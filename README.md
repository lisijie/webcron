# webcron
------------

一个定时任务管理器，基于Go语言和beego框架开发。用于统一管理项目中的定时任务，提供可视化配置界面、执行日志记录、邮件通知等功能，无需依赖*unix下的crontab服务。

## 项目背景

开发此项目是为了解决本人所在公司的PHP项目中定时任务繁多，使用crontab不好管理的问题。我所在项目的定时任务也是PHP编写的，属于整个项目的一部分，我希望能有一个系统可以统一配置这些定时任务，并且可以查看每次任务的执行情况，任务执行完成或失败能够自动邮件提醒开发人员，因此做了这个项目。

## 功能特点

* 统一管理多种定时任务。
* 秒级定时器，使用crontab的时间表达式。
* 可随时暂停任务。
* 记录每次任务的执行结果。
* 执行结果邮件通知。

## 界面截图

![webcron](https://raw.githubusercontent.com/lisijie/webcron/master/screenshot.png)


## 安装说明

系统需要安装Go和MySQL。

获取源码

	$ go get github.com/lisijie/webcron
	
打开配置文件 conf/app.conf，修改相关配置。
	

创建数据库webcron，再导入install.sql

	$ mysql -u username -p -D webcron < install.sql

运行
	
	$ ./webcron
	或
	$ nohup ./webcron 2>&1 > error.log &
	设为后台运行

访问： 

http://localhost:8000

帐号：admin
密码：admin888