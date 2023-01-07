// $PF_IGNORE$
package routes

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/graph"
)

func graphqlRoutes(as *app.State, r *router.Router, logger util.Logger) {
	r.GET("/graphql", controller.GraphiQL)
	r.POST("/graphql", graph.Auth(fasthttpadaptor.NewFastHTTPHandler(as.Services.GQL)))
}
