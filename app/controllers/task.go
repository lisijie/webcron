package controllers

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/cron"
	"github.com/lisijie/webcron/app/libs"
	"github.com/lisijie/webcron/app/models"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	BaseController
}

// 任务列表
func (this *TaskController) List() {
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	result, count := models.TaskGetList(page, this.pageSize)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.TaskName
		row["cron_spec"] = v.CronSpec

		e := cron.GetEntryById(v.Id)
		if e != nil {
			row["next_time"] = beego.Date(e.Next, "Y-m-d H:i:s")
			row["prev_time"] = beego.Date(e.Prev, "Y-m-d H:i:s")
			row["running"] = true
		} else {
			row["next_time"] = "-"
			row["prev_time"] = "-"
			row["running"] = false
		}
		list[k] = row
	}

	this.Data["pageTitle"] = "任务列表"
	this.Data["list"] = list
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.UrlFor("TaskController.List"), true).ToString()
	this.display()
}

// 添加任务
func (this *TaskController) Add() {

	if this.isPost() {
		task := new(models.Task)
		task.UserId = this.userId
		task.GroupId = 0
		task.TaskName = this.GetString("task_name")
		task.Concurrent, _ = this.GetInt("concurrent")
		task.CronSpec = this.GetString("cron_spec")
		task.Command = this.GetString("command")

		if _, err := models.TaskAdd(task); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
	}

	this.display()
}

// 编辑任务
func (this *TaskController) Edit() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		task.TaskName = strings.TrimSpace(this.GetString("task_name"))
		task.Concurrent, _ = this.GetInt("concurrent")
		task.CronSpec = strings.TrimSpace(this.GetString("cron_spec"))
		task.Command = strings.TrimSpace(this.GetString("command"))
		if err := task.Update(); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		this.ajaxMsg("", MSG_OK)
	}

	this.Data["task"] = task
	this.display()
}

// 任务执行日志列表
func (this *TaskController) Logs() {
	taskId, _ := this.GetInt("id")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	task, err := models.TaskGetById(taskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	result, count := models.TaskLogGetList(task.Id, page, this.pageSize)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["status"] = v.Status
		list[k] = row
	}

	this.Data["pageTitle"] = "任务执行日志"
	this.Data["list"] = list
	this.Data["task"] = task
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.UrlFor("TaskController.Logs", "id", taskId), true).ToString()
	this.display()
}

// 查看日志详情
func (this *TaskController) ViewLog() {
	id, _ := this.GetInt("id")

	taskLog, err := models.TaskLogGetById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	task, err := models.TaskGetById(taskLog.TaskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	data := make(map[string]interface{})
	data["id"] = taskLog.Id
	data["output"] = taskLog.Output
	data["error"] = taskLog.Error
	data["start_time"] = beego.Date(time.Unix(taskLog.CreateTime, 0), "Y-m-d H:i:s")
	data["process_time"] = float64(taskLog.ProcessTime) / 1000
	data["ouput_size"] = libs.SizeFormat(float64(len(taskLog.Output)))
	data["status"] = taskLog.Status

	this.Data["task"] = task
	this.Data["data"] = data
	this.display()
}

// 批量操作日志
func (this *TaskController) LogBatch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			models.TaskLogDelById(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 批量操作
func (this *TaskController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "active":
			if task, err := models.TaskGetById(id); err == nil {
				job, err := cron.NewJobFromTask(task)
				if err == nil {
					cron.AddJob(task.CronSpec, job)
				}
			}
		case "pause":
			cron.RemoveJob(id)
		case "delete":
			models.TaskDel(id)
		}
	}

	this.ajaxMsg("", MSG_OK)
}

// 启动任务
func (this *TaskController) Start() {
	id, _ := this.GetInt("id")

	task, err := models.TaskGetById(id)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	job, err := cron.NewJobFromTask(task)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	cron.AddJob(task.CronSpec, job)

	this.ajaxMsg("", MSG_OK)
}

// 暂停任务
func (this *TaskController) Pause() {
	id, _ := this.GetInt("id")

	cron.RemoveJob(id)

	this.ajaxMsg("", MSG_OK)
}
