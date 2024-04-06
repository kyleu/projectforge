package app

import (
	"context"

	{{{ .ServicesImports }}}
)

{{{ .ServicesDefinition }}}

//nolint:revive
func initCoreServices(ctx context.Context, st *State{{{ if .HasModule "audit" }}}, auditSvc *audit.Service{{{ end }}}, logger util.Logger) CoreServices {
	return {{{ .ServicesConstructor }}}
}
