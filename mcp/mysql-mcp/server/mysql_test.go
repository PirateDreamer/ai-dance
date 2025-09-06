package main

import (
	"testing"
)

func TestConnectMysql(t *testing.T) {
	_, err := ConnectMysql(ConnectMysqlParam{DBName: "mario", Host: "localhost", Pass: "root", Port: "3306", User: "root"})
	if err != nil {
		t.Error(err)
	}
}

func TestQuery(t *testing.T) {
	var err error
	if Mysql, err = ConnectMysql(ConnectMysqlParam{DBName: "mario", Host: "localhost", Pass: "root", Port: "3306", User: "root"}); err != nil {
		t.Error(err)
	}
	var res QueryResult
	if res, err = Query("DESCRIBE serve", "table_structure"); err != nil {
		t.Error(err)
	}
	t.Log(res.Data)
}
