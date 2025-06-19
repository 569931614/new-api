# 授权码外部接口使用文档

## 目录
- [概述](#概述)
- [接口列表](#接口列表)
- [机器码绑定接口](#机器码绑定接口)
- [授权码验证接口](#授权码验证接口)
- [安全机制说明](#安全机制说明)
- [客户端实现示例](#客户端实现示例)
- [机器码获取方法](#机器码获取方法)
- [错误处理](#错误处理)
- [最佳实践](#最佳实践)

---

## 概述

本文档详细说明了授权码系统的两个外部接口的使用方法。这些接口允许第三方应用程序安全地绑定和验证授权码，采用了挑战-响应机制确保安全性。

### 接口特点
- **无需认证**：接口可直接调用，无需登录或API密钥
- **安全可靠**：采用HMAC-SHA256挑战-响应机制防止篡改
- **机器绑定**：支持授权码与特定设备绑定
- **状态管理**：完整的授权码生命周期管理

---

## 接口列表

| 接口名称 | 方法 | 地址 | 用途 |
|---------|------|------|------|
| 机器码绑定 | POST | `/api/auth/bind` | 将授权码与机器码绑定 |
| 授权码验证 | POST | `/api/auth/validate` | 验证授权码有效性 |
| 获取渠道列表 | POST | `/api/auth/channels` | 根据授权码获取可用渠道列表 |

---

## 机器码绑定接口

### 基本信息
- **接口地址：** `POST /api/auth/bind`
- **Content-Type：** `application/json`
- **用途：** 将授权码与特定的机器码绑定，绑定后授权码状态变为"激活"
- **使用场景：** 用户首次在设备上使用授权码时

### 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| machine_code | string | 是 | 机器码（设备唯一标识） |

### 请求示例

#### cURL
```bash
curl -X POST http://your-domain/api/auth/bind \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "9WQrAHZcsvOwydLj",
    "machine_code": "DESKTOP-ABC123-12345"
  }'
```

#### JavaScript
```javascript
async function bindMachineCode(authCode, machineCode) {
  const response = await fetch('/api/auth/bind', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      auth_code: authCode,
      machine_code: machineCode
    })
  });

  const result = await response.json();
  return result;
}

// 使用示例
try {
  const result = await bindMachineCode('9WQrAHZcsvOwydLj', 'DESKTOP-ABC123-12345');
  if (result.success) {
    console.log('绑定成功');
  } else {
    console.error('绑定失败:', result.message);
  }
} catch (error) {
  console.error('请求失败:', error);
}
```

#### Python
```python
import requests

def bind_machine_code(auth_code, machine_code, base_url="http://your-domain"):
    """绑定机器码"""
    response = requests.post(f'{base_url}/api/auth/bind', json={
        'auth_code': auth_code,
        'machine_code': machine_code
    })

    return response.json()

# 使用示例
try:
    result = bind_machine_code('9WQrAHZcsvOwydLj', 'DESKTOP-ABC123-12345')
    if result['success']:
        print('绑定成功')
    else:
        print(f'绑定失败: {result["message"]}')
except Exception as e:
    print(f'请求失败: {e}')
```

### 响应格式

#### 成功响应
```json
{
  "success": true,
  "message": "机器码绑定成功"
}
```

#### 失败响应
```json
{
  "success": false,
  "message": "错误描述"
}
```

### 常见错误

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| 授权码不存在 | 提供的授权码无效 | 检查授权码是否正确 |
| 授权码状态不允许绑定机器码 | 授权码已被使用或禁用 | 联系管理员检查授权码状态 |
| 授权码已过期 | 授权码超过有效期 | 联系管理员更新授权码 |
| 该机器码已被其他授权码绑定 | 机器码已被占用 | 使用不同的机器码或联系管理员 |

---

## 授权码验证接口

### 基本信息
- **接口地址：** `POST /api/auth/validate`
- **Content-Type：** `application/json`
- **用途：** 验证授权码是否有效，采用挑战-响应机制确保安全
- **使用场景：** 每次需要验证用户权限时

### 验证流程

授权码验证采用两步验证流程：

1. **第一步：获取挑战值**
   - 发送授权码和机器码
   - 服务器返回挑战值和时间戳

2. **第二步：提交验证响应**
   - 计算挑战值的SHA256哈希
   - 提交挑战值和响应进行验证

### 第一步：获取挑战值

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| machine_code | string | 是 | 机器码 |

#### 请求示例
```bash
curl -X POST http://your-domain/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "9WQrAHZcsvOwydLj",
    "machine_code": "DESKTOP-ABC123-12345"
  }'
```

#### 响应示例
```json
{
  "success": true,
  "message": "请完成验证挑战",
  "challenge": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
  "timestamp": 1703123456,
  "expires_in": 300
}
```

### 第二步：提交验证响应

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| machine_code | string | 是 | 机器码 |
| challenge | string | 是 | 从第一步获取的挑战值 |
| response | string | 是 | 挑战值的SHA256哈希 |

#### 响应计算方法
```javascript
// JavaScript
const response = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(challenge))
  .then(buffer => Array.from(new Uint8Array(buffer))
    .map(b => b.toString(16).padStart(2, '0')).join(''));
```

```python
# Python
import hashlib
response = hashlib.sha256(challenge.encode()).hexdigest()
```

#### 成功响应
```json
{
  "success": true,
  "message": "授权码验证成功",
  "data": {
    "user_type": 1,
    "is_bot": false,
    "wx_auto_x_code": "wxautox_code_here",
    "expired_time": 1703123456
  }
}
```

#### 返回数据说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| user_type | int | 用户类型（1:普通用户, 10:管理员, 100:超级管理员） |
| is_bot | bool | 是否为机器人账户 |
| wx_auto_x_code | string | WxAutoX码 |
| expired_time | int | 过期时间戳（-1表示永不过期） |

---

## 获取渠道列表接口

### 基本信息
- **接口地址：** `POST /api/auth/channels`
- **Content-Type：** `application/json`
- **用途：** 根据授权码获取关联分组下的可用渠道列表
- **使用场景：** 需要获取特定授权码可访问的渠道信息时

### 验证流程

与授权码验证接口相同，采用两步验证流程：

1. **第一步：获取挑战值**
   - 发送授权码和机器码
   - 服务器返回挑战值和时间戳

2. **第二步：提交验证响应**
   - 计算挑战值的SHA256哈希
   - 提交挑战值和响应获取渠道列表

### 第一步：获取挑战值

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| machine_code | string | 是 | 机器码 |

#### 请求示例
```bash
curl -X POST http://your-domain/api/auth/channels \
  -H "Content-Type: application/json" \
  -d '{
    "auth_code": "9WQrAHZcsvOwydLj",
    "machine_code": "DESKTOP-ABC123-12345"
  }'
```

#### 响应示例
```json
{
  "success": true,
  "message": "请完成验证挑战",
  "challenge": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
  "timestamp": 1703123456,
  "expires_in": 300
}
```

### 第二步：提交验证响应

#### 请求参数
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| auth_code | string | 是 | 授权码 |
| machine_code | string | 是 | 机器码 |
| challenge | string | 是 | 从第一步获取的挑战值 |
| response | string | 是 | 挑战值的SHA256哈希 |

#### 成功响应
```json
{
  "success": true,
  "message": "获取渠道列表成功",
  "data": {
    "group": "default",
    "channels": [
      {
        "id": 1,
        "name": "OpenAI官方",
        "type": 1,
        "status": 1,
        "models": ["gpt-3.5-turbo", "gpt-4"],
        "group": "default",
        "priority": 0,
        "weight": 100
      },
      {
        "id": 2,
        "name": "Claude API",
        "type": 15,
        "status": 1,
        "models": ["claude-3-sonnet", "claude-3-haiku"],
        "group": "default",
        "priority": 0,
        "weight": 80
      }
    ],
    "total": 2
  }
}
```

#### 返回数据说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| group | string | 授权码关联的分组名称 |
| channels | array | 渠道列表 |
| total | int | 渠道总数 |

#### 渠道信息字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int | 渠道ID |
| name | string | 渠道名称 |
| type | int | 渠道类型 |
| status | int | 渠道状态（1:启用） |
| models | array | 支持的模型列表 |
| group | string | 渠道分组 |
| priority | int | 优先级 |
| weight | int | 权重 |

### 特殊情况

#### 无分组授权码
如果授权码没有设置分组，将返回空的渠道列表：

```json
{
  "success": true,
  "message": "获取渠道列表成功",
  "data": {
    "group": "",
    "channels": [],
    "total": 0
  }
}
```

#### 分组无渠道
如果分组下没有可用渠道，也会返回空列表：

```json
{
  "success": true,
  "message": "获取渠道列表成功",
  "data": {
    "group": "test_group",
    "channels": [],
    "total": 0
  }
}
```

---

## 安全机制说明

### 挑战-响应验证原理

1. **挑战生成**
   ```
   challenge = HMAC-SHA256(auth_code:machine_code:timestamp, server_secret)
   ```

2. **响应计算**
   ```
   response = SHA256(challenge)
   ```

3. **验证过程**
   - 服务器重新生成挑战值进行比对
   - 验证响应的正确性
   - 允许5分钟的时间窗口

### 安全特性

- **防篡改**：挑战值使用服务器密钥签名，无法伪造
- **防重放**：挑战值包含时间戳，有5分钟有效期
- **设备绑定**：授权码与特定机器码绑定
- **双重验证**：挑战值完整性 + 响应正确性

---

## 客户端实现示例

### JavaScript 完整实现

```javascript
class AuthCodeClient {
  constructor(baseUrl = '') {
    this.baseUrl = baseUrl;
  }

  /**
   * 绑定机器码
   */
  async bindMachineCode(authCode, machineCode) {
    const response = await fetch(`${this.baseUrl}/api/auth/bind`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        auth_code: authCode,
        machine_code: machineCode
      })
    });

    const result = await response.json();
    if (!result.success) {
      throw new Error(result.message);
    }

    return result;
  }

  /**
   * 验证授权码
   */
  async validateAuthCode(authCode, machineCode) {
    try {
      // 第一步：获取挑战
      const challengeResponse = await fetch(`${this.baseUrl}/api/auth/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          auth_code: authCode,
          machine_code: machineCode
        })
      });

      const challengeData = await challengeResponse.json();
      if (!challengeData.success) {
        throw new Error(challengeData.message);
      }

      // 第二步：计算响应
      const challenge = challengeData.challenge;
      const encoder = new TextEncoder();
      const data = encoder.encode(challenge);
      const hashBuffer = await crypto.subtle.digest('SHA-256', data);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

      // 第三步：提交验证
      const validationResponse = await fetch(`${this.baseUrl}/api/auth/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          auth_code: authCode,
          machine_code: machineCode,
          challenge: challenge,
          response: response
        })
      });

      const validationData = await validationResponse.json();
      if (!validationData.success) {
        throw new Error(validationData.message);
      }

      return validationData.data;

    } catch (error) {
      throw new Error(`验证失败: ${error.message}`);
    }
  }

  /**
   * 获取渠道列表
   */
  async getChannels(authCode, machineCode) {
    try {
      // 第一步：获取挑战
      const challengeResponse = await fetch(`${this.baseUrl}/api/auth/channels`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          auth_code: authCode,
          machine_code: machineCode
        })
      });

      const challengeData = await challengeResponse.json();
      if (!challengeData.success) {
        throw new Error(challengeData.message);
      }

      // 第二步：计算响应
      const challenge = challengeData.challenge;
      const encoder = new TextEncoder();
      const data = encoder.encode(challenge);
      const hashBuffer = await crypto.subtle.digest('SHA-256', data);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      const response = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

      // 第三步：提交验证
      const channelsResponse = await fetch(`${this.baseUrl}/api/auth/channels`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          auth_code: authCode,
          machine_code: machineCode,
          challenge: challenge,
          response: response
        })
      });

      const channelsData = await channelsResponse.json();
      if (!channelsData.success) {
        throw new Error(channelsData.message);
      }

      return channelsData.data;

    } catch (error) {
      throw new Error(`获取渠道列表失败: ${error.message}`);
    }
  }

  /**
   * 获取机器码（示例实现）
   */
  getMachineCode() {
    // 这里可以根据实际需求实现机器码获取逻辑
    const userAgent = navigator.userAgent;
    const platform = navigator.platform;
    const language = navigator.language;

    const machineInfo = `${userAgent}-${platform}-${language}`;

    // 简单哈希（生产环境建议使用更复杂的算法）
    let hash = 0;
    for (let i = 0; i < machineInfo.length; i++) {
      const char = machineInfo.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // 转换为32位整数
    }

    return Math.abs(hash).toString(16);
  }
}

// 使用示例
const authClient = new AuthCodeClient('http://localhost:3000');

async function example() {
  try {
    const authCode = '9WQrAHZcsvOwydLj';
    const machineCode = authClient.getMachineCode();

    // 首次使用：绑定机器码
    await authClient.bindMachineCode(authCode, machineCode);
    console.log('机器码绑定成功');

    // 验证授权码
    const userInfo = await authClient.validateAuthCode(authCode, machineCode);
    console.log('验证成功:', userInfo);

    // 获取渠道列表
    const channelsInfo = await authClient.getChannels(authCode, machineCode);
    console.log('渠道信息:', channelsInfo);
    console.log(`分组: ${channelsInfo.group}`);
    console.log(`可用渠道数: ${channelsInfo.total}`);

    // 根据用户信息决定权限
    if (userInfo.user_type >= 10) {
      console.log('管理员权限');
    } else {
      console.log('普通用户权限');
    }

  } catch (error) {
    console.error('操作失败:', error.message);
  }
}
```

### Python 完整实现

```python
import hashlib
import requests
import json
import platform
import uuid

class AuthCodeClient:
    def __init__(self, base_url=""):
        self.base_url = base_url

    def bind_machine_code(self, auth_code, machine_code):
        """绑定机器码"""
        response = requests.post(f'{self.base_url}/api/auth/bind', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })

        result = response.json()
        if not result['success']:
            raise Exception(result['message'])

        return result

    def validate_auth_code(self, auth_code, machine_code):
        """验证授权码"""
        # 第一步：获取挑战
        challenge_response = requests.post(f'{self.base_url}/api/auth/validate', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })

        challenge_data = challenge_response.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])

        # 第二步：计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()

        # 第三步：提交验证
        validation_response = requests.post(f'{self.base_url}/api/auth/validate', json={
            'auth_code': auth_code,
            'machine_code': machine_code,
            'challenge': challenge,
            'response': response
        })

        validation_data = validation_response.json()
        if not validation_data['success']:
            raise Exception(validation_data['message'])

        return validation_data['data']

    def get_channels(self, auth_code, machine_code):
        """获取渠道列表"""
        # 第一步：获取挑战
        challenge_response = requests.post(f'{self.base_url}/api/auth/channels', json={
            'auth_code': auth_code,
            'machine_code': machine_code
        })

        challenge_data = challenge_response.json()
        if not challenge_data['success']:
            raise Exception(challenge_data['message'])

        # 第二步：计算响应
        challenge = challenge_data['challenge']
        response = hashlib.sha256(challenge.encode()).hexdigest()

        # 第三步：提交验证
        channels_response = requests.post(f'{self.base_url}/api/auth/channels', json={
            'auth_code': auth_code,
            'machine_code': machine_code,
            'challenge': challenge,
            'response': response
        })

        channels_data = channels_response.json()
        if not channels_data['success']:
            raise Exception(channels_data['message'])

        return channels_data['data']

    def get_machine_code(self):
        """获取机器码"""
        # 获取系统信息
        hostname = platform.node()
        system = platform.system()
        machine = platform.machine()
        processor = platform.processor()

        # 组合机器信息
        machine_info = f"{hostname}-{system}-{machine}-{processor}"

        # 生成哈希
        return hashlib.md5(machine_info.encode()).hexdigest()

# 使用示例
def main():
    client = AuthCodeClient("http://localhost:3000")

    try:
        auth_code = "9WQrAHZcsvOwydLj"
        machine_code = client.get_machine_code()

        print(f"机器码: {machine_code}")

        # 首次使用：绑定机器码
        bind_result = client.bind_machine_code(auth_code, machine_code)
        print("机器码绑定成功")

        # 验证授权码
        user_info = client.validate_auth_code(auth_code, machine_code)
        print(f"验证成功: {user_info}")

        # 获取渠道列表
        channels_info = client.get_channels(auth_code, machine_code)
        print(f"渠道信息: {channels_info}")
        print(f"分组: {channels_info['group']}")
        print(f"可用渠道数: {channels_info['total']}")

        # 根据用户信息决定权限
        if user_info['user_type'] >= 10:
            print("管理员权限")
        else:
            print("普通用户权限")

        if user_info['is_bot']:
            print("机器人账户")

        if user_info['expired_time'] != -1:
            import datetime
            expire_date = datetime.datetime.fromtimestamp(user_info['expired_time'])
            print(f"过期时间: {expire_date}")
        else:
            print("永不过期")

    except Exception as e:
        print(f"操作失败: {e}")

if __name__ == "__main__":
    main()
```

---

## 机器码获取方法

### Windows 环境

#### PowerShell
```powershell
# 获取系统信息
$computerName = $env:COMPUTERNAME
$processor = (Get-WmiObject -Class Win32_Processor).Name
$motherboard = (Get-WmiObject -Class Win32_BaseBoard).SerialNumber
$machineCode = "$computerName-$processor-$motherboard"

# 生成哈希
$hash = [System.Security.Cryptography.MD5]::Create()
$bytes = [System.Text.Encoding]::UTF8.GetBytes($machineCode)
$hashBytes = $hash.ComputeHash($bytes)
$hashString = [System.BitConverter]::ToString($hashBytes) -replace '-'
Write-Output $hashString.ToLower()
```

#### C#
```csharp
using System;
using System.Management;
using System.Security.Cryptography;
using System.Text;

public static string GetMachineCode()
{
    string computerName = Environment.MachineName;
    string processor = "";
    string motherboard = "";

    // 获取处理器信息
    using (ManagementObjectSearcher searcher = new ManagementObjectSearcher("SELECT * FROM Win32_Processor"))
    {
        foreach (ManagementObject obj in searcher.Get())
        {
            processor = obj["Name"].ToString();
            break;
        }
    }

    // 获取主板信息
    using (ManagementObjectSearcher searcher = new ManagementObjectSearcher("SELECT * FROM Win32_BaseBoard"))
    {
        foreach (ManagementObject obj in searcher.Get())
        {
            motherboard = obj["SerialNumber"].ToString();
            break;
        }
    }

    string machineInfo = $"{computerName}-{processor}-{motherboard}";

    // 生成MD5哈希
    using (MD5 md5 = MD5.Create())
    {
        byte[] bytes = Encoding.UTF8.GetBytes(machineInfo);
        byte[] hash = md5.ComputeHash(bytes);
        return BitConverter.ToString(hash).Replace("-", "").ToLower();
    }
}
```

### Linux 环境

#### Bash
```bash
#!/bin/bash

# 获取系统信息
hostname=$(hostname)
cpu_info=$(cat /proc/cpuinfo | grep "model name" | head -1 | cut -d: -f2 | xargs)
machine_id=$(cat /etc/machine-id 2>/dev/null || cat /var/lib/dbus/machine-id 2>/dev/null || echo "unknown")

# 组合信息
machine_info="$hostname-$cpu_info-$machine_id"

# 生成MD5哈希
echo -n "$machine_info" | md5sum | cut -d' ' -f1
```

### macOS 环境

#### Swift
```swift
import Foundation
import IOKit

func getMachineCode() -> String {
    let host = Host.current()
    let hostname = host.name ?? "unknown"

    // 获取硬件UUID
    let platformExpert = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching("IOPlatformExpertDevice"))
    let serialNumberAsCFString = IORegistryEntryCreateCFProperty(platformExpert, kIOPlatformUUIDKey, kCFAllocatorDefault, 0)
    let serialNumber = serialNumberAsCFString?.takeUnretainedValue() as? String ?? "unknown"

    let machineInfo = "\(hostname)-\(serialNumber)"

    // 生成MD5哈希
    let data = machineInfo.data(using: .utf8)!
    let hash = Insecure.MD5.hash(data: data)
    return hash.map { String(format: "%02hhx", $0) }.joined()
}
```

---

## 错误处理

### 错误分类

#### 1. 网络错误
```javascript
try {
  const result = await authClient.validateAuthCode(authCode, machineCode);
} catch (error) {
  if (error.message.includes('fetch')) {
    console.error('网络连接失败，请检查网络设置');
  }
}
```

#### 2. 业务逻辑错误
```javascript
const errorHandlers = {
  '授权码不存在': () => {
    console.log('请检查授权码是否正确输入');
    // 提示用户重新输入授权码
  },

  '授权码无效或已过期': () => {
    console.log('授权码已过期，请联系管理员获取新的授权码');
    // 跳转到联系页面或显示联系信息
  },

  '机器码不匹配': () => {
    console.log('此授权码已绑定其他设备，无法在当前设备使用');
    // 提示用户联系管理员或使用其他授权码
  },

  '该机器码已被其他授权码绑定': () => {
    console.log('当前设备已绑定其他授权码');
    // 提示用户使用已绑定的授权码
  },

  '验证挑战失败': () => {
    console.log('验证失败，请重试');
    // 自动重试或提示用户重试
  }
};

try {
  const result = await authClient.validateAuthCode(authCode, machineCode);
} catch (error) {
  const handler = errorHandlers[error.message];
  if (handler) {
    handler();
  } else {
    console.error('未知错误:', error.message);
  }
}
```

### 重试机制

```javascript
class AuthCodeClient {
  async validateAuthCodeWithRetry(authCode, machineCode, maxRetries = 3) {
    for (let i = 0; i < maxRetries; i++) {
      try {
        return await this.validateAuthCode(authCode, machineCode);
      } catch (error) {
        if (i === maxRetries - 1) {
          throw error; // 最后一次重试失败，抛出错误
        }

        // 可重试的错误
        if (error.message.includes('验证挑战失败') ||
            error.message.includes('网络')) {
          console.log(`第${i + 1}次验证失败，正在重试...`);
          await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1))); // 递增延迟
          continue;
        }

        throw error; // 不可重试的错误，直接抛出
      }
    }
  }
}
```

---

## 最佳实践

### 1. 安全建议

- **HTTPS通信**：生产环境必须使用HTTPS协议
- **机器码保护**：机器码应该基于硬件特征生成，难以伪造
- **错误信息**：不要在客户端暴露过多的错误详情
- **日志记录**：记录关键操作的日志，便于审计

### 2. 性能优化

- **缓存机制**：验证成功后可以缓存结果一段时间
- **连接复用**：使用HTTP连接池减少连接开销
- **超时设置**：设置合理的请求超时时间

```javascript
class AuthCodeClient {
  constructor(baseUrl = '', options = {}) {
    this.baseUrl = baseUrl;
    this.cache = new Map();
    this.cacheTimeout = options.cacheTimeout || 300000; // 5分钟缓存
  }

  async validateAuthCodeCached(authCode, machineCode) {
    const cacheKey = `${authCode}-${machineCode}`;
    const cached = this.cache.get(cacheKey);

    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data;
    }

    const result = await this.validateAuthCode(authCode, machineCode);
    this.cache.set(cacheKey, {
      data: result,
      timestamp: Date.now()
    });

    return result;
  }
}
```

### 3. 用户体验

- **进度提示**：显示验证进度，特别是网络较慢时
- **友好错误**：提供用户友好的错误提示
- **自动重试**：对于临时性错误自动重试
- **离线处理**：考虑网络断开时的处理方案

```javascript
// 用户友好的验证函数
async function validateWithUI(authCode, machineCode) {
  // 显示加载状态
  showLoading('正在验证授权码...');

  try {
    const result = await authClient.validateAuthCodeWithRetry(authCode, machineCode);
    hideLoading();
    showSuccess('验证成功！');
    return result;
  } catch (error) {
    hideLoading();

    // 用户友好的错误提示
    const friendlyMessages = {
      '授权码不存在': '授权码无效，请检查输入是否正确',
      '授权码无效或已过期': '授权码已过期，请联系管理员',
      '机器码不匹配': '此授权码已在其他设备使用',
      '验证挑战失败': '验证失败，请重试'
    };

    const friendlyMessage = friendlyMessages[error.message] || '验证失败，请稍后重试';
    showError(friendlyMessage);

    throw error;
  }
}
```

### 4. 监控和日志

```javascript
class AuthCodeClient {
  constructor(baseUrl = '', options = {}) {
    this.baseUrl = baseUrl;
    this.logger = options.logger || console;
  }

  async validateAuthCode(authCode, machineCode) {
    const startTime = Date.now();

    try {
      this.logger.info('开始验证授权码', { authCode: authCode.substring(0, 4) + '****' });

      const result = await this._doValidate(authCode, machineCode);

      this.logger.info('授权码验证成功', {
        duration: Date.now() - startTime,
        userType: result.user_type
      });

      return result;
    } catch (error) {
      this.logger.error('授权码验证失败', {
        error: error.message,
        duration: Date.now() - startTime
      });

      throw error;
    }
  }
}
```

---

## 附录

### 状态码说明

| 状态码 | 名称 | 说明 |
|--------|------|------|
| 1 | 启用 | 授权码可正常使用 |
| 2 | 禁用 | 授权码被管理员禁用 |
| 3 | 已使用 | 授权码已被使用（一次性） |
| 4 | 待激活 | 授权码等待机器码绑定 |
| 5 | 激活 | 授权码已绑定机器码并激活 |

### 用户类型说明

| 类型值 | 名称 | 权限级别 |
|--------|------|----------|
| 1 | 普通用户 | 基础权限 |
| 10 | 管理员 | 管理权限 |
| 100 | 超级管理员 | 最高权限 |

### 常见问题

**Q: 机器码绑定后可以更改吗？**
A: 已激活的授权码不能更改机器码，需要联系管理员重置。

**Q: 验证失败后多久可以重试？**
A: 没有重试限制，但建议间隔1-2秒后重试。

**Q: 授权码可以在多个设备上使用吗？**
A: 不可以，每个授权码只能绑定一个机器码。

**Q: 挑战值的有效期是多久？**
A: 挑战值有效期为5分钟。

**Q: 如何处理网络不稳定的情况？**
A: 建议实现重试机制，对于网络错误自动重试2-3次。

**Q: 机器码如何生成？**
A: 建议基于硬件特征（CPU、主板序列号等）生成，确保唯一性和稳定性。

**Q: 是否支持批量验证？**
A: 当前版本不支持批量验证，每次只能验证一个授权码。

### 技术支持

如果在使用过程中遇到问题，请联系技术支持团队：

- **邮箱**：support@example.com
- **文档更新**：本文档会根据系统更新持续维护
- **版本兼容性**：本文档适用于 v1.0+ 版本

---

*最后更新时间：2024年12月*
