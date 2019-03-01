package controllers

import (
	"encoding/json"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"github.com/clakeboy/golib/utils"
	"github.com/gin-gonic/gin"
	"gonotify/models"
	"gonotify/queue"
	"time"
)

//控制器
type TaskController struct {
	c *gin.Context
}

func NewTaskController(c *gin.Context) *TaskController {
	return &TaskController{c: c}
}

//查询
func (m *TaskController) ActionQuery(args []byte) (*ckdb.QueryResult, error) {
	var params struct {
		TaskName string `json:"task_name"`
		Once     string `json:"once"`
		Page     int    `json:"page"`
		Number   int    `json:"number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	var where []q.Matcher
	if params.TaskName != "" {
		where = append(where, q.Re("TaskName", params.TaskName))
	}

	if params.Once != "all" {
		once := utils.YN(params.Once == "loop", false, true).(bool)
		where = append(where, q.Eq("Once", once))
	}

	model := models.NewTaskModel(nil)
	res, err := model.Query(params.Page, params.Number, where...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//添加
func (m *TaskController) ActionInsert(args []byte) error {
	var params struct {
		TaskName     string `json:"task_name"`
		TimeRule     string `json:"time_rule"`
		Once         bool   `json:"once"`
		Disable      bool   `json:"disable"`
		NotifyUrl    string `json:"notify_url"`
		NotifyMethod string `json:"notify_method"`
		NotifyData   string `json:"notify_data"`
		NotifyNumber int    `json:"notify_number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	model := models.NewTaskModel(nil)

	data := &models.TaskData{
		TaskName:     params.TaskName,
		TimeRule:     params.TimeRule,
		Once:         params.Once,
		NotifyUrl:    params.NotifyUrl,
		NotifyMethod: params.NotifyMethod,
		NotifyData:   params.NotifyData,
		NotifyNumber: params.NotifyNumber,
		Source:       "System",
		CreatedDate:  int(time.Now().Unix()),
	}

	err = model.Save(data)
	if err != nil {
		return utils.Error("添加任务出错!", err)
	}
	//添加到任务执行管理
	//service.Notify.AddNotify(data)
	if !data.Disable {
		queue.Notify.AddNotify(data)
	}
	return nil
}

//删除
func (m *TaskController) ActionDelete(args []byte) error {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	model := models.NewTaskModel(nil)
	return model.Delete(q.Eq("Id", params.Id))
}

//使用ID得到一条记录
func (m *TaskController) ActionFind(args []byte) (*models.TaskData, error) {
	var params struct {
		Id int `json:"id"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return nil, err
	}

	model := models.NewTaskModel(nil)
	data, err := model.GetById(params.Id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//修改
func (m *TaskController) ActionUpdate(args []byte) error {
	var params struct {
		Id           int    `json:"id"`
		TaskName     string `json:"task_name"`
		TimeRule     string `json:"time_rule"`
		Once         bool   `json:"once"`
		Disable      bool   `json:"disable"`
		NotifyUrl    string `json:"notify_url"`
		NotifyMethod string `json:"notify_method"`
		NotifyData   string `json:"notify_data"`
		NotifyNumber int    `json:"notify_number"`
	}

	err := json.Unmarshal(args, &params)
	if err != nil {
		return err
	}

	model := models.NewTaskModel(nil)
	data, err := model.GetById(params.Id)
	if err != nil {
		return err
	}

	data.TaskName = params.TaskName
	data.TimeRule = params.TimeRule
	data.Once = params.Once
	data.Disable = params.Disable
	data.NotifyUrl = params.NotifyUrl
	data.NotifyMethod = params.NotifyMethod
	data.NotifyData = params.NotifyData
	data.NotifyNumber = params.NotifyNumber

	err = model.Update(data)
	if err != nil {
		return err
	}
	queue.Notify.RemoveNotify(data)
	if !data.Disable {
		queue.Notify.AddNotify(data)
	}

	return nil
}
