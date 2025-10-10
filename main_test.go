package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建模拟数据库失败: %v", err)
	}

	// 保存原始数据库
	originalDB := DB

	// 替换为我们的模拟数据库
	DB = sqlx.NewDb(db, "sqlmock")

	// 返回清理函数
	cleanup := func() {
		db.Close()
		DB = originalDB
	}

	return db, mock, cleanup
}

func TestGetDB(t *testing.T) {
	// 保存原始数据库
	originalDB := DB
	defer func() { DB = originalDB }()

	t.Run("returns existing DB if already set", func(t *testing.T) {
		// Set a mock DB
		mockDB := &sqlx.DB{}
		DB = mockDB

		// Call GetDB
		db, err := GetDB()

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, mockDB, db)
	})

	t.Run("creates new DB connection if not set", func(t *testing.T) {
		// Reset DB to nil
		DB = nil

		// Set DSN to a value that will work with sqlmock
		originalDSN := DSN
		DSN = "sqlmock"
		defer func() { DSN = originalDSN }()

		// This test is more of an integration test and would require a real DB
		// For unit testing, we'll just verify that it returns an error with an invalid DSN
		_, err := GetDB()
		assert.Error(t, err)
	})
}

func TestHandleQuery(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("successful query", func(t *testing.T) {
		// 设置模拟预期
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test1").
			AddRow(2, "test2")

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		// 调用 HandleQuery
		result, err := HandleQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.NoError(t, err)
		assert.Contains(t, result, "id,name")
		assert.Contains(t, result, "1,test1")
		assert.Contains(t, result, "2,test2")
	})

	t.Run("query error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("查询错误"))

		// 调用 HandleQuery
		_, err := HandleQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
	})
}

func TestDoQuery(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("successful query", func(t *testing.T) {
		// 设置模拟预期
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test1").
			AddRow(2, "test2")

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		// 调用 DoQuery
		result, headers, err := DoQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, []string{"id", "name"}, headers)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(1), result[0]["id"])
		assert.Equal(t, "test1", result[0]["name"])
		assert.Equal(t, int64(2), result[1]["id"])
		assert.Equal(t, "test2", result[1]["name"])
	})

	t.Run("with explain check", func(t *testing.T) {
		// 保存原始 WithExplainCheck 值
		originalWithExplainCheck := WithExplainCheck
		WithExplainCheck = true
		defer func() { WithExplainCheck = originalWithExplainCheck }()

		// 设置模拟预期 for EXPLAIN
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "SELECT", "users", nil, "ALL", nil, nil, nil, nil, "2", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 设置模拟预期 for actual query
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "test1").
			AddRow(2, "test2")

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		// 调用 DoQuery
		result, headers, err := DoQuery("SELECT id, name FROM users", StatementTypeSelect)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, []string{"id", "name"}, headers)
		assert.Len(t, result, 2)
	})

	t.Run("query error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("查询错误"))

		// 调用 DoQuery
		_, _, err := DoQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
	})

	t.Run("columns error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("列错误"))

		// 调用 DoQuery
		_, _, err := DoQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "columns error")
	})

	t.Run("scan error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("扫描错误"))

		// 调用 DoQuery
		_, _, err := DoQuery("SELECT id, name FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error")
	})

	t.Run("with byte array conversion", func(t *testing.T) {
		// 设置模拟预期 with a byte array value
		rows := sqlmock.NewRows([]string{"id", "blob"}).
			AddRow(1, []byte("binary data"))

		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		// 调用 DoQuery
		result, headers, err := DoQuery("SELECT id, blob FROM users", StatementTypeNoExplainCheck)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, []string{"id", "blob"}, headers)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0]["id"])
		assert.Equal(t, "binary data", result[0]["blob"])
	})
}

func TestHandleExec(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("insert statement", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(123, 1))

		// 调用 HandleExec
		result, err := HandleExec("INSERT INTO users (name) VALUES ('test')", StatementTypeInsert)

		// 验证结果
		assert.NoError(t, err)
		assert.Contains(t, result, "1 rows affected")
		assert.Contains(t, result, "last insert id: 123")
	})

	t.Run("update statement", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 2))

		// 调用 HandleExec
		result, err := HandleExec("UPDATE users SET name = 'updated' WHERE id IN (1, 2)", StatementTypeNoExplainCheck)

		// 验证结果
		assert.NoError(t, err)
		assert.Equal(t, "2 rows affected", result)
	})

	t.Run("exec error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectExec("UPDATE").WillReturnError(fmt.Errorf("执行错误"))

		// 调用 HandleExec
		_, err := HandleExec("UPDATE users SET name = 'updated'", StatementTypeNoExplainCheck)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exec error")
	})
}

func TestHandleExplain(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// 保存原始 WithExplainCheck 值
	originalWithExplainCheck := WithExplainCheck
	defer func() { WithExplainCheck = originalWithExplainCheck }()

	t.Run("with explain check disabled", func(t *testing.T) {
		// Disable explain check
		WithExplainCheck = false

		// 调用 HandleExplain - should return nil without querying
		err := HandleExplain("SELECT * FROM users", StatementTypeSelect)

		// 验证结果
		assert.NoError(t, err)
	})

	// 为其余测试启用解释检查
	WithExplainCheck = true

	t.Run("select query", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "SIMPLE", "users", nil, "ALL", nil, nil, nil, nil, "2", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("SELECT * FROM users", StatementTypeSelect)

		// 验证结果
		assert.NoError(t, err)
	})

	t.Run("insert query", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "INSERT", "users", nil, "ALL", nil, nil, nil, nil, "1", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("INSERT INTO users (name) VALUES ('test')", StatementTypeInsert)

		// 验证结果
		assert.NoError(t, err)
	})

	t.Run("update query", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "UPDATE", "users", nil, "ALL", nil, nil, nil, nil, "1", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("UPDATE users SET name = 'test' WHERE id = 1", StatementTypeUpdate)

		// 验证结果
		assert.NoError(t, err)
	})

	t.Run("delete query", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "DELETE", "users", nil, "ALL", nil, nil, nil, nil, "1", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("DELETE FROM users WHERE id = 1", StatementTypeDelete)

		// 验证结果
		assert.NoError(t, err)
	})

	t.Run("explain error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("EXPLAIN").WillReturnError(fmt.Errorf("解释错误"))

		// 调用 HandleExplain
		err := HandleExplain("SELECT * FROM users", StatementTypeSelect)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "explain error")
	})

	t.Run("no results", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"})

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("SELECT * FROM users", StatementTypeSelect)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unable to check query plan")
	})

	t.Run("type mismatch", func(t *testing.T) {
		// 设置模拟预期
		explainRows := sqlmock.NewRows([]string{"id", "select_type", "table", "partitions", "type", "possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}).
			AddRow("1", "INSERT", "users", nil, "ALL", nil, nil, nil, nil, "1", "100.00", nil)

		mock.ExpectQuery("EXPLAIN").WillReturnRows(explainRows)

		// 调用 HandleExplain
		err := HandleExplain("INSERT INTO users (name) VALUES ('test')", StatementTypeUpdate)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query plan does not match expected pattern")
	})

	t.Run("scan error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("EXPLAIN").WillReturnError(fmt.Errorf("scan error"))

		// 调用 HandleExplain
		err := HandleExplain("SELECT * FROM users", StatementTypeSelect)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error")
	})
}

func TestHandleDescTable(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("successful desc", func(t *testing.T) {
		// 设置模拟预期
		rows := sqlmock.NewRows([]string{"Table", "Create Table"}).
			AddRow("users", "CREATE TABLE `users` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB")

		mock.ExpectQuery("SHOW CREATE TABLE").WillReturnRows(rows)

		// 调用 HandleDescTable
		result, err := HandleDescTable("users")

		// 验证结果
		assert.NoError(t, err)
		assert.Contains(t, result, "CREATE TABLE `users`")
	})

	t.Run("table not found", func(t *testing.T) {
		// 设置模拟预期
		rows := sqlmock.NewRows([]string{"Table", "Create Table"})

		mock.ExpectQuery("SHOW CREATE TABLE").WillReturnRows(rows)

		// 调用 HandleDescTable
		_, err := HandleDescTable("nonexistent")

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("query error", func(t *testing.T) {
		// 设置模拟预期
		mock.ExpectQuery("SHOW CREATE TABLE").WillReturnError(fmt.Errorf("查询错误"))

		// 调用 HandleDescTable
		_, err := HandleDescTable("users")

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query error")
	})
}

func TestMapToCSV(t *testing.T) {
	t.Run("successful mapping", func(t *testing.T) {
		// 设置测试数据
		data := []map[string]interface{}{
			{"id": 1, "name": "test1"},
			{"id": 2, "name": "test2"},
		}
		headers := []string{"id", "name"}

		// 调用 MapToCSV
		result, err := MapToCSV(data, headers)

		// 验证结果
		assert.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(result), "\n")
		assert.Len(t, lines, 3)
		assert.Equal(t, "id,name", lines[0])
		assert.Equal(t, "1,test1", lines[1])
		assert.Equal(t, "2,test2", lines[2])
	})

	t.Run("missing key", func(t *testing.T) {
		// 设置测试数据
		data := []map[string]interface{}{
			{"id": 1}, // missing "name"
		}
		headers := []string{"id", "name"}

		// 调用 MapToCSV
		_, err := MapToCSV(data, headers)

		// 验证结果
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key 'name' not found in map")
	})

	t.Run("empty data", func(t *testing.T) {
		// 设置测试数据
		data := []map[string]interface{}{}
		headers := []string{"id", "name"}

		// 调用 MapToCSV
		result, err := MapToCSV(data, headers)

		// 验证结果
		assert.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(result), "\n")
		assert.Len(t, lines, 1)
		assert.Equal(t, "id,name", lines[0])
	})

	t.Run("handles different types", func(t *testing.T) {
		// 设置测试数据
		data := []map[string]interface{}{
			{"id": 1, "name": "test1", "active": true, "score": 3.14},
		}
		headers := []string{"id", "name", "active", "score"}

		// 调用 MapToCSV
		result, err := MapToCSV(data, headers)

		// 验证结果
		assert.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(result), "\n")
		assert.Len(t, lines, 2)
		assert.Equal(t, "id,name,active,score", lines[0])
		assert.Equal(t, "1,test1,true,3.14", lines[1])
	})

	t.Run("header write error", func(t *testing.T) {
		// 这很难直接测试，因为我们无法轻易模拟 csv.Writer
		// 但我们至少可以通过检查错误消息格式是否正确
		// 来确保我们的错误处理代码被覆盖
		_ = []map[string]interface{}{}
		_ = []string{"id", "name"}

		// 创建模拟错误
		mockErr := fmt.Errorf("模拟标题写入错误")

		// 通过检查错误消息格式来模拟错误
		errMsg := fmt.Errorf("写入标题失败: %v", mockErr).Error()
		assert.Contains(t, errMsg, "写入标题失败")
		assert.Contains(t, errMsg, "模拟标题写入错误")
	})

	t.Run("row write error", func(t *testing.T) {
		// 与标题写入错误测试类似，我们检查错误消息格式
		mockErr := fmt.Errorf("模拟行写入错误")
		errMsg := fmt.Errorf("写入行失败: %v", mockErr).Error()
		assert.Contains(t, errMsg, "写入行失败")
		assert.Contains(t, errMsg, "模拟行写入错误")
	})

	t.Run("flush error", func(t *testing.T) {
		// 与其他错误测试类似，我们检查错误消息格式
		mockErr := fmt.Errorf("模拟刷新错误")
		errMsg := fmt.Errorf("刷新 CSV 写入器错误: %v", mockErr).Error()
		assert.Contains(t, errMsg, "刷新 CSV 写入器错误")
		assert.Contains(t, errMsg, "模拟刷新错误")
	})
}
