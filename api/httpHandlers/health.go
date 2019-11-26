package httpHandlers

import (
	"github.com/gin-gonic/gin"
)

// health-check
func Health(c *gin.Context) {
	resp := map[string]string{"status": "ok"}
	c.JSON(200, resp)
	c.Abort()
}
