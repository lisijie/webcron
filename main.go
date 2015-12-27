package main

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/controllers"
	"github.com/lisijie/webcron/app/jobs"
	_ "github.com/lisijie/webcron/app/mail"
	"github.com/lisijie/webcron/app/models"
)

func main() {
	models.Init()
	jobs.InitJobs()

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/help", &controllers.HelpController{}, "*:Index")

	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.GroupController{})

	beego.SessionOn = true
	beego.Run()
}
