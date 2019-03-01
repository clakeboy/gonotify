package controllers

import (
	"fmt"
	"github.com/clakeboy/golib/httputils"
	"github.com/gin-gonic/gin"
	"gonotify/queue"
)

type ServerManagementController struct {
	c *gin.Context
}

func NewServerManagementController(con *gin.Context) *ServerManagementController {
	return &ServerManagementController{c: con}
}

//得到服务状态
func (s *ServerManagementController) ActionStatus() {
	str := fmt.Sprintf("%s", queue.Notify.Status())
	s.c.String(200, str)
}

func (s *ServerManagementController) ActionSessionSet() {
	val := s.c.DefaultQuery("s", "this is default session string")
	session := s.c.MustGet("session").(*httputils.HttpSession)
	session.Set("clake", val)
	s.c.String(200, "session set done")
}

func (s *ServerManagementController) ActionNotify(args []byte) {
	fmt.Println(string(args))
	s.c.String(200, "success")
}
