package models

import (
	"fmt"
	"github.com/astaxie/beego"
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

func TaskLogGetList(page, pageSize int, filters ...interface{}) ([]*TaskLog, int64) {
	offset := (page - 1) * pageSize

	logs := make([]*TaskLog, 0)

	query := orm.NewOrm().QueryTable(TableName("task_log"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

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

func TaskLogDelByTaskId(taskId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task_log")).Filter("task_id", taskId).Delete()
}

func TaskLogCount(taskId int) int64 {
	count, err := orm.NewOrm().QueryTable(TableName("task_log")).Filter("task_id", taskId).Count()
	if err != nil {
		beego.Error(fmt.Printf("查询失败原因 %s", err.Error()))
		return 0
	}
	return count
}
