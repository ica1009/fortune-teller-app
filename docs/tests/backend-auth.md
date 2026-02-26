# 后端：注册/登录接口测试

## 环境

- 后端地址：`http://localhost:10081`（Docker）或 `http://localhost:8081`（本地）
- 需先启动 PostgreSQL 与后端服务。

---

## 用例 1：注册新用户

| 项目 | 说明 |
|------|------|
| 接口 | `POST /api/auth/register` |
| 请求体 | `{"username":"testuser","password":"pass1234"}` |
| 预期 | 201，`{"message":"registered"}` |
| 通过标准 | 状态码 201。 |

---

## 用例 2：注册重复用户名

| 项目 | 说明 |
|------|------|
| 接口 | `POST /api/auth/register` |
| 请求体 | 与已存在用户相同的 username |
| 预期 | 400，`{"error":"username already exists"}` |
| 通过标准 | 状态码 400。 |

---

## 用例 3：登录成功

| 项目 | 说明 |
|------|------|
| 接口 | `POST /api/auth/login` |
| 请求体 | 已注册用户的 username 与正确 password |
| 预期 | 200，`{"token":"<JWT>"}` |
| 通过标准 | 状态码 200，token 为合法 JWT。 |

---

## 用例 4：登录失败

| 项目 | 说明 |
|------|------|
| 接口 | `POST /api/auth/login` |
| 请求体 | 错误密码或不存在用户名 |
| 预期 | 401，`{"error":"invalid username or password"}` |
| 通过标准 | 状态码 401。 |

---

## 快速验证

```bash
curl -s -X POST http://localhost:10081/api/auth/register -H "Content-Type: application/json" -d '{"username":"u1","password":"pass1234"}'
curl -s -X POST http://localhost:10081/api/auth/login   -H "Content-Type: application/json" -d '{"username":"u1","password":"pass1234"}'
```
