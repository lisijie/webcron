# webcron

一个使用Go语言和beego框架开发的定时任务管理器。用于统一管理项目中的定时任务，提供可视化配置界面、执行日志记录、邮件通知等功能。


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