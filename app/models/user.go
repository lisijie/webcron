package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/lisijie/webcron/app/libs"
	"errors"
	"fmt"
	"gopkg.in/ldap.v2"
)

const (
	AUTH_LCOAL = 0
	AUTH_LDAP  = 1
	STATUS_NORMAL = 0
	STATUS_FORBIDDEN = 1
	ROLE_MANAGER = 1
	ROLE_GUEST = 0
)
type User struct {
	Id        int    `orm:"pk;auto;unique;column(id)" json:"id"`
	Account   string `orm:"column(account);size(64);unique" json:"account"`    //账户
	Password  string `orm:"column(password);size(32)" json:"password"`         //密码
	Salt      string `orm:"column(salt);size(10)" json:"salt"`                 //密码盐
	UserName  string `orm:"column(user_name);size(20)" json:"user_name"`       //昵称
	Email     string `orm:"column(email);size(50)" json:"email"`               //邮箱
	LastLogin int64  `orm:"column(last_login);type(bigint)" json:"last_login"` //最后登录时间
	LastIp    string `orm:"column(last_ip);size(15)" json:"last_ip"`           //最后登录IP
	Status    int    `orm:"column(status);type(int);default(0)" json:"status"` //状态，0正常 1禁用
	Auth      int    `orm:"column(auth);type(int);default(0) json:"auth"`      //0 Local 1  LDAP
	Role      int    `orm:"column(role);type(int);default(0) json:"role"`      //角色，0 Guest  1 Manager
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

func UserGetByAccount(account string) (*User, error) {
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("account", account).One(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserUpdate(user *User, fields ...string) error {
	_, err := orm.NewOrm().Update(user, fields...)
	return err
}


func UserGetList(page, pageSize int) ([]*User, int64) {
	offset := (page - 1) * pageSize

	list := make([]*User, 0)

	query := orm.NewOrm().QueryTable(TableName("user"))
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)

	return list, total
}

func UserDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("user")).Filter("id", id).Delete()
	return err
}


func CheckPassword(account ,password string)(*User, error){
	u := new(User)

	err := orm.NewOrm().QueryTable(TableName("user")).Filter("account", account).One(u)
	if err == nil {
		if u.Status == STATUS_FORBIDDEN {
			return u,errors.New("该帐号已禁用")
		}
		if u.Auth == AUTH_LCOAL{
			if u.Password == libs.Md5([]byte(password+u.Salt)) {
				return u, nil
			}else{
				return u,errors.New("帐号或密码错误")
			}
		}
	}
	err = u.LDAPLogin(account,password)
	if err != nil {
		return  u, err
	}
	return u,nil
}


func (u *User) LDAPLogin(account,password string) (error) {

	if beego.AppConfig.DefaultBool("ldap.enable", false) == false {
		return errors.New("没有启用LDAP服务器")
	}
	addr := fmt.Sprintf("%s:%d", beego.AppConfig.String("ldap.host"), beego.AppConfig.DefaultInt("ldap.port", 3268))
	lc, err := ldap.Dial("tcp", addr)
	if err != nil {
		return errors.New("无法连接到LDAP服务器")
	}
	defer lc.Close()
	err = lc.Bind(beego.AppConfig.String("ldap.user"), beego.AppConfig.String("ldap.password"))
	if err != nil {
		return errors.New("第一次LDAP绑定失败")
	}
	searchRequest := ldap.NewSearchRequest(
		beego.AppConfig.String("ldap.base"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		//修改objectClass通过配置文件获取值
		fmt.Sprintf("(&(%s)(%s=%s))", beego.AppConfig.String("ldap.filter"), beego.AppConfig.String("ldap.attribute"), account),
		[]string{"dn", "mail", "displayName"},
		nil,
	)
	searchResult, err := lc.Search(searchRequest)
	if err != nil {
		return errors.New("LDAP搜索失败")
	}
	if len(searchResult.Entries) != 1 {
		return errors.New("LDAP用户不存在或者多于一个")
	}
	userdn := searchResult.Entries[0].DN
	err = lc.Bind(userdn, password)
	if err != nil {
		return errors.New("用户密码错误")
	}
	if u.Id == 0 {
		u.Email = searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap.mail"))
		u.UserName = searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap.name"))
		u.Account = account
		u.Salt = string(utils.RandomCreateBytes(10))
		u.Password = libs.Md5([]byte(password + u.Salt))
		u.Status = STATUS_NORMAL
		u.Auth = AUTH_LDAP
		u.Role = ROLE_GUEST
		_,err := UserAdd(u)
		if err != nil {
			return errors.New(fmt.Sprint("自动注册LDAP用户错误:%s",err.Error()))
		}
	}else{
		u.Email = searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap.mail"))
		u.UserName = searchResult.Entries[0].GetAttributeValue(beego.AppConfig.String("ldap.name"))
		u.Update()
	}
	return nil
}
