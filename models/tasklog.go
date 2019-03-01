package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/clakeboy/golib/ckdb"
)

var logdb *storm.DB

type TaskLogData struct {
	Id          int    `storm:"id,increment" json:"id"` //主键,自增长
	TaskId      int    `json:"task_id" storm:"index"`   //任务ID
	TaskName    string `json:"task_name"`               //任务名称
	ExecReceive string `json:"exec_receive"`            //执行通知后的回复信息 success 为成功
	ExecError   string `json:"exec_error"`              //执行时的错误信息
	ExecTime    int    `json:"exec_time"`               //执行时间
	CreatedDate int    `json:"created_date"`            //创建时间
}

//表名
type TaskLogModel struct {
	Table string `json:"table"` //表名
	storm.Node
}

func NewTaskLogModel(db *storm.DB) *TaskLogModel {
	if db == nil {
		if logdb == nil {
			logdb, _ = storm.Open("./db/tasklog.db")
		}
		db = logdb
	}

	return &TaskLogModel{
		Table: "task_log",
		Node:  db.From("task_log"),
	}
}

//通过ID拿到记录
func (t *TaskLogModel) GetById(id int) (*TaskLogData, error) {
	data := &TaskLogData{}
	err := t.One("id", id, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//查询条件得到任务数据列表
func (t *TaskLogModel) Query(page, number int, where ...q.Matcher) (*ckdb.QueryResult, error) {
	var list []TaskLogData
	count, err := t.Select(where...).Count(new(TaskLogData))
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
func (t *TaskLogModel) List(page, number int, where ...q.Matcher) ([]TaskLogData, error) {
	var list []TaskLogData
	err := t.Select(where...).Limit(number).Skip((page - 1) * number).Reverse().Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
