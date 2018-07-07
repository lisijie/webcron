package models

import (
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id        int    `orm:"pk;auto;unique;column(id)" json:"id"`
	UserName  string `orm:"column(user_name);size(20)" json:"user_name"`       //用户名
	Password  string `orm:"column(password);size(32)" json:"password"`         //密码
	Salt      string `orm:"column(salt);size(10)" json:"salt"`                 //密码盐
	Email     string `orm:"column(email);size(50)" json:"email"`               //邮箱
	LastLogin int64  `orm:"column(last_login);type(bigint)" json:"last_login"` //最后登录时间
	LastIp    string `orm:"column(last_ip);size(15)" json:"last_ip"`           //最后登录IP
	Status    int    `orm:"column(status);type(int);default(0)" json:"status"` //状态，0正常 -1禁用
}

func (u *User) TableName() string {
	return TableName("user")
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func UserAdd(user *User) (int64, error) {
	return orm.NewOrm().Insert(user)
}

func UserGetById(id int) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserGetByName(userName string) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("user_name", userName).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}
