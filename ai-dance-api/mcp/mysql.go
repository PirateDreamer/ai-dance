package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Mysql *gorm.DB

type ConnectMysqlParam struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DBName string `json:"db_name"`
}

func ConnectMysql(param ConnectMysqlParam) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", param.User, param.Pass, param.Host, param.Port, param.DBName)
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return
	}
	return
}

type QueryResult struct {
	Data any
}

func Query(sql, queryType string) (res QueryResult, err error) {
	switch queryType {
	case "table_list":
		// sql 检查
		var tables []string
		if err = Mysql.Raw(sql).Scan(&tables).Error; err != nil {
			return
		}
		res = QueryResult{
			Data: tables,
		}
	case "table_structure":
		// sql 检查
		var ddl []map[string]any
		if err = Mysql.Raw(sql).Scan(&ddl).Error; err != nil {
			return
		}
		res = QueryResult{
			Data: ddl,
		}
	case "table_data":
		// sql 检查
		var data []map[string]any
		if err = Mysql.Raw(sql).Scan(&data).Error; err != nil {
			return
		}
		res = QueryResult{
			Data: data,
		}
	default:
		err = fmt.Errorf("不支持的查询类型: %s", queryType)
	}
	return
}
