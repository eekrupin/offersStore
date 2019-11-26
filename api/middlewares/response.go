package middlewares

import (
	//"encoding/json"
	"github.com/gin-gonic/gin"
	//"highloadcup/travels/config"
	//"net/http"
	//"runtime/debug"
	//"strings"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// after request

		/*		if c.Request.Method == http.MethodPost {
				meta := c.MustGet(config.KeyMeta).(json.RawMessage)
				result, _ := json.Marshal(c.MustGet(config.KeyResponse))
				var response interface{}

				isError := len(c.Errors) > 0
				if isError {
					stack := strings.Split(string(debug.Stack()), "\n")
					response = config.ResponseError{
						Success:  0,
						Envelope: config.Envelope{Meta: meta},
						Error: config.RError{
							Message:    result,
							StackTrace: stack,
						},
					}
				} else {
					response = config.ResponseSuccess{
						Success:  1,
						Envelope: config.Envelope{Meta: meta},
						Data:     result,
					}
				}

				c.JSON(http.StatusOK, response)
			}*/
	}
}
