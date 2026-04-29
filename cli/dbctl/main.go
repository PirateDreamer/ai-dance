package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

// MySQLConn 单个数据库连接配置
type MySQLConn struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	Env      string `yaml:"-"` // 从 server.env 继承
}

type Config struct {
	Server struct {
		Env string `yaml:"env"`
	} `yaml:"server"`
	Databases map[string]MySQLConn `yaml:"databases"` // 命名连接
}

// resolveConn 根据名称查找连接配置；空名称取第一个
func resolveConn(cfg *Config, name string) (string, *MySQLConn, error) {
	if len(cfg.Databases) == 0 {
		return "", nil, fmt.Errorf("配置中没有任何数据库连接，请检查 config.yaml")
	}
	// 未指定名称：取第一个（仅一个时自动选中）
	if name == "" {
		if len(cfg.Databases) > 1 {
			return "", nil, fmt.Errorf("存在多个连接，请用 -n <连接名> 指定，使用 -list 查看")
		}
		for n, mc := range cfg.Databases {
			if mc.Env == "" {
				mc.Env = cfg.Server.Env
			}
			return n, &mc, nil
		}
	}
	if mc, ok := cfg.Databases[name]; ok {
		if mc.Env == "" {
			mc.Env = cfg.Server.Env
		}
		return name, &mc, nil
	}
	return "", nil, fmt.Errorf("连接 '%s' 不存在，使用 -list 查看可用连接", name)
}

func buildDSN(mc *MySQLConn) string {
	if mc.Charset == "" {
		mc.Charset = "utf8mb4"
	}
	if mc.Port == "" {
		mc.Port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mc.UserName, mc.Password, mc.Host, mc.Port, mc.Database, mc.Charset)
}

// printConnList 列出所有可用连接（隐藏密码）
func printConnList(cfg *Config) {
	type ConnInfo struct {
		Name     string `json:"name"`
		Database string `json:"database"`
		Env      string `json:"env"`
	}
	var list []ConnInfo

	env := cfg.Server.Env
	if env == "" {
		env = "unknown"
	}

	for name, mc := range cfg.Databases {
		e := mc.Env
		if e == "" {
			e = env
		}
		list = append(list, ConnInfo{
			Name:     name,
			Database: mc.Database,
			Env:      e,
		})
	}

	out, _ := json.MarshalIndent(list, "", "  ")
	fmt.Println(string(out))
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, nil
}

func main() {
	configPath := flag.String("c", "./config.yaml", "配置文件路径")
	connName := flag.String("n", "", "连接名称（对应 config.yaml 中 databases 下的 key）")
	query := flag.String("e", "", "要执行的 SQL 语句")
	listConn := flag.Bool("l", false, "列出所有可用的数据库连接（不显示密码）")
	flag.Parse()

	if *query == "" && !*listConn {
		fmt.Fprintln(os.Stderr, "用法:")
		fmt.Fprintln(os.Stderr, "  dbctl -l                              列出所有可用连接")
		fmt.Fprintln(os.Stderr, "  dbctl -n <连接名> -e \"<SQL>\"         执行SQL")
		fmt.Fprintln(os.Stderr, "  dbctl -e \"<SQL>\"                    使用默认连接执行SQL")
		os.Exit(1)
	}

	// 加载配置
	cfg, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// 列出连接
	if *listConn {
		printConnList(cfg)
		return
	}

	// 解析连接
	resolvedName, mc, err := resolveConn(cfg, *connName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	dsn := buildDSN(mc)

	// 环境安全检查
	sqlUpper := strings.ToUpper(strings.TrimSpace(*query))
	isWrite := isWriteSQL(sqlUpper)
	if mc.Env == "prod" && isWrite {
		fmt.Fprintln(os.Stderr, "⚠️  当前为生产环境(prod)，拒绝执行写操作！")
		os.Exit(1)
	}

	// 危险语句拦截
	if isDangerousSQL(sqlUpper) {
		fmt.Fprintln(os.Stderr, "🚫 检测到危险SQL（DROP/TRUNCATE/ALTER），已拦截！")
		os.Exit(1)
	}

	// 无 WHERE 的写操作拦截
	if isWrite && !strings.Contains(sqlUpper, "WHERE") && !strings.HasPrefix(sqlUpper, "INSERT") {
		fmt.Fprintln(os.Stderr, "🚫 UPDATE/DELETE 必须包含 WHERE 条件！")
		os.Exit(1)
	}

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接数据库失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "数据库 Ping 失败: %v\n", err)
		os.Exit(1)
	}

	// 打印连接信息到 stderr（不暴露密码）
	fmt.Fprintf(os.Stderr, "[%s] %s:%s/%s (env=%s)\n", resolvedName, mc.Host, mc.Port, mc.Database, mc.Env)

	// 执行 SQL
	if isWrite {
		execWrite(db, *query)
	} else {
		execRead(db, *query)
	}
}

func isWriteSQL(upper string) bool {
	return strings.HasPrefix(upper, "INSERT") ||
		strings.HasPrefix(upper, "UPDATE") ||
		strings.HasPrefix(upper, "DELETE") ||
		strings.HasPrefix(upper, "REPLACE") ||
		strings.HasPrefix(upper, "BEGIN") ||
		strings.HasPrefix(upper, "START")
}

func isDangerousSQL(upper string) bool {
	return strings.Contains(upper, "DROP TABLE") ||
		strings.Contains(upper, "DROP DATABASE") ||
		strings.Contains(upper, "TRUNCATE") ||
		strings.Contains(upper, "ALTER TABLE")
}

func execRead(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "查询失败: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	if len(cols) == 0 {
		fmt.Println("(空结果集)")
		return
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			fmt.Fprintf(os.Stderr, "扫描行失败: %v\n", err)
			os.Exit(1)
		}
		row := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	// 输出 JSON
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	fmt.Fprintf(os.Stderr, "\n共 %d 行\n", len(results))
}

func execWrite(db *sql.DB, query string) {
	result, err := db.Exec(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "执行失败: %v\n", err)
		os.Exit(1)
	}
	affected, _ := result.RowsAffected()
	lastID, _ := result.LastInsertId()
	fmt.Printf("执行成功: affected_rows=%d, last_insert_id=%d\n", affected, lastID)
}
