package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
)

func main() {

	if err := LoadConfig(); err != nil {
		panic(err)
	}

	s := server.NewMCPServer(
		"mcp-mysql-server",
		"0.0.1",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	connectTool := mcp.NewTool("connect",
		mcp.WithDescription("连接数据库"),
		mcp.WithString("connect_name",
			mcp.Required(),
			mcp.Description("主机"),
		),
	)

	executeTool := mcp.NewTool("execute",
		mcp.WithDescription("执行sql，支持执行更新、插入、删除数据，支持DDL（建表、改表）"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("SQL语句"),
		),
	)

	queryTool := mcp.NewTool("query",
		mcp.WithDescription("查询数据，支持查询数据库中的表列表、表中的数据、以及表的结构"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("SQL语句"),
		),
		mcp.WithString("query_type",
			mcp.Required(),
			mcp.Description("查询类型"),
			mcp.Enum("table_list", "table_structure", "table_data"),
		),
	)

	s.AddTool(connectTool, func(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
		var connectName string
		if connectName, err = request.RequireString("connect_name"); err != nil {
			return
		}
		prefix := fmt.Sprintf("mysql.%s", connectName)
		if Mysql, err = ConnectMysql(ConnectMysqlParam{
			Host:   viper.GetString(prefix + ".host"),
			Pass:   viper.GetString(prefix + ".pass"),
			Port:   viper.GetString(prefix + ".port"),
			User:   viper.GetString(prefix + ".user"),
			DBName: viper.GetString(prefix + ".db"),
		}); err != nil {
			return
		}

		res = mcp.NewToolResultText("连接成功")

		return
	})

	s.AddTool(executeTool, func(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
		var sql string
		if sql, err = request.RequireString("sql"); err != nil {
			return
		}

		result := Mysql.Exec(sql)
		if result.Error != nil {
			err = result.Error
			return
		}

		res = mcp.NewToolResultText(fmt.Sprintf("执行成功，影响行数：%d", result.RowsAffected))
		return
	})

	s.AddTool(queryTool, func(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
		var sql, queryType string
		if sql, err = request.RequireString("sql"); err != nil {
			return
		}
		if queryType, err = request.RequireString("query_type"); err != nil {
			return
		}
		var result QueryResult
		if result, err = Query(sql, queryType); err != nil {
			return
		}
		var resBytes []byte
		if resBytes, err = json.Marshal(result.Data); err != nil {
			return
		}
		res = mcp.NewToolResultText(string(resBytes))
		return
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
