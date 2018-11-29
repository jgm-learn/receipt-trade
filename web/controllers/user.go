package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	this.TplName = "user.html"
}
func (this *UserController) GetUserId() {
	user := new(models.User)
	userAddr := this.GetString("userAddr")
	user.PublicKey = userAddr

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	this.Data["json"] = &user
	this.ServeJSON()
}

func (this *UserController) GetUserNonce() {
	user := new(models.User)
	userAddr := this.GetString("userAddr")
	user.PublicKey = userAddr

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	this.Data["json"] = &user
	this.ServeJSON()
}

func (this *UserController) GetFunds() {
	//orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30)
	//o := orm.NewOrm()

	user := new(models.User)
	pK := this.GetString("pucKey")
	user.PublicKey = pK

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	fmt.Println(user)

	userFunds := new(models.UserFunds)
	userFunds.UserId = user.UserId
	o.Read(userFunds) //根据用户id查询资金信息
	fmt.Printf("%d\n", userFunds.TotalFunds)

	this.Data["json"] = &userFunds
	this.ServeJSON()
}

func (this *UserController) GetReceipt() {
	//orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30) //连接数据库
	//o := orm.NewOrm()

	pK := this.GetString("pucKey")
	user := new(models.User)
	user.PublicKey = pK

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	fmt.Printf("%s receipt:\n", user.UserName)

	var userReceipts []models.UserReceipt
	_, err = o.QueryTable("UserReceipt").Filter("UserId", user.UserId).All(&userReceipts)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &userReceipts
	this.ServeJSON()
}

func (this *UserController) ListTrade() {
	fmt.Printf("<---------------------------->\n")
	fmt.Printf("user.go ListTrade() 挂牌执行输出如下：\n")
	this.TplName = "user.html"
	var orderSell models.OrderSell //用户仓单结构体

	body := this.Ctx.Input.RequestBody //获取http数据
	fmt.Printf("前端传来数据如下：\n")
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &orderSell); err == nil {
		orderSell.Insert() //写入数据库
	} else {
		fmt.Println("json Unmarshal 出错：")
		fmt.Println(err)
	}

	num, err := o.QueryTable("user").Filter("user_id", orderSell.UserId).Update(orm.Params{
		"nonce": orderSell.NonceSell})
	if err == nil {
		fmt.Printf("更新user表 nonce字段成功 id = %d\n", num)
	} else {
		fmt.Printf("更新user表 nonce字段失败 %v", err)
	}
	fmt.Printf("<---------------------------->\n")
}
