package middleware

import (
	"awesomeProject/initializer"
	"awesomeProject/models"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func UserAuth(c *gin.Context) {

	// Get the cookie of the request

	tokenString, err := c.Cookie("store")
	if err != nil {
		c.HTML(http.StatusOK, "home.html", gin.H{})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validation

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && token != nil {

		// check the exp

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// find the user with token subject

		var user models.User

		initializer.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			location := url.URL{Path: "/login"}
			c.Redirect(http.StatusFound, location.RequestURI())

			c.AbortWithStatus(http.StatusUnauthorized)
			return

		}

		if user.Status == "blocked" {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"error": "Admin Blocked You",
			})
			//return
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// atach to request

		c.Set("user", user)

		// continue

		c.Next()
		// fmt.Println(claims["exp"], claims["sub"])
	} else {
		c.HTML(http.StatusOK, "login.html", gin.H{})

		//c.AbortWithStatus(http.StatusUnauthorized)
	}

}
