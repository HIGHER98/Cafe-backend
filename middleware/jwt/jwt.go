package jwt

import (
	"net/http"

	"cafe/pkg/e"
	"cafe/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = 20002
				default:
					code = 20001
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
