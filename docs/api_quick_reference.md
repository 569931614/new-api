# API 快速参考

## 基础信息
- **基础URL**: `http://localhost:3000`
- **Content-Type**: `application/json`

## 接口概览

| 接口 | 方法 | 路径 | 功能 |
|------|------|------|------|
| 机器码绑定 | POST | `/api/auth/bind` | 绑定机器码到授权码 |
| 授权码验证 | POST | `/api/auth/validate` | 验证授权码有效性 |
| 获取渠道列表 | GET | `/api/auth/channels` | 获取可用渠道 |

## 快速示例

### 1. 机器码绑定
```bash
curl -X POST http://localhost:3000/api/auth/bind \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "ABC123DEF456",
    "machine_code": "MACHINE001"
  }'
```

### 2. 授权码验证
```bash
curl -X POST http://localhost:3000/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "ABC123DEF456",
    "machine_code": "MACHINE001"
  }'
```

### 3. 获取渠道列表
```bash
curl -X GET "http://localhost:3000/api/auth/channels?auth_code=ABC123DEF456"
```

## 响应格式

### 成功响应
```json
{
  "success": true,
  "message": "操作成功",
  "data": { /* 具体数据 */ }
}
```

### 失败响应
```json
{
  "success": false,
  "message": "错误描述"
}
```

## 多分组功能 🆕

### 分组格式
- **授权码分组**: `"vip,premium,enterprise"`
- **渠道分组**: `"vip,premium"`
- **匹配规则**: 任一分组匹配即可访问

### 渠道响应示例
```json
{
  "success": true,
  "data": {
    "channels": [
      {
        "id": 1,
        "name": "OpenAI官方",
        "type": 1,
        "models": "gpt-3.5-turbo,gpt-4",
        "group": "vip,premium"
      }
    ],
    "auth_groups": ["vip", "premium"]
  }
}
```

## 状态码
- `1`: 启用
- `4`: 待激活
- `5`: 激活
- `2`: 禁用
- `3`: 已使用

## 错误处理
- `200`: 成功
- `400`: 参数错误
- `401`: 未授权
- `500`: 服务器错误

## 集成步骤
1. **绑定机器码** → 激活授权码
2. **验证授权码** → 检查有效性
3. **获取渠道** → 获取可用服务
4. **使用服务** → 调用具体渠道API
