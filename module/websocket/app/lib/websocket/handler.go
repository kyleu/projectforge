package websocket

import (
	"context"

	"{{{ .Package }}}/app/util"
)

type Handler func(ctx context.Context, s *Service, conn *Connection, svc string, cmd string, param []byte, logger util.Logger) error

func EchoHandler(_ context.Context, s *Service, conn *Connection, svc string, cmd string, param []byte, logger util.Logger) error {
	logger.Debugf("echoing [%s] message command [%s] from connection [%s]: %s", svc, cmd, conn.ID, string(param))
	return s.WriteMessage(conn.ID, &Message{Channel: svc, From: &conn.ID, Cmd: cmd, Param: param}, logger)
}

func LogHandler(_ context.Context, _ *Service, conn *Connection, svc string, cmd string, param []byte, logger util.Logger) error {
	logger.Debugf("received [%s] message command [%s] from connection [%s]: %s", svc, cmd, conn.ID, string(param))
	return nil
}
