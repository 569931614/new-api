# 外部接口文档

## 概述

本文档描述了系统提供的外部接口，主要用于授权码验证、机器码绑定和渠道获取等功能。

## 基础信息

- **基础URL**: `http://localhost:3000`
- **Content-Type**: `application/json`
- **字符编码**: UTF-8

## 接口列表

### 1. 机器码绑定接口

#### 接口描述
用于将机器码与授权码进行绑定，绑定后授权码状态变为"激活"。

#### 请求信息
- **URL**: `/api/auth/bind`
- **方法**: `POST`
- **认证**: 无需认证

#### 请求参数
```json
{
  "auth_code": "string",     // 授权码（必填）
  "machine_code": "string"   // 机器码（必填）
}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "success": true,
  "message": "机器码绑定成功",
  "data": {
    "auth_code": "ABC123DEF456",
    "machine_code": "MACHINE001",
    "status": 5,
    "bind_time": 1703123456
  }
}
```

**失败响应 (400)**:
```json
{
  "success": false,
  "message": "授权码不存在或已被使用"
}
```

#### 状态码说明
- `200`: 绑定成功
- `400`: 请求参数错误或业务逻辑错误
- `500`: 服务器内部错误

---

### 2. 授权码验证接口

#### 接口描述
验证授权码和机器码的有效性，用于客户端登录验证。

#### 请求信息
- **URL**: `/api/auth/validate`
- **方法**: `POST`
- **认证**: 无需认证

#### 请求参数
```json
{
  "auth_code": "string",     // 授权码（必填）
  "machine_code": "string"   // 机器码（必填）
}
```

#### 响应示例
**成功响应 (200)**:
```json
{
  "success": true,
  "message": "验证成功",
  "data": {
    "valid": true,
    "auth_code": "ABC123DEF456",
    "name": "测试授权码",
    "user_type": 1,
    "expired_time": 1735660800,
    "status": 5,
    "is_bot": false,
    "wx_auto_x_code": "WX001",
    "groups": ["vip", "premium"]
  }
}
```

**失败响应 (400)**:
```json
{
  "success": false,
  "message": "授权码验证失败",
  "data": {
    "valid": false
  }
}
```

#### 字段说明
- `valid`: 验证是否通过
- `user_type`: 用户类型 (1: 普通用户, 10: 管理员, 100: 超级管理员)
- `expired_time`: 过期时间戳，-1表示永不过期
- `status`: 状态 (1: 启用, 4: 待激活, 5: 激活)
- `is_bot`: 是否为机器人账户
- `groups`: 分组列表（支持多个分组）

---

### 3. 获取渠道列表接口 🆕

#### 接口描述
根据授权码获取可用的渠道列表，支持多分组查询。

#### 请求信息
- **URL**: `/api/auth/channels`
- **方法**: `GET`
- **认证**: 无需认证

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| business_type | int | 否 | 业务类型过滤（1:对话, 2:应用, 3:工作流），不传则返回所有类型 |

**示例URL**:
- 获取所有渠道：`/api/auth/channels?auth_code=ABC123DEF456`
- 获取对话类型渠道：`/api/auth/channels?auth_code=ABC123DEF456&business_type=1`
- 获取应用类型渠道：`/api/auth/channels?auth_code=ABC123DEF456&business_type=2`
- 获取工作流类型渠道：`/api/auth/channels?auth_code=ABC123DEF456&business_type=3`

#### 响应示例
**成功响应 (200)**:
```json
{
  "success": true,
  "message": "获取成功",
  "data": {
    "channels": [
      {
        "id": 1,
        "name": "OpenAI官方",
        "type": 1,
        "business_type": 1,
        "status": 1,
        "models": ["gpt-3.5-turbo", "gpt-4"],
        "group": "vip,premium",
        "priority": 0,
        "weight": 100
      },
      {
        "id": 2,
        "name": "Claude API",
        "type": 14,
        "business_type": 1,
        "status": 1,
        "models": ["claude-3-sonnet", "claude-3-opus"],
        "group": "premium",
        "priority": 0,
        "weight": 80
      },
      {
        "id": 3,
        "name": "Dify工作流",
        "type": 37,
        "business_type": 3,
        "status": 1,
        "models": ["dify-workflow"],
        "group": "premium",
        "priority": 0,
        "weight": 60
      }
    ],
    "total": 3,
    "auth_groups": ["vip", "premium"]
  }
}
```

**无权限响应 (200)**:
```json
{
  "success": true,
  "message": "获取成功",
  "data": {
    "channels": [],
    "total": 0,
    "auth_groups": []
  }
}
```

**失败响应 (400)**:
```json
{
  "success": false,
  "message": "授权码验证失败"
}
```

#### 字段说明
- `channels`: 可用渠道列表
- `total`: 渠道总数
- `auth_groups`: 授权码关联的分组列表
- `type`: 渠道技术类型 (1: OpenAI, 14: Claude, 37: Dify, 等)
- `business_type`: 业务类型 (1: 对话, 2: 应用, 3: 工作流)
- `models`: 支持的模型列表（数组格式）
- `group`: 渠道所属分组（逗号分隔，支持多分组）

---

## 业务类型过滤功能 🆕

### 业务类型说明
系统支持三种业务类型的渠道：
- **1 - 对话类型**: 用于聊天对话的渠道，如ChatGPT、Claude等
- **2 - 应用类型**: 用于特定应用功能的渠道，如图像生成、语音合成等
- **3 - 工作流类型**: 用于复杂工作流的渠道，如Dify工作流等

### 过滤使用方法
通过在请求URL中添加 `business_type` 参数来过滤特定类型的渠道：

```bash
# 获取所有类型渠道
GET /api/auth/channels?auth_code=ABC123

# 只获取对话类型渠道
GET /api/auth/channels?auth_code=ABC123&business_type=1

# 只获取应用类型渠道
GET /api/auth/channels?auth_code=ABC123&business_type=2

# 只获取工作流类型渠道
GET /api/auth/channels?auth_code=ABC123&business_type=3
```

### 应用场景
- **客户端分类显示**: 根据业务类型在不同界面展示相应渠道
- **功能模块隔离**: 不同功能模块只获取对应类型的渠道
- **性能优化**: 减少不必要的渠道数据传输

---

## 多分组功能说明 🆕

### 分组匹配逻辑
1. **授权码分组**: 支持多个分组，用逗号分隔存储，如 `"vip,premium,enterprise"`
2. **渠道分组**: 支持多个分组，用逗号分隔存储，如 `"vip,premium"`
3. **匹配规则**: 授权码的任一分组与渠道的任一分组匹配即可访问该渠道

### 示例场景
```
授权码分组: "vip,premium"
渠道A分组: "vip"           -> ✅ 可访问（vip匹配）
渠道B分组: "premium,pro"   -> ✅ 可访问（premium匹配）
渠道C分组: "enterprise"    -> ❌ 不可访问（无匹配分组）
渠道D分组: ""              -> ❌ 不可访问（渠道无分组）
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误或业务逻辑错误 |
| 401 | 未授权访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 常见错误信息

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| "授权码不存在" | 提供的授权码无效 | 检查授权码是否正确 |
| "授权码已过期" | 授权码超过有效期 | 联系管理员更新授权码 |
| "机器码不匹配" | 机器码与绑定的不一致 | 使用正确的机器码 |
| "授权码已被禁用" | 授权码状态为禁用 | 联系管理员启用授权码 |
| "该机器码已被其他授权码绑定" | 机器码重复绑定 | 使用未绑定的机器码 |

## 使用示例

### JavaScript示例
```javascript
// 绑定机器码
async function bindMachineCode(authCode, machineCode) {
  const response = await fetch('/api/auth/bind', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      auth_code: authCode,
      machine_code: machineCode
    })
  });
  
  const result = await response.json();
  return result;
}

// 验证授权码
async function validateAuth(authCode, machineCode) {
  const response = await fetch('/api/auth/validate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      auth_code: authCode,
      machine_code: machineCode
    })
  });
  
  const result = await response.json();
  return result;
}

// 获取可用渠道
async function getChannels(authCode, businessType = null) {
  let url = `/api/auth/channels?auth_code=${authCode}`;
  if (businessType) {
    url += `&business_type=${businessType}`;
  }

  const response = await fetch(url, {
    method: 'GET'
  });

  const result = await response.json();
  return result;
}

// 获取对话类型渠道
async function getChatChannels(authCode) {
  return getChannels(authCode, 1);
}

// 获取应用类型渠道
async function getAppChannels(authCode) {
  return getChannels(authCode, 2);
}

// 获取工作流类型渠道
async function getWorkflowChannels(authCode) {
  return getChannels(authCode, 3);
}
```

### Python示例
```python
import requests
import json

def bind_machine_code(auth_code, machine_code):
    url = "http://localhost:3000/api/auth/bind"
    data = {
        "auth_code": auth_code,
        "machine_code": machine_code
    }
    
    response = requests.post(url, json=data)
    return response.json()

def validate_auth(auth_code, machine_code):
    url = "http://localhost:3000/api/auth/validate"
    data = {
        "auth_code": auth_code,
        "machine_code": machine_code
    }
    
    response = requests.post(url, json=data)
    return response.json()

def get_channels(auth_code, business_type=None):
    url = f"http://localhost:3000/api/auth/channels?auth_code={auth_code}"
    if business_type:
        url += f"&business_type={business_type}"

    response = requests.get(url)
    return response.json()

def get_chat_channels(auth_code):
    """获取对话类型渠道"""
    return get_channels(auth_code, 1)

def get_app_channels(auth_code):
    """获取应用类型渠道"""
    return get_channels(auth_code, 2)

def get_workflow_channels(auth_code):
    """获取工作流类型渠道"""
    return get_channels(auth_code, 3)
```

---

## 最佳实践

### 1. 安全建议
- **HTTPS**: 生产环境必须使用HTTPS协议
- **参数验证**: 客户端应验证所有输入参数
- **错误处理**: 妥善处理所有可能的错误情况
- **重试机制**: 实现合理的重试机制，避免频繁请求

### 2. 性能优化
- **缓存策略**: 客户端可缓存渠道列表，减少重复请求
- **连接复用**: 使用HTTP连接池提高性能
- **超时设置**: 设置合理的请求超时时间

### 3. 集成流程建议
```
1. 应用启动 -> 验证授权码
2. 首次使用 -> 绑定机器码
3. 定期检查 -> 验证授权码有效性
4. 获取服务 -> 获取可用渠道列表
5. 异常处理 -> 处理各种错误情况
```

## 状态码详解

### 授权码状态
- `1`: 启用 - 可正常使用
- `2`: 禁用 - 已被管理员禁用
- `3`: 已使用 - 一次性授权码已被使用
- `4`: 待激活 - 需要绑定机器码激活
- `5`: 激活 - 已绑定机器码，可正常使用

### 用户类型
- `1`: 普通用户 - 基础权限
- `10`: 管理员 - 管理权限
- `100`: 超级管理员 - 最高权限

## 故障排查

### 常见问题
1. **连接超时**: 检查网络连接和服务器状态
2. **授权失败**: 验证授权码和机器码是否正确
3. **无可用渠道**: 检查授权码分组配置
4. **接口返回500**: 查看服务器日志排查问题

### 调试技巧
- 使用curl命令测试接口
- 检查请求和响应的JSON格式
- 查看服务器日志获取详细错误信息

---

## 更新日志

### v1.2.0 (2025-06-19)
- 🆕 新增业务类型过滤功能
- 🆕 支持按业务类型（对话、应用、工作流）获取渠道
- 🔧 优化渠道数据结构，添加business_type字段
- 📝 更新接口文档和使用示例

### v1.1.0 (2025-06-19)
- 🆕 新增多分组支持功能
- 🆕 优化渠道获取接口，支持多分组匹配
- 🔧 改进分组匹配逻辑
- 📝 更新接口文档

### v1.0.0 (2025-06-18)
- 🎉 初始版本发布
- ✅ 机器码绑定接口
- ✅ 授权码验证接口
- ✅ 渠道获取接口
