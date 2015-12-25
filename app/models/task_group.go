package models

import (
	"github.com/astaxie/beego/orm"
)

type TaskGroup struct {
	Id         int
	UserId     int
	GroupName  string
	CreateTime int64
}

func (t *TaskGroup) TableName() string {
	return TableName("task_group")
}

func TaskGroupAdd(obj *TaskGroup) (int64, error) {
	return orm.NewOrm().Insert(obj)
}

func TaskGroupGetById(id int) (*TaskGroup, error) {
	obj := &TaskGroup{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func TaskGroupDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("task_group")).Filter("id", id).Delete()
	return err
}
