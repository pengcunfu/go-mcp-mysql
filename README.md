# go-mcp-mysql
[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/pengcunfu/go-mcp-mysql)](https://archestra.ai/mcp-catalog/pengcunfu__go-mcp-mysql)

## 概述

零负担、开箱即用的模型上下文协议（MCP）服务器，用于与 MySQL 交互和自动化操作。无需 Node.js 或 Python 环境。该服务器提供对 MySQL 数据库和表进行 CRUD 操作的工具，以及只读模式以防止意外的写入操作。您还可以通过添加 `--with-explain-check` 标志让 MCP 服务器在执行查询前使用 `EXPLAIN` 语句检查查询计划。

请注意，这是一个正在开发中的项目，可能还不适合生产环境使用。

## 安装

1. 获取最新的 [发布版本](https://github.com/pengcunfu/go-mcp-mysql/releases) 并将其放入您的 `$PATH` 或您可以轻松访问的地方。

2. 或者如果您已安装 Go，可以从源码构建：

```sh
go install -v github.com/pengcunfu/go-mcp-mysql@latest
```

## 使用方法

### 方法 A：使用命令行参数

```json
{
  "mcpServers": {
    "mysql": {
      "command": "go-mcp-mysql",
      "args": [
        "--host", "localhost",
        "--user", "root",
        "--pass", "password",
        "--port", "3306",
        "--db", "mydb"
      ]
    }
  }
}
```

### 方法 B：使用 DSN 和自定义选项

```json
{
  "mcpServers": {
    "mysql": {
      "command": "go-mcp-mysql",
      "args": [
        "--dsn", "username:password@tcp(localhost:3306)/mydb?parseTime=true&loc=Local"
      ]
    }
  }
}
```

更多详情请参考 [MySQL DSN](https://github.com/go-sql-driver/mysql#dsn-data-source-name)。

注意：对于将二进制文件放在 `$PATH` 之外的用户，您需要将 `go-mcp-mysql` 替换为二进制文件的完整路径：例如：如果您将二进制文件放在 **Downloads** 文件夹中，您可以使用以下路径：

```json
{
  "mcpServers": {
    "mysql": {
      "command": "C:\\Users\\<username>\\Downloads\\go-mcp-mysql.exe",
      "args": [
        ...
      ]
    }
  }
}
```

### 可选标志

- 添加 `--read-only` 标志以启用只读模式。在此模式下，只有以 `list`、`read_` 和 `desc_` 开头的工具可用。添加此标志后请确保刷新/重启 MCP 服务器。
- 默认情况下，CRUD 查询将首先使用 `EXPLAIN ?` 语句执行，以检查生成的查询计划是否符合预期模式。添加 `--with-explain-check` 标志以禁用此行为。

## 工具

### 模式工具

1. `list_database`

    - 列出 MySQL 服务器中的所有数据库。
    - 参数：无
    - 返回：匹配的数据库名称列表。

2. `list_table`

    - 列出 MySQL 服务器中的所有表。
    - 参数：
        - `name`：如果提供，列出具有指定名称的表，与 SQL `SHOW TABLES LIKE '%name%'` 相同。否则，列出所有表。
    - 返回：匹配的表名称列表。

3. `create_table`

    - 在 MySQL 服务器中创建新表。
    - 参数：
        - `query`：创建表的 SQL 查询。
    - 返回：受影响的行数。

4. `alter_table`

    - 修改 MySQL 服务器中的现有表。LLM 被告知不要删除现有表或列。
    - 参数：
        - `query`：修改表的 SQL 查询。
    - 返回：受影响的行数。

5. `desc_table`

    - 描述表的结构。
    - 参数：
        - `name`：要描述的表名。
    - 返回：表的结构。

### 数据工具

1. `read_query`

    - 执行只读 SQL 查询。
    - 参数：
        - `query`：要执行的 SQL 查询。
    - 返回：查询结果。

2. `write_query`

    - 执行写入 SQL 查询。
    - 参数：
        - `query`：要执行的 SQL 查询。
    - 返回：受影响的行数，最后插入 ID：<last_insert_id>。

3. `update_query`

    - 执行更新 SQL 查询。
    - 参数：
        - `query`：要执行的 SQL 查询。
    - 返回：受影响的行数。

4. `delete_query`

    - 执行删除 SQL 查询。
    - 参数：
        - `query`：要执行的 SQL 查询。
    - 返回：受影响的行数。

## 许可证

MIT
