package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"github.com/lisijie/webcron/app/jobs"
	"github.com/lisijie/webcron/app/libs"
	"github.com/lisijie/webcron/app/models"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MainController struct {
	BaseController
}

// 首页
func (this *MainController) Index() {
	this.Data["pageTitle"] = "系统概况"

	// 即将执行的任务
	entries := jobs.GetEntries(30)
	jobList := make([]map[string]interface{}, len(entries))
	for k, v := range entries {
		row := make(map[string]interface{})
		job := v.Job.(*jobs.Job)
		row["task_id"] = job.GetId()
		row["task_name"] = job.GetName()
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
	}

	// 最近执行的日志
	logs, _ := models.TaskLogGetList(1, 20)

	// 批量查询task得到taskMap
	taskIds := make([]int, 0)
	for _, v := range logs {
		taskIds = append(taskIds, v.TaskId)
	}
	taskIds = models.MakeIntSliceUnique(taskIds)
	tasks,_ := models.TaskGetList(1, 100, "id__in", taskIds)
	var taskMap = map[int]*models.Task{}
	for _, task := range tasks {
		taskMap[task.Id] = task
	}

	recentLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, ok := taskMap[v.TaskId]
		taskName := ""
		if ok {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["output"] = beego.Substr(v.Output, 0, 100)
		row["status"] = v.Status
		recentLogs[k] = row
	}

	// 最近执行失败的日志
	logs, _ = models.TaskLogGetList(1, 20, "status__lt", 0)

	// 批量查询task得到taskMap
	taskIds = make([]int, 0)
	for _, v := range logs {
		taskIds = append(taskIds, v.TaskId)
	}
	taskIds = models.MakeIntSliceUnique(taskIds)
	tasks,_ = models.TaskGetList(1, 100, "id__in", taskIds)
	taskMap = map[int]*models.Task{}
	for _, task := range tasks {
		taskMap[task.Id] = task
	}

	errLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, ok := taskMap[v.TaskId]
		taskName := ""
		if ok {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = beego.Substr(v.Error, 0, 100)
		row["status"] = v.Status
		errLogs[k] = row
	}

	this.Data["recentLogs"] = recentLogs
	this.Data["errLogs"] = errLogs
	this.Data["jobs"] = jobList
	this.Data["cpuNum"] = runtime.NumCPU()
	this.display()
}

// 个人信息
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user, _ := models.UserGetById(this.userId)

	if this.isPost() {
		flash := beego.NewFlash()
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("密码长度必须大于6位")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("两次输入的密码不一致")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		flash.Success("修改成功！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "个人信息"
	this.Data["user"] = user
	this.display()
}

// 登录
func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix()
				models.UserUpdate(user)

				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}

				this.redirect(beego.URLFor("TaskController.List"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}

	this.TplName = "main/login.html"
}

// 退出登录
func (this *MainController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.redirect(beego.URLFor("MainController.Login"))
}

// 获取系统时间
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}
