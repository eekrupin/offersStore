package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func TracerMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan(c.Request.URL.Path)
		defer span.Finish()

		c.Set("span", span)

		c.Next()
	}
}
