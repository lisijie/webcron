package main

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/controllers"
	"github.com/lisijie/webcron/app/models"
)

func main() {
	models.Init()

	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")

	beego.AutoRouter(&controllers.TaskController{})

	beego.SessionOn = true
	beego.Run()
}
