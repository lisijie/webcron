package models

import (
	"github.com/astaxie/beego/orm"
)

type TaskLog struct {
	Id          int
	TaskId      int
	Output      string
	Error       string
	Status      int
	ProcessTime int
	CreateTime  int64
}

func (t *TaskLog) TableName() string {
	return TableName("task_log")
}

func TaskLogAdd(t *TaskLog) (int64, error) {
	return orm.NewOrm().Insert(t)
}

func TaskLogGetList(taskId, page, pageSize int) ([]*TaskLog, int64) {
	offset := (page - 1) * pageSize

	logs := make([]*TaskLog, 0)

	query := orm.NewOrm().QueryTable(TableName("task_log"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&logs)

	return logs, total
}

func TaskLogGetById(id int) (*TaskLog, error) {
	obj := &TaskLog{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func TaskLogDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task_log")).Filter("id", id).Delete()
	return err
}
