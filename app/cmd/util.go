// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/kirsle/configdir"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/util"
)

var (
	_initialized = false
	_buildInfo   *app.BuildInfo
	_flags       = &Flags{}
	_logger      util.Logger
)

type Flags struct {
	Address   string
	Port      uint16
	ConfigDir string
	Debug     bool
}

func (f *Flags) Addr() string {
	if f.Port == 0 {
		return f.Address
	}
	return fmt.Sprintf("%s:%d", f.Address, f.Port)
}

func (f *Flags) Clone(port uint16) *Flags {
	return &Flags{
		Address:   f.Address,
		Port:      port,
		ConfigDir: f.ConfigDir,
		Debug:     f.Debug,
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
	if _flags.ConfigDir == "" {
		_flags.ConfigDir = configdir.LocalConfig(util.AppName)
		_ = configdir.MakePath(_flags.ConfigDir)
	}
	err := util.InitAcronyms()
	if err != nil {
		return err
	}
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	_logger = l
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

var (
	maxHeaderSize = 1024 * 256
	maxBodySize   = 1024 * 1024 * 128
)

func serve(name string, listener net.Listener, h fasthttp.RequestHandler) error {
	x := &fasthttp.Server{Handler: h, Name: name, ReadBufferSize: maxHeaderSize, NoDefaultServerHeader: true, MaxRequestBodySize: maxBodySize}
	if err := x.Serve(listener); err != nil {
		return errors.Wrap(err, "unable to run http server")
	}
	return nil
}

func listenandserve(name string, addr string, port uint16, h fasthttp.RequestHandler) (uint16, error) {
	p, l, err := listen(addr, port)
	if err != nil {
		return p, err
	}
	err = serve(name, l, h)
	if err != nil {
		return p, err
	}
	return 0, nil
}
