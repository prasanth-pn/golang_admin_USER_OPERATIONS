package controllers

import (
	"awesomeProject/initializer"
	"awesomeProject/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/url"
	"os"
	"time"
)

func AdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminlogin.html", gin.H{
		"content": "This is an index page...",
	})
}

func AdminHome(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"content": "this is an index page",
		"message": user,
	})
}

func AdminLoginSubmit(c *gin.Context) {
	// Get the email and password from req body

	var body struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required,min=5"`
	}
	fmt.Println(c.Bind(&body))
	if c.Bind(&body) != nil {
		c.HTML(http.StatusBadRequest, "adminlogin.html", gin.H{
			"error": "Invalid Inputs Please Check Inputsqq",
		})
		return
	}

	// Look up reqested user

	var user models.Admin

	// fmt.Print("\n\n email :", body.Email, "\npassword :", body.Password, "\n\n")

	// It equals to : SELECT * FROM users WHERE email = requested email;
	initializer.DB.First(&user, "username = ?", body.Email)

	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "adminlogin.html", gin.H{
			"error": "invalid user name or password e ",
		})
		return
	}
	if user.Password != body.Password {
		c.HTML(http.StatusBadRequest, "adminlogin.html", gin.H{
			"error": "invalid user name or password  p",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		fmt.Println(err)
		return
	}

	// set cookie and Send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin", tokenString, 3600*24*30, "", "", false, true)

	location := url.URL{Path: "/show-users"}
	c.Redirect(http.StatusFound, location.RequestURI())

}
func AdminLogout(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	tokenString, err := c.Cookie("admin")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin", tokenString, -1, "", "", false, true)
	location := url.URL{Path: "/admin"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func AdminShowUser(c *gin.Context) {
	//user, _ := c.Get("user")

	var users []models.User
	fmt.Println(users)

	// Get all records
	initializer.DB.Order("email").Where("status = ?", "unblocked").Find(&users)
	fmt.Println(users, "hhhhhhh")
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	//	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "adminpanel.html", gin.H{
		"content": "This is an index page...",
		"users":   users,
	})
}
func AdminShowBlockedUser(c *gin.Context) {
	//user, _ := c.Get("user")

	var users []models.User
	fmt.Println(users)

	// Get all records
	initializer.DB.Order("email").Where("status = ?", "blocked").Find(&users)
	fmt.Println(users, "hhhhhhh")
	// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

	c.HTML(http.StatusOK, "adminblockedusers.html", gin.H{
		"content": "This is an index page...",
		"users":   users,
	})

}

func Adminuserblock(c *gin.Context) {
	userid := c.Param("id")
	fmt.Println(userid)
	var user models.User
	initializer.DB.Model(&user).Where("id=?", userid).Update("Status", "blocked")
	location := url.URL{Path: "/show-users"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
func Unblock(c *gin.Context) {
	userid := c.Param("id")
	fmt.Println(userid)
	var user models.User
	initializer.DB.Model(&user).Where("id=?", userid).Update("Status", "unblocked")
	location := url.URL{Path: "/blocked-user"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
func Delete(c *gin.Context) {
	userid := c.Param("id")
	var user models.User
	initializer.DB.Unscoped().Where("id = ?", userid).Delete(&user)
	// DELETE FROM users WHERE id = 10;
	//initializer.DB.Model(&user).Where("id=?", userid).Update("Status", "unblocked")
	location := url.URL{Path: "/show-users"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
