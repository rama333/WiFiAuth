package restapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ResponseBadRequest(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseInternalserverError(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}

func ResponseStatusNotFound(message string, c *gin.Context) {

	resp := APIResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}

	c.JSON(http.StatusOK, resp)
}
