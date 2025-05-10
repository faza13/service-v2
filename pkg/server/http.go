package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service/app/controllers/restapi"
	"service/app/middlewares"
	"service/config"
	"service/pkg/otel"
)

type groupsHandlers func(group *gin.RouterGroup, restApi *restapi.Restapi, mid *middlewares.Middlewares)

func RunHTTPServer(
	ctx context.Context,
	cfg *config.Router,
	tracer trace.Tracer,
	restApi *restapi.Restapi,
	mid *middlewares.Middlewares,
	groupsHandlers ...groupsHandlers,
) {
	r := gin.New()

	//tracing otel middleware
	if tracer != nil {
		r.Use(func(c *gin.Context) {
			otelCtx := otel.InjectTracing(c, tracer, "")
			name := fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.URL.Path)
			ctxStart, span := otel.AddSpan(otelCtx, name)
			defer span.End()
			span.SetAttributes(attribute.KeyValue{
				Key:   "tes",
				Value: attribute.Value{},
			})

			c.Request = c.Request.WithContext(ctxStart)

			c.Next()
		})
	}

	setupMiddlewares(r)

	group := r.Group(cfg.Prefix)
	for i := range groupsHandlers {
		groupsHandlers[i](group, restApi, mid)
	}

	s := http.Server{
		Addr:                         fmt.Sprintf(":%s", cfg.Port),
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
		//BaseContext: func(net.Listener) context.Context {
		//	return nil
		//},
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func setupMiddlewares(route *gin.Engine) {

}

func guidMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
