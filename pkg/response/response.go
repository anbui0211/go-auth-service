package response

import (
	"github.com/gin-gonic/gin"
)

type ResponseMessageData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RMiddlewareError(c *gin.Context, code int, msg string) {
	c.JSON(code, ResponseMessageData{
		Code:    code,
		Message: msg,
	})
	c.Abort()
}
