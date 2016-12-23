package controllers

import (
	"UserManaServer/control"
	"UserManaServer/process"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type LoginController struct { //登录
	beego.Controller
}

type DenyController struct { //登录
	beego.Controller
}
type UserController struct { //用户界面
	beego.Controller
}
type RegisController struct { //用户注册
	beego.Controller
}
type ModifyFController struct { //修改自己的信息
	beego.Controller
}
type ModifyRController struct { //修改别人的信息
	beego.Controller
}

type AuthController struct { //权限管理
	beego.Controller
}
type BillboardController struct { //公告管理
	beego.Controller
}
type UserMgnController struct { //用户管理
	beego.Controller
}
type ModifyController struct { //修改信息管理
	beego.Controller
}

func (this *DenyController) Get() {
	this.TplName = "deny.html"
}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}
func (this *LoginController) Post() {
	name := this.Input().Get("name") //实际上是id
	pwd := this.Input().Get("pwd")
	m := make(map[string]string)
	m["user_id"] = name
	m["pwd"] = pwd
	m = control.Login(m)
	if m["tag"] == "111" { //登陆成功
		userId, _ := strconv.Atoi(name)
		token := process.NewToken(userId, pwd)
		process.G_nameId_token[userId] = token
		process.G_token_nameId[token] = userId

		this.Data["Tag"] = false
		this.SetSession("token", token)
		this.SetSession("user_id", name)
		this.SetSession("role", control.GetRole(userId))
		this.Redirect("/user", 302)

	} else { //用户名或密码错误
		this.Data["Msg"] = "用户名或密码错误！！"
		this.Data["Tag"] = true
		this.TplName = "login.html"
	}
}

func (this *UserController) Get() {

	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("user_id")
	strid, _ := data.(string)
	userId, _ := strconv.Atoi(strid)
	data = this.GetSession("role")
	role, _ := data.(string)

	if !process.CheckToken(userId, token) {
		this.Redirect("/login", 302)
	} else {
		var auths string
		if role == "sup_admin" {
			auths = "0,1,2,3"
		} else if role == "admin" {
			auths = "1,3"
		} else {
			auths = "3"
		}
		_, m := control.MViewUsermsg(strid)
		subauth := strings.Split(auths, ",")
		this.Data["Auths"] = subauth
		this.Data["id"] = strid
		this.Data["name"] = m["name"]
		if m["sex"] == "male" {
			this.Data["sex"] = "男"
		} else {
			this.Data["sex"] = "女"
		}

		this.Data["age"] = m["age"]

		r := m["role"]
		if r == "admin" {
			this.Data["role"] = "普通管理员"
		} else if r == "sup_admin" {
			this.Data["role"] = "超级管理员"
		} else {
			this.Data["role"] = "普通用户"
		}
		this.Data["intro"] = m["intro"]
		fmt.Println("0成功读取完成公告")
		b, m := control.MViewBillboard()
		fmt.Println("1成功读取完成公告")
		if b {
			last_times := strings.Split(m["last_times"], "%s9k@!")
			names := strings.Split(m["names"], "%s9k@!")
			billboard_ids := strings.Split(m["billboard_ids"], "%s9k@!")
			contents := strings.Split(m["contents"], "%s9k@!")
			user_ids := strings.Split(m["user_ids"], "%s9k@!")

			for i := 0; i < len(billboard_ids); i++ {
				if len(billboard_ids[i]) > 0 {
					this.Data["exist"+strconv.Itoa(i+1)] = true
					this.Data["last_times"+strconv.Itoa(i+1)] = last_times[i]
					this.Data["names"+strconv.Itoa(i+1)] = names[i]
					this.Data["billboard_ids"+strconv.Itoa(i+1)] = billboard_ids[i]
					this.Data["contents"+strconv.Itoa(i+1)] = contents[i]
					this.Data["user_ids"+strconv.Itoa(i+1)] = user_ids[i]
				}
			}
		}
		this.TplName = "user.html"
	}
}

func (this *AuthController) Get() {
	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("role")
	R, _ := data.(string)
	if token == "" {
		this.Redirect("/login", 302)
	} else if R != "sup_admin" {
		this.Redirect("/deny", 302)
	} else {
		data = this.GetSession("name")
		name, _ := data.(string)
		this.Data["name"] = name
		this.TplName = "auth.html"
	}
}
func (this *AuthController) Post() {

	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("role")
	R, _ := data.(string)
	if token == "" {
		this.Ctx.WriteString("0")
	} else if R != "sup_admin" {
		this.Ctx.WriteString("0")
	}

	op := this.Input().Get("op")
	id := this.Input().Get("id")
	role := this.Input().Get("role")
	fmt.Println(op, "--", id, "--", role)
	if op == "search" {
		b, m := control.MViewUsermsg(id)
		if b {
			this.Ctx.WriteString(m["role"] + "##0%7" + m["name"])
		} else {
			this.Ctx.WriteString("##0%7")
		}
	} else if op == "modify" {
		b2 := control.MDispatchAuth(id, role)
		if b2 {
			this.Ctx.WriteString("1")
		} else {
			this.Ctx.WriteString("0")
		}
	}
}

func (this *BillboardController) Get() {

	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("role")
	R, _ := data.(string)
	if token == "" {
		this.Redirect("/login", 302)
	} else if R != "admin" && R != "sup_admin" {
		this.Redirect("/deny", 302)
	}

	b, m := control.MViewBillboard()
	if b {
		last_times := strings.Split(m["last_times"], "%s9k@!")
		names := strings.Split(m["names"], "%s9k@!")
		billboard_ids := strings.Split(m["billboard_ids"], "%s9k@!")
		contents := strings.Split(m["contents"], "%s9k@!")
		user_ids := strings.Split(m["user_ids"], "%s9k@!")
		L := len(billboard_ids)
		fmt.Println("长度：", L)
		for i := 0; i < L; i++ {
			if len(billboard_ids[i]) > 0 {
				this.Data["exist"+strconv.Itoa(i+1)] = true
				this.Data["last_times"+strconv.Itoa(i+1)] = last_times[i]
				this.Data["names"+strconv.Itoa(i+1)] = names[i]
				this.Data["billboard_ids"+strconv.Itoa(i+1)] = billboard_ids[i]
				this.Data["contents"+strconv.Itoa(i+1)] = contents[i]
				this.Data["user_ids"+strconv.Itoa(i+1)] = user_ids[i]
			}
		}
	}
	this.TplName = "billboard.html"
}
func (this *BillboardController) Post() {

	data1 := this.GetSession("role")
	R, _ := data1.(string)
	if R != "admin" && R != "sup_admin" {
		this.Ctx.WriteString("0")
	}

	op := this.Input().Get("op")
	content := this.Input().Get("content")
	id := this.Input().Get("id")
	data := this.GetSession("token")
	token, _ := data.(string)
	iid := process.G_token_nameId[token]
	strid := strconv.Itoa(iid)
	ok, m := control.MViewUsermsg(strid)
	if !ok {
		this.Ctx.WriteString("0")
		this.Redirect("/login", 302)
	}
	if op == "add" {
		b, _ := control.MAddBillboard("", m["name"], content, strid)
		if b {
			this.Ctx.WriteString("1")
		} else {
			this.Ctx.WriteString("0")
		}
		return
	} else if op == "del" {
		b2 := control.MDeleteBillboard(id)
		if b2 {
			this.Ctx.WriteString("1")
		} else {
			this.Ctx.WriteString("0")
		}
	}
	return
}
func (this *UserMgnController) Get() {

	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("role")
	R, _ := data.(string)
	if token == "" {
		this.Redirect("/login", 302)
	} else if R != "sup_admin" {
		this.Redirect("/deny", 302)
	}

	this.TplName = "userMgn.html"
	//	this.Ctx.WriteString(s)
}
func (this *UserMgnController) Post() {

	data := this.GetSession("token")
	token, _ := data.(string)
	data = this.GetSession("role")
	R, _ := data.(string)
	if token == "" {
		this.Ctx.WriteString("0")
	} else if R != "sup_admin" {
		this.Ctx.WriteString("0")
	}

	id := this.Input().Get("id")
	op := this.Input().Get("op")
	if op == "search" {
		b, m := control.MViewUsermsg(id)
		if b {
			msg := ""
			msg += m["name"] + "##%07"
			if m["sex"] == "female" {
				msg += "女" + "##%07"
			} else {
				msg += "男" + "##%07"
			}
			msg += m["age"] + "##%07"
			msg += m["role"] + "##%07"
			msg += m["intro"] + "##%07"
			this.Ctx.WriteString(msg)
		} else {
			this.Ctx.WriteString("##%07")
		}
	} else if op == "del" {
		b := control.MDeleteUser(id)
		if b {
			this.Ctx.WriteString("1")
		} else {
			this.Ctx.WriteString("0")
		}
	}

	//	this.Ctx.WriteString(s)
}

func (this *ModifyController) Get() {
	//判断令牌和权限
	//假设都通过了
	this.Data["name"] = "name"
	this.Data["isfemale"] = true
	this.Data["age"] = 16
	this.Data["record"] = "name66666666"
	this.TplName = "modify.html"
	//	this.Ctx.WriteString(s)
}
func (this *ModifyController) Post() {
	/*name := this.Input().Get("name")
	age := this.Input().Get("age")
	sex := this.Input().Get("sex")
	record := this.Input().Get("record")*/
	//可以成功接收了
	this.Ctx.WriteString("1")
	return
}
func (this *RegisController) Post() {
	name := this.Input().Get("name")
	pwd := this.Input().Get("pwd")
	age := this.Input().Get("age")
	sex := this.Input().Get("sex")
	intro := this.Input().Get("intro")

	m := make(map[string]string)
	m["name"] = name
	m["pwd"] = pwd
	m["age"] = age
	m["sex"] = sex
	m["intro"] = intro
	m = control.AddUser(m)

	if m["tag"] == "101" {
		this.Ctx.WriteString(m["user_id"])
	} else {
		this.Ctx.WriteString("0")
	}
	//添加用户
	//返回提示信息：0开头代表失败，否则代表用户的账号
}
func (this *RegisController) Get() {
	this.TplName = "regis.html"
}

func (this *ModifyFController) Get() {

	data := this.GetSession("token")
	token, _ := data.(string)
	iid := process.G_token_nameId[token]
	//获取id	（int）

	b, m := control.MViewUsermsg(strconv.Itoa(iid))
	if b {
		this.Data["name"] = m["name"]
		this.Data["isfemale"] = m["sex"] == "female"
		this.Data["age"] = m["age"]
		this.Data["intro"] = m["intro"]
	} else {
		this.Redirect("/login", 302)
	}
	fmt.Println("")
	this.TplName = "modify-ourself.html"
}
func (this *ModifyFController) Post() {
	name := this.Input().Get("name")
	age := this.Input().Get("age")
	sex := this.Input().Get("sex")
	intro := this.Input().Get("intro")
	fmt.Println("name:", name)
	data := this.GetSession("token")
	token, _ := data.(string)
	iid := process.G_token_nameId[token]

	b := control.MModifyUserBasicmsg(strconv.Itoa(iid), name, sex, age, intro)
	if b {
		this.Ctx.WriteString("1")
	} else {
		this.Ctx.WriteString("0")
	}
	return
}

func (this *ModifyRController) Get() {
	this.TplName = "modify-other.html"
}
