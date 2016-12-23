package control

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

const (
	CDispatchAuth    = "00" //为用户分配权限
	CAddUser         = "01" //增加用户
	CDeleteUser      = "02" //删除用户
	CAddBillboard    = "03" //增加公告
	CViewBillboard   = "04" //查看公告
	CModifyBillboard = "05" //修改公告
	CDeleteBillboard = "06" //删除公告
	CViewUsermsg     = "07" //查看指定用户所有信息
	CModifyUsermsg   = "08" //修改指定用户所有信息
	CViewMsg         = "09" //查看自己所有信息
	CModifyMsg       = "10" //修改自己所有信息
	CLogin           = "11" //登陆
	CModifyPwd       = "12" //修改密码
)

var db *sql.DB

var G_op_func map[string]func(map[string]string) map[string]string = make(map[string]func(map[string]string) map[string]string)

func init() {
	SetupDB("root", "qq5566", "localhost", "3306", "webserver")
	G_op_func[CDispatchAuth] = DispatchAuth
	G_op_func[CAddUser] = AddUser
	G_op_func[CDeleteUser] = DelUser
	G_op_func[CAddBillboard] = AddBillboard
	G_op_func[CViewBillboard] = ViewBillboard
	G_op_func[CModifyBillboard] = ModifyBillboard
	G_op_func[CDeleteBillboard] = DelBillboard
	G_op_func[CViewUsermsg] = ViewMsg
	G_op_func[CModifyUsermsg] = ModifyMsg
	G_op_func[CViewMsg] = ViewMsg
	G_op_func[CModifyMsg] = ModifyMsg
	G_op_func[CLogin] = Login
	G_op_func[CModifyPwd] = ModifyPwd
}

func AddUser(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	b, m2 := MAddUser(arg_info["name"], arg_info["pwd"], arg_info["sex"], arg_info["age"], arg_info["intro"])
	if b {
		m["tag"] = "1" + CAddUser
		m["user_id"] = m2["user_id"]
	} else {
		m["tag"] = "0" + CAddUser
	}
	return m
}

func DelUser(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	b := MDeleteUser(arg_info["user_id"])
	if b {
		m["tag"] = "1" + CDeleteUser
	} else {
		m["tag"] = "0" + CDeleteUser
	}
	return m
}

func DispatchAuth(arg_info map[string]string) map[string]string {
	b := MDispatchAuth(arg_info["user_id"], arg_info["role"])
	m := make(map[string]string)
	if b {
		m["tag"] = "1" + CDispatchAuth
	} else {
		m["tag"] = "0" + CDispatchAuth
	}
	return m
}

func ModifyUserMsg(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	b := MModifyUsermsg(arg_info["user_id"], arg_info["name"], arg_info["sex"], arg_info["pwd"], arg_info["age"], arg_info["intro"])
	if b {
		m["tag"] = "1" + CModifyUsermsg
	} else {
		m["tag"] = "0" + CModifyUsermsg
	}
	return m
}

func AddBillboard(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//MAddBillboard(last_time, name, content, suser_id string) (bool, map[string]string)
	b, m2 := MAddBillboard(arg_info["last_time"], arg_info["name"], arg_info["content"], arg_info["user_id"])
	if b {
		m["tag"] = "1" + CAddBillboard
		m["billboard_id"] = m2["billboard_id"]
	} else {
		m["tag"] = "0" + CAddBillboard
	}
	return m
}

func ViewBillboard(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//MViewBillboard() (bool, map[string]string)
	//contents,last_times,names,billboard_ids,user_ids
	b, m2 := MViewBillboard()
	if b {
		m["tag"] = "1" + CViewBillboard
		m["contents"] = m2["contents"]
		m["last_times"] = m2["last_times"]
		m["names"] = m2["names"]
		m["billboard_ids"] = m2["billboard_ids"]
		m["user_ids"] = m2["user_ids"]
	} else {
		m["tag"] = "0" + CViewBillboard
	}
	return m
}

func ModifyBillboard(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	// MModifyBillboard(sid, content string) bool
	//contents,last_times,names,billboard_ids,user_ids
	b := MModifyBillboard(arg_info["user_id"], arg_info["content"])
	if b {
		m["tag"] = "1" + CModifyBillboard
	} else {
		m["tag"] = "0" + CModifyBillboard
	}
	return m
}

func DelBillboard(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	// MDeleteBillboard(sid string) bool
	//contents,last_times,names,billboard_ids,user_ids
	b := MDeleteBillboard(arg_info["billboard_id"])
	if b {
		m["tag"] = "1" + CDeleteBillboard
	} else {
		m["tag"] = "0" + CDeleteBillboard
	}
	return m
}

func ViewMsg(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//  MViewUsermsg(suser_id string) (bool, map[string]string)
	//name sex	pwd  age  intro
	b, m2 := MViewUsermsg(arg_info["user_id"])
	if b {
		m["tag"] = "1" + CViewUsermsg
		m["name"] = m2["name"]
		m["sex"] = m2["sex"]
		m["pwd"] = m2["pwd"]
		m["age"] = m2["age"]
		m["intro"] = m2["intro"]
	} else {
		m["tag"] = "0" + CViewUsermsg
	}
	return m
}

func ModifyMsg(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//  MModifyUsermsg(suser_id, name, sex, pwd, age, intro string) bool
	//name sex	pwd  age  intro
	b := MModifyUsermsg(arg_info["user_id"], arg_info["name"], arg_info["sex"], arg_info["pwd"], arg_info["age"], arg_info["intro"])
	if b {
		m["tag"] = "1" + CModifyUsermsg
	} else {
		m["tag"] = "0" + CModifyUsermsg
	}
	return m
}

func Login(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//  MLogin(suser_id, password string) (bool, map[string]string)
	//name sex	pwd  age  intro
	b, m2 := MLogin(arg_info["user_id"], arg_info["pwd"])
	if b {
		fmt.Println("ok")
		m["tag"] = "1" + CLogin
		m["auth"] = m2["auth"]
	} else {
		fmt.Println("wrong")
		m["tag"] = "0" + CLogin
	}
	return m
}
func ModifyPwd(arg_info map[string]string) map[string]string {
	m := make(map[string]string)
	//   MModifyPwd(suser_id, newpwd string) bool
	//name sex	pwd  age  intro
	b := MModifyPwd(arg_info["user_id"], arg_info["pwd"])
	if b {
		m["tag"] = "1" + CModifyPwd
	} else {
		m["tag"] = "0" + CModifyPwd
	}
	return m
}

// AuthorizationManagement project models.go
//########################################数据库API层#############################################################
//#
//# 数据库访问层
//# 由业务逻辑层调用
//# 由于数据的合法性和存在性已经在业务逻辑层判断过，故此接口不再进行检验，直接增删改查
//##########################################################################################################

//-------------------------------------------------------------------------------------------------------
//****************************************权限号常量*********************************************************
//--------------------------------------------------------------------------------------------------------

//*****************************************************************************************************************
//***********************************数据库连接及配置*******************************************************************
//*********************************************************************************************************

//连接数据库
func SetupDB(root, password, ip, port, database string) {
	var err error
	db, err = sql.Open("mysql", root+":"+password+"@tcp("+ip+":"+port+")/"+database+"?charset=utf8&parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	/*cr_table := "create table if not exists users(user_id int auto_increment primary key, user_name varchar(64),user_password varchar(64),state tinyint default 0)"
	stmt, err := db.Prepare(cr_table)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}*/
}

//-----------------------------------判断用户是否存在----------------------------------------------------
//若存在，返回true和用户名字
//若不存在，返回false和空字符串
func IsUserExist(sid string) bool {
	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println(err)
	}
	name := ""
	err = db.QueryRow("select name from users where user_id=?", iid).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			fmt.Println(err)
		}
	}

	return true
}

//----------------------------------------------判断公告是否存在---------------------------------------------
func IsBillboardExist(sid string) bool {
	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(sid)
	if err != nil {
		return false
	}
	name := ""
	err = db.QueryRow("select name from billboard where billboard_id=?", iid).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			fmt.Println(err)
			return true
		}
	}

	return true
}

//*****************************************************************************************************************
//*********************************************数据库接口************************************************************
//*****************************************************************************************************************

//-----------------------------------------为用户分配权限-----------------------------------------------------------------

func MDispatchAuth(sid, role string) bool {

	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println(err)
	}

	/*//若该用户不存在，返回失败
	if !models.IsUserExist(a["name"]) {
		b["tag"] = "0"
		b["reason"] = "该用户不存在"
		return b
	}*/

	stmt, err := db.Prepare("UPDATE user_role SET role=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(role, iid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true

}

//------------------------------------增加用户，初始角色默认为普通用户-----------------------------------------------
func MAddUser(name, pwd, sex, age, intro string) (bool, map[string]string) {
	var b = map[string]string{} //用于返回
	//插入users表
	stmt, err := db.Prepare("INSERT INTO users (name,pwd,age,sex,intro)VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	res, err := stmt.Exec(name, pwd, age, sex, intro)

	if err != nil {
		fmt.Println(err)
		return false, b
	}

	lastid, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return false, b
	}

	b["user_id"] = strconv.Itoa(int(lastid))

	//插入user_role表，默认初始角色为普通用户
	stmt, err = db.Prepare("INSERT INTO user_role (user_id,role)VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	_, err = stmt.Exec(lastid, "user")
	if err != nil {
		fmt.Println(err)
		return false, b
	}

	return true, b

}

//----------------------------------------删除用户---------------------------------------------
func MDeleteUser(sid string) bool {
	//fmt.Println("stept 1")
	fmt.Println("传进来的id：", sid)
	//将string类型的sid转换为int的iid

	//从users表中删除
	stmt, err := db.Prepare("DELETE FROM users WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(sid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//从user_role表中删除
	stmt, err = db.Prepare("DELETE FROM user_role WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(sid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//------------------------------------增加公告-----------------------------------------------------------
func MAddBillboard(last_time, name, content, suser_id string) (bool, map[string]string) {

	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(suser_id)
	if err != nil {
		fmt.Println(err)
	}

	var b = map[string]string{} //用于返回
	//插入billboard表
	stmt, err := db.Prepare("INSERT INTO billboard (name,content,user_id)VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	res, err := stmt.Exec(name, content, iid)
	if err != nil {
		fmt.Println(err)
		return false, b
	}

	//获取公告id
	lastid, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return false, b
	}

	//转为字符串
	b["billboard_id"] = strconv.Itoa(int(lastid))

	return true, b
}

//----------------------------------------查看公告(前三条)-----------------------------------------------------------
func MViewBillboard() (bool, map[string]string) {

	var b = map[string]string{} //用于返回
	b["last_times"] = ""
	b["names"] = ""
	b["billboard_ids"] = ""
	b["user_ids"] = ""
	var (
		last_time string
		name      string
		id        int
		content   string
		user_id   int
	)
	rows, err := db.Query("select last_time,name,billboard_id,content,user_id from billboard order by billboard_id desc limit 3")
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	defer rows.Close()
	isEmpty := true
	for rows.Next() {
		isEmpty = false
		err := rows.Scan(&last_time, &name, &id, &content, &user_id)
		if err != nil {
			fmt.Println(err)
			return false, b
		}
		last_time = strings.Replace(last_time, "T", " ", -1)
		last_time = strings.Replace(last_time, "Z", "", -1)
		b["last_times"] += (last_time + "%s9k@!")
		b["names"] += (name + "%s9k@!")
		b["billboard_ids"] += (strconv.Itoa(id) + "%s9k@!")
		b["contents"] += (content + "%s9k@!")
		b["user_ids"] += (strconv.Itoa(user_id) + "%s9k@!")

	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	if !isEmpty {
		b["last_times"] = b["last_times"][0 : len(b["last_times"])-6]
		b["names"] = b["names"][0 : len(b["names"])-6]
		b["billboard_ids"] = b["billboard_ids"][0 : len(b["billboard_ids"])-6]
		b["contents"] = b["contents"][0 : len(b["contents"])-6]
		b["user_ids"] = b["user_ids"][0 : len(b["user_ids"])-6]
	}
	return true, b
}

//-----------------------------------------------修改公告-------------------------------------------------------
func MModifyBillboard(sid, content string) bool {

	iid, err := strconv.Atoi(sid)
	if err != nil {
		return false
	}

	stmt, err := db.Prepare("UPDATE billboard SET content=?,last_time=? WHERE billboard_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}

	last_time := time.Now().String()[0:19]
	//fmt.Println(last_time)
	_, err = stmt.Exec(content, last_time, iid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//-----------------------------------------删除公告---------------------------------------------------------
func MDeleteBillboard(sid string) bool {

	//将传进来的string的sid转换为int的iid
	iid, err := strconv.Atoi(sid)
	if err != nil {
		return false
	}

	stmt, err := db.Prepare("DELETE FROM billboard WHERE billboard_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(iid)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//------------------------------------查看指定用户所有信息-------------------------------------------------
func MViewUsermsg(suser_id string) (bool, map[string]string) {

	var b = map[string]string{} //用于返回
	var err error
	var (
		bname  = ""
		bsex   = ""
		bage   = ""
		bpwd   = ""
		brole  = ""
		bintro = ""
	)
	err = db.QueryRow("select name,pwd,sex,age,intro from users where user_id=?", suser_id).Scan(&bname, &bpwd, &bsex, &bage, &bintro)
	if err != nil {
		fmt.Println("查询基本信息出错" + err.Error())
		return false, b
	}
	err = db.QueryRow("select role from user_role where user_id=?", suser_id).Scan(&brole)
	if err != nil {
		fmt.Println("查询角色时出错")
		return false, b
	}
	b["name"] = bname
	b["sex"] = bsex
	b["age"] = bage
	b["pwd"] = bpwd
	b["role"] = brole
	b["intro"] = bintro
	fmt.Println(b)

	return true, b
}

//-------------------------------------修改指定用户所有信息-------------------------------------------------

func MModifyUsermsg(suser_id, name, sex, pwd, age, intro string) bool {

	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(suser_id)
	if err != nil {
		fmt.Println(err)
	}

	//修改user表
	stmt, err := db.Prepare("UPDATE users SET name=?,sex=?,pwd=?,age=?,intro=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = stmt.Exec(name, sex, pwd, age, intro, iid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	/*//修改user_role表的name属性
	stmt, err = db.Prepare("UPDATE user_role SET user=? WHERE user=?")
	if err != nil {
		log.Fatal(err)
		return false
	}

	_, err = stmt.Exec(newname, oldname)
	if err != nil {
		log.Fatal(err)
		return false
	}*/
	return true
}

func MModifyUserBasicmsg(suser_id, name, sex, age, intro string) bool {

	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(suser_id)
	if err != nil {
		fmt.Println(err)
	}

	//修改user表
	stmt, err := db.Prepare("UPDATE users SET name=?,sex=?,age=?,intro=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = stmt.Exec(name, sex, age, intro, iid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//--------------------------------------------登录---------------------------------------------------------------
func MLogin(suser_id, password string) (bool, map[string]string) {
	fmt.Println(suser_id, password)
	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(suser_id)
	if err != nil {
		fmt.Println(err)
	}

	var b = map[string]string{} //用于返回

	var user_password string
	err = db.QueryRow("select pwd from users where user_id = ?", iid).Scan(&user_password)
	if err != nil {
		fmt.Println(err)
		return false, b
	}
	if password != user_password {
		return false, b
	}
	b["auth"] = GetAuth(iid)
	return true, b
}

//-----------------------------------------修改密码--------------------------------------------------
func MModifyPwd(suser_id, newpwd string) bool {

	//将string类型的sid转换为int的iid
	iid, err := strconv.Atoi(suser_id)
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := db.Prepare("UPDATE users SET pwd=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = stmt.Exec(newpwd, iid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetAuth(id int) string {
	var role string
	var err error
	err = db.QueryRow("select role from user_role where user_id = ?", id).Scan(&role)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var auth string
	err = db.QueryRow("select auth from role_auth where role = ?", role).Scan(&auth)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return auth
}

func GetRole(id int) string {
	var role string
	var err error
	err = db.QueryRow("select role from user_role where user_id = ?", id).Scan(&role)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return role
}
