package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	//"highloadcup/travels/config"
	//"highloadcup/travels/services/loggerService"
	//"net/http"
)

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request

		c.Set("requestId", uuid.NewV4().String())

		/*		var req config.Request

				if c.Request.Method == http.MethodPost {
					err := c.ShouldBindJSON(&req)
					if err != nil {
						loggerService.GetMainLogger().Error(c, err.Error())

						c.JSON(http.StatusOK, map[string]string{"error": err.Error()})

						c.Abort()

						return
					}

					c.Set(config.KeyMeta, req.Meta)
					c.Set(config.KeyRequest, req.Data)
					loggerService.GetMainLogger().Info(c, string(req.Meta))
					loggerService.GetMainLogger().Info(c, string(req.Data))
				}*/

		c.Next()
	}
}
