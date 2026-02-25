# 高级算命 · Fortune Teller App

前后端分离的「高级算命」占卜应用：Go 后端提供随机运势 API，React 前端提供占卜界面。

## 技术栈

- **后端**: Go 1.22，标准库 HTTP + CORS，端口 `8081`
- **前端**: React 18 + Vite 5 + TypeScript，端口 `5174`
- **运行**: 支持本地直接运行或 Docker Compose 一键启动

## 本地直接运行

### 后端

```bash
cd backend
go mod tidy
go run ./cmd/server
```

服务监听 `http://localhost:8081`。

### 前端

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 `http://localhost:5174`。开发模式下请求会通过 Vite 代理到后端 `/api`。

## Docker 一键启动

在项目根目录执行：

```bash
docker compose up --build
```

- 后端 API: http://localhost:8081  
- 前端页面: http://localhost:5174  

若端口 8081 或 5174 已被占用，可修改 `docker-compose.yml` 中 `ports` 映射（如 `9081:8081`、`9075:80`），并相应调整前端构建参数 `VITE_API_BASE`（例如 `http://localhost:9081`）。

## 脚本启动（可选）

```bash
chmod +x start.sh
./start.sh
```

`start.sh` 会以后台方式启动后端与前端（需已安装 Go 与 Node），具体见脚本内容。

## API 说明

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET | `/api/categories` | 运势类别列表 |
| GET | `/api/fortune` | 随机抽一条运势 |
| GET | `/api/fortune?category=love` | 按类别（love/career/health/wealth/general）抽签 |

## 后续

可在此仓库继续开发、提交并推送到 GitHub；如需修改端口或技术栈，可参考上述说明调整。
