# SaaS Monitor 部署指南

## 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [Docker Compose 部署](#docker-compose-部署)
- [手动部署](#手动部署)
- [环境变量](#环境变量)
- [健康检查](#健康检查)
- [故障排除](#故障排除)
- [生产环境建议](#生产环境建议)

## 前置要求

### 软件要求
- Docker 20.10+
- Docker Compose 1.29+ 或 Docker Compose V2
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

### 网络要求
- 端口 3000: 前端访问
- 端口 8080: 后端 API
- 端口 5432: PostgreSQL (可选，如使用外部数据库)
- 端口 6379: Redis (可选，如使用外部 Redis)

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd sass-monitor
```

### 2. 配置环境

复制环境变量模板:
```bash
cp .env.example .env
```

编辑 `.env` 文件,修改数据库密码等敏感信息:
```bash
vim .env
```

复制配置文件模板:
```bash
cp backend/configs/config.yaml.example backend/configs/config.yaml
```

编辑 `backend/configs/config.yaml`,配置数据库连接信息:
```bash
vim backend/configs/config.yaml
```

### 3. 一键部署

```bash
./deploy.sh deploy
```

### 4. 访问系统

- 前端: http://localhost:3000
- 后端 API: http://localhost:8080
- 健康检查: http://localhost:8080/health

**默认登录账号:**
- 用户名: `admin`
- 密码: `admin123`

## 配置说明

### 后端配置 (backend/configs/config.yaml)

```yaml
server:
  port: 8080
  mode: release
  jwt_secret: "your-jwt-secret-key"
  jwt_expire_hours: 24

databases:
  # SaaS 监控数据库
  saas_monitor:
    type: postgres
    host: "postgres"  # Docker 内部使用服务名
    port: 5432
    user: "postgres"
    password: "${POSTGRES_PASSWORD}"
    database: "saas_monitor"

  # Light Admin 业务数据库 (只读监控)
  light_admin:
    type: postgres
    host: "外部数据库地址"
    port: 35432
    user: "postgres"
    password: "数据库密码"
    database: "light_admin"
    readonly: true

clickhouse:
  - name: "traces"
    host: "外部 ClickHouse 地址"
    port: 39000
    user: "default"
    password: "ClickHouse 密码"
    database: "light_traces"

redis:
  host: "redis"  # Docker 内部使用服务名
  port: "6379"
  password: "${REDIS_PASSWORD}"
  database: 0
```

### 环境变量配置 (.env)

```bash
# PostgreSQL 配置
POSTGRES_PASSWORD=your_secure_password
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# Redis 配置
REDIS_PASSWORD=your_redis_password
REDIS_HOST=redis
REDIS_PORT=6379

# ClickHouse 配置 (外部数据库)
CLICKHOUSE_HOST=192.168.2.81
CLICKHOUSE_PORT=39000
CLICKHOUSE_PASSWORD=your_clickhouse_password
```

## Docker Compose 部署

### 服务架构

```
┌─────────────────────────────────────────────┐
│           SaaS Monitor System              │
├─────────────────────────────────────────────┤
│  Frontend (Nginx + React)                   │
│  ├─ Port: 3000                              │
│  └─ 静态文件服务 + API 代理                 │
├─────────────────────────────────────────────┤
│  Backend (Go + Gin)                         │
│  ├─ Port: 8080                              │
│  └─ REST API + 监控采集                     │
├─────────────────────────────────────────────┤
│  PostgreSQL 15                              │
│  ├─ Port: 5432                              │
│  └─ 监控数据存储                            │
├─────────────────────────────────────────────┤
│  Redis 7                                    │
│  ├─ Port: 6379                              │
│  └─ 会话缓存                                │
└─────────────────────────────────────────────┘
```

### 部署命令

#### 完整部署
```bash
./deploy.sh deploy
```

#### 仅构建镜像
```bash
./deploy.sh build
```

#### 启动服务
```bash
./deploy.sh start
```

#### 停止服务
```bash
./deploy.sh stop
```

#### 重启服务
```bash
./deploy.sh restart
```

#### 查看状态
```bash
./deploy.sh status
```

#### 查看日志
```bash
./deploy.sh logs

# 或查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
```

#### 清理数据
```bash
./deploy.sh clean
```

### 手动 Docker Compose 命令

```bash
# 构建并启动
docker-compose up -d --build

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v

# 查看日志
docker-compose logs -f

# 进入容器
docker-compose exec backend sh
docker-compose exec frontend sh
```

## 手动部署

### 后端部署

#### 1. 构建后端
```bash
cd backend
go mod tidy
go build -o sass-monitor ./cmd/api
```

#### 2. 配置数据库
创建数据库:
```sql
CREATE DATABASE saas_monitor;
```

#### 3. 运行后端
```bash
./sass-monitor
```

### 前端部署

#### 1. 安装依赖
```bash
cd frontend
npm install
```

#### 2. 构建前端
```bash
npm run build
```

#### 3. 使用 Nginx 部署
```bash
# 复制构建产物到 Nginx 目录
cp -r build/* /usr/share/nginx/html/

# 配置 Nginx 反向代理
# 参考 frontend/nginx.conf
```

## 环境变量

### 必需的环境变量

| 变量名 | 说明 | 示例 |
|--------|------|------|
| POSTGRES_PASSWORD | PostgreSQL 密码 | `secure_password` |
| REDIS_PASSWORD | Redis 密码 | `redis_password` |

### 可选的环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| POSTGRES_HOST | PostgreSQL 主机 | `postgres` |
| POSTGRES_PORT | PostgreSQL 端口 | `5432` |
| REDIS_HOST | Redis 主机 | `redis` |
| REDIS_PORT | Redis 端口 | `6379` |
| BACKEND_PORT | 后端端口 | `8080` |
| FRONTEND_PORT | 前端端口 | `3000` |
| GIN_MODE | Gin 运行模式 | `release` |

## 健康检查

### 后端健康检查
```bash
curl http://localhost:8080/health
```

预期响应:
```json
{
  "status": "healthy",
  "timestamp": "2025-11-24T10:00:00Z",
  "version": "1.0.0"
}
```

### 前端健康检查
```bash
curl http://localhost:3000/health
```

### Docker 健康检查
```bash
docker-compose ps
```

所有服务的 `STATUS` 应显示 `Up (healthy)`。

## 故障排除

### 问题 1: 容器启动失败

**症状:** 容器状态为 `Exited` 或 `Restarting`

**解决方案:**
```bash
# 查看日志
docker-compose logs backend

# 检查配置文件
cat backend/configs/config.yaml

# 确保数据库连接正确
docker-compose exec postgres psql -U postgres -d saas_monitor
```

### 问题 2: 前端无法连接后端

**症状:** 前端页面显示网络错误

**解决方案:**
```bash
# 检查后端是否运行
curl http://localhost:8080/health

# 检查 Nginx 配置
docker-compose exec frontend cat /etc/nginx/conf.d/default.conf

# 检查网络连接
docker-compose exec frontend ping backend
```

### 问题 3: 数据库连接失败

**症状:** 后端日志显示数据库连接错误

**解决方案:**
```bash
# 检查 PostgreSQL 状态
docker-compose exec postgres pg_isready

# 检查数据库密码
docker-compose exec postgres psql -U postgres -c "SELECT 1"

# 重启数据库
docker-compose restart postgres
```

### 问题 4: 端口冲突

**症状:** 容器启动时提示端口被占用

**解决方案:**
```bash
# 检查端口占用
lsof -i :3000
lsof -i :8080

# 修改 docker-compose.yml 中的端口映射
# 或停止占用端口的进程
```

## 生产环境建议

### 1. 安全配置

- ✅ 修改所有默认密码
- ✅ 使用强密码和密钥
- ✅ 启用 HTTPS (使用 Nginx + Let's Encrypt)
- ✅ 配置防火墙规则
- ✅ 定期更新依赖和镜像

### 2. 性能优化

- ✅ 调整数据库连接池大小
- ✅ 配置 Redis 最大内存
- ✅ 启用 Nginx Gzip 压缩
- ✅ 配置静态资源缓存
- ✅ 使用 CDN 加速静态资源

### 3. 监控和日志

- ✅ 配置日志轮转
- ✅ 集成监控告警系统
- ✅ 配置日志聚合 (如 ELK Stack)
- ✅ 设置资源限制 (CPU、内存)

### 4. 备份策略

```bash
# 备份 PostgreSQL
docker-compose exec postgres pg_dump -U postgres saas_monitor > backup.sql

# 备份 Redis
docker-compose exec redis redis-cli --rdb /data/dump.rdb

# 备份配置文件
tar -czf config-backup.tar.gz backend/configs/
```

### 5. 资源限制

编辑 `docker-compose.yml` 添加资源限制:

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 6. 使用外部数据库

对于生产环境,建议使用独立的数据库服务:

1. 注释掉 `docker-compose.yml` 中的 `postgres` 和 `redis` 服务
2. 在 `backend/configs/config.yaml` 中配置外部数据库地址
3. 确保网络连通性和防火墙规则

## 更新部署

### 1. 拉取最新代码
```bash
git pull origin main
```

### 2. 重新构建镜像
```bash
./deploy.sh build
```

### 3. 重启服务
```bash
./deploy.sh restart
```

### 4. 验证更新
```bash
./deploy.sh status
curl http://localhost:8080/health
```

## 卸载

```bash
# 停止并删除所有服务和数据
./deploy.sh clean

# 删除镜像
docker rmi sass-monitor-backend sass-monitor-frontend

# 删除网络
docker network rm sass-monitor_sass-monitor-network
```

## 技术支持

如有问题或建议,请提交 Issue 或联系开发团队。

---

**文档版本:** 1.0.0
**最后更新:** 2025-11-24
