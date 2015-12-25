package models

import (
	"github.com/astaxie/beego/orm"
)

type Task struct {
	Id         int
	UserId     int
	GroupId    int
	TaskName   string
	CronSpec   string
	Concurrent int
	Command    string
	CreateTime int64
}

func (t *Task) TableName() string {
	return TableName("task")
}

func (t *Task) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskAdd(task *Task) (int64, error) {
	return orm.NewOrm().Insert(task)
}

func TaskGetList(page, pageSize int) ([]*Task, int64) {
	offset := (page - 1) * pageSize

	tasks := make([]*Task, 0)

	query := orm.NewOrm().QueryTable(TableName("task"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&tasks)

	return tasks, total
}

func TaskGetById(id int) (*Task, error) {
	task := &Task{
		Id: id,
	}

	err := orm.NewOrm().Read(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func TaskDel(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task")).Filter("id", id).Delete()
	return err
}
