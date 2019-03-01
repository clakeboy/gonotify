package router

import (
	"github.com/gin-gonic/gin"
	"gonotify/controllers"
)

func GetController(controller_name string, c *gin.Context) interface{} {
	switch controller_name {
	case "server":
		return controllers.NewServerManagementController(c)
	case "task":
		return controllers.NewTaskController(c)
	case "log":
		return controllers.NewTaskLogController(c)
	case "account":
		return controllers.NewAccountController(c)
	case "login":
		return controllers.NewLoginController(c)
	default:
		return nil
	}
}
