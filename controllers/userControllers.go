package controllers


import (
	"awesomeProject/initializer"
	"awesomeProject/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func UserHome(c *gin.Context) {
	user, _ := c.Get("user")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"message": user,
	})
}
func UserLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"message": "ghh",
	})
}

func UserRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"message": "register success",
	})
	return
}
func Logout(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	tokenString, err := c.Cookie("store")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("store", tokenString, -1, "", "", false, true)
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())

}
func UserRegisterDb(c *gin.Context) {
	var data struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Status   string
	}
	if err := c.Bind(&data); err != nil {
		fmt.Println("hhhhhhhhhhhhh")
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "invalid input please check inputs",
		})
		return
	}
	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user := models.User{Email: data.Email, Password: string(hash), Status: "unblocked"}

	result := initializer.DB.Create(&user) //pass the data to create
	if result.Error != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "failed to create already registered",
		})
		return
	}
	//respond
	c.HTML(http.StatusOK, "login.html", gin.H{
		"data": "this is login page",
	})
	return
}

func UserAuth(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	// Get the email and password from req body

	var Data struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if c.Bind(&Data) != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid Inputs Please Check Inputs",
		})
		return
	}

	fmt.Println(Data)

	// Look up requested user

	var user models.User

	// It equals to : SELECT * FROM users WHERE email = requested email;
	initializer.DB.First(&user, "Email = ?", Data.Email)

	if user.ID == 0 {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "invalid user name or password ",
		})
		return
	}

	// Compare sent password with user password hash

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Data.Password))
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "invalid user name or password ",
		})
		return
	}

	// Create token
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
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

	fmt.Println(user.Status)

	// set cookie and Send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("store", tokenString, 3600*24*30, "", "", false, true)

	if user.Status == "blocked" {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"error": "you are blocked by admin",
		})

	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"content": "This is an index page...",
		"token":   tokenString,
		"message": Data,
	})

}
func Getloginsubmit(c *gin.Context) {
	c.Header("Cache-Control", "no-cache,no-store,must-revalidate")
	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())

}
