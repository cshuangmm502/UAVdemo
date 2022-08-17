package web

import (
	"UAVdemo/web/controller"
	"fmt"
	"net/http"
)

//func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
//	group.handle("GET", relativePath, handlers)
//	group.handle("POST", relativePath, handlers)
//	group.handle("PUT", relativePath, handlers)
//	group.handle("PATCH", relativePath, handlers)
//	group.handle("HEAD", relativePath, handlers)
//	group.handle("OPTIONS", relativePath, handlers)
//	group.handle("DELETE", relativePath, handlers)
//	group.handle("CONNECT", relativePath, handlers)
//	group.handle("TRACE", relativePath, handlers)
//	return group.returnObj()
//}

func WebStart(app *controller.Application)  {


	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 指定第一次打开系统进入的页面
	http.HandleFunc("/", app.Home)
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/backToHome", app.BackToHome)			// 返回首页
	// 登陆
	http.HandleFunc("/loginPage", app.LoginView)
	http.HandleFunc("/login", app.Login)

	//// 查询
	//http.HandleFunc("/queryPage", app.QueryPage)		// 转至查询信息页面
	//http.HandleFunc("/findFoodByID", app.FindFoodByID)	// 根据id查询并转至查询结果页面
	//
	//http.HandleFunc("/addFoodPage", app.AddFoodPage) // 显示添加信息页面
	//http.HandleFunc("/addFood", app.AddFood)         // 提交修改请求并跳转添加成功提示页面
	//
	//http.HandleFunc("/addStorageUnitPage", app.AddStorageUnitPage) // 显示添加信息页面
	//http.HandleFunc("/addStorageUnit", app.AddStorageUnit)         // 提交修改请求并跳转添加成功提示页面
	//
	//http.HandleFunc("/addTransporterPage", app.AddTransporterPage) // 显示添加信息页面
	//http.HandleFunc("/addTransporter", app.AddTransporter)         // 提交修改请求并跳转添加成功提示页面
	//
	//http.HandleFunc("/addDealerPage", app.AddDealerPage) // 显示添加信息页面
	//http.HandleFunc("/addDealerInfo", app.AddDealerInfo)         // 提交修改请求并跳转添加成功提示页面
	//
	//http.HandleFunc("/addCustomerPage", app.AddCustomerInfoPage) // 显示添加信息页面
	//http.HandleFunc("/addCustomerInfo", app.AddCustomerInfo)         // 提交修改请求并跳转添加成功提示页面

	fmt.Println("---------------------------------------------")
	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}