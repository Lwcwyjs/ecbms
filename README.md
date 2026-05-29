# 边缘计算盒子管理系统 (ECBMS)

基于Ubuntu系统的边缘计算盒子Web管理系统，提供网络配置、资源监控功能，方便用户对边缘计算盒子进行远程管理和维护。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: SQLite3
- **认证**: JWT
- **系统监控**: gopsutil

### 前端
- **框架**: Vue 3 (Composition API)
- **UI组件库**: Element Plus
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **图表**: ECharts

## 项目结构

```
ecbms/
├── backend/                 # 后端代码
│   ├── cmd/
│   │   └── server/         # 服务入口
│   │       └── main.go
│   ├── internal/
│   │   ├── handler/        # API处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   └── service/        # 业务逻辑
│   ├── pkg/
│   │   └── utils/          # 工具函数
│   ├── config/             # 配置
│   ├── go.mod
│   └── go.sum
├── frontend/               # 前端代码
│   ├── src/
│   │   ├── api/            # API接口
│   │   ├── components/     # 组件
│   │   ├── router/         # 路由
│   │   ├── store/          # 状态管理
│   │   ├── utils/          # 工具函数
│   │   ├── views/          # 页面视图
│   │   └── main.js
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
└── docs/                   # 文档
```

## 功能特性

### 1. 系统概览 (Dashboard)
- 主机名、CPU/内存/磁盘使用率实时展示
- 系统信息展示（操作系统、CPU型号、内存等）
- 快捷操作（刷新数据、重启系统、关闭系统）

### 2. 网络配置 (Network)
- 网络接口列表查看
- DHCP/静态IP配置
- Ping网络测试
- DNS配置查看

### 3. 资源监控 (Monitor)
- CPU使用率趋势图
- 内存使用率趋势图
- 磁盘使用率仪表盘
- 网络流量趋势图
- 实时数据刷新（3秒间隔）

### 4. 系统日志 (Logs)
- 系统操作日志记录
- 日志级别分类（信息/警告/错误）
- 日志类型分类（认证/网络/系统/用户）

### 5. 用户管理 (Users)
- 用户列表查看
- 添加/删除用户
- 角色权限控制（管理员/普通用户）
- 默认管理员账号: `admin` / `admin123`

### 6. 个人设置 (Settings)
- 修改密码
- 账户信息查看

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- Ubuntu 20.04+ (推荐)

### 后端启动

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

### 生产部署

#### 后端编译
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o ecbms-server cmd/server/main.go
```

#### 前端构建
```bash
cd frontend
npm run build
```

## API接口

### 认证接口
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息
- `POST /api/auth/change-password` - 修改密码

### 系统接口
- `GET /api/system/info` - 获取系统信息
- `GET /api/system/stats` - 获取系统实时状态
- `GET /api/system/logs` - 获取系统日志
- `POST /api/system/reboot` - 重启系统
- `POST /api/system/shutdown` - 关闭系统

### 网络接口
- `GET /api/network/interfaces` - 获取网络接口列表
- `GET /api/network/configs` - 获取网络配置历史
- `POST /api/network/configure` - 保存网络配置
- `POST /api/network/apply` - 应用网络配置
- `POST /api/network/ping` - Ping测试
- `GET /api/network/dns` - 获取DNS配置

### 用户管理接口 (需管理员权限)
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `DELETE /api/users/:id` - 删除用户

## 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | `8080` | 服务端口 |
| `GIN_MODE` | `debug` | Gin运行模式 (debug/release) |
| `DB_PATH` | `./data/ecbms.db` | SQLite数据库路径 |
| `JWT_SECRET` | `ecbms-secret-key-2024` | JWT密钥 |
| `JWT_EXPIRE_HOURS` | `24` | JWT过期时间(小时) |

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| `admin` | `admin123` | 管理员 |

## 注意事项

1. **生产环境**请务必修改默认密码和JWT密钥
2. 网络配置应用功能需要Linux系统权限
3. 系统重启/关机功能需要root权限
4. 建议使用systemd管理后台服务

## License

MIT
"# ecbms" 
