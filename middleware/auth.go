package middleware

import (
	"net/http"
	"userapi/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token is not provided."})
		return
	}

	username, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Token is not authorized."})
		return
	}

	context.Set("username", username)
	context.Next()
}
