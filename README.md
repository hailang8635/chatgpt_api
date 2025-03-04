这个仓库是一个关于微信公众号对接ChatGPT的后台接口项目，主要实现了微信公众号与ChatGPT的对接功能。以下是对该仓库的详细介绍：

### 项目结构
```
chatgpt_api
├── .gitignore
├── README.md
├── chatgpt_api.go
├── go.mod
├── go.sum
├── domain
│   └── domain.go
├── wechat_server
│   ├── gpt_api_from_browser.go
│   ├── gpt_api_from_curl.go
│   ├── gpt_api_from_wechat.go
│   ├── gpt_api_from_wechat_v01.go
│   ├── item_helper.go
│   └── model_api_route.go
├── api_from_ai
│   ├── api_error.go
│   ├── chatgpt_api_test.go
│   ├── handler_chatgpt_api.go
│   ├── handler_deepseek_api.go
│   └── handler_glm_api.go
├── utils
│   ├── db_utils_gorm.go
│   ├── dbutils
│   │   ├── gorm_mysql_demo.go
│   │   └── jmoiron_sqlx_demo.go
│   └── string_utils.go
├── http_test
├── img
├── config
│   ├── ban_words.txt
│   ├── config_files
│   └── init_properties.go
```

### 主要文件及功能说明
- **README.md**：项目的说明文档，包含功能说明、使用说明、编译命令等信息。
- **chatgpt_api.go**：项目的入口文件。
- **go.mod 和 go.sum**：Go模块的配置文件，用于管理项目的依赖。
- **domain/domain.go**：定义了项目中使用的一些数据结构，如`KeywordAndAnswerItem`等。
- **wechat_server 目录**：包含处理微信相关请求的代码，如`chatHandlerWithDB`函数用于处理微信的聊天请求。
- **api_from_ai 目录**：包含与AI接口交互的代码，如处理DeepSeek、ChatGPT、GLM等API的请求。
- **utils 目录**：包含一些工具类的代码，如数据库操作、字符串处理等。其中`dbutils`子目录包含了使用`gorm`和`sqlx`进行数据库操作的示例代码。
- **config 目录**：包含项目的配置文件，如`ban_words.txt`用于存储禁止的词汇，`init_properties.go`用于初始化项目的配置。

### 主要功能模块
- **数据库操作**：
  - 使用`gorm`和`sqlx`进行数据库操作，如创建表、插入数据、查询数据等。
  - 包含了对`UserInfo`和`KeywordAndAnswerItem`等数据结构的数据库操作。
- **微信请求处理**：
  - 处理微信公众号的请求，如校验微信平台、读取用户请求数据、处理聊天请求等。
  - 支持GPT-4（TODO）。
- **AI接口交互**：
  - 与DeepSeek、ChatGPT、GLM等AI接口进行交互，发送请求并处理响应。

### 使用说明
1. 设置SDK：在IDE中设置GOROOT和Go Modules。
2. 初始化项目：运行`go mod init chatgpt_api`。
3. 编译项目：根据目标平台的操作系统和体系架构，使用不同的编译命令，如：
   - Mac:
     - `CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build chatgpt_api.go`
     - `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build chatgpt_api.go`
     - `CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build chatgpt_api.go`
4. 发布到国际云：将编译后的文件发布到国际云上，暴露API接口。
5. 配置微信后台：将API接口地址配置到微信后台。

通过以上步骤，就可以实现微信公众号与ChatGPT的对接。