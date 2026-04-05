package models

import (
	"fmt"
	"go-admin-api/common/global"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

func InitDb(db *gorm.DB) (err error) {
	filePath := "config/db.sql"
	err = ExecSql(db, filePath)
	if global.Driver == "postgres" {
		filePath = "config/pg.sql"
		err = ExecSql(db, filePath)
	}
	return err
}

func ExecSql(db *gorm.DB, filePath string) error {
	sql, err := Ioutil(filePath)
	if err != nil {
		fmt.Println("failed to read database initialization script, reason:", err.Error())
		return err
	}
	sqlList := strings.Split(sql, ";")
	for i := 0; i < len(sqlList)-1; i++ {
		if strings.Contains(sqlList[i], "--") {
			fmt.Println(sqlList[i])
			continue
		}
		sql := strings.Replace(sqlList[i]+";", "\n", "", -1)
		sql = strings.TrimSpace(sql)
		if err = db.Exec(sql).Error; err != nil {
			log.Printf("error sql: %s", sql)
			if !strings.Contains(err.Error(), "Query was empty") {
				return err
			}
		}
	}
	return nil
}

func Ioutil(filePath string) (string, error) {
	if contents, err := os.ReadFile(filePath); err == nil {
		result := strings.Replace(string(contents), "\n", "", 1)
		fmt.Println("Use ioutil.ReadFile to read a file:", result)
		return result, nil
	} else {
		return "", err
	}
}
