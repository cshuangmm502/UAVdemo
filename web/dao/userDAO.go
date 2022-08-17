package dao

import (
	"UAVdemo/web/model"
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

//Create user table
func CreateTableWithUser() {
	sqlStr := `CREATE TABLE IF NOT EXISTS user (
				id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
				username VARCHAR (64),
				propic VARCHAR (64),
				PASSWORD VARCHAR (64),
				role varchar(64),
				phone VARCHAR (64),
				STATUS varchar(64),
				createtime VARCHAR (64)
			);
			alter table user default character set utf8;
			alter table user change username username varchar(64) character set utf8;
			alter table user change propic propic varchar(64) character set utf8;
			alter table user change role role varchar(64) character set utf8;
			alter table user change STATUS status varchar(64) character set utf8;`
	Exec(sqlStr)
	fmt.Println("---------------------------------------------")
	fmt.Println("user table created")
}

func TimeStampToData(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func MD5(str string) string {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5str
}

func CreateStaffInUser() {
	un := "s1"
	psw := MD5("1")
	role := "员工1"
	tel := 1777
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,phone,status,createtime) values (?,?,?,?,?,?)",
		un, psw, role, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}
func CreateStaff2InUser() {
	un := "s2"
	psw := MD5("1")
	role := "员工2"
	tel := 1777
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,phone,status,createtime) values (?,?,?,?,?,?)",
		un, psw, role, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}
func CreateStaff3InUser() {
	un := "s3"
	psw := MD5("1")
	role := "员工3"
	tel := 1777
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,phone,status,createtime) values (?,?,?,?,?,?)",
		un, psw, role, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}
func CreateStaff4InUser() {
	un := "s4"
	psw := MD5("1")
	role := "员工4"
	tel := 1777
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,phone,status,createtime) values (?,?,?,?,?,?)",
		un, psw, role, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}
func CreateStaff5InUser() {
	un := "s5"
	psw := MD5("1")
	role := "员工5"
	tel := 1777
	st := "正常"
	ct := TimeStampToData(time.Now().Unix())

	_, _ = Exec("insert into user(username,password,role,phone,status,createtime) values (?,?,?,?,?,?)",
		un, psw, role, tel, st, ct)

	fmt.Println("---------------------------------------------")
	fmt.Println("User created")
}


//按条件查询
func QueryUserWightCon(con string) int {
	sqlStr := fmt.Sprintf("select id from user %s", con)
	fmt.Println(sqlStr)
	row := QueryRowDB(sqlStr)

	fmt.Println("Row is:", row)
	id := 0
	row.Scan(&id)

	log.Println(id)

	fmt.Println("查到的id为", id)
	return id
}


// 通过 username 和 password 查找 User全部信息
func FindUserByUsernameAndPassword(username string, password string) (user *model.User) {

	var id int
	var role string //0 普通， 1 管理员

	var phone string
	var status string // 0 正常状态， 1 删除
	var createtime string

	sqlStr := fmt.Sprintf("select id, role,  phone, status, createtime from user where username='%s' and password='%s'", username, password)
	row := QueryRowDB(sqlStr)
	_ = row.Scan(&id, &role,  &phone, &status, &createtime)

	user = &model.User{
		Id:         id,
		Username:   username,
		Password:   password,
		Role:       role,
		Phone:      phone,
		Status:     status,
		Createtime: createtime,
	}
	return
}
