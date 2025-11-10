# go-mcp-mysql

<div align="center">

[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/pengcunfu/go-mcp-mysql)](https://archestra.ai/mcp-catalog/pengcunfu__go-mcp-mysql)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/pengcunfu/go-mcp-mysql)](https://goreportcard.com/report/github.com/pengcunfu/go-mcp-mysql)
[![GitHub release](https://img.shields.io/github/release/pengcunfu/go-mcp-mysql.svg)](https://github.com/pengcunfu/go-mcp-mysql/releases)

**零负担、开箱即用的 MySQL MCP 服务器**

[English](README_EN.md) | 简体中文

</div>

---

## ✨ 特性

- 🚀 **零依赖**：无需 Node.js 或 Python 环境，单一可执行文件
- 🔒 **安全可靠**：支持只读模式，防止意外写入操作
- 📊 **查询优化**：可选的 EXPLAIN 检查，优化查询性能
- 🛠️ **完整 CRUD**：支持数据库和表的完整生命周期管理
- ⚡ **高性能**：基于 Go 语言开发，性能卓越
- 🎯 **易于使用**：支持命令行参数和 DSN 两种配置方式

> ⚠️ **注意**：本项目正在积极开发中，建议在生产环境使用前进行充分测试。

## 📦 安装

### 方式一：下载预编译二进制文件

从 [Releases 页面](https://github.com/pengcunfu/go-mcp-mysql/releases) 下载适合您操作系统的最新版本，并将其放入 `$PATH` 或您可以轻松访问的位置。

### 方式二：从源码构建

如果您已安装 Go 1.21 或更高版本：

```bash
go install -v github.com/pengcunfu/go-mcp-mysql@latest
```

## 🚀 快速开始

### 配置方式 A：使用命令行参数

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

### 配置方式 B：使用 DSN 连接字符串

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

> 💡 **提示**：更多 DSN 配置选项请参考 [MySQL DSN 文档](https://github.com/go-sql-driver/mysql#dsn-data-source-name)。

### 使用绝对路径

如果二进制文件不在 `$PATH` 中，需要使用完整路径。例如，Windows 用户可以这样配置：

```json
{
  "mcpServers": {
    "mysql": {
      "command": "C:\\Users\\<username>\\Downloads\\go-mcp-mysql.exe",
      "args": ["--host", "localhost", "--user", "root", "--pass", "password"]
    }
  }
}
```

## ⚙️ 配置选项

### 可选标志

| 标志 | 说明 |
|------|------|
| `--read-only` | 启用只读模式，仅允许 `list`、`read_` 和 `desc_` 开头的工具，防止数据修改 |
| `--with-explain-check` | 在执行 CRUD 查询前使用 `EXPLAIN` 检查查询计划，帮助优化性能 |

> 📌 **注意**：修改标志后需要重启 MCP 服务器才能生效。

## 🛠️ 可用工具

### 数据库模式管理

#### `list_database`
列出 MySQL 服务器中的所有数据库。
- **参数**：无
- **返回**：数据库名称列表

#### `list_table`
列出 MySQL 服务器中的所有表。
- **参数**：
  - `name`（可选）：表名过滤条件，等同于 `SHOW TABLES LIKE '%name%'`
- **返回**：匹配的表名称列表

#### `create_table`
在 MySQL 服务器中创建新表。
- **参数**：
  - `query`：CREATE TABLE SQL 语句
- **返回**：受影响的行数

#### `alter_table`
修改现有表结构（不支持删除表或列）。
- **参数**：
  - `query`：ALTER TABLE SQL 语句
- **返回**：受影响的行数

#### `desc_table`
查看表结构详情。
- **参数**：
  - `name`：表名
- **返回**：表的结构信息

### 数据操作

#### `read_query`
执行只读 SQL 查询（SELECT）。
- **参数**：
  - `query`：SELECT SQL 语句
- **返回**：查询结果集

#### `write_query`
执行写入 SQL 查询（INSERT）。
- **参数**：
  - `query`：INSERT SQL 语句
- **返回**：受影响的行数和最后插入的 ID

#### `update_query`
执行更新 SQL 查询（UPDATE）。
- **参数**：
  - `query`：UPDATE SQL 语句
- **返回**：受影响的行数

#### `delete_query`
执行删除 SQL 查询（DELETE）。
- **参数**：
  - `query`：DELETE SQL 语句
- **返回**：受影响的行数

## 🤝 贡献

欢迎贡献！如果您有任何想法、建议或发现了 bug，请：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📝 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

## 🔗 相关链接

- [GitHub 仓库](https://github.com/pengcunfu/go-mcp-mysql)
- [问题反馈](https://github.com/pengcunfu/go-mcp-mysql/issues)
- [MCP 协议文档](https://modelcontextprotocol.io/)
- [MySQL 驱动文档](https://github.com/go-sql-driver/mysql)

## 👤 作者

**pengcunfu**

- GitHub: [@pengcunfu](https://github.com/pengcunfu)

---

<div align="center">

如果这个项目对您有帮助，请给它一个 ⭐️！

</div>
