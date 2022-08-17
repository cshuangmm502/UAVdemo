package controller

import (
	"UAVdemo/web/dao"
	"UAVdemo/web/model"
	"UAVdemo/web/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//进入首页
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin==false{
		ShowView(w,r,"AccountRelated/login.html",data)
	}else {
		ShowView(w, r, "index.html", data)
	}
}

// 返回首页
func (app *Application) BackToHome(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "index.html", data)
}



// 随机数字6位
func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 进入登录界面
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		ShowView(w, r, "index.html", data)
		return
	}
	ShowView(w, r, "AccountRelated/login.html", data)
}

// 用户登录
func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Sess         *model.Session
		FailedLogin  bool
		IsLogin      bool
		IsSuperAdmin bool
		IsAdmin      bool
		IsUser       bool
		IsStaff1      bool
		IsStaff2      bool
		IsStaff3      bool
		IsStaff4     bool
		IsStaff5     bool
		Msg          string
	}{
		Sess:         nil,
		FailedLogin:  false,
		IsLogin:      false,
		IsSuperAdmin: false,
		IsAdmin:      false,
		IsUser:       false,
		IsStaff1:      false,
		IsStaff2:      false,
		IsStaff3:      false,
		IsStaff4:      false,
		IsStaff5:      false,
		Msg:          "",
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("默认参数已就绪")
	fmt.Println("---------------------------------------------")
	//获取表格信息
	username := r.FormValue("loginName")
	Password := r.FormValue("password")
	password := utils.MD5(Password)
	fmt.Println("---------------------------------------------")
	fmt.Println("前端表格读取完成")

	//返回完整的用户信息
	user := dao.FindUserByUsernameAndPassword(username, password)
	fmt.Println("---------------------------------------------")
	fmt.Println("用户", user.Username, "查询结果已传回，正在核查")

	if user.Id == 0 {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println("用户名或密码错误，登陆失败，以未登录状态返回首页")
		data.Msg = "用户名或密码错误"
		ShowView(w, r, "AccountRelated/login.html", data)
		return

	} else if user.Status == "异常" {
		data.FailedLogin = true
		fmt.Println("---------------------------------------------")
		fmt.Println(user.Role, user.Username, "账户受限，登陆失败，以未登录状态返回首页")
		data.Msg = user.Role + user.Username + "账户受限，登陆失败，请联系管理员"
		ShowView(w, r, "index.html", data)
		return

	} else if user.Status == "正常" {
		uuid := utils.CreateUUID()
		session := &model.Session{
			SessionID:  uuid,
			UserID:     user.Id,
			UserName:   user.Username,
			PassWord:   user.Password,
			Role:       user.Role,
			Phone:      user.Phone,
			Status:     user.Status,
			CreateTime: user.Createtime,
		}

		_ = dao.AddSession(session)

		fmt.Println("---------------------------------------------")
		fmt.Println("Session已设置")

		cookie := http.Cookie{
			Name:     "user",
			Value:    uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		fmt.Println("---------------------------------------------")
		fmt.Println("Cookie已送往浏览器")
		if user.Role == "员工1" {
			data.IsStaff1 = true
		}
		if user.Role == "员工2" {
			data.IsStaff2 = true
		}
		if user.Role == "员工3" {
			data.IsStaff3 = true
		}
		if user.Role == "员工4" {
			data.IsStaff4 = true
		}
		if user.Role == "员工5" {
			data.IsStaff5 = true
		}
		data.IsLogin = true
		data.Sess = session
		fmt.Println("---------------------------------------------")
		fmt.Println("默认参数已更新")

		ShowView(w, r, "index.html", data)

	}
}

// 退出登陆
func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	data := utils.DeleteSession(r)
	fmt.Println(data.Msg)
	ShowView(w, r, "AccountRelated/login.html", data)
}

func (app *Application) About(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "PublicOption/about.html", data)
}

