// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

// Lib starts the application as a library, returning the actual TCP port the server is listening on (as an int32 to make interop easier).
func Lib(path string) int32 {
	if _buildInfo == nil {
		_buildInfo = &app.BuildInfo{Version: "0.2.17"}
	}
	f := &Flags{Address: "0.0.0.0", Port: 0, ConfigDir: path}

	if err := initIfNeeded(); err != nil {
		panic(errors.WithStack(errors.Wrap(err, "error initializing application")))
	}

	r, logger, err := loadServer(f, _logger)
	if err != nil {
		panic(errors.WithStack(err))
	}

	port, listener, err := listen(f.Address, f.Port)
	if err != nil {
		panic(errors.WithStack(err))
	}

	logger.Infof("%s library started on port [%d]", util.AppName, port)

	go func() {
		e := serve(util.AppKey, listener, r)
		if e != nil {
			panic(errors.WithStack(e))
		}
	}()

	return int32(port)
}
