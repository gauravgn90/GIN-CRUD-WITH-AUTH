package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Custom logger middleware

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s %d bytes \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.BodySize,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
