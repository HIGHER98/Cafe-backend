package jwt

import (
	"encoding/json"
	"net/http"

	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Token struct {
	Token string `json:"token"`
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var t Token
		code := e.SUCCESS
		data, err := c.GetRawData()
		if err != nil {
			logging.Warn(err)
		}
		json.Unmarshal(data, &t)
		token := t.Token
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
