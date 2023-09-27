// Package telemetry - Content managed by Project Forge, see [projectforge.md] for details.
package telemetry

import (
	"context"
	"io"

	"projectforge.dev/projectforge/app/util"
)

func RunProcess(ctx context.Context, cmd string, path string, in io.Reader, out io.Writer, er io.Writer, logger util.Logger) (int, error) {
	exec, _ := util.StringSplit(cmd, ' ', true)
	_, span, _ := StartSpan(ctx, "process:"+exec, logger)
	defer span.Complete()
	span.Attribute("cmd", cmd)
	span.Attribute("path", path)
	return util.RunProcess(cmd, path, in, out, er)
}

func RunProcessSimple(ctx context.Context, cmd string, path string, logger util.Logger) (int, string, error) {
	exec, _ := util.StringSplit(cmd, ' ', true)
	_, span, _ := StartSpan(ctx, "process-simple:"+exec, logger)
	defer span.Complete()
	span.Attribute("cmd", cmd)
	span.Attribute("path", path)
	return util.RunProcessSimple(cmd, path)
}
