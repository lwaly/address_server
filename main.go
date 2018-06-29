package main

import (
	_ "AddressServer/routers"
	"fmt"

	"AddressServer/controllers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	if 0 != controllers.Init() {
		fmt.Printf("fail to init")
		return
	}
	beego.Run()
}
