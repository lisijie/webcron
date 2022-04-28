package main

import (
	"github.com/astaxie/beego"
	"html/template"
	"net/http"
	"os"
	"webcron/app/controllers"
	"webcron/app/jobs"
	_ "webcron/app/mail"
	"webcron/app/models"
)

const VERSION = "1.1.1"

func main() {
	_, e := os.Stat("/srun3/etc/srun4-webcron.conf")
	if e == nil || os.IsNotExist(e) {
		_ = beego.LoadAppConfig("ini", "/srun3/etc/srun4-webcron.conf")
	}

	models.Init()
	jobs.InitJobs()

	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})

	// 生产环境不输出debug日志
	if beego.AppConfig.String("runmode") == "prod" {
		beego.SetLevel(beego.LevelInformational)
		beego.SetViewsPath("/srun3/www/srun4-webcron/views")
		beego.SetStaticPath("/static", "/srun3/www/srun4-webcron/static")
	}
	beego.AppConfig.Set("version", VERSION)

	// 路由设置
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/gettime", &controllers.MainController{}, "*:GetTime")
	beego.Router("/help", &controllers.HelpController{}, "*:Index")
	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.GroupController{})

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
