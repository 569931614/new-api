# 根据授权码获取API密钥 - 使用示例

## 概述

本文档提供了如何使用新增的 `GET /api/auth/api_key` 接口根据授权码获取绑定的API密钥的详细示例。

## 前置条件

1. 授权码必须是**激活状态**（status = 5）
2. 授权码必须已经**绑定了API密钥**（token_id > 0）
3. 绑定的API密钥必须是**启用状态**（status = 1）

## 快速开始

### 1. 基本调用

```bash
curl -X GET "http://your-domain/api/auth/api_key?auth_code=9WQrAHZcsvOwydLj"
```

### 2. 成功响应示例

```json
{
  "success": true,
  "message": "获取API密钥成功",
  "data": {
    "token_id": 1,
    "token_name": "我的API密钥",
    "api_key": "sk-1234567890abcdef1234567890abcdef",
    "status": 1,
    "expired_time": 1703123456,
    "remain_quota": 1000000,
    "unlimited_quota": false,
    "group": "default",
    "auth_code_info": {
      "code": "9WQrAHZcsvOwydLj",
      "name": "测试授权码",
      "user_type": 1,
      "is_bot": false,
      "group": "default"
    }
  }
}
```

## 编程语言示例

### JavaScript/Node.js

```javascript
async function getApiKeyByAuthCode(authCode) {
  try {
    const response = await fetch(`/api/auth/api_key?auth_code=${encodeURIComponent(authCode)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    const result = await response.json();
    
    if (result.success) {
      console.log('API密钥:', result.data.api_key);
      console.log('剩余配额:', result.data.remain_quota);
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('获取API密钥失败:', error.message);
    throw error;
  }
}

// 使用示例
getApiKeyByAuthCode('9WQrAHZcsvOwydLj')
  .then(data => {
    console.log('获取成功:', data);
    // 使用API密钥进行后续操作
    useApiKey(data.api_key);
  })
  .catch(error => {
    console.error('操作失败:', error);
  });
```

### Python

```python
import requests

def get_api_key_by_auth_code(auth_code, base_url="http://your-domain"):
    """根据授权码获取API密钥"""
    try:
        response = requests.get(f'{base_url}/api/auth/api_key', params={
            'auth_code': auth_code
        })
        
        result = response.json()
        
        if result['success']:
            print(f"API密钥: {result['data']['api_key']}")
            print(f"剩余配额: {result['data']['remain_quota']}")
            return result['data']
        else:
            raise Exception(result['message'])
            
    except Exception as e:
        print(f"获取API密钥失败: {e}")
        raise

# 使用示例
try:
    api_key_data = get_api_key_by_auth_code('9WQrAHZcsvOwydLj')
    print(f"获取成功: {api_key_data}")
    
    # 使用API密钥进行后续操作
    api_key = api_key_data['api_key']
    # use_api_key(api_key)
    
except Exception as e:
    print(f"操作失败: {e}")
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

type ApiKeyResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    struct {
        TokenId        int    `json:"token_id"`
        TokenName      string `json:"token_name"`
        ApiKey         string `json:"api_key"`
        Status         int    `json:"status"`
        ExpiredTime    int64  `json:"expired_time"`
        RemainQuota    int64  `json:"remain_quota"`
        UnlimitedQuota bool   `json:"unlimited_quota"`
        Group          string `json:"group"`
        AuthCodeInfo   struct {
            Code     string `json:"code"`
            Name     string `json:"name"`
            UserType int    `json:"user_type"`
            IsBot    bool   `json:"is_bot"`
            Group    string `json:"group"`
        } `json:"auth_code_info"`
    } `json:"data"`
}

func GetApiKeyByAuthCode(authCode, baseURL string) (*ApiKeyResponse, error) {
    // 构建URL
    u, err := url.Parse(baseURL + "/api/auth/api_key")
    if err != nil {
        return nil, err
    }
    
    q := u.Query()
    q.Set("auth_code", authCode)
    u.RawQuery = q.Encode()
    
    // 发送请求
    resp, err := http.Get(u.String())
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // 解析JSON
    var result ApiKeyResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    
    if !result.Success {
        return nil, fmt.Errorf(result.Message)
    }
    
    return &result, nil
}

func main() {
    authCode := "9WQrAHZcsvOwydLj"
    baseURL := "http://your-domain"
    
    result, err := GetApiKeyByAuthCode(authCode, baseURL)
    if err != nil {
        fmt.Printf("获取API密钥失败: %v\n", err)
        return
    }
    
    fmt.Printf("API密钥: %s\n", result.Data.ApiKey)
    fmt.Printf("剩余配额: %d\n", result.Data.RemainQuota)
    fmt.Printf("授权码信息: %+v\n", result.Data.AuthCodeInfo)
    
    // 使用API密钥进行后续操作
    // useApiKey(result.Data.ApiKey)
}
```

## 错误处理

### 常见错误及处理方法

```javascript
const errorHandlers = {
  '授权码参数不能为空': () => {
    console.log('请提供有效的授权码参数');
  },
  
  '授权码不存在': () => {
    console.log('授权码无效，请检查输入');
  },
  
  '授权码无效或已过期': () => {
    console.log('授权码已过期，请联系管理员');
  },
  
  '授权码未激活': () => {
    console.log('请先激活授权码（绑定机器码）');
  },
  
  '授权码未绑定API密钥': () => {
    console.log('授权码未绑定API密钥，请联系管理员');
  },
  
  '绑定的API密钥不存在或已被禁用': () => {
    console.log('关联的API密钥已失效，请联系管理员');
  }
};

async function getApiKeyWithErrorHandling(authCode) {
  try {
    return await getApiKeyByAuthCode(authCode);
  } catch (error) {
    const handler = errorHandlers[error.message];
    if (handler) {
      handler();
    } else {
      console.error('未知错误:', error.message);
    }
    throw error;
  }
}
```

## 实际应用场景

### 1. 自动化脚本

```python
#!/usr/bin/env python3
"""
自动化脚本：根据授权码获取API密钥并执行任务
"""

import requests
import sys

def main():
    if len(sys.argv) != 2:
        print("用法: python script.py <授权码>")
        sys.exit(1)
    
    auth_code = sys.argv[1]
    
    try:
        # 获取API密钥
        api_key_data = get_api_key_by_auth_code(auth_code)
        api_key = api_key_data['api_key']
        
        print(f"成功获取API密钥，剩余配额: {api_key_data['remain_quota']}")
        
        # 使用API密钥执行任务
        execute_tasks_with_api_key(api_key)
        
    except Exception as e:
        print(f"脚本执行失败: {e}")
        sys.exit(1)

def execute_tasks_with_api_key(api_key):
    """使用API密钥执行具体任务"""
    # 这里实现具体的业务逻辑
    print(f"使用API密钥 {api_key[:8]}**** 执行任务...")

if __name__ == "__main__":
    main()
```

### 2. 配置管理

```javascript
class ConfigManager {
  constructor() {
    this.apiKey = null;
    this.authCode = null;
  }
  
  async initialize(authCode) {
    this.authCode = authCode;
    
    try {
      const apiKeyData = await getApiKeyByAuthCode(authCode);
      this.apiKey = apiKeyData.api_key;
      
      console.log('配置初始化成功');
      return true;
    } catch (error) {
      console.error('配置初始化失败:', error.message);
      return false;
    }
  }
  
  getApiKey() {
    if (!this.apiKey) {
      throw new Error('API密钥未初始化，请先调用 initialize()');
    }
    return this.apiKey;
  }
  
  async refreshApiKey() {
    if (!this.authCode) {
      throw new Error('授权码未设置');
    }
    
    const apiKeyData = await getApiKeyByAuthCode(this.authCode);
    this.apiKey = apiKeyData.api_key;
    
    console.log('API密钥已刷新');
    return this.apiKey;
  }
}

// 使用示例
const config = new ConfigManager();

async function setupApplication() {
  const authCode = process.env.AUTH_CODE || '9WQrAHZcsvOwydLj';
  
  if (await config.initialize(authCode)) {
    console.log('应用程序配置完成');
    // 开始应用程序逻辑
    startApplication();
  } else {
    console.error('应用程序初始化失败');
    process.exit(1);
  }
}
```

## 安全建议

1. **HTTPS通信**: 生产环境必须使用HTTPS协议
2. **密钥保护**: 获取到的API密钥应安全存储，避免明文保存
3. **访问控制**: 限制对此接口的访问频率
4. **日志记录**: 记录API密钥获取操作的审计日志
5. **定期检查**: 定期验证API密钥的有效性

## 注意事项

- 此接口不需要身份验证，但需要有效的激活授权码
- API密钥获取后建议缓存使用，避免频繁请求
- 如果API密钥配额不足，请联系管理员充值
- 授权码状态变更可能影响API密钥的获取

---

*文档版本: v1.0*  
*最后更新: 2025年6月*
