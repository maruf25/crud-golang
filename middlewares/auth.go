package middlewares

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/models"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(allowedRole []models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				c.JSON(404, gin.H{
					"message": "Error no cookies",
				})
				c.Abort()
				return
			}
		}

		tokenString := tokenCookie.Value
		claim := &config.JWTClaim{}

		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
			return config.SecretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthotized",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthotized",
			})
			c.Abort()
			return
		}

		if allow := slices.Contains(allowedRole, models.Role(claim.Role)); !allow {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unauthotized",
			})
			c.Abort()
			return
		}

		c.Set("email", claim.Email)
		c.Set("userId", claim.UserId)

		c.Next()
	}
}
