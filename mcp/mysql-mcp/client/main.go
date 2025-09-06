package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// 这里的路径是上面编译的mcp-server可执行文件
	mcpClient, err := client.NewStdioMCPClient("D:/workspace/mcp-project/mcp-mysql/server/mcp-mysql-server.exe", []string{})
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Mcp Mysql Client",
		Version: "1.0.0",
	}

	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("初始化成功，服务器信息: %s %s\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)

	// 调用工具
	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolRequest.Params.Name = "connect"
	toolRequest.Params.Arguments = map[string]any{
		"connect_name": "local",
	}

	result, err := mcpClient.CallTool(ctx, toolRequest)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Content[0].(mcp.TextContent).Text)
}
