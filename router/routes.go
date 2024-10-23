package router

import (
	"userapi/handler"
	"userapi/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {

	authenticated := route.Group("/")
	authenticated.Use(middleware.Authenticate)

	//signup
	route.POST("/signup", handler.UserSignUp)

	//login
	route.POST("/login", handler.Userlogin)

	//add userinfo
	authenticated.POST("/adduser", handler.AddUser)

	//get all userinfo
	route.GET("/getuser", handler.GetAllUser)

	//get user by id
	route.GET("/getuser/:id", handler.GetUserById)

	//update user by id
	authenticated.POST("/updateuser/:id", handler.UpdateUser)

	//delete user by id
	authenticated.DELETE("/deleteuser/:id", handler.DeleteUser)
}
