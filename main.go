package main

import (
	"github.com/astaxie/beego"
	"./app/controllers"
	"./app/jobs"
	_ "./app/mail"
	"./app/models"
	"html/template"
	"net/http"
	"fmt"
	"os/exec"
	"os"
	"path/filepath"
)

const VERSION = "1.0.0"

func main() {
	models.Init()
	jobs.InitJobs()
	fmt.Println(">>>>>>>>:",beego.BConfig.WebConfig.ViewsPath);
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
	}
	beego.AppConfig.Set("version", VERSION)
	//设置视图模板路径
	beego.SetViewsPath(GetAPPRootPath()+"/views/")
	//设置静态资源路径
	beego.SetStaticPath("/static",GetAPPRootPath()+"/static")

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
func GetAPPRootPath()string{
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}
