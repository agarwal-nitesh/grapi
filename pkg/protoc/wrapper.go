package protoc

import (
	"context"
	"os"
	"path/filepath"

	"github.com/izumin5210/clicontrib/pkg/clog"
	"github.com/izumin5210/gex/pkg/tool"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/utils/exec"

	"github.com/izumin5210/grapi/pkg/clui"
	"github.com/izumin5210/grapi/pkg/grapicmd/util/fs"
)

type Wrapper interface {
	Exec(context.Context) error
}

type wrapperImpl struct {
	cfg      *Config
	fs       afero.Fs
	ui       clui.UI
	execer   exec.Interface
	toolRepo tool.Repository
	rootDir  string
	binDir   string
}

func NewWrapper(cfg *Config, fs afero.Fs, execer exec.Interface, ui clui.UI, toolRepo tool.Repository, rootDir, binDir string) Wrapper {
	return &wrapperImpl{
		cfg:      cfg,
		fs:       fs,
		ui:       ui,
		execer:   execer,
		toolRepo: toolRepo,
		rootDir:  rootDir,
		binDir:   binDir,
	}
}

func (e *wrapperImpl) Exec(ctx context.Context) (err error) {
	e.ui.Section("Execute protoc")

	e.ui.Subsection("Install plugins")
	err = errors.WithStack(e.installPlugins(ctx))
	if err != nil {
		return
	}

	e.ui.Subsection("Execute protoc")
	err = errors.WithStack(e.execProtocAll(ctx))

	return
}

func (e *wrapperImpl) installPlugins(ctx context.Context) error {
	return errors.WithStack(e.toolRepo.BuildAll(ctx))
}

func (e *wrapperImpl) execProtocAll(ctx context.Context) error {
	protoFiles, err := e.cfg.ProtoFiles(e.fs, e.rootDir)
	if err != nil {
		return errors.WithStack(err)
	}

	var errs []error
	for _, path := range protoFiles {
		err = e.execProtoc(ctx, path)
		relPath, _ := filepath.Rel(e.rootDir, path)
		if err == nil {
			e.ui.ItemSuccess(relPath)
		} else {
			errs = append(errs, err)
			e.ui.ItemFailure(relPath)
		}
	}

	if len(errs) > 0 {
		clog.Error("failed to execute protoc", "errors", errs)
		return errors.New("failed to execute protoc")
	}

	return nil
}

func (e *wrapperImpl) execProtoc(ctx context.Context, protoPath string) error {
	outDir, err := e.cfg.OutDirOf(e.rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = fs.CreateDirIfNotExists(e.fs, outDir); err != nil {
		return errors.WithStack(err)
	}

	cmds, err := e.cfg.Commands(e.rootDir, protoPath)
	if err != nil {
		return errors.WithStack(err)
	}

	path := e.binDir + string(filepath.ListSeparator) + os.Getenv("PATH")
	env := append(os.Environ(), "PATH="+path)

	for _, args := range cmds {
		cmd := e.execer.CommandContext(ctx, args[0], args[1:]...)
		cmd.SetEnv(env)
		cmd.SetDir(e.rootDir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return errors.Wrapf(err, "failed to execute command: %v\n%s", args, string(out))
		}
	}

	return nil
}
