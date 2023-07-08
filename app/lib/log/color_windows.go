package log

import (
	"os"

	"golang.org/x/sys/windows"
)

var _ = func() any {
	println("blergh start")
	m := &master{}
	m.initStdios()
	println("blergh stop")
	return nil
}()

type master struct {
	in     windows.Handle
	inMode uint32

	out     windows.Handle
	outMode uint32

	err     windows.Handle
	errMode uint32
}

func (m *master) initStdios() {
	m.in = windows.Handle(os.Stdin.Fd())
	if err := windows.GetConsoleMode(m.in, &m.inMode); err == nil {
		_ = windows.SetConsoleMode(m.in, m.inMode|windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
	}

	m.out = windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(m.out, &m.outMode); err == nil {
		if err := windows.SetConsoleMode(m.out, m.outMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil {
			m.outMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		} else {
			_ = windows.SetConsoleMode(m.out, m.outMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		}
	}

	m.err = windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(m.err, &m.errMode); err == nil {
		if err := windows.SetConsoleMode(m.err, m.errMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil {
			m.errMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		} else {
			_ = windows.SetConsoleMode(m.err, m.errMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		}
	}
}
