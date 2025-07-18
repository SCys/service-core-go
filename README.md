# Service Core Go

一个高性能的Go服务核心库，提供了常用的工具函数和基础设施组件。

## 功能特性

- **HTTP服务**: 基于Fiber框架的HTTP服务工具
- **数据库操作**: 基于pgx的PostgreSQL数据库操作工具
- **日志系统**: 基于logrus的日志系统，支持文件轮转
- **配置管理**: 基于ini文件的配置管理
- **JSON处理**: 高性能的JSON序列化和反序列化
- **工具函数**: 常用的工具函数集合

## 主要改进

### 1. 修复了测试问题
- 修复了 `basic_fields_test.go` 中指针类型字段处理的问题
- 改进了 `utils_test.go` 中的 `IsZero` 函数测试
- 所有测试现在都能正常通过

### 2. 代码质量优化
- 修复了 `fiber.go` 中的拼写错误 (`X-Forward-For` → `X-Forwarded-For`)
- 改进了 `utils.go` 中的不安全转换函数，添加了安全警告注释
- 优化了 `logging.go` 中的日志配置，使路径可配置
- 改进了 `config.go` 中的全局变量管理，使用线程安全的单例模式

### 3. 新增功能
- 添加了 `IsZeroValue` 函数，支持更通用的零值检查
- 改进了错误处理和日志记录

## 使用方法

### 初始化配置
```go
err := core.InitGlobalConfig("config.ini")
if err != nil {
    log.Fatal(err)
}
```

### 使用日志
```go
core.I("信息日志")
core.W("警告日志")
core.E("错误日志", err)
```

### 使用数据库操作
```go
// 定义模型
type User struct {
    core.BasicFields
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func (u User) TableName() string {
    return "users"
}

// 插入数据
user := User{
    BasicFields: core.NewBasicFields(),
    Name: "张三",
    Age:  25,
}
err := core.PGXInsert(ctx, db, &user)
```

### 使用HTTP工具
```go
app := fiber.New()
app.Use(core.FiberBasicInfo())

app.Get("/api/users", func(c *fiber.Ctx) error {
    // 获取客户端IP
    ip := core.FiberIP(c)
    
    // 返回JSON响应
    return core.FiberJSON(c, core.H{"users": users})
})
```

## 依赖管理

项目使用Go 1.24.0，主要依赖包括：
- `github.com/gofiber/fiber/v2` - HTTP框架
- `github.com/jackc/pgx/v5` - PostgreSQL驱动
- `github.com/sirupsen/logrus` - 日志库
- `github.com/bytedance/sonic` - 高性能JSON库

## 测试

运行所有测试：
```bash
go test -v
```

## 注意事项

1. `String` 和 `Bytes` 函数使用了 `unsafe.Pointer`，仅适用于临时使用
2. 日志文件默认保存在 `./log/main.log`，可通过 `core.LogPath` 变量配置
3. 配置管理使用单例模式，确保线程安全

## 许可证

MIT License 