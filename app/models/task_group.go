package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type TaskGroup struct {
	Id          int    `orm:"pk;auto;unique;column(id)" json:"id"`
	UserId      int    `orm:"column(user_id);type(int);default(0);index" json:"user_id"` //用户ID
	GroupName   string `orm:"column(group_name);size(50)" json:"group_name"`             //组名
	Description string `orm:"column(description);size(255)" json:"description"`          //说明
	CreateTime  int64  `orm:"column(create_time);type(bigint)" json:"create_time"`
}

func (t *TaskGroup) TableName() string {
	return TableName("task_group")
}

func (t *TaskGroup) Update(fields ...string) error {
	if t.GroupName == "" {
		return fmt.Errorf("组名不能为空")
	}
	if _, err := orm.NewOrm().Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func TaskGroupAdd(obj *TaskGroup) (int64, error) {
	if obj.GroupName == "" {
		return 0, fmt.Errorf("组名不能为空")
	}
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

func TaskGroupGetList(page, pageSize int) ([]*TaskGroup, int64) {
	offset := (page - 1) * pageSize

	list := make([]*TaskGroup, 0)

	query := orm.NewOrm().QueryTable(TableName("task_group"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}
