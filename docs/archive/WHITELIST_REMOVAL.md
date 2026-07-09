# 白名单功能移除说明

## 变更内容

**日期**: 2026-07-05  
**版本**: v0.2.0

### 移除功能
- ❌ 白名单管理（添加/移除/查询）
- ❌ `whitelist` 数据表
- ❌ 相关 API 端点

### 新增功能
- ✅ **在线玩家实时监控**
- ✅ 自动从游戏日志解析玩家信息
- ✅ 在线人数统计
- ✅ 在线时长计算
- ✅ 历史玩家记录
- ✅ 踢出玩家功能
- ✅ 自动刷新（10秒间隔）

## 为什么移除白名单？

### 原因分析

1. **功能定位不符**
   - 星露谷多人游戏通常是小圈子（朋友/家人）
   - 服务器地址不公开，陌生人无法得知
   - 白名单在私人服务器场景下没有实际意义

2. **安全性问题**
   - 游戏基于角色名识别，玩家可随意改名
   - 无法真正验证玩家身份
   - IP 地址不可靠（动态 IP、共用 IP）
   - 容易被冒充

3. **用户体验**
   - 增加额外的管理负担
   - 朋友加入前需要提前配置
   - 如果忘记添加会导致连接失败

### 新功能优势

**在线玩家监控更实用**：
- ✅ 实时看到谁在玩
- ✅ 自动记录玩家信息
- ✅ 在线时长统计
- ✅ 历史玩家追踪
- ✅ 无需提前配置

## 技术实现

### 数据库变更

**旧表**（已删除）:
```sql
CREATE TABLE whitelist (
    id INTEGER PRIMARY KEY,
    player_name TEXT NOT NULL UNIQUE,
    player_id TEXT,
    added_at DATETIME
);
```

**新表**:
```sql
CREATE TABLE online_players (
    id INTEGER PRIMARY KEY,
    player_name TEXT NOT NULL,
    connected_at DATETIME,
    last_seen DATETIME,
    is_online BOOLEAN
);
```

### API 变更

**移除的端点**:
- `POST /api/v1/players/whitelist` - 添加白名单
- `DELETE /api/v1/players/whitelist/:id` - 移除白名单

**新增的端点**:
- `POST /api/v1/players/refresh` - 手动刷新玩家列表
- `DELETE /api/v1/players/cleanup` - 清理历史记录

**修改的端点**:
- `GET /api/v1/players` - 现在返回在线玩家列表 + 统计信息

### 日志解析

新功能通过解析游戏日志文件来获取玩家信息：

```go
// 匹配玩家加入
joinPattern := regexp.MustCompile(`(\w+) joined the game`)

// 匹配玩家离开
leftPattern := regexp.MustCompile(`(\w+) left the game`)
```

## 前端变更

### 页面改造

**Players.vue** 完全重写：
- 移除白名单输入框
- 新增统计卡片（在线人数/历史人数）
- 在线玩家列表显示：
  - 玩家名称
  - 在线状态徽章
  - 加入时间
  - 在线时长
- 自动刷新（10秒）
- 手动刷新按钮
- 清理历史记录功能

### 样式优化

- 统计数据卡片
- 玩家条目悬停效果
- 在线徽章显示
- 空状态提示

## 迁移指南

### 对现有用户的影响

**如果你已经在使用旧版本**：

1. **数据迁移**
   - 白名单数据不会自动迁移
   - 升级后白名单表将被删除
   - 建议升级前导出白名单数据（如需要）

2. **API 兼容性**
   - 旧的白名单 API 返回 410 Gone
   - 前端需要更新到新版本

3. **功能替代**
   - 使用在线玩家监控替代白名单
   - 踢出功能保留，可以管理不欢迎的玩家

### 升级步骤

```bash
# 1. 备份数据库（可选）
cp data/stardew-panel.db data/stardew-panel.db.backup

# 2. 更新代码
git pull

# 3. 重新构建
docker-compose down
docker-compose up -d --build

# 4. 数据库会自动创建新表
```

## 未来计划

如果确实需要访问控制，可以考虑：

1. **Steam ID 验证**（高优先级）
   - 绑定 Steam 账号
   - 更安全的身份验证

2. **密码保护**（中优先级）
   - 服务器级别密码
   - 简单但有效

3. **邀请码系统**（低优先级）
   - 生成一次性邀请码
   - 适合半公开服务器

## 反馈

如果你对这个改动有任何意见或建议，欢迎提交 Issue：
https://github.com/LoganLazy/stardew-panel/issues

---

**总结**：白名单功能被更实用的"在线玩家监控"取代，更符合私人服务器的使用场景。
