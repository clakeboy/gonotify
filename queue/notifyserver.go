package queue

import (
	"fmt"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/components/task"
	"github.com/clakeboy/golib/utils"
	"gonotify/models"
	"time"
)

const (
	NotifyGet  = "GET"
	NotifyPost = "POST"
	NotifyJson = "JSON"
)

var Notify *NotifyServer

type NotifyServer struct {
	task *task.Management
}

func NewNotifyServer() *NotifyServer {
	return &NotifyServer{
		task: task.NewManagement(),
	}
}

//开始服务
func (n *NotifyServer) Start() {
	if n.task.Status() == "Running" {
		return
	}
	model := models.NewTaskModel(nil)
	//所有循环任务
	model.Select(q.Eq("Once", false), q.Eq("Disable", false)).Each(new(models.TaskData), func(row interface{}) error {
		n.AddNotify(row.(*models.TaskData))
		return nil
	})
	//所有未执行成功的单次任务
	model.Select(q.Eq("Once", true), q.Eq("IsExecute", false)).Each(new(models.TaskData), func(row interface{}) error {
		n.AddNotify(row.(*models.TaskData))
		return nil
	})

	n.task.Start()
}

//结束服务
func (n *NotifyServer) Stop() {
	n.task.Stop()
	n.task.ClearTask()
}

//服务状态
func (n *NotifyServer) Status() string {
	return n.task.Status()
}

//添加一个通知任务
func (n *NotifyServer) AddNotify(data *models.TaskData) {
	n.task.AddTaskString(data.TimeRule, n.execNotify, n.callbackNotify, data)
}

//删除一个通知任务
func (n *NotifyServer) RemoveNotify(data *models.TaskData) {
	n.task.RemoveForeach(func(item *task.Item) bool {
		itemData := item.Args[0].(*models.TaskData)
		return itemData.Id == data.Id
	})
}

func (n *NotifyServer) execNotify(item *task.Item) bool {
	taskData, ok := item.Args[0].(*models.TaskData)
	if !ok {
		return false
	}

	httpClient := utils.NewHttpClient()
	httpClient.SetTimeout(time.Second * 5)
	var res []byte
	var err error
	switch taskData.NotifyMethod {
	case NotifyGet:
		res, err = httpClient.Get(taskData.NotifyUrl)
	case NotifyPost:
		postData := utils.M{}
		if taskData.NotifyData != "" {
			err = postData.ParseJsonString(taskData.NotifyData)
			if err != nil {
				break
			}
		}
		res, err = httpClient.Post(taskData.NotifyUrl, postData)
	case NotifyJson:
		res, err = httpClient.PostJsonString(taskData.NotifyUrl, taskData.NotifyData)
	default:
		return false
	}

	if string(res) == "success" {
		if taskData.Once {
			taskData.IsExecute = true
		}
	}

	taskLog := &models.TaskLogData{
		TaskId:      taskData.Id,
		TaskName:    taskData.TaskName,
		ExecReceive: string(res),
		ExecError:   fmt.Sprintf("%v", err),
		ExecTime:    int(time.Now().Unix()),
		CreatedDate: int(time.Now().Unix()),
	}

	logModel := models.NewTaskLogModel(nil)
	logModel.Save(taskLog)

	taskData.NotifiedNumber++

	return taskData.Once
}

func (n *NotifyServer) callbackNotify(item *task.Item) {
	taskData, ok := item.Args[0].(*models.TaskData)
	if !ok {
		return
	}

	//如果只执行一次,并且已经执行功能就删除这个任务
	if taskData.Once && taskData.IsExecute {
		n.task.RemoveTask(item)
	} else if taskData.Once {
		//未执行成功
		//如果执行次数已经达到最大执行次数也删除这个任务执行
		if taskData.NotifiedNumber >= taskData.NotifyNumber {
			n.task.RemoveTask(item)
			taskData.IsExecute = true
		}
	}

	model := models.NewTaskModel(nil)
	model.Update(taskData)
}
