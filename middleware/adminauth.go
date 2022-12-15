package middleware

import (
	"awesomeProject/initializer"
	"awesomeProject/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

func AdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("admin")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//decode validation
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpexted signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil

	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token != nil {
		//check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var admin models.Admin
		initializer.DB.First(&admin, claims["sub"])
		if admin.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("admin", admin)

		//continue

		c.Next()

	} else {
		c.HTML(http.StatusOK, "adminpanel.html", gin.H{})
	}
}
