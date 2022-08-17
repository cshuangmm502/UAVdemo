package dao

import (
	"fmt"
)

//Create action table
func CreateTableWithAction() {
	sqlStr := `CREATE TABLE IF NOT EXISTS action (
				action_id BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL,
				user_id BIGINT NOT NULL,
				action_type VARCHAR (64) NOT NULL,
				action_target_name VARCHAR (64) NOT NULL,
				action_target_id VARCHAR (64) NOT NULL,
				action_time VARCHAR (64) NOT NULL,
				FOREIGN KEY (user_id) REFERENCES user (id)
			);
			alter table action default character set utf8;
			alter table action change action_type action_type varchar(64) character set utf8;
			alter table action change action_target_name action_target_name varchar(64) character set utf8;
			alter table action change action_target_id action_target_id varchar(64) character set utf8;
			alter table action change action_time action_time varchar(64) character set utf8;`
	Exec(sqlStr)
	fmt.Println("---------------------------------------------")
	fmt.Println("Action table created")
}

