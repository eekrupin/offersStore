package httpHandlers

import (
	"fmt"
	"github.com/eekrupin/offersStore/modules"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sync"
	"time"
)

var loadLocker sync.Mutex

func LoadShareFile(c *gin.Context) {
	defer c.Abort()
	LoadSource(c, modules.FileToCassandra)
}

func LoadUrlFile(c *gin.Context) {
	defer c.Abort()
	LoadSource(c, modules.HttpFileToCassandra)
}

func LoadSource(c *gin.Context, loader func(url string) (offers int, err error)) {
	loadLocker.Lock()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in LoadSource", r)
			resp := map[string]string{"status": "error: " + fmt.Sprint(r)}
			c.JSON(404, resp)
		}
	}()
	defer loadLocker.Unlock()

	source, ok := RequestStringParam(c, "source")
	if !ok || source == "" {
		c.JSON(400, map[string]string{"status": "param 'source' must be fill"})
		return
	}
	start := time.Now()
	offers, err := loader(source)
	if err != nil {
		c.JSON(404, map[string]string{"status": "error: " + err.Error()})
		return
	}
	resp := map[string]string{"status": "ok", "offers": strconv.Itoa(offers), "durationSeconds": fmt.Sprintf("%f", time.Since(start).Seconds())}
	c.JSON(200, resp)
	return
}
