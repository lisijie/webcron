package models

import (
	"github.com/astaxie/beego/orm"
)

type TaskLog struct {
	Id          int    `orm:"pk;auto;unique;column(id)" json:"id"`
	TaskId      int    `orm:"column(task_id);type(int);default(0);index" json:"task_id"`     //任务ID
	Output      string `orm:"column(output);type(text)" json:"output"`                       //任务输出
	Error       string `orm:"column(error);type(text)" json:"error"`                         //错误信息
	Status      int    `orm:"column(status);type(int);default(0)" json:"status"`             //状态
	ProcessTime int    `orm:"column(process_time);type(int);default(0)" json:"process_time"` //消耗时间/毫秒
	CreateTime  int64  `orm:"column(create_time);type(bigint)" json:"create_time"`
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

func TaskLogGetErrorListTop(page, pageSize int) ([]*TaskLog, int64) {
	offset := (page - 1) * pageSize - 1

	var logs []*TaskLog
	sql := "SELECT * FROM "+TableName("task_log") +
		" WHERE id IN (SELECT MAX(id) FROM  "+TableName("task_log") + " GROUP BY task_id) " +
		" AND status = -1 " +
		" ORDER BY id DESC LIMIT ?,?"
	total,_ := orm.NewOrm().Raw(sql,offset,pageSize).QueryRows(&logs)

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
