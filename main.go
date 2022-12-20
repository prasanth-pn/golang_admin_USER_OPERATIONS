package main

import (
	"awesomeProject/controllers"
	"awesomeProject/initializer"
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectDatabase()
	initializer.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*.html")

	r.GET("/", middleware.UserAuth, controllers.UserHome)
	r.GET("/login", controllers.UserLogin)
	r.GET("/register", controllers.UserRegister)
	r.POST("/register-submit", controllers.UserRegisterDb)
	r.POST("/login-submit", controllers.UserAuth)
	r.GET("/login-submit", controllers.Getloginsubmit)
	r.GET("/logout", controllers.Logout)
	//admin
	r.GET("/show-users", middleware.AdminAuth, controllers.AdminShowUser)
	r.GET("/admin", controllers.AdminHome)
	r.GET("/admin-logout", controllers.AdminLogout)
	r.GET("/admin-login", controllers.AdminLogin)
	r.POST("/admin-submit", controllers.AdminLoginSubmit)
	r.GET("/blockuser/:id", controllers.Adminuserblock)
	r.GET("/unblock/:id", controllers.Unblock)
	r.GET("delete/:id", controllers.Delete)
	r.GET("/blocked-user", middleware.AdminAuth, controllers.AdminShowBlockedUser)
	//sdfsdfplsmjkl
	r.Run()
}
