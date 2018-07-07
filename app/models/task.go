package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	TASK_SUCCESS = 0  // 任务执行成功
	TASK_ERROR   = -1 // 任务执行出错
	TASK_TIMEOUT = -2 // 任务执行超时
)

type Task struct {
	Id           int    `orm:"pk;auto;unique;column(id)" json:"id"`
	UserId       int    `orm:"column(user_id);type(int);default(0);index" json:"user_id"`       //用户ID
	GroupId      int    `orm:"column(group_id);type(int);default(0);index" json:"group_id"`     //分组ID
	TaskName     string `orm:"column(task_name);size(50)" json:"task_name"`                     //任务名称
	TaskType     int    `orm:"column(task_type);type(int);default(0)" json:"task_type"`         //任务类型
	Description  string `orm:"column(description);size(255)" json:"description"`                //任务描述
	CronSpec     string `orm:"column(cron_spec);size(100)" json:"cron_spec"`                    //时间表达式
	Concurrent   int    `orm:"column(concurrent);type(int);default(0)" json:"concurrent"`       //是否只允许一个实例
	Command      string `orm:"column(command);type(text)" json:"command"`                       //命令详情
	Status       int    `orm:"column(status);type(int);default(0)" json:"status"`               //0停用 1启用
	Notify       int    `orm:"column(notify);type(int);default(0)" json:"notify"`               //通知设置
	NotifyEmail  string `orm:"column(notify_email);size(2000)" json:"notify_email"`             //通知人列表
	Timeout      int    `orm:"column(timeout);type(int);default(0)" json:"timeout"`             //超时设置
	ExecuteTimes int    `orm:"column(execute_times);type(int);default(0)" json:"execute_times"` //累计执行次数
	PrevTime     int64  `orm:"column(prev_time);type(bigint)" json:"prev_time"`                 //上次执行时间
	CreateTime   int64  `orm:"column(create_time);type(bigint)" json:"create_time"`
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
	if task.TaskName == "" {
		return 0, fmt.Errorf("TaskName字段不能为空")
	}
	if task.CronSpec == "" {
		return 0, fmt.Errorf("CronSpec字段不能为空")
	}
	if task.Command == "" {
		return 0, fmt.Errorf("Command字段不能为空")
	}
	if task.CreateTime == 0 {
		task.CreateTime = time.Now().Unix()
	}
	return orm.NewOrm().Insert(task)
}

func TaskGetList(page, pageSize int, filters ...interface{}) ([]*Task, int64) {
	offset := (page - 1) * pageSize

	tasks := make([]*Task, 0)

	query := orm.NewOrm().QueryTable(TableName("task"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&tasks)

	return tasks, total
}

func TaskResetGroupId(groupId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("task")).Filter("group_id", groupId).Update(orm.Params{
		"group_id": 0,
	})
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
