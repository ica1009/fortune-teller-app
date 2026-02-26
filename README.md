# 高级算命 · Fortune Teller App

前后端分离的「高级算命」占卜应用：Go 后端提供随机运势 API 与可选 AI 占卜（OpenAI 兼容 / DeepSeek），React 前端提供占卜界面。

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

- 后端 API: http://localhost:10081  
- 前端页面: http://localhost:10075  

（宿主机端口 10081/10075；若仍冲突，可改 `docker-compose.yml` 中 `ports` 与 `VITE_API_BASE`。）

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
| GET | `/api/fortune?category=love` | 按类别抽签 |
| POST | `/api/fortune/ai` | AI 占卜，请求体可选 `{"category":"love"}`；需配置 `OPENAI_API_KEY`，否则 503 |

## AI 占卜（可选）

- 与 [rain-flower-calendar](https://github.com/ica1009/rain-flower-calendar) 共用同一套 Key：通过环境变量 **`OPENAI_API_KEY`**（及可选 **`OPENAI_BASE_URL`**，默认 `https://api.deepseek.com`）配置。
- 本地运行后端时：`export OPENAI_API_KEY=your-key` 后再 `go run ./cmd/server`。
- Docker：复制 `.env.example` 为 `.env` 并填写 `OPENAI_API_KEY`，然后 `docker compose up --build`；或直接 `OPENAI_API_KEY=your-key docker compose up --build`。
- 未配置时前端「AI 占卜」会提示暂不可用，普通抽签不受影响。

## 后续

可在此仓库继续开发、提交并推送到 GitHub；如需修改端口或技术栈，可参考上述说明调整。
