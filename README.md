# Cloudflare DDNS 工具

一个基于Go语言的Cloudflare动态DNS（DDNS）工具，用于自动更新Cloudflare DNS记录以指向当前公网IP地址。

## 功能特性

- 🔄 **自动IP检测**：通过ip.sb服务获取当前公网IP地址
- 📝 **DNS记录管理**：自动创建或更新Cloudflare DNS A记录
- ⏱️ **定时检查**：可配置的IP检查间隔时间
- 🐳 **Docker支持**：提供Docker镜像便于部署
- 🔧 **灵活配置**：支持命令行参数和环境变量配置
- 🛡️ **错误处理**：完善的错误处理和重试机制

## 项目结构

```
cloudflare/
├── main.go          # 主程序入口
├── go.mod           # Go模块依赖管理
├── go.sum           # 依赖校验文件
├── Dockerfile       # Docker构建配置
├── Makefile         # 构建脚本
└── README.md        # 项目说明文档
```

## 快速开始

### 1. 构建项目

```bash
# 使用Makefile构建
make build

# 或直接使用go build
go build -o bin/cloudflare main.go
```

### 2. 运行程序

```bash
# 使用命令行参数
./bin/cloudflare \
  -CF_API_EMAIL=your-email@example.com \
  -CF_API_KEY=your-api-key \
  -CF_ZONE_NAME=example.com \
  -CF_DNS_NAME=subdomain \
  -CF_DNS_TTL=120 \
  -IP_CHECK_DURATION=5m \
  -CF_TIMEOUT=30s

# 或使用环境变量
export CF_API_EMAIL=your-email@example.com
export CF_API_KEY=your-api-key
export CF_ZONE_NAME=example.com
export CF_DNS_NAME=subdomain
export CF_DNS_TTL=120
export IP_CHECK_DURATION=5m
export CF_TIMEOUT=30s

./bin/cloudflare
```

### 3. 使用Docker

```bash
# 构建Docker镜像
docker build -t cloudflare-ddns .

# 运行容器
docker run -d \
  -e CF_API_EMAIL=your-email@example.com \
  -e CF_API_KEY=your-api-key \
  -e CF_ZONE_NAME=example.com \
  -e CF_DNS_NAME=subdomain \
  -e CF_DNS_TTL=120 \
  -e IP_CHECK_DURATION=5m \
  -e CF_TIMEOUT=30s \
  cloudflare-ddns
```

## 配置参数

### 必需参数

| 参数 | 环境变量 | 描述 | 默认值 |
|------|----------|------|--------|
| `-CF_API_KEY` | `CF_API_KEY` | Cloudflare API密钥 | 无 |
| `-CF_API_EMAIL` | `CF_API_EMAIL` | Cloudflare账户邮箱 | 无 |
| `-CF_ZONE_NAME` | `CF_ZONE_NAME` | Cloudflare域名区域名称 | 无 |
| `-CF_DNS_NAME` | `CF_DNS_NAME` | DNS记录名称（子域名） | 无 |

### 可选参数

| 参数 | 环境变量 | 描述 | 默认值 |
|------|----------|------|--------|
| `-CF_DNS_TTL` | `CF_DNS_TTL` | DNS记录的TTL值（秒） | 0（自动） |
| `-IP_CHECK_DURATION` | `IP_CHECK_DURATION` | IP检查间隔时间 | 1分钟 |
| `-CF_TIMEOUT` | `CF_TIMEOUT` | API调用超时时间 | 30秒 |

## 工作原理

1. **IP地址获取**：程序通过访问 `https://api-ipv4.ip.sb/ip` 获取当前公网IPv4地址
2. **Cloudflare API认证**：使用提供的API密钥和邮箱进行认证
3. **DNS记录查询**：查找指定域名和子域名的DNS A记录
4. **记录管理**：
   - 如果记录不存在，创建新的A记录
   - 如果记录存在且IP地址不同，更新记录
   - 如果记录存在且IP地址相同，跳过更新
5. **定时循环**：按照配置的间隔时间重复执行上述过程

## 代码结构说明

### 主要函数

#### `main()`
程序入口点，负责：
- 解析命令行参数和环境变量
- 设置默认值
- 启动定时循环执行DDNS更新

#### `DDNS(ctx, cfAPIKey, cfAPIEmail, zoneName, dnsName, ttl)`
核心DDNS逻辑，负责：
- 获取当前IP地址
- 初始化Cloudflare API客户端
- 查询DNS记录
- 创建或更新DNS记录

#### `getIP(ctx)`
通过 `https://api-ipv4.ip.sb/ip` 获取当前公网IPv4地址：

- 校验 HTTP 状态码必须为 200
- 使用 `io.LimitReader` 限制响应体最大 4KB，避免异常响应占用过多内存
- 响应为纯文本 IP（如 `180.172.40.161`），经 `strings.TrimSpace` 后使用 `net.ParseIP` 校验必须为合法 IPv4

## 依赖项

- **cloudflare-go** (v0.64.0): Cloudflare官方Go SDK
- **golang.org/x/net**: Go网络库，用于上下文支持

## 开发指南

### 代码格式化

```bash
make fmt
```

### 清理构建文件

```bash
make clean
```

### 运行测试

```bash
go run main.go [参数]
```

## 部署建议

### 1. 获取Cloudflare API密钥

1. 登录Cloudflare控制台
2. 进入"我的个人资料" > "API令牌"
3. 创建新的API令牌，选择"编辑区域DNS"权限

### 2. 系统服务部署（Linux）

创建systemd服务文件 `/etc/systemd/system/cloudflare-ddns.service`：

```ini
[Unit]
Description=Cloudflare DDNS Service
After=network.target

[Service]
Type=simple
User=ddns
Environment=CF_API_EMAIL=your-email@example.com
Environment=CF_API_KEY=your-api-key
Environment=CF_ZONE_NAME=example.com
Environment=CF_DNS_NAME=subdomain
Environment=CF_DNS_TTL=120
Environment=IP_CHECK_DURATION=5m
Environment=CF_TIMEOUT=30s
WorkingDirectory=/opt/cloudflare-ddns
ExecStart=/opt/cloudflare-ddns/cloudflare
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 3. 容器化部署

使用提供的Dockerfile构建镜像，可通过Docker Compose或Kubernetes进行编排。

## 故障排除

### 常见问题

1. **认证失败**
   - 检查API密钥和邮箱是否正确
   - 确认API令牌具有DNS编辑权限

2. **DNS记录未更新**
   - 检查域名区域名称是否正确
   - 确认子域名配置正确
   - 查看程序日志获取详细错误信息

3. **IP获取失败**
   - 检查网络连接
   - 确认api-ipv4.ip.sb服务可用性

### 日志说明

程序会输出以下类型的日志信息：
- `get IPInfo success`: IP地址获取成功
- `DNS Record Created`: 创建新的DNS记录
- `DNS Record Not Change`: IP地址未变化，跳过更新
- `update DNS record success`: 成功更新DNS记录
- `[FAIL]`: 操作失败，显示错误信息

## 安全注意事项

1. **API密钥保护**：不要将API密钥提交到版本控制系统
2. **最小权限原则**：使用仅具有DNS编辑权限的API令牌
3. **网络隔离**：在生产环境中限制网络访问权限
4. **日志管理**：定期清理日志文件，避免敏感信息泄露

## 许可证

本项目采用MIT许可证。详见LICENSE文件（如有）。

## 贡献指南

欢迎提交Issue和Pull Request来改进这个项目。

## 联系方式

- 项目维护者：Sinute <sinute@outlook.com>
- 问题反馈：通过GitHub Issues提交