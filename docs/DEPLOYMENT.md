# 部署指南

本指南提供多种部署 Go Standards MCP Server 的方式。

## 部署方式

- [本地开发部署](#本地开发部署)
- [Docker 部署](#docker-部署)
- [Docker Compose 部署](#docker-compose-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [systemd 服务部署](#systemd-服务部署)

---

## 本地开发部署

### 前置要求

- Go 1.21+
- golangci-lint（可选）

### 步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/MOONL0323/go-standards-mcp-server.git
   cd go-standards-mcp-server
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **安装 golangci-lint**
   ```bash
   # macOS
   brew install golangci-lint

   # Linux
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

   # Windows (使用 PowerShell)
   iwr -useb https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | iex
   ```

4. **构建**
   ```bash
   make build
   ```

5. **运行**
   ```bash
   # stdio 模式
   ./bin/mcp-server

   # HTTP 模式
   ./bin/mcp-server --mode http --port 8080
   ```

---

## Docker 部署

### 构建镜像

```bash
docker build -t go-standards-mcp-server:latest .
```

### 运行容器

#### stdio 模式（不适用于 Docker）

stdio 模式用于本地集成，不适合 Docker 部署。

#### HTTP 模式

```bash
docker run -d \
  --name mcp-server \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs:ro \
  -v $(pwd)/reports:/app/reports \
  go-standards-mcp-server:latest \
  --mode http --port 8080
```

### 环境变量配置

```bash
docker run -d \
  --name mcp-server \
  -p 8080:8080 \
  -e MCP_SERVER_MODE=http \
  -e MCP_SERVER_PORT=8080 \
  -e MCP_LOG_LEVEL=info \
  go-standards-mcp-server:latest
```

---

## Docker Compose 部署

### 基础部署

创建 `docker-compose.yml`（项目已包含）：

```yaml
version: '3.8'

services:
  mcp-server:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_MODE=http
      - SERVER_PORT=8080
      - LOG_LEVEL=info
    volumes:
      - ./configs:/app/configs:ro
      - ./reports:/app/reports
    restart: unless-stopped
```

### 启动服务

```bash
# 启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止
docker-compose down
```

### 完整部署（包含 Redis 和 PostgreSQL）

使用项目提供的完整 `docker-compose.yml`：

```bash
# 启动所有服务
docker-compose up -d

# 查看所有服务状态
docker-compose ps

# 查看特定服务日志
docker-compose logs -f mcp-server

# 停止所有服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

---

## Kubernetes 部署

### 创建部署配置

创建 `k8s/deployment.yaml`：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mcp-server
  labels:
    app: mcp-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: mcp-server
  template:
    metadata:
      labels:
        app: mcp-server
    spec:
      containers:
      - name: mcp-server
        image: go-standards-mcp-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_MODE
          value: "http"
        - name: SERVER_PORT
          value: "8080"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 创建服务配置

创建 `k8s/service.yaml`：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mcp-server
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
  selector:
    app: mcp-server
```

### 创建 ConfigMap

创建 `k8s/configmap.yaml`：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mcp-server-config
data:
  default.yaml: |
    server:
      mode: http
      port: 8080
    log:
      level: info
    # ... 其他配置
```

### 部署到 Kubernetes

```bash
# 创建命名空间
kubectl create namespace mcp-server

# 应用配置
kubectl apply -f k8s/configmap.yaml -n mcp-server
kubectl apply -f k8s/deployment.yaml -n mcp-server
kubectl apply -f k8s/service.yaml -n mcp-server

# 查看状态
kubectl get pods -n mcp-server
kubectl get svc -n mcp-server

# 查看日志
kubectl logs -f deployment/mcp-server -n mcp-server

# 扩缩容
kubectl scale deployment/mcp-server --replicas=5 -n mcp-server
```

---

## systemd 服务部署

### 创建服务文件

创建 `/etc/systemd/system/mcp-server.service`：

```ini
[Unit]
Description=Go Standards MCP Server
After=network.target

[Service]
Type=simple
User=mcp
Group=mcp
WorkingDirectory=/opt/mcp-server
ExecStart=/opt/mcp-server/bin/mcp-server --config /opt/mcp-server/configs/default.yaml
Restart=always
RestartSec=10

# 环境变量
Environment="MCP_SERVER_MODE=http"
Environment="MCP_SERVER_PORT=8080"
Environment="MCP_LOG_LEVEL=info"

# 安全加固
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/mcp-server/reports /opt/mcp-server/tmp

[Install]
WantedBy=multi-user.target
```

### 安装和启动

```bash
# 创建用户
sudo useradd -r -s /bin/false mcp

# 创建目录
sudo mkdir -p /opt/mcp-server/{bin,configs,reports,tmp}

# 复制文件
sudo cp bin/mcp-server /opt/mcp-server/bin/
sudo cp -r configs /opt/mcp-server/

# 设置权限
sudo chown -R mcp:mcp /opt/mcp-server

# 重载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start mcp-server

# 设置开机自启
sudo systemctl enable mcp-server

# 查看状态
sudo systemctl status mcp-server

# 查看日志
sudo journalctl -u mcp-server -f
```

---

## 生产环境配置建议

### 1. 资源限制

```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

### 2. 日志配置

```yaml
log:
  level: warn
  output: /var/log/mcp-server/app.log
  format: json
```

### 3. 数据库配置

使用 PostgreSQL 而不是 SQLite：

```yaml
storage:
  type: postgres
  postgres:
    host: postgres.example.com
    port: 5432
    user: mcp_user
    password: ${POSTGRES_PASSWORD}  # 使用环境变量
    database: mcp_server
    sslmode: require
```

### 4. 缓存配置

启用 Redis 缓存：

```yaml
cache:
  enabled: true
  type: redis
  redis:
    addr: redis.example.com:6379
    password: ${REDIS_PASSWORD}
    db: 0
```

### 5. 监控配置

- 集成 Prometheus 指标导出
- 配置健康检查端点
- 设置告警规则

### 6. 安全配置

- 使用 HTTPS/TLS
- 启用认证和授权
- 限制并发连接数
- 配置防火墙规则

---

## 性能调优

### 1. 并发限制

```yaml
analyzer:
  concurrent_limit: 20  # 根据 CPU 核心数调整
```

### 2. 超时设置

```yaml
analyzer:
  timeout: 600  # 10 分钟，根据项目规模调整
```

### 3. 连接池

```yaml
storage:
  postgres:
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 5m
```

---

## 监控和维护

### 健康检查

```bash
# HTTP 模式
curl http://localhost:8080/health

# 使用 MCP 工具
# 发送 health_check 工具调用
```

### 日志轮转

配置 logrotate（`/etc/logrotate.d/mcp-server`）：

```
/var/log/mcp-server/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 mcp mcp
    sharedscripts
    postrotate
        systemctl reload mcp-server
    endscript
}
```

### 备份

```bash
# 备份配置
tar -czf mcp-server-config-$(date +%Y%m%d).tar.gz configs/

# 备份数据库（PostgreSQL）
pg_dump -h localhost -U mcp_user mcp_server > mcp-server-$(date +%Y%m%d).sql
```

---

## 故障排除

### 服务无法启动

1. 检查日志
2. 验证配置文件
3. 确认端口未被占用
4. 检查文件权限

### 性能问题

1. 增加并发限制
2. 优化数据库查询
3. 启用缓存
4. 检查资源使用情况

### 连接问题

1. 检查网络配置
2. 验证防火墙规则
3. 确认服务正在运行
4. 检查负载均衡器配置

---

## 更新和回滚

### 更新服务

```bash
# 构建新版本
make build

# 停止服务
sudo systemctl stop mcp-server

# 备份当前版本
sudo cp /opt/mcp-server/bin/mcp-server /opt/mcp-server/bin/mcp-server.bak

# 部署新版本
sudo cp bin/mcp-server /opt/mcp-server/bin/

# 启动服务
sudo systemctl start mcp-server

# 检查状态
sudo systemctl status mcp-server
```

### 回滚

```bash
# 停止服务
sudo systemctl stop mcp-server

# 恢复备份
sudo cp /opt/mcp-server/bin/mcp-server.bak /opt/mcp-server/bin/mcp-server

# 启动服务
sudo systemctl start mcp-server
```

---

## 相关资源

- [配置参考](../configs/default.yaml)
- [API 文档](API.md)
- [故障排除指南](TROUBLESHOOTING.md)
- [性能调优指南](PERFORMANCE.md)

---

**注意**: 部署到生产环境前，请务必进行充分的测试！
