package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

const (
	userName = "root"
	password = "123456"
	ip       = "192.168.174.182"
	port     = "3306"
	dbName   = "food"
)

var db *sql.DB

func InitMysql() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8&multiStatements=true"}, "")

	if db == nil {
		fmt.Println("---------------------------------------------")
		fmt.Println("Database Connected")
		db, _ = sql.Open("mysql", path)
		DeleteTable()
		CreateTableWithUser()
		CreateTableWithSession()
		CreateTableWithAction()
		CreateStaff5InUser()
		CreateStaff4InUser()
		CreateStaff2InUser()
		CreateStaff3InUser()
		CreateStaffInUser()
	}
}

//查询
func QueryRowDB(sqlStr string) *sql.Row {
	return db.QueryRow(sqlStr)
}

//操作数据库
func Exec(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

func DeleteTable() {
	sqlStr := `SET FOREIGN_KEY_CHECKS = 0;
			DROP TABLE IF EXISTS user;
			SET FOREIGN_KEY_CHECKS = 1;
			DROP TABLE IF EXISTS session;
			DROP TABLE IF EXISTS action;
			DROP TABLE IF EXISTS tea;`
	fmt.Println("---------------------------------------------")
	fmt.Println("table deleted")
	_, _ = Exec(sqlStr)
}
