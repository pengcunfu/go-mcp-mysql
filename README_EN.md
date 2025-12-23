# go-mcp-mysql

<div align="center">

[![Trust Score](https://archestra.ai/mcp-catalog/api/badge/quality/pengcunfu/go-mcp-mysql)](https://archestra.ai/mcp-catalog/pengcunfu__go-mcp-mysql)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/pengcunfu/go-mcp-mysql)](https://goreportcard.com/report/github.com/pengcunfu/go-mcp-mysql)
[![GitHub release](https://img.shields.io/github/release/pengcunfu/go-mcp-mysql.svg)](https://github.com/pengcunfu/go-mcp-mysql/releases)

**Zero-dependency, ready-to-use MySQL MCP Server**

English | [简体中文](README.md)

</div>

---

## ✨ Features

- 🚀 **Zero Dependencies**: No Node.js or Python required, single executable binary
- 🔒 **Safe & Secure**: Read-only mode support to prevent accidental write operations
- 📊 **Query Optimization**: Optional EXPLAIN check for query performance optimization
- 🛠️ **Full CRUD**: Complete lifecycle management for databases and tables
- ⚡ **High Performance**: Built with Go for exceptional performance
- 🎯 **Easy to Use**: Support for both command-line arguments and DSN configuration

> ⚠️ **Note**: This project is under active development. Please test thoroughly before using in production.

## 📦 Installation

### Option 1: Download Pre-built Binary

Download the latest release for your operating system from the [Releases page](https://github.com/pengcunfu/go-mcp-mysql/releases) and place it in your `$PATH` or an easily accessible location.

### Option 2: Build from Source

If you have Go 1.21 or higher installed:

```bash
go install -v github.com/pengcunfu/go-mcp-mysql@latest
```

## 🚀 Quick Start

### Configuration Method A: Using Command-Line Arguments

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

### Configuration Method B: Using DSN Connection String

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

> 💡 **Tip**: For more DSN configuration options, refer to the [MySQL DSN Documentation](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

### Using Absolute Path

If the binary is not in your `$PATH`, use the full path. For example, Windows users can configure it like this:

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

## ⚙️ Configuration Options

### Optional Flags

| Flag | Description |
|------|-------------|
| `--read-only` | Enable read-only mode, allowing only tools starting with `list`, `read_`, and `desc_` to prevent data modification |
| `--with-explain-check` | Use `EXPLAIN` to check query plans before executing CRUD queries for performance optimization |

> 📌 **Note**: You need to restart the MCP server after changing flags for them to take effect.

## 🛠️ Available Tools

### Database Schema Management

#### `list_database`
List all databases in the MySQL server.
- **Parameters**: None
- **Returns**: List of database names

#### `list_table`
List all tables in the MySQL server.
- **Parameters**:
  - `name` (optional): Table name filter, equivalent to `SHOW TABLES LIKE '%name%'`
- **Returns**: List of matching table names

#### `create_table`
Create a new table in the MySQL server.
- **Parameters**:
  - `query`: CREATE TABLE SQL statement
- **Returns**: Number of affected rows

#### `alter_table`
Modify existing table structure (does not support dropping tables or columns).
- **Parameters**:
  - `query`: ALTER TABLE SQL statement
- **Returns**: Number of affected rows

#### `desc_table`
View table structure details.
- **Parameters**:
  - `name`: Table name
- **Returns**: Table structure information

#### `use_database`
Select the current database to use. Executes a `USE database` statement to switch databases.
- **Parameters**:
  - `name`: Database name to use
- **Returns**: Operation result message

> 💡 **Tip**: If you don't specify the `--db` parameter during configuration, you can use this tool to select a database after connecting.

### Data Operations

#### `read_query`
Execute read-only SQL queries (SELECT).
- **Parameters**:
  - `query`: SELECT SQL statement
- **Returns**: Query result set

#### `write_query`
Execute write SQL queries (INSERT).
- **Parameters**:
  - `query`: INSERT SQL statement
- **Returns**: Number of affected rows and last insert ID

#### `update_query`
Execute update SQL queries (UPDATE).
- **Parameters**:
  - `query`: UPDATE SQL statement
- **Returns**: Number of affected rows

#### `delete_query`
Execute delete SQL queries (DELETE).
- **Parameters**:
  - `query`: DELETE SQL statement
- **Returns**: Number of affected rows

## 🤝 Contributing

Contributions are welcome! If you have any ideas, suggestions, or find bugs, please:

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the [Apache License 2.0](LICENSE).

## 🔗 Related Links

- [GitHub Repository](https://github.com/pengcunfu/go-mcp-mysql)
- [Issue Tracker](https://github.com/pengcunfu/go-mcp-mysql/issues)
- [MCP Protocol Documentation](https://modelcontextprotocol.io/)
- [MySQL Driver Documentation](https://github.com/go-sql-driver/mysql)

## 👤 Author

**pengcunfu**

- GitHub: [@pengcunfu](https://github.com/pengcunfu)

---

<div align="center">

If this project helps you, please give it a ⭐️!

</div>
