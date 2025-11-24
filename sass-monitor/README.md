# Sass后台监控系统

企业级多数据库监控后台系统，用于监控light_admin系统的PostgreSQL和ClickHouse资源使用情况。

## 系统架构

### 技术栈
- **后端**: Go + Gin + GORM + JWT认证
- **前端**: React + TypeScript + Ant Design
- **数据库连接**: PostgreSQL + ClickHouse + Redis
- **UI设计**: 浅蓝色主题，支持白天/黑夜模式切换

### 监控目标
- **PostgreSQL**: light_admin数据库（只读监控）
- **ClickHouse**: traces、logs、llm、rum四个数据库
- **Redis**: 缓存系统监控

### 核心功能
1. **管理员认证**: JWT令牌认证，角色权限控制
2. **组织监控**: 按organization_id统计资源使用情况
3. **多数据库监控**: PostgreSQL、ClickHouse、Redis状态监控
4. **实时仪表板**: 资源使用情况可视化展示
5. **告警系统**: 自定义告警规则和通知
6. **主题切换**: 支持白天/黑夜模式

## 项目结构

```
sass-monitor/
├── backend/                    # Go后端服务
│   ├── cmd/api/               # 应用入口
│   ├── internal/              # 内部模块
│   │   ├── auth/             # 认证逻辑
│   │   ├── database/         # 数据库管理
│   │   ├── handlers/         # HTTP处理器
│   │   ├── middleware/       # 中间件
│   │   ├── models/           # 数据模型
│   │   └── services/         # 业务服务
│   ├── pkg/                  # 公共包
│   │   ├── config/           # 配置管理
│   │   └── utils/            # 工具函数
│   └── configs/              # 配置文件
├── frontend/                  # React前端
│   ├── src/
│   │   ├── components/       # 组件
│   │   ├── pages/           # 页面
│   │   ├── services/        # API服务
│   │   ├── utils/           # 工具函数
│   │   └── styles/          # 样式文件
│   └── public/              # 静态资源
└── docs/                    # 文档
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- PostgreSQL 12+
- ClickHouse 22+
- Redis 6+

### 后端启动

1. 克隆项目
```bash
git clone <repository-url>
cd sass-monitor
```

2. 安装Go依赖
```bash
cd backend
go mod tidy
```

3. 配置数据库连接
编辑 `configs/config.yaml` 文件，配置数据库连接信息。

4. 创建监控数据库
```sql
CREATE DATABASE saas_monitor;
```

5. 启动服务
```bash
go run cmd/api/main.go
```

服务将在 http://localhost:8080 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
# 或
yarn install
```

3. 启动开发服务器
```bash
npm start
# 或
yarn start
```

前端将在 http://localhost:3000 启动

## API文档

### 认证接口

#### 登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

#### 获取用户信息
```http
GET /api/v1/profile
Authorization: Bearer <token>
```

### 监控接口

#### 获取仪表板概览
```http
GET /api/v1/dashboard/overview
Authorization: Bearer <token>
```

#### 获取组织列表
```http
GET /api/v1/dashboard/organizations?page=1&page_size=20&search=keyword
Authorization: Bearer <token>
```

#### 获取组织指标
```http
GET /api/v1/dashboard/organizations/{org_id}/metrics
Authorization: Bearer <token>
```

#### 获取数据库状态
```http
GET /api/v1/dashboard/database-status
Authorization: Bearer <token>
```

#### 获取监控指标
```http
GET /api/v1/monitoring/metrics?database_type=postgresql&metric_type=storage
Authorization: Bearer <token>
```

## 配置说明

### 数据库配置
```yaml
databases:
  saas_monitor:
    type: postgres
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "saas_monitor"

  light_admin:
    type: postgres
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "light_admin"
    readonly: true

clickhouse:
  - name: "traces"
    host: "localhost"
    port: 9000
    user: "default"
    password: ""
    database: "light_traces"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  database: 0
```

### 监控配置
```yaml
monitoring:
  collect_interval: 5        # 数据采集间隔（分钟）
  retention_days: 30         # 数据保留天数
  alerts:
    enabled: true
    cpu_threshold: 80
    memory_threshold: 85
    disk_threshold: 90
```

## 监控指标

### PostgreSQL指标
- 连接数
- 数据库大小
- 表大小
- 查询性能
- 活跃会话数

### ClickHouse指标
- 数据库大小
- 表数量
- 行数统计
- 查询响应时间
- 压缩率

### Redis指标
- 内存使用率
- 连接数
- 键值对数量
- 命中率
- 响应时间

### 组织维度指标
- 用户数量
- 订阅数量
- 存储使用量
- API调用次数
- 数据查询量

## 告警系统

### 告警类型
- **系统告警**: CPU、内存、磁盘使用率
- **数据库告警**: 连接数、响应时间、错误率
- **业务告警**: 用户数量、订阅状态

### 告警级别
- **info**: 信息提示
- **warning**: 警告
- **critical**: 严重告警

### 通知方式
- 邮件通知
- 短信通知
- Webhook通知

## 部署

### Docker部署

1. 构建镜像
```bash
# 后端镜像
cd backend
docker build -t sass-monitor-backend .

# 前端镜像
cd ../frontend
docker build -t sass-monitor-frontend .
```

2. Docker Compose部署
```bash
docker-compose up -d
```

### 生产环境配置

1. 环境变量配置
2. 数据库连接池优化
3. 日志配置
4. 监控和告警配置
5. 安全配置

## 开发指南

### 添加新的监控指标

1. 在 `internal/models/admin.go` 中定义指标模型
2. 在 `internal/services/` 中实现数据采集逻辑
3. 在 `internal/handlers/` 中添加API接口
4. 在前端添加展示组件

### 自定义告警规则

1. 在 `models.AlertRule` 中定义规则结构
2. 实现告警检测逻辑
3. 配置通知方式
4. 添加告警历史记录

## 贡献

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 发起Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交Issue或联系开发团队。