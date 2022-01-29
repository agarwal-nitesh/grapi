// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package testing

import (
	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/izumin5210/grapi/pkg/svcgen"
)

// Injectors from wire.go:

func NewTestApp(command *gencmd.Command, wrapper protoc.Wrapper, ui cli.UI) (*svcgen.App, error) {
	ctx := gencmd.ProvideCtx(command)
	grapicmdCtx := gencmd.ProvideGrapiCtx(ctx)
	rootDir := grapicmd.ProvideRootDir(grapicmdCtx)
	config := grapicmd.ProvideProtocConfig(grapicmdCtx)
	grapicmdConfig := grapicmd.ProvideConfig(grapicmdCtx)
	builder := svcgen.ProvideParamsBuilder(rootDir, config, grapicmdConfig)
	app := &svcgen.App{
		ProtocWrapper: wrapper,
		ParamsBuilder: builder,
	}
	return app, nil
}
