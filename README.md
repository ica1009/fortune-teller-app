# 高级算命 · Fortune Teller App

前后端分离的「高级算命」占卜应用：Go 后端提供随机运势 API 与可选 AI 占卜（OpenAI 兼容 / DeepSeek），React 前端提供占卜界面；支持用户注册/登录（用户名+密码，PostgreSQL，JWT）。

## 技术栈

- **后端**: Go，标准库 HTTP + CORS，PostgreSQL 用户存储，bcrypt + JWT
- **前端**: React 18 + Vite 5 + TypeScript
- **运行**: Docker Compose 一键启动（含 DB），或本地起 DB + 后端 + 前端

## Docker 一键启动（推荐）

需设置 `JWT_SECRET`（生产务必改为强随机）：

```bash
export JWT_SECRET=your-secret-key
docker compose up -d
```

- 前端: http://localhost:10075  
- 后端 API: http://localhost:10081  
- 数据库: PostgreSQL 宿主机端口 5434  

未登录会先进入登录/注册页，登录后可抽签占卜。

## 本地开发（需 PostgreSQL）

若只起数据库容器、本机跑后端与前端：

```bash
docker compose up -d db
# 等几秒后
export DATABASE_URL="postgres://app:appsecret@localhost:5434/fortune_teller?sslmode=disable"
export JWT_SECRET=dev-secret
cd backend && go run ./cmd/server
```

另开终端：

```bash
cd frontend && npm install && npm run dev
```

浏览器打开 http://localhost:5174。开发模式下 Vite 代理 `/api` 到后端。

## API 说明

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查 |
| GET | `/api/categories` | 运势类别列表 |
| GET | `/api/fortune` | 随机抽一条运势 |
| GET | `/api/fortune?category=love` | 按类别抽签 |
| POST | `/api/auth/register` | 注册（body: username, password） |
| POST | `/api/auth/login` | 登录（body: username, password；返回 token） |
| POST | `/api/fortune/ai` | AI 占卜，请求体可选 `{"category":"love"}`；需配置 `OPENAI_API_KEY`，否则 503 |

## AI 占卜（可选）

- 与 [rain-flower-calendar](https://github.com/ica1009/rain-flower-calendar) 共用同一套 Key：通过环境变量 **`OPENAI_API_KEY`**（及可选 **`OPENAI_BASE_URL`**，默认 `https://api.deepseek.com`）配置。
- 本地运行后端时：`export OPENAI_API_KEY=your-key` 后再 `go run ./cmd/server`。
- Docker：复制 `.env.example` 为 `.env` 并填写 `OPENAI_API_KEY`，然后 `docker compose up --build`；或直接 `OPENAI_API_KEY=your-key docker compose up --build`。
- 未配置时前端「AI 占卜」会提示暂不可用，普通抽签不受影响。

## 后续

可在此仓库继续开发、提交并推送到 GitHub；如需修改端口或技术栈，可参考上述说明调整。
