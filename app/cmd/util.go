package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/kirsle/configdir"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/log"
	"github.com/kyleu/projectforge/app/util"
)

var (
	_initialized = false
	_buildInfo   *app.BuildInfo
	_flags       = &Flags{}
	_logger      *zap.SugaredLogger
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

func initIfNeeded() error {
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
		return port, nil, errors.Wrap(err, fmt.Sprintf("unable to listen on port [%d]", port))
	}
	if port == 0 {
		addr := l.Addr().String()
		_, portStr := util.SplitStringLast(addr, ':', true)
		actualPort, err := strconv.Atoi(portStr)
		if err != nil {
			return 0, nil, errors.Wrap(err, "invalid port ["+portStr+"]")
		}
		port = uint16(actualPort)
	}
	return port, l, nil
}

func serve(name string, listener net.Listener, h fasthttp.RequestHandler) error {
	x := &fasthttp.Server{Handler: h, Name: name, ReadBufferSize: 65536, NoDefaultServerHeader: true}
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
