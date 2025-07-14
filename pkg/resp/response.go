package resp

import (
	"github.com/gin-gonic/gin"
	"github.com/teakingwang/gin-demo/pkg/errs"
	"net/http"
)

type HTTPResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    errs.CodeSuccess,
		Message: errs.Message(errs.CodeSuccess),
		Data:    data,
	})
}

func WriteError(c *gin.Context, appErr *errs.AppError) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
		Data:    nil,
	})
}
