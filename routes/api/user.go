package api

import (
	"github.com/gin-gonic/gin"
	"service/app/controllers/restapi"
	"service/app/middlewares"
)

func NewUserApi(r *gin.RouterGroup, rest *restapi.Restapi, mid *middlewares.Middlewares) {
	api := r.Group("/user")
	noAuth := api.Group("/data")
	{
		noAuth.GET("", rest.UserHandler.List)

		noAuthUser := noAuth.Group("/user")
		{
			noAuthUser.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello World",
				})
			})
		}
	}
}
