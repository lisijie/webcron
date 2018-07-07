package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lifei6671/mindoc/conf"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Init() {
	fmt.Println("Initializing...")

	adapter := beego.AppConfig.String("db.adapter")
	if adapter == "mysql" {
		dbhost := beego.AppConfig.String("db.host")
		dbport := beego.AppConfig.String("db.port")
		dbuser := beego.AppConfig.String("db.user")
		dbpassword := beego.AppConfig.String("db.password")
		dbname := beego.AppConfig.String("db.name")
		timezone := beego.AppConfig.String("db.timezone")
		if dbport == "" {
			dbport = "3306"
		}
		dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
		if timezone != "" {
			dsn = dsn + "&loc=" + url.QueryEscape(timezone)
		}
		err := orm.RegisterDataBase("default", "mysql", dsn)

		if err != nil {
			beego.Error("注册默认数据库失败:", err)
		}
	} else if adapter == "sqlite3" {
		orm.DefaultTimeLoc = time.UTC
		database := beego.AppConfig.String("db.database")
		if strings.HasPrefix(database, "./") {
			database = filepath.Join(conf.WorkingDirectory, string(database[1:]))
		}

		dbPath := filepath.Dir(database)
		os.MkdirAll(dbPath, 0777)

		err := orm.RegisterDataBase("default", "sqlite3", database)

		if err != nil {
			beego.Error("注册默认数据库失败:", err)
		}
	} else {
		beego.Error("不支持的数据库类型.")
		os.Exit(1)
	}

	orm.RegisterModel(new(User), new(Task), new(TaskGroup), new(TaskLog))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err.Error())
		os.Exit(1)
	}

	_, err = UserGetByName("admin")
	if err == orm.ErrNoRows {
		usr := new(User)
		usr.UserName = "admin"
		usr.Email = "admin@example.com"
		usr.Password = "7fef6171469e80d32c0559f88b377245"
		usr.Salt = ""
		usr.Status = 0

		_, err := UserAdd(usr)
		if err != nil {
			panic("User.Add => " + err.Error())
			os.Exit(0)
		}

	}

	fmt.Println("Initializing Successfully!")
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
