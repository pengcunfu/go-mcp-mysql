package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	StatementTypeNoExplainCheck = ""
	StatementTypeSelect         = "SELECT"
	StatementTypeInsert         = "INSERT"
	StatementTypeUpdate         = "UPDATE"
	StatementTypeDelete         = "DELETE"
)

var (
	Host string
	User string
	Pass string
	Port int
	Db   string

	DSN string

	ReadOnly         bool
	WithExplainCheck bool

	DB *sqlx.DB
)

type ExplainResult struct {
	Id           *string `db:"id"`
	SelectType   *string `db:"select_type"`
	Table        *string `db:"table"`
	Partitions   *string `db:"partitions"`
	Type         *string `db:"type"`
	PossibleKeys *string `db:"possible_keys"`
	Key          *string `db:"key"`
	KeyLen       *string `db:"key_len"`
	Ref          *string `db:"ref"`
	Rows         *string `db:"rows"`
	Filtered     *string `db:"filtered"`
	Extra        *string `db:"Extra"`
}

type ShowCreateTableResult struct {
	Table       string `db:"Table"`
	CreateTable string `db:"Create Table"`
}

func main() {
	flag.StringVar(&Host, "host", "localhost", "MySQL 主机名")
	flag.StringVar(&User, "user", "root", "MySQL 用户名")
	flag.StringVar(&Pass, "pass", "", "MySQL 密码")
	flag.IntVar(&Port, "port", 3306, "MySQL 端口")
	flag.StringVar(&Db, "db", "", "MySQL 数据库")

	flag.StringVar(&DSN, "dsn", "", "MySQL DSN")

	flag.BoolVar(&ReadOnly, "read-only", false, "启用只读模式")
	flag.BoolVar(&WithExplainCheck, "with-explain-check", false, "执行前使用 `EXPLAIN` 检查查询计划")
	flag.Parse()

	if len(DSN) == 0 {
		DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", User, Pass, Host, Port, Db)
	}

	s := server.NewMCPServer(
		"go-mcp-mysql",
		"0.1.0",
	)

	// 模式工具
	listDatabaseTool := mcp.NewTool(
		"list_database",
		mcp.WithDescription("列出 MySQL 服务器中的所有数据库"),
	)

	listTableTool := mcp.NewTool(
		"list_table",
		mcp.WithDescription("列出 MySQL 服务器中的所有表"),
	)

	createTableTool := mcp.NewTool(
		"create_table",
		mcp.WithDescription("在 MySQL 服务器中创建新表。确保为每个列和表本身添加了适当的注释"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("创建表的 SQL 查询"),
		),
	)

	alterTableTool := mcp.NewTool(
		"alter_table",
		mcp.WithDescription("修改 MySQL 服务器中的现有表。确保为每个修改的列更新了注释。不要删除表或现有列！"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("修改表的 SQL 查询"),
		),
	)

	descTableTool := mcp.NewTool(
		"desc_table",
		mcp.WithDescription("描述表的结构"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("要描述的表名"),
		),
	)

	useDatabaseTool := mcp.NewTool(
		"use_database",
		mcp.WithDescription("选择当前使用的数据库。执行 USE database 语句切换数据库"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("要使用的数据库名"),
		),
	)

	// 数据工具
	readQueryTool := mcp.NewTool(
		"read_query",
		mcp.WithDescription("执行只读 SQL 查询。在编写 WHERE 条件之前确保了解表结构。如有必要请先调用 `desc_table`"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("要执行的 SQL 查询"),
		),
	)

	writeQueryTool := mcp.NewTool(
		"write_query",
		mcp.WithDescription("执行写入 SQL 查询。执行查询前确保了解表结构。确保数据类型与列定义匹配"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("要执行的 SQL 查询"),
		),
	)

	updateQueryTool := mcp.NewTool(
		"update_query",
		mcp.WithDescription("执行更新 SQL 查询。执行查询前确保了解表结构。确保始终有 WHERE 条件。如有必要请先调用 `desc_table`"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("要执行的 SQL 查询"),
		),
	)

	deleteQueryTool := mcp.NewTool(
		"delete_query",
		mcp.WithDescription("执行删除 SQL 查询。执行查询前确保了解表结构。确保始终有 WHERE 条件。如有必要请先调用 `desc_table`"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("要执行的 SQL 查询"),
		),
	)

	s.AddTool(listDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery("SHOW DATABASES", StatementTypeNoExplainCheck)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(listTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery("SHOW TABLES", StatementTypeNoExplainCheck)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	if !ReadOnly {
		s.AddTool(createTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeNoExplainCheck)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(alterTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeNoExplainCheck)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	s.AddTool(descTableTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleDescTable(request.Params.Arguments["name"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(useDatabaseTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleUseDatabase(request.Params.Arguments["name"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	s.AddTool(readQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := HandleQuery(request.Params.Arguments["query"].(string), StatementTypeSelect)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(result), nil
	})

	if !ReadOnly {
		s.AddTool(writeQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeInsert)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(updateQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeUpdate)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if !ReadOnly {
		s.AddTool(deleteQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			result, err := HandleExec(request.Params.Arguments["query"].(string), StatementTypeDelete)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			return mcp.NewToolResultText(result), nil
		})
	}

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}

func GetDB() (*sqlx.DB, error) {
	if DB != nil {
		return DB, nil
	}

	db, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		return nil, fmt.Errorf("建立数据库连接失败: %v", err)
	}

	DB = db

	return DB, nil
}

func HandleQuery(query, expect string) (string, error) {
	result, headers, err := DoQuery(query, expect)
	if err != nil {
		return "", err
	}

	s, err := MapToCSV(result, headers)
	if err != nil {
		return "", err
	}

	return s, nil
}

func DoQuery(query, expect string) ([]map[string]interface{}, []string, error) {
	db, err := GetDB()
	if err != nil {
		return nil, nil, err
	}

	if len(expect) > 0 {
		if err := HandleExplain(query, expect); err != nil {
			return nil, nil, err
		}
	}

	rows, err := db.Queryx(query)
	if err != nil {
		return nil, nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	result := []map[string]interface{}{}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return nil, nil, err
		}

		resultRow := map[string]interface{}{}
		for i, col := range cols {
			switch v := row[i].(type) {
			case []byte:
				resultRow[col] = string(v)
			default:
				resultRow[col] = v
			}
		}
		result = append(result, resultRow)
	}

	return result, cols, nil
}

func HandleExec(query, expect string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	if len(expect) > 0 {
		if err := HandleExplain(query, expect); err != nil {
			return "", err
		}
	}

	result, err := db.Exec(query)
	if err != nil {
		return "", err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	switch expect {
	case StatementTypeInsert:
		li, err := result.LastInsertId()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%d rows affected, last insert id: %d", ra, li), nil
	default:
		return fmt.Sprintf("%d rows affected", ra), nil
	}
}

func HandleExplain(query, expect string) error {
	if !WithExplainCheck {
		return nil
	}

	db, err := GetDB()
	if err != nil {
		return err
	}

	rows, err := db.Queryx(fmt.Sprintf("EXPLAIN %s", query))
	if err != nil {
		return err
	}

	result := []ExplainResult{}
	for rows.Next() {
		var row ExplainResult
		if err := rows.StructScan(&row); err != nil {
			return err
		}
		result = append(result, row)
	}

	if len(result) != 1 {
		return fmt.Errorf("无法检查查询计划，拒绝执行")
	}

	match := false
	switch expect {
	case StatementTypeInsert:
		fallthrough
	case StatementTypeUpdate:
		fallthrough
	case StatementTypeDelete:
		if *result[0].SelectType == expect {
			match = true
		}
	default:
		// for SELECT type query, the select_type will be multiple values
		// here we check if it's not INSERT, UPDATE or DELETE
		match = true
		for _, typ := range []string{StatementTypeInsert, StatementTypeUpdate, StatementTypeDelete} {
			if *result[0].SelectType == typ {
				match = false
				break
			}
		}
	}

	if !match {
		return fmt.Errorf("查询计划不符合预期模式，拒绝执行")
	}

	return nil
}

func HandleDescTable(name string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	rows, err := db.Queryx(fmt.Sprintf("SHOW CREATE TABLE %s", name))
	if err != nil {
		return "", err
	}

	result := []ShowCreateTableResult{}
	for rows.Next() {
		var row ShowCreateTableResult
		if err := rows.StructScan(&row); err != nil {
			return "", err
		}
		result = append(result, row)
	}

	if len(result) == 0 {
		return "", fmt.Errorf("表 %s 不存在", name)
	}

	return result[0].CreateTable, nil
}

func HandleUseDatabase(name string) (string, error) {
	db, err := GetDB()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(fmt.Sprintf("USE `%s`", name))
	if err != nil {
		return "", fmt.Errorf("切换数据库失败: %v", err)
	}

	return fmt.Sprintf("已成功切换到数据库: %s", name), nil
}

func MapToCSV(m []map[string]interface{}, headers []string) (string, error) {
	var csvBuf strings.Builder
	writer := csv.NewWriter(&csvBuf)

	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("写入标题失败: %v", err)
	}

	for _, item := range m {
		row := make([]string, len(headers))
		for i, header := range headers {
			value, exists := item[header]
			if !exists {
				return "", fmt.Errorf("在映射中未找到键 '%s'", header)
			}
			row[i] = fmt.Sprintf("%v", value)
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("写入行失败: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("刷新 CSV 写入器错误: %v", err)
	}

	return csvBuf.String(), nil
}
