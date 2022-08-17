package main

import (
	"UAVdemo/web"
	"UAVdemo/web/cliInit"
	"UAVdemo/web/controller"
	"UAVdemo/web/dao"

	//"UAVdemo/web/dao"
)

func main() {

	//DB Conn
	dao.InitMysql()

	//Web
	app := controller.Application{
		cliInit.CliInit(),
	}

	defer cliInit.SDK.Close()

	web.WebStart(&app)
}
