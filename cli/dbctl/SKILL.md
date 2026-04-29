---
name: mysql-operation
description: 查询和修改项目 MySQL 数据库内容。当用户需要查看数据库表数据、执行 SQL 查询、更新/插入/删除数据、生成gorm model、或排查数据问题时使用此技能。触发词：查数据库、改数据、SQL查询、表数据、数据修复。
---

# 数据库查询与修改

## CLI 工具

**命令执行**: `dbctl`
**配置**: `E:\Workspace\me\ai-dance\cli\dbctl\config.yaml`dbctl 专用，与项目 config 独立）

### 第一步：查看可用连接

首次使用时，先列出所有可用连接：

```bash
dbctl -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -l
```

输出示例（不含密码）：

```json
[
  {"name": "school_idea", "database": "school_idea", "env": "dev"}
]
```

### 第二步：通过连接名执行 SQL

```bash
# 使用默认连接（仅一个连接时自动选中）
dbctl -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -e "<SQL>"

# 指定连接名称
dbctl -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -n "<NAME>" -e "<SQL>"
```

**参数说明**:

- `-c` 配置文件路径，每次必须带上这个参数且参数必须是 -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml
- `-n` 连接名称（对应 config 中 `databases` 下的 key，仅一个时可省略）
- `-e` 要执行的 SQL 语句
- `-l` 列出所有可用连接

### 输出格式

- SELECT 查询：JSON 数组，stderr 显示连接信息和行数
- 写操作：输出 `affected_rows` 和 `last_insert_id`

## 内置安全防护

1. **生产环境拦截** — `env=prod` 时自动拒绝写操作
2. **危险语句拦截** — 拦截 `DROP TABLE`、`DROP DATABASE`、`TRUNCATE`、`ALTER TABLE`
3. **无 WHERE 拦截** — 拦截无 WHERE 的 `UPDATE`/`DELETE`

## 操作规范

### 查询（直接执行）

```bash
dbctl -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -e "SHOW TABLES;"
dbctl -c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -e "DESCRIBE sc_exam_record;"
dbctl --c E:\Workspace\me\ai-dance\cli\dbctl\config.yaml -e "SELECT id, user_id FROM sc_exam_record WHERE user_id=123 LIMIT 20;"
```

### 修改（需用户确认）

1. **先 SELECT 确认** 影响范围
2. **告知用户** 影响行数
3. **用户确认后** 再执行写操作

# 

## sql生成model规则

### 生成的tag规则

gorm tag带上字段名称，字段类型，默认值、以及字段注释

字段带上json tag 并且名称是小驼峰命名

### 数据格式转化规则

myql bigint 转化为 int64

mysql int  转化为 int32

mysql varchar 转化为 string

mysql datetime（created_at和updated_at 字段转化为time.Time, deleted_at字段转化为gorm.DeletedAt）转化为 *time.Time

其他的mysql类型转化的时候问我一下