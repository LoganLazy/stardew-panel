# 🔐 认证系统实现 v0.3.0

**日期**: 2026-07-05  
**版本**: v0.3.0  
**类型**: 安全增强 - 密码保护

---

## ✅ 实现内容

### 🔐 完整的认证系统

为了保护公网访问的安全，添加了完整的密码保护功能：

---

## 📊 新增功能

### 1. 用户登录 ✨
- ✅ 登录页面（用户名 + 密码）
- ✅ Token 认证机制
- ✅ 会话管理（24小时有效期）
- ✅ 自动跳转

### 2. 密码修改 ✨
- ✅ 修改密码页面
- ✅ 旧密码验证
- ✅ 新密码强度检查
- ✅ 修改后自动退出

### 3. 访问控制 ✨
- ✅ 所有 API 需要认证
- ✅ Token 自动附加
- ✅ 401 自动跳转登录页
- ✅ 路由守卫保护

### 4. 默认账号 ✨
- ✅ 首次启动自动创建
- ✅ 控制台显示默认账号
- ✅ 提示立即修改密码

---

## 🏗️ 技术架构

### 后端（Go）

#### 1. 数据库表
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

#### 2. 认证服务（service/auth.go）
```go
✅ Login()                // 登录
✅ ChangePassword()       // 修改密码
✅ InitDefaultUser()      // 初始化默认用户
✅ ValidateToken()        // 验证 Token
✅ Logout()               // 退出登录
```

**特性**：
- 使用 bcrypt 加密密码
- 生成 64 位随机 Token
- 内存会话存储（简化实现）
- 24 小时 Token 过期

#### 3. 认证处理器（handler/auth.go）
```go
✅ Login()                // POST /auth/login
✅ Logout()               // POST /auth/logout
✅ ChangePassword()       // POST /auth/change-password
✅ CheckAuth()            // GET /auth/check
✅ AuthMiddleware()       // 认证中间件
```

#### 4. 路由配置
```go
// 无需认证的路由
/health                    // 健康检查
/api/v1/auth/login         // 登录

// 需要认证的路由（所有其他 API）
/api/v1/server/*
/api/v1/mods/*
/api/v1/players/*
/api/v1/saves/*
/api/v1/logs
/api/v1/install/*
```

---

### 前端（Vue）

#### 1. 登录页面（Login.vue）
```vue
✅ 用户名输入
✅ 密码输入（隐藏）
✅ 登录按钮
✅ 错误提示
✅ 默认账号提示
```

**设计**：
- 温暖的渐变背景
- 居中的登录框
- 清晰的表单
- 友好的提示信息

#### 2. 设置页面（Settings.vue）
```vue
✅ 修改密码表单
✅ 账户信息显示
✅ 退出登录按钮
```

#### 3. API 拦截器（api/index.js）
```javascript
// 请求拦截器
✅ 自动附加 Token（Authorization: Bearer XXX）

// 响应拦截器
✅ 401 自动跳转登录页
✅ 自动清除过期 Token
```

#### 4. 路由守卫（router/index.js）
```javascript
✅ 检查登录状态
✅ 未登录自动跳转登录页
✅ 已登录访问登录页跳转首页
```

#### 5. 导航更新（App.vue）
```vue
✅ 登录页隐藏导航
✅ 添加"设置"链接
✅ 响应式布局
```

---

## 🔐 安全特性

### 1. 密码安全
```
✅ bcrypt 加密（成本因子 10）
✅ 明文密码从不保存
✅ 密码强度检查（最少 6 位）
```

### 2. Token 安全
```
✅ 64 位随机 Token
✅ 24 小时自动过期
✅ 每次登录生成新 Token
✅ 退出登录立即失效
```

### 3. 传输安全
```
✅ Token 通过 Header 传输
✅ 不在 URL 中暴露
✅ HTTPS 推荐（生产环境）
```

### 4. 会话管理
```
✅ 内存会话存储
✅ 自动过期清理
✅ 单点登录（一个 Token）
```

---

## 📝 使用说明

### 首次启动

1. **启动服务器**
   ```bash
   cd server && go run main.go
   ```

2. **控制台输出**
   ```
   ✅ 默认用户已创建
      用户名: admin
      密码: admin123
      ⚠️  请立即登录后修改密码！
   ```

3. **访问面板**
   ```
   http://localhost:8080
   自动跳转到 /login
   ```

4. **登录**
   - 用户名: `admin`
   - 密码: `admin123`

5. **立即修改密码**
   - 点击导航栏"⚙️ 设置"
   - 修改密码
   - 重新登录

---

### 修改密码

1. **进入设置页面**
   - 点击导航栏"⚙️ 设置"

2. **填写表单**
   - 当前密码: 输入旧密码
   - 新密码: 至少 6 位
   - 确认新密码: 再次输入

3. **提交**
   - 点击"修改密码"
   - 确认修改
   - 自动退出登录
   - 使用新密码重新登录

---

### 退出登录

1. **进入设置页面**
2. **点击"退出登录"**
3. **自动跳转到登录页**

---

## 🚨 安全建议

### 生产环境必做

1. **立即修改默认密码**
   ```
   ❌ 不要使用 admin123
   ✅ 使用强密码（8位以上，包含字母数字）
   ```

2. **启用 HTTPS**
   ```bash
   # 使用 Nginx 反向代理
   server {
       listen 443 ssl;
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

3. **限制访问 IP**
   ```nginx
   # Nginx 配置
   allow 192.168.1.0/24;
   deny all;
   ```

4. **定期修改密码**
   ```
   建议：每 3 个月修改一次
   ```

5. **使用防火墙**
   ```bash
   # 只允许特定端口
   ufw allow 443/tcp
   ufw enable
   ```

---

## 📊 文件变更

### 后端新增
```
server/models/auth.go           ✨ 认证模型
server/service/auth.go          ✨ 认证服务
server/handler/auth.go          ✨ 认证处理器
server/database/database.go     🔧 添加 users 表
server/router/router.go         🔧 添加认证路由和中间件
server/main.go                  🔧 初始化默认用户
server/go.mod                   🔧 添加 golang.org/x/crypto
```

### 前端新增
```
web/src/views/Login.vue         ✨ 登录页面
web/src/views/Settings.vue      ✨ 设置页面
web/src/api/api.js              🔧 添加认证 API
web/src/api/index.js            🔧 添加 Token 拦截器
web/src/router/index.js         🔧 添加路由守卫
web/src/App.vue                 🔧 导航更新
```

---

## 🎯 后续改进（可选）

### 短期
1. 记住密码功能
2. 密码找回功能
3. 登录日志记录

### 中期
4. 多用户支持
5. 角色权限系统
6. 会话持久化（Redis）

### 长期
7. 双因素认证（2FA）
8. OAuth 登录
9. API Key 支持

---

## ✅ 测试清单

- [x] 登录功能正常
- [x] 未登录自动跳转
- [x] Token 正确附加
- [x] 401 自动处理
- [x] 修改密码功能
- [x] 退出登录功能
- [x] 默认用户创建
- [x] 密码加密存储
- [x] 路由守卫工作

---

## 🎊 总结

### v0.3.0 - 安全增强

**新增功能**：
- ✅ 完整的认证系统
- ✅ 登录/退出功能
- ✅ 密码修改功能
- ✅ Token 认证机制
- ✅ 路由守卫保护

**安全提升**：
- 🔒 所有 API 需要认证
- 🔒 密码 bcrypt 加密
- 🔒 Token 自动过期
- 🔒 401 自动处理

**用户体验**：
- ✨ 温暖的登录界面
- ✨ 清晰的设置页面
- ✨ 友好的提示信息
- ✨ 自动跳转

---

**现在可以安全地将面板暴露到公网了！** 🚀
