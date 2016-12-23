package routers

import (
	"UserManaServer/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/deny", &controllers.DenyController{})
	beego.Router("/regis", &controllers.RegisController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/auth", &controllers.AuthController{})
	beego.Router("/billboard", &controllers.BillboardController{})
	beego.Router("/userMgn", &controllers.UserMgnController{})
	beego.Router("/modify_ourself", &controllers.ModifyFController{})
	beego.Router("/modify_other", &controllers.ModifyRController{})
}
