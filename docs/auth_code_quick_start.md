# 授权码接口快速开始

## 简介

本文档提供授权码外部接口的快速使用指南，帮助开发者快速集成授权码验证功能。

## 两个核心接口

### 1. 绑定机器码
```
POST /api/auth/bind
```

### 2. 验证授权码
```
POST /api/auth/validate
```

### 3. 获取渠道列表
```
POST /api/auth/channels
```

## 快速集成示例

### JavaScript 版本

```javascript
// 1. 创建客户端类
class AuthClient {
  constructor(baseUrl = '') {
    this.baseUrl = baseUrl;
  }

  // 绑定机器码
  async bind(authCode, machineCode) {
    const response = await fetch(`${this.baseUrl}/api/auth/bind`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });
    return response.json();
  }

  // 验证授权码
  async validate(authCode, machineCode) {
    // 第一步：获取挑战
    const challengeResp = await fetch(`${this.baseUrl}/api/auth/validate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });
    
    const challengeData = await challengeResp.json();
    if (!challengeData.success) throw new Error(challengeData.message);

    // 第二步：计算响应
    const challenge = challengeData.challenge;
    const encoder = new TextEncoder();
    const data = encoder.encode(challenge);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

    // 第三步：提交验证
    const validateResp = await fetch(`${this.baseUrl}/api/auth/validate`, {
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

  // 获取渠道列表
  async getChannels(authCode, machineCode) {
    // 第一步：获取挑战
    const challengeResp = await fetch(`${this.baseUrl}/api/auth/channels`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ auth_code: authCode, machine_code: machineCode })
    });

    const challengeData = await challengeResp.json();
    if (!challengeData.success) throw new Error(challengeData.message);

    // 第二步：计算响应
    const challenge = challengeData.challenge;
    const encoder = new TextEncoder();
    const data = encoder.encode(challenge);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

    // 第三步：提交验证
    const channelsResp = await fetch(`${this.baseUrl}/api/auth/channels`, {
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

  // 获取简单机器码
  getMachineCode() {
    const info = `${navigator.userAgent}-${navigator.platform}`;
    let hash = 0;
    for (let i = 0; i < info.length; i++) {
      hash = ((hash << 5) - hash) + info.charCodeAt(i);
      hash = hash & hash;
    }
    return Math.abs(hash).toString(16);
  }
}

// 2. 使用示例
async function main() {
  const client = new AuthClient('http://localhost:3000');
  const authCode = 'your_auth_code_here';
  const machineCode = client.getMachineCode();

  try {
    // 首次使用：绑定机器码
    await client.bind(authCode, machineCode);
    console.log('绑定成功');

    // 验证授权码
    const userInfo = await client.validate(authCode, machineCode);
    console.log('验证成功:', userInfo);

    // 获取渠道列表
    const channelsInfo = await client.getChannels(authCode, machineCode);
    console.log('渠道信息:', channelsInfo);
    console.log(`分组: ${channelsInfo.group}, 渠道数: ${channelsInfo.total}`);

    // 检查权限
    if (userInfo.user_type >= 10) {
      console.log('管理员权限');
    } else {
      console.log('普通用户权限');
    }
  } catch (error) {
    console.error('操作失败:', error.message);
  }
}

// 调用
main();
```

### Python 版本

```python
import hashlib
import requests
import platform

class AuthClient:
    def __init__(self, base_url=""):
        self.base_url = base_url
    
    def bind(self, auth_code, machine_code):
        """绑定机器码"""
        response = requests.post(f'{self.base_url}/api/auth/bind', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })
        return response.json()
    
    def validate(self, auth_code, machine_code):
        """验证授权码"""
        # 第一步：获取挑战
        challenge_resp = requests.post(f'{self.base_url}/api/auth/validate', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })
        
        challenge_data = challenge_resp.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])
        
        # 第二步：计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()
        
        # 第三步：提交验证
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
        """获取渠道列表"""
        # 第一步：获取挑战
        challenge_resp = requests.post(f'{self.base_url}/api/auth/channels', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })

        challenge_data = challenge_resp.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])

        # 第二步：计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()

        # 第三步：提交验证
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

    def get_machine_code(self):
        """获取简单机器码"""
        info = f"{platform.node()}-{platform.system()}-{platform.machine()}"
        return hashlib.md5(info.encode()).hexdigest()

# 使用示例
def main():
    client = AuthClient("http://localhost:3000")
    auth_code = "your_auth_code_here"
    machine_code = client.get_machine_code()
    
    try:
        # 首次使用：绑定机器码
        bind_result = client.bind(auth_code, machine_code)
        print("绑定成功")
        
        # 验证授权码
        user_info = client.validate(auth_code, machine_code)
        print(f"验证成功: {user_info}")

        # 获取渠道列表
        channels_info = client.get_channels(auth_code, machine_code)
        print(f"渠道信息: {channels_info}")
        print(f"分组: {channels_info['group']}, 渠道数: {channels_info['total']}")

        # 检查权限
        if user_info['user_type'] >= 10:
            print("管理员权限")
        else:
            print("普通用户权限")
            
    except Exception as e:
        print(f"操作失败: {e}")

if __name__ == "__main__":
    main()
```

## 常见错误处理

```javascript
// 错误处理示例
try {
  const result = await client.validate(authCode, machineCode);
} catch (error) {
  switch (error.message) {
    case '授权码不存在':
      alert('授权码无效，请检查输入');
      break;
    case '授权码无效或已过期':
      alert('授权码已过期，请联系管理员');
      break;
    case '机器码不匹配':
      alert('此授权码已在其他设备使用');
      break;
    case '该机器码已被其他授权码绑定':
      alert('当前设备已绑定其他授权码');
      break;
    default:
      alert('验证失败，请重试');
  }
}
```

## 状态说明

| 状态 | 说明 |
|------|------|
| 1 - 启用 | 可正常使用 |
| 2 - 禁用 | 被管理员禁用 |
| 3 - 已使用 | 已被使用（一次性） |
| 4 - 待激活 | 等待机器码绑定 |
| 5 - 激活 | 已绑定并激活 |

## 用户类型

| 类型 | 说明 |
|------|------|
| 1 | 普通用户 |
| 10 | 管理员 |
| 100 | 超级管理员 |

## 注意事项

1. **首次使用**：需要先调用绑定接口
2. **机器码唯一**：每个机器码只能绑定一个授权码
3. **HTTPS**：生产环境建议使用HTTPS
4. **错误重试**：网络错误可以重试，业务错误不建议重试
5. **缓存结果**：验证成功后可以缓存结果5-10分钟

## 完整文档

详细的接口文档请参考：[授权码外部接口使用文档](./auth_code_api.md)
