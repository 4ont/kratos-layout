# First of all

**执行 make init**
它会帮你安装好protobuf相关的工具和插件

**在你的IDE里执行全局的替换（refactor）将`kratos-layout`修改为你的工程名。**

工程内的提供了一个简单的示例: api/example. 其中定义一个sample服务。你可以参照它进行需求开发。
当你的接口能够工作后，可以删除api/example并重新生成接口代码: `make api`, 然后清理example的其他代码使工程通过编译：`make build`
小结一下工作的步骤：
- 完成你的第一个接口
- 删除 api/example
- 执行`make build`
- 修复编译错误
- 重复上面两个步骤直到编译成功

# 框架说明 
该项目使用了 [kratos][1] 微服务框架进行开发。
本项目运用 protobuf 接口描述语言定义对外提供的服务接口和数据模型，
提供 HTTP 和 gRPC 协议的微服务接口。
此外，protobuf 还被运用于错误码定义和调用内部网关 API 接口。
本项目还使用了编译时依赖注入技术，通过构造函数注入，降低编写复杂的应用初始化代码的难度。
本项目实现了基于 [opentelemetry][3] 和 [jaeger][4] 的链路追踪，方便应用排错和性能优化。

本项目参考了 DDD 架构分层思路，将应用内部组织为单向依赖的四层结构：
- 服务及事件层
- 业务层
- 数据聚合及微服务访问代理层
- 数据访层
其中业务层是关键，主要的业务逻辑应该在这里分派，业务层和数据层通过interface与model的定义将biz和repo层的依赖关系解耗，
当项目的代码规模变大时你会发现这个设计的好处的。
 
# Kratos Project Template

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```

### 生成 gRPC 服务客户端和服务端代码

每当我们更新 api protobuf 文件时，我们需要使用代码生成工具生成客户端。
当我们增加服务或方法时，还需要生成服务端代码。
我们可以利用 kratos 命令行工具帮助我们实现这个目标。
生成 gRPC 客户端代码，请执行以下命令：

    cd <repo-root-dir>
    kratos proto client api/<path>/<filename>.proto


生成 gRPC 服务端代码，请执行以下命令：

    cd <repo-root-dir>
	kratos proto server api/<path>/<filename>.proto -t internal/service/

### 服务注册
kratos框架默认支持grpc和http两种服务暴露方式。一个service要向外提供服务前需要先向框架进行注册：

编辑 internal/server/grpc.go

```go
    ...
srv := grpc.NewServer(opts...)
probeapi.RegisterProbeServer(srv, probe)
```

编辑 internal/server/http.go

```go
    ...
srv := http.NewServer(opts...)
probeapi.RegisterProbeHTTPServer(srv, probe)
```

## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```


## Config
支持多种配置数据源：
- 本地配置文件
- aws appconfig

### local config
通过命令行参数来指定配置源：
```
Usage of ./bin/kratos-layout:
  -conf string
        config path, eg: -conf config.yaml (default "../../configs")
  -conf-src string
        config source, eg: -conf-src aws-appconfig (default "file")

```
conf-src 默认为: file, 将 conf参数指定的路径读取配置文件

当conf-src为 aws-appconfig 时，会从以下环境变量读取aws appconfig的元配置信息以便连接到appconfig服务
- AWS_REGION=ap-southeast-1
- AWS_ROLE_ARN=arn:aws:iam::302759042447:role/eks_aws_access_role
- AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token
- AWS_ENVIRONMENT=t2etu44
- AWS_CONFIGURATION_ID=1a37jep
- AWS_APPLICATION_ID=o7i5bvr

AWS_WEB_IDENTITY_TOKEN_FILE指定的aws identity token内容在本地的保存路径, **根据实际情况修改以上示例的值**


---

[1]: https://go-kratos.dev
[2]: https://git.ucloudadmin.com/devops/virtualenv
[3]: https://opentelemetry.io/
[4]: https://www.jaegertracing.io/