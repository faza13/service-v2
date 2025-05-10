package api

import (
	"github.com/gin-gonic/gin"
	"service/app/controllers/restapi"
	"service/app/middlewares"
)

func NewPermissionApi(r *gin.RouterGroup, restApi *restapi.Restapi, mid *middlewares.Middlewares) {
	api := r.Group("/permission")
	noAuth := api.Group("/noauth")
	{
		noAuth.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "ini permission",
			})
		})

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
