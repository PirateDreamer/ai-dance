package main

import (
	"fmt"
	"regexp"
	"strings"
)

type RequestInfo struct {
	Method      string
	Url         string
	BaseUrl     string
	Uri         string
	Body        map[string]string
	Headers     map[string]string
	QueryParams map[string]string
	PathParams  map[string]string
}

// 解析curl
func parseCurl(curl string) (reqInfo RequestInfo, err error) {
	// 初始化结构体
	reqInfo = RequestInfo{
		Body:        make(map[string]string),
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		PathParams:  make(map[string]string),
		Method:      "GET", // 默认方法
	}

	// 去除首尾空格
	curl = strings.TrimSpace(curl)

	// 检查是否以 curl 开头
	if !strings.HasPrefix(curl, "curl") {
		err = fmt.Errorf("invalid curl command")
		return
	}

	// 移除 curl 前缀
	command := strings.TrimPrefix(curl, "curl")
	command = strings.TrimSpace(command)

	// 解析各种参数
	tokens := tokenize(command)

	// 解析URL和其他参数
	err = parseTokens(tokens, &reqInfo)
	if err != nil {
		return
	}

	// 如果URL存在，解析URL组件
	if reqInfo.Url != "" {
		parseUrlComponents(&reqInfo)
	}

	return
}

// 简单的词法分析器，将命令行字符串分解为标记
func tokenize(command string) []string {
	var tokens []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for i, char := range command {
		switch char {
		case '"', '\'':
			if !inQuote {
				inQuote = true
				quoteChar = char
			} else if quoteChar == char && (i == 0 || command[i-1] != '\\') {
				inQuote = false
				quoteChar = rune(0)
			} else {
				current.WriteRune(char)
			}
		case ' ', '\t':
			if inQuote {
				current.WriteRune(char)
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteRune(char)
		}
	}

	// 添加最后一个标记
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// 解析标记并填充请求信息
func parseTokens(tokens []string, reqInfo *RequestInfo) error {
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token {
		case "-X", "--request":
			if i+1 < len(tokens) {
				reqInfo.Method = strings.ToUpper(tokens[i+1])
				i++
			}
		case "-H", "--header":
			if i+1 < len(tokens) {
				header := tokens[i+1]
				if parts := strings.SplitN(header, ":", 2); len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					reqInfo.Headers[key] = value
				}
				i++
			}
		case "-d", "--data", "--data-raw", "--data-binary":
			if i+1 < len(tokens) {
				parseData(tokens[i+1], reqInfo)
				// 如果没有指定方法，默认为POST
				if reqInfo.Method == "GET" {
					reqInfo.Method = "POST"
				}
				i++
			}
		case "-u", "--user":
			if i+1 < len(tokens) {
				reqInfo.Headers["Authorization"] = "Basic " + tokens[i+1]
				i++
			}
		case "-A", "--user-agent":
			if i+1 < len(tokens) {
				reqInfo.Headers["User-Agent"] = tokens[i+1]
				i++
			}
		case "-e", "--referer":
			if i+1 < len(tokens) {
				reqInfo.Headers["Referer"] = tokens[i+1]
				i++
			}
		case "-b", "--cookie":
			if i+1 < len(tokens) {
				reqInfo.Headers["Cookie"] = tokens[i+1]
				i++
			}
		default:
			// 检查是否是URL
			if isUrl(token) {
				reqInfo.Url = token
			}
		}
	}

	return nil
}

// 检查字符串是否为URL
func isUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// 解析数据参数
func parseData(data string, reqInfo *RequestInfo) {
	// 处理JSON数据
	if strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}") {
		reqInfo.Headers["Content-Type"] = "application/json"
		// 对于简单实现，我们只标记内容类型
		return
	}

	// 处理表单数据
	if strings.Contains(data, "=") {
		reqInfo.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		pairs := strings.Split(data, "&")
		for _, pair := range pairs {
			if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
				reqInfo.Body[kv[0]] = kv[1]
			}
		}
	}
}

// 解析URL组件
func parseUrlComponents(reqInfo *RequestInfo) {
	// 提取基础URL和路径
	url := reqInfo.Url

	// 使用正则表达式解析URL
	re := regexp.MustCompile(`^(https?://[^/]+)(/.*)?$`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		reqInfo.BaseUrl = matches[1]
		if len(matches) > 2 {
			reqInfo.Uri = matches[2]
		} else {
			reqInfo.Uri = "/"
		}
	}

	// 解析查询参数
	if strings.Contains(reqInfo.Uri, "?") {
		parts := strings.Split(reqInfo.Uri, "?")
		reqInfo.Uri = parts[0]

		if len(parts) > 1 && parts[1] != "" {
			queryString := parts[1]
			pairs := strings.Split(queryString, "&")
			for _, pair := range pairs {
				if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
					reqInfo.QueryParams[kv[0]] = kv[1]
				} else if len(kv) == 1 {
					reqInfo.QueryParams[kv[0]] = ""
				}
			}
		}
	}
}

func main() {
	curl := `curl --location --request POST 'https://schoolidea.seakoi.net/api/lottery/runLottery' \
--header 'Authorization-Token;' \
--header 'Content-Type: application/json' \
--data-raw '{
    "activity_id": 0,
    "lottery_count": 0
}'`
	reqInfo, err := parseCurl(curl)
	if err != nil {
		panic(err)
	}
	fmt.Println(reqInfo)
}
