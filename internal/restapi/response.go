package restapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIResponse struct {
	Error   int32  `json:"error"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func ResponseOkRequest(message string, c *gin.Context) {
	resp := APIResponse{
		Code:    http.StatusOK,
		Error:   0,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseBadRequest(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusBadRequest,
		Error:   1,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseInternalserverError(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusInternalServerError,
		Error:   1,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseStatusNotFound(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusNotFound,
		Error:   1,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}
