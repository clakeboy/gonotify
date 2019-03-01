package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
	"gonotify/common"
)

//任务
type TaskData struct {
	Id             int    `json:"id" storm:"id,increment"`  //主键,自增长
	TaskName       string `json:"task_name" storm:"unique"` //任务名称
	TimeRule       string `json:"time_rule"`                //任务时间规则
	Once           bool   `json:"once"`                     //是否只执行一次
	IsExecute      bool   `json:"is_execute"`               //是否已经执行,Once 为 true 才有效
	Disable        bool   `json:"disable"`                  //是否禁用,BOOL值,每次都必须更新,不然就会还原为false
	NotifyUrl      string `json:"notify_url"`               //通知地址
	NotifyMethod   string `json:"notify_method"`            //通知方法 GET,POST,JSON
	NotifyData     string `json:"notify_data"`              //通知内容, 为JSON字符串
	NotifyNumber   int    `json:"notify_number"`            //通知次数
	NotifiedNumber int    `json:"notified_number"`          //已通知次数
	Source         string `json:"source"`                   //任务来源
	CreatedDate    int    `json:"created_date"`             //创建时间 unix 时间截
}

//表名
type TaskModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewTaskModel(db *storm.DB) *TaskModel {
	if db == nil {
		db = common.BDB
	}

	return &TaskModel{
		Table: "task",
		Node:  db.From("task"),
	}
}

//通过ID拿到记录
func (t *TaskModel) GetById(id int) (*TaskData, error) {
	data := &TaskData{}
	err := t.One("Id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//通过Name 拿到记录
func (t *TaskModel) GetByName(name string) (*TaskData, error) {
	data := &TaskData{}
	err := t.One("TaskName", name, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//查询条件得到任务数据列表
func (t *TaskModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []TaskData
	count, err := t.Select(where...).Count(new(TaskData))
	if err != nil {
		return nil, err
	}
	err = t.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return &ckdb.QueryResult{
		Count: count,
		List:  list,
	}, nil
}

//查询条件得到任务数据列表
func (t *TaskModel) List(page, number int, where ...q.Matcher) ([]TaskData, error) {
	var list []TaskData
	err := t.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//删除记录
func (t *TaskModel) Delete(where ...q.Matcher) error {
	return t.Select(where...).Delete(new(TaskData))
}
