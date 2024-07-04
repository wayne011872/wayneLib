package wayneLib

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErr "github.com/wayne011872/api-toolkit/errors"
	"github.com/wayne011872/log"
)

var (
	ServerErrorHandler = func(c *gin.Context, service string, err error) {
		if err == nil {
			return
		}

		l := log.GetByGinCtx(c)
		if l != nil {
			l.WarnPkg(err)
		}
		if apiErr, ok := err.(apiErr.ApiError); ok {
			c.AbortWithStatusJSON(apiErr.GetStatus(),
				map[string]interface{}{
					"status":  apiErr.GetStatus(),
					"error":   apiErr.Error(),
					"service": service,
				})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				map[string]interface{}{
					"status":  http.StatusInternalServerError,
					"title":   err.Error(),
					"service": service,
				})
		}
	}
)
