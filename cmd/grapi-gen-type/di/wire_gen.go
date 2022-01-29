// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/protoc"
)

// Injectors from wire.go:

func NewApp(command *gencmd.Command) (*App, error) {
	ctx := gencmd.ProvideCtx(command)
	grapicmdCtx := gencmd.ProvideGrapiCtx(ctx)
	config := grapicmd.ProvideProtocConfig(grapicmdCtx)
	fs := grapicmd.ProvideFS(grapicmdCtx)
	executor := grapicmd.ProvideExec(grapicmdCtx)
	io := grapicmd.ProvideIO(grapicmdCtx)
	ui := cli.UIInstance(io)
	rootDir := grapicmd.ProvideRootDir(grapicmdCtx)
	gexConfig := protoc.ProvideGexConfig(fs, executor, io, rootDir)
	repository, err := protoc.ProvideToolRepository(gexConfig)
	if err != nil {
		return nil, err
	}
	wrapper := protoc.NewWrapper(config, fs, executor, ui, repository, rootDir)
	app := &App{
		Protoc: wrapper,
	}
	return app, nil
}
