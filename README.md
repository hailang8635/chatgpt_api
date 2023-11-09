### 功能说明

微信公众号对接ChatGPT的后台接口



### 使用说明

golang编译后发布到国际云上后暴露API接口，把地址配置到微信后台使用


```
go mod init chatgpt_api
```

## go build

```
### GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
### GOARCH：目标平台的体系架构（386、amd64、arm）
### 当CGO_ENABLED=1， 进行编译时， 会将文件中引用libc的库（比如常用的net包），以动态链接的方式生成目标文件。
### 当CGO_ENABLED=0， 进行编译时， 则会把在目标文件中未定义的符号（外部函数）一起链接到可执行文件中。
  
SET GOOS=linux
SET GOARCH=arm
SET CGO_ENABLED=0

Mac:
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build go_main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build go_main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build go_main.go

```

### 提交记录
