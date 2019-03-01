package controllers

import (
	"encoding/json"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"github.com/gin-gonic/gin"
	"gonotify/models"
)

//控制器
type TaskLogController struct {
	c *gin.Context
}

func NewTaskLogController(c *gin.Context) *TaskLogController {
	return &TaskLogController{c: c}
}

//查询
func (m *TaskLogController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		TaskId int `json:"task_id"`
		Page   int `json:"page"`
		Number int `json:"number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	var where []q.Matcher
	if params.TaskId != 0 {
		where = append(where, q.Eq("TaskId", params.TaskId))
	}

	model := models.NewTaskLogModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}
