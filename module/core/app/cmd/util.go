package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kirsle/configdir"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/util"
)

var (
	_initialized = false
	_buildInfo   *app.BuildInfo
	_flags       = &Flags{}
)

type Flags struct {
	Address    string
	Port       uint16
	ConfigDir  string
	WorkingDir string
	Debug      bool
}

func (f *Flags) Addr() string {
	if f.Port == 0 {
		return f.Address
	}
	return fmt.Sprintf("%s:%d", f.Address, f.Port)
}

func (f *Flags) Clone(port uint16) *Flags {
	return &Flags{
		Address:    f.Address,
		Port:       port,
		ConfigDir:  f.ConfigDir,
		WorkingDir: f.WorkingDir,
		Debug:      f.Debug,
	}
}

var initMu sync.Mutex

func initIfNeeded() error {
	initMu.Lock()
	defer initMu.Unlock()

	if _initialized {
		return nil
	}
	if _buildInfo == nil {
		return errors.New("no build info")
	}
	if _flags.WorkingDir != "" && _flags.WorkingDir != "." {
		if err := os.Chdir(_flags.WorkingDir); err != nil {
			return errors.Wrapf(err, "failed to change working directory to [%s]", _flags.WorkingDir)
		}
	}
	if _flags.ConfigDir == "" {
		if x := util.GetEnv(util.AppKey + "_working_directory"); x != "" {
			_flags.ConfigDir = x
		} else {
			_flags.ConfigDir = configdir.LocalConfig(util.AppName)
		}
	}
	err := util.InitAcronyms({{{ .Acronyms }}})
	if err != nil {
		return err
	}
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	util.ConfigDir = _flags.ConfigDir
	util.DEBUG = _flags.Debug
	util.RootLogger = l
	_initialized = true
	return nil
}

func listen(address string, port uint16) (uint16, net.Listener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return port, nil, errors.Wrapf(err, "unable to listen on port [%d]", port)
	}
	if port == 0 {
		addr := l.Addr().String()
		_, portStr := util.StringSplitLast(addr, ':', true)
		actualPort, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return 0, nil, errors.Wrapf(err, "invalid port [%s]", portStr)
		}
		port = uint16(actualPort)
	}
	return port, l, nil
}

var maxHeaderSize = 1024 * 256

func serve(listener net.Listener, h http.Handler) error {
	x := &http.Server{Handler: h, MaxHeaderBytes: maxHeaderSize, ReadHeaderTimeout: time.Minute}
	if err := x.Serve(listener); err != nil {
		return errors.Wrap(err, "unable to run http server")
	}
	return nil
}

func listenandserve(addr string, port uint16, h http.Handler) (uint16, error) {
	p, l, err := listen(addr, port)
	if err != nil {
		return p, err
	}
	err = serve(l, h)
	if err != nil {
		return p, err
	}
	return 0, nil
}
