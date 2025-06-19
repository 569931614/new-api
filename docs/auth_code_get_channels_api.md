# 授权码获取渠道列表 GET 接口文档

## 概述

新增了一个简化的 GET 接口，用于根据授权码获取关联分组下的可用渠道列表。这个接口相比原有的 POST 接口更加简单，不需要复杂的两步验证流程。

## 接口信息

- **接口地址：** `GET /api/auth/channels`
- **请求方式：** GET
- **Content-Type：** 无需设置（URL参数）
- **认证要求：** 无需认证
- **用途：** 根据授权码获取关联分组下的可用渠道列表

## 请求参数

### URL 参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |

### 请求示例

```bash
# 基本请求
curl -X GET "http://localhost:3000/api/auth/channels?auth_code=YOUR_AUTH_CODE"

# PowerShell 请求
Invoke-WebRequest -Uri "http://localhost:3000/api/auth/channels?auth_code=YOUR_AUTH_CODE" -Method GET
```

## 响应格式

### 成功响应

```json
{
  "success": true,
  "message": "获取渠道列表成功",
  "data": {
    "auth_groups": ["group1", "group2"],
    "channels": [
      {
        "id": 1,
        "name": "渠道名称",
        "type": 1,
        "business_type": 1,
        "status": 1,
        "models": ["gpt-3.5-turbo", "gpt-4"],
        "group": "group1",
        "priority": 0,
        "weight": 100
      }
    ],
    "total": 1
  }
}
```

### 错误响应

#### 授权码参数缺失
```json
{
  "success": false,
  "message": "授权码参数不能为空"
}
```

#### 授权码不存在
```json
{
  "success": false,
  "message": "授权码不存在"
}
```

#### 授权码无效或已过期
```json
{
  "success": false,
  "message": "授权码无效或已过期"
}
```

#### 授权码未激活
```json
{
  "success": false,
  "message": "授权码未激活"
}
```

## 响应字段说明

### data 字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| auth_groups | array | 授权码关联的分组列表 |
| channels | array | 可用渠道列表 |
| total | int | 渠道总数 |

### channels 数组中的字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 渠道ID |
| name | string | 渠道名称 |
| type | int | 渠道类型（技术类型：OpenAI、Claude等） |
| business_type | int | 业务类型（1:对话, 2:应用, 3:工作流） |
| status | int | 渠道状态（1:启用, 2:禁用） |
| models | array | 支持的模型列表 |
| group | string | 渠道所属分组 |
| priority | int | 优先级 |
| weight | int | 权重 |

## 使用说明

### 授权码状态要求

- 授权码必须存在于系统中
- 授权码必须处于激活状态（status = 5）
- 授权码不能过期（expired_time = -1 或大于当前时间戳）

### 分组关联

- 授权码必须关联至少一个分组
- 只返回授权码关联分组下的渠道
- 支持多分组授权码（用逗号分隔）

### 安全性

- 接口不暴露敏感信息（如API密钥等）
- 只返回渠道的基本信息
- 无需机器码验证（简化版接口）

## 与原有 POST 接口的区别

| 特性 | GET 接口 | POST 接口 |
|------|----------|-----------|
| 请求方式 | GET | POST |
| 参数传递 | URL参数 | JSON Body |
| 验证流程 | 单步验证 | 两步验证（挑战-响应） |
| 机器码验证 | 不需要 | 需要 |
| 安全性 | 较低 | 较高 |
| 使用复杂度 | 简单 | 复杂 |

## 测试示例

### JavaScript 测试

```javascript
async function getChannels(authCode) {
  try {
    const response = await fetch(`http://localhost:3000/api/auth/channels?auth_code=${authCode}`);
    const data = await response.json();
    
    if (data.success) {
      console.log('渠道列表:', data.data.channels);
      console.log('关联分组:', data.data.auth_groups);
    } else {
      console.error('获取失败:', data.message);
    }
  } catch (error) {
    console.error('请求失败:', error);
  }
}

// 使用示例
getChannels('YOUR_AUTH_CODE');
```

### Python 测试

```python
import requests

def get_channels(auth_code):
    try:
        response = requests.get(f'http://localhost:3000/api/auth/channels?auth_code={auth_code}')
        data = response.json()
        
        if data['success']:
            print('渠道列表:', data['data']['channels'])
            print('关联分组:', data['data']['auth_groups'])
        else:
            print('获取失败:', data['message'])
    except Exception as e:
        print('请求失败:', e)

# 使用示例
get_channels('YOUR_AUTH_CODE')
```

## 注意事项

1. **授权码状态**：确保授权码处于激活状态
2. **分组配置**：授权码必须关联有效的分组
3. **渠道状态**：只返回启用状态的渠道
4. **缓存考虑**：建议客户端适当缓存结果，避免频繁请求
5. **错误处理**：务必处理各种错误情况

## 更新日志

- **2025-06-19**：新增 GET 接口，提供简化的渠道列表获取方式
