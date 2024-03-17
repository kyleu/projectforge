package app

import (
	"context"

	{{{ .ServicesImports }}}
)

{{{ .ServicesDefinition }}}

func initCoreServices(ctx context.Context, st *State{{{ if .HasModule "audit" }}}, auditSvc *audit.Service{{{ end }}}, logger util.Logger) CoreServices {
	return {{{ .ServicesConstructor }}}
}
