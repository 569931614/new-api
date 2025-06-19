# 授权码系统文档

## 文档概览

本目录包含授权码系统的完整文档，包括外部接口使用指南、快速开始教程和测试示例。

## 文档列表

### 📚 主要文档

1. **[授权码外部接口使用文档](./auth_code_api.md)**
   - 完整的接口文档
   - 详细的安全机制说明
   - 多语言客户端实现示例
   - 最佳实践和错误处理

2. **[快速开始指南](./auth_code_quick_start.md)**
   - 简化版使用指南
   - 快速集成示例
   - 常见问题解答

3. **[测试示例](./auth_code_test_examples.md)**
   - cURL 测试命令
   - Postman 测试配置
   - JavaScript/Python 测试脚本
   - 测试检查清单

## 功能概述

### 🔐 核心功能

- **机器码绑定**：将授权码与特定设备绑定
- **安全验证**：采用挑战-响应机制防止篡改
- **状态管理**：完整的授权码生命周期管理
- **多平台支持**：支持 Windows、Linux、macOS

### 🛡️ 安全特性

- **HMAC-SHA256 签名**：防止挑战值被篡改
- **时间窗口限制**：防止重放攻击
- **设备绑定**：一码一机，防止跨设备使用
- **双重验证**：挑战完整性 + 响应正确性

## 接口列表

| 接口 | 方法 | 地址 | 用途 |
|------|------|------|------|
| 绑定机器码 | POST | `/api/auth/bind` | 将授权码与机器码绑定 |
| 验证授权码 | POST | `/api/auth/validate` | 验证授权码有效性 |
| 获取渠道列表 | POST | `/api/auth/channels` | 根据授权码获取可用渠道列表 |

## 快速开始

### 1. 绑定机器码

```bash
curl -X POST http://your-domain/api/auth/bind \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "your_auth_code",
    "machine_code": "your_machine_code"
  }'
```

### 2. 验证授权码

```bash
# 第一步：获取挑战
curl -X POST http://your-domain/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "your_auth_code",
    "machine_code": "your_machine_code"
  }'

# 第二步：计算响应并提交验证
# response = SHA256(challenge)
curl -X POST http://your-domain/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "your_auth_code",
    "machine_code": "your_machine_code",
    "challenge": "challenge_from_step1",
    "response": "sha256_hash_of_challenge"
  }'
```

## 客户端示例

### JavaScript

```javascript
class AuthClient {
  async bind(authCode, machineCode) {
    const response = await fetch('/api/auth/bind', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });
    return response.json();
  }

  async validate(authCode, machineCode) {
    // 获取挑战
    const challengeResp = await fetch('/api/auth/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });
    
    const challengeData = await challengeResp.json();
    if (!challengeData.success) throw new Error(challengeData.message);

    // 计算响应
    const challenge = challengeData.challenge;
    const encoder = new TextEncoder();
    const data = encoder.encode(challenge);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

    // 提交验证
    const validateResp = await fetch('/api/auth/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        auth_code: authCode,
        machine_code: machineCode,
        challenge: challenge,
        response: response
      })
    });

    const validateData = await validateResp.json();
    if (!validateData.success) throw new Error(validateData.message);
    return validateData.data;
  }

  async getChannels(authCode, machineCode) {
    // 获取挑战
    const challengeResp = await fetch('/api/auth/channels', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });

    const challengeData = await challengeResp.json();
    if (!challengeData.success) throw new Error(challengeData.message);

    // 计算响应
    const challenge = challengeData.challenge;
    const encoder = new TextEncoder();
    const data = encoder.encode(challenge);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

    // 提交验证
    const channelsResp = await fetch('/api/auth/channels', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        auth_code: authCode,
        machine_code: machineCode,
        challenge: challenge,
        response: response
      })
    });

    const channelsData = await channelsResp.json();
    if (!channelsData.success) throw new Error(channelsData.message);
    return channelsData.data;
  }
}
```

### Python

```python
import hashlib
import requests

class AuthClient:
    def __init__(self, base_url=""):
        self.base_url = base_url
    
    def bind(self, auth_code, machine_code):
        response = requests.post(f'{self.base_url}/api/auth/bind', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })
        return response.json()
    
    def validate(self, auth_code, machine_code):
        # 获取挑战
        challenge_resp = requests.post(f'{self.base_url}/api/auth/validate', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })
        
        challenge_data = challenge_resp.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])
        
        # 计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()
        
        # 提交验证
        validate_resp = requests.post(f'{self.base_url}/api/auth/validate', json={
            'auth_code': auth_code,
            'machine_code': machine_code,
            'challenge': challenge,
            'response': response
        })
        
        validate_data = validate_resp.json()
        if not validate_data['success']:
            raise Exception(validate_data['message'])

        return validate_data['data']

    def get_channels(self, auth_code, machine_code):
        # 获取挑战
        challenge_resp = requests.post(f'{self.base_url}/api/auth/channels', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })

        challenge_data = challenge_resp.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])

        # 计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()

        # 提交验证
        channels_resp = requests.post(f'{self.base_url}/api/auth/channels', json={
            'auth_code': auth_code,
            'machine_code': machine_code,
            'challenge': challenge,
            'response': response
        })

        channels_data = channels_resp.json()
        if not channels_data['success']:
            raise Exception(channels_data['message'])

        return channels_data['data']
```

## 状态说明

| 状态码 | 名称 | 说明 |
|--------|------|------|
| 1 | 启用 | 授权码可正常使用 |
| 2 | 禁用 | 授权码被管理员禁用 |
| 3 | 已使用 | 授权码已被使用（一次性） |
| 4 | 待激活 | 授权码等待机器码绑定 |
| 5 | 激活 | 授权码已绑定机器码并激活 |

## 用户类型

| 类型值 | 名称 | 权限级别 |
|--------|------|----------|
| 1 | 普通用户 | 基础权限 |
| 10 | 管理员 | 管理权限 |
| 100 | 超级管理员 | 最高权限 |

## 常见问题

### Q: 如何获取授权码？
A: 授权码由管理员在管理面板中创建和分发。

### Q: 机器码如何生成？
A: 建议基于硬件特征（CPU、主板序列号等）生成，确保唯一性和稳定性。

### Q: 授权码可以在多个设备上使用吗？
A: 不可以，每个授权码只能绑定一个机器码。

### Q: 验证失败后如何处理？
A: 检查错误信息，对于网络错误可以重试，对于业务错误需要根据具体情况处理。

### Q: 挑战值的有效期是多久？
A: 挑战值有效期为5分钟。

## 技术支持

如果在使用过程中遇到问题：

1. 首先查看 [完整文档](./auth_code_api.md)
2. 尝试使用 [测试示例](./auth_code_test_examples.md) 进行调试
3. 检查网络连接和服务器状态
4. 联系技术支持团队

## 更新日志

- **v1.0.0** (2024-12-19)
  - 初始版本发布
  - 支持机器码绑定和验证功能
  - 实现挑战-响应安全机制
  - 提供多语言客户端示例

---

*文档最后更新时间：2024年12月19日*
