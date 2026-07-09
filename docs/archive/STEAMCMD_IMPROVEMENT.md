# SteamCMD 功能改进

## 🔐 更新内容

### 问题
之前的 SteamCMD 实现使用匿名登录：
```go
"+login", "anonymous"  // ❌ 无法下载需要购买的游戏
```

**问题**：
- 星露谷物语需要账号才能下载
- 匿名登录会失败
- 用户无法使用 SteamCMD 功能

### 解决方案
现在支持账号密码登录：
```go
"+login", username, password  // ✅ 使用用户提供的账号
```

---

## 🆕 新功能

### 前端改进

**Setup.vue 新增输入框**：
```html
<input v-model="steamUsername" placeholder="your_steam_username" />
<input type="password" v-model="steamPassword" placeholder="••••••••" />
```

**特性**：
- ✅ 密码输入框（type="password"）
- ✅ 按钮禁用逻辑（必须填写账号密码）
- ✅ 使用后自动清空密码
- ✅ 详细的警告提示

### 后端改进

**service/install.go**：
- ✅ 接收用户名和密码参数
- ✅ 参数验证（不能为空）
- ✅ 详细的错误识别：
  - 密码错误
  - 用户名不存在
  - Steam Guard 验证
  - 未购买游戏

**handler/install.go**：
- ✅ 验证必填参数
- ✅ 返回友好错误信息

**models/models.go**：
- ✅ 新增 `SteamUsername` 字段
- ✅ 新增 `SteamPassword` 字段

---

## 🔒 安全性说明

### 密码处理

**后端**：
- ✅ 密码仅用于调用 SteamCMD
- ✅ 不保存到数据库
- ✅ 不记录到日志（日志只显示用户名）
- ✅ 仅在内存中临时存储

**前端**：
- ✅ 使用 `type="password"` 隐藏输入
- ✅ 使用后立即清空
- ✅ 不保存到 localStorage
- ✅ 明确提示用户密码不会被保存

### 开源项目的安全性

**为什么不用担心**：
1. **自托管**：用户在自己的服务器上部署
2. **开源**：代码完全透明，可以审查
3. **本地处理**：密码只在用户自己的服务器处理
4. **不上传**：密码不会发送到第三方服务器

**vs 商业产品**：
- ❌ 商业产品可能上传密码到云端
- ❌ 闭源无法验证安全性
- ✅ 开源自托管最安全

---

## 🚨 重要提示

### 界面警告信息

```
⚠️ 重要提示：
- 需要 Steam 账号并且已购买星露谷物语
- 密码仅在服务器本地使用，不会被保存或上传
- 如果启用了 Steam Guard，需要先在服务器上手动登录 SteamCMD
- 下载速度取决于网络环境，国内可能较慢
- 建议使用代理或 VPN 以获得更好的下载速度
```

### Steam Guard 处理

**如果账号启用了 Steam Guard**：

1. 用户会收到错误：`需要 Steam Guard 验证码`
2. 解决方法：
   ```bash
   # 在服务器上手动运行
   steamcmd
   Steam> login your_username
   # 输入密码和验证码
   # 完成验证后，面板即可使用
   ```

3. 验证信息会保存在 SteamCMD 中，之后面板可以直接使用

---

## 📊 错误处理

### 后端错误识别

```go
if strings.Contains(outputStr, "Invalid Password") {
    return "Steam 密码错误"
}
if strings.Contains(outputStr, "Invalid Username") {
    return "Steam 用户名不存在"
}
if strings.Contains(outputStr, "Steam Guard") {
    return "需要 Steam Guard 验证码"
}
if strings.Contains(outputStr, "No subscription") {
    return "该账号未购买星露谷物语"
}
```

### 前端友好提示

```javascript
if (errorMsg.includes('密码错误')) {
  alert('❌ Steam 密码错误，请检查后重试')
}
else if (errorMsg.includes('Steam Guard')) {
  alert('❌ 需要 Steam Guard 验证码\n\n请先在服务器上手动登录...')
}
else if (errorMsg.includes('未购买')) {
  alert('❌ 该账号未购买星露谷物语\n\n需要拥有游戏才能下载服务端')
}
```

---

## 🎯 使用流程

### 正常流程

1. 用户在 Setup 页面选择 "SteamCMD 下载"
2. 输入 Steam 用户名
3. 输入 Steam 密码
4. 选择安装路径
5. 点击"开始下载"
6. 等待 5-15 分钟
7. 下载完成，密码自动清空

### Steam Guard 流程

1. 用户首次使用遇到 Steam Guard 错误
2. SSH 登录服务器
3. 运行 `steamcmd`
4. 手动登录并输入验证码
5. 完成验证
6. 回到面板继续使用

---

## ✅ 测试清单

- [x] 输入账号密码能正常提交
- [x] 密码框类型为 password
- [x] 必填验证正常工作
- [x] 后端接收参数正确
- [x] 错误信息正确识别
- [x] 密码不显示在日志中
- [x] 使用后密码自动清空
- [x] 警告信息正确显示

---

## 📝 API 变更

### 请求

**旧版**：
```json
{
  "path": "/opt/stardew-server",
  "install_smapi": true
}
```

**新版**：
```json
{
  "path": "/opt/stardew-server",
  "steam_username": "your_username",
  "steam_password": "your_password",
  "install_smapi": true
}
```

### 响应

**成功**：
```json
{
  "success": true,
  "installation": {
    "game_path": "/opt/stardew-server",
    "install_method": "steamcmd",
    "has_smapi": true,
    "version": "1.6.9"
  }
}
```

**失败**：
```json
{
  "error": "Steam 密码错误"
}
```

---

## 🎊 总结

SteamCMD 功能现在**完全可用**！

**改进**：
- ✅ 支持账号密码登录
- ✅ 详细的错误提示
- ✅ 安全的密码处理
- ✅ 友好的用户体验
- ✅ Steam Guard 指引

**适用场景**：
- 服务器没有游戏文件
- 想要干净的安装
- 需要自动化部署

**推荐度**：⭐⭐⭐（如果有 Steam 账号和良好网络）
