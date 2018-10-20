package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/libs"
	"github.com/lisijie/webcron/app/models"
	"strconv"
	"strings"
	"github.com/astaxie/beego/utils"
	"time"
)

type UserController struct {
	BaseController
}

func (this *UserController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	result, count := models.UserGetList(page, this.pageSize)


	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["Id"] = v.Id
		row["Account"] = v.Account
		row["UserName"] = v.UserName
		row["Email"] = v.Email
		row["LastLogin"] = beego.Date(time.Unix(v.LastLogin, 0), "Y-m-d H:i:s")
		row["LastIp"] = v.LastIp
		row["Status"] = v.Status
		row["Auth"] = v.Auth
		row["Role"] = v.Role
		list[k] = row
	}



	this.Data["pageTitle"] = "用户列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("UserController.List"), true).ToString()
	this.display()
}

func (this *UserController) Add() {
	if this.isPost() {
		user := new(models.User)
		user.Id = 0
		user.Role, _ = this.GetInt("role")
		user.Auth = models.AUTH_LCOAL
		user.Status, _ = this.GetInt("status")
		user.Account = strings.TrimSpace(this.GetString("account"))
		user.UserName = strings.TrimSpace(this.GetString("username"))
		user.Email = strings.TrimSpace(this.GetString("email"))
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if len(password1) < 6 {
			this.ajaxMsg("密码长度必须大于6位", MSG_ERR)
		} else if password2 != password1 {
			this.ajaxMsg("两次输入的密码不一致", MSG_ERR)
		} else {
			user.Salt = string(utils.RandomCreateBytes(10))
			user.Password = libs.Md5([]byte(password1 + user.Salt))
		}
		_, err := models.UserAdd(user)
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "添加用户"
	this.display()
}

func (this *UserController) Edit() {
	id, _ := this.GetInt("id")

	user, err := models.UserGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		if user.Id > 1 {
			user.Role, _ = this.GetInt("role")
			user.Status, _ = this.GetInt("status")
		}
		if user.Auth == models.AUTH_LCOAL {
			user.UserName = strings.TrimSpace(this.GetString("username"))
			user.Email = strings.TrimSpace(this.GetString("email"))
		}
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				this.ajaxMsg("密码长度必须大于6位", MSG_ERR)
			} else if password2 != password1 {
				this.ajaxMsg("两次输入的密码不一致", MSG_ERR)
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
			}
		}
		err := user.Update()
		if err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}
		this.ajaxMsg("", MSG_OK)
	}

	this.Data["pageTitle"] = "编辑用户"
	this.Data["pageUser"] = user
	this.display()
}

func (this *UserController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的用户", MSG_ERR)
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			models.UserDelById(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}
