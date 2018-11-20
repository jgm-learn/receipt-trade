package main

import (
	"github.com/astaxie/beego"
	_ "receipt-trade/web/routers"
)

func main() {
	beego.Run()
}
