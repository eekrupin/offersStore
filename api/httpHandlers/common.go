package httpHandlers

import (
	"github.com/gin-gonic/gin"
)

func RequestStringParam(c *gin.Context, param string) (request string, ok bool) {
	values := c.Request.URL.Query()
	_, ok = values[param]
	str := values.Get(param)
	if ok && str == "" {
		c.Keys["err"] = true
		return "", false
	} else {
		return str, ok
	}
}
