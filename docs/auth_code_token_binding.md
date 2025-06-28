# 授权码API密钥绑定功能

## 功能概述

本功能为授权码系统添加了API密钥绑定能力，允许管理员将授权码与特定的API密钥进行关联，实现更精细的权限控制和资源管理。

## 功能特性

### 1. 数据模型扩展
- 在 `AuthCode` 结构体中添加了 `TokenId` 字段
- 支持授权码与API密钥的一对一绑定关系
- 兼容现有数据，未绑定的授权码 `token_id` 为 0

### 2. 后端API扩展

#### 新增接口
- `GET /api/auth_code/available_tokens` - 获取可用的API密钥列表

#### 修改接口
- `POST /api/auth_code/` - 创建授权码时支持绑定API密钥
- `PUT /api/auth_code/` - 更新授权码时支持修改绑定的API密钥
- `POST /api/auth_code/batch` - 批量创建授权码时支持绑定API密钥

### 3. 前端界面增强

#### 授权码编辑表单
- 添加API密钥选择下拉框
- 显示可用的API密钥列表（脱敏显示）
- 支持选择"不绑定API密钥"选项

#### 批量创建表单
- 添加API密钥选择功能
- 批量创建的所有授权码将绑定到同一个API密钥

#### 授权码列表表格
- 新增"绑定密钥"列，显示绑定状态
- 在详情弹窗中显示绑定的API密钥信息

## 使用说明

### 1. 创建带API密钥绑定的授权码

1. 进入授权码管理页面
2. 点击"添加授权码"按钮
3. 填写基本信息（授权码、名称、描述等）
4. 在"绑定API密钥"下拉框中选择要绑定的API密钥
5. 点击"提交"完成创建

### 2. 批量创建带API密钥绑定的授权码

1. 进入授权码管理页面
2. 点击"批量生成"按钮
3. 设置生成数量和基本配置
4. 在"绑定API密钥"下拉框中选择要绑定的API密钥
5. 点击"生成"完成批量创建

### 3. 修改授权码的API密钥绑定

1. 在授权码列表中找到目标授权码
2. 点击"编辑"按钮
3. 在"绑定API密钥"下拉框中选择新的API密钥或选择"不绑定API密钥"
4. 点击"提交"保存修改

### 4. 查看绑定状态

- 在授权码列表的"绑定密钥"列中查看绑定状态
- 点击"查看"按钮在弹窗中查看详细的绑定信息

## 技术实现

### 数据库变更
```sql
-- 为 auth_codes 表添加 token_id 字段
ALTER TABLE auth_codes ADD COLUMN token_id INTEGER DEFAULT 0;
CREATE INDEX idx_auth_codes_token_id ON auth_codes(token_id);
```

### API接口

#### 获取可用API密钥列表
```http
GET /api/auth_code/available_tokens
Authorization: Bearer <admin_token>

Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "测试密钥",
      "key": "sk-1234****5678"
    }
  ]
}
```

#### 创建授权码（带API密钥绑定）
```http
POST /api/auth_code/
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "code": "AUTH123456",
  "name": "测试授权码",
  "description": "用于测试的授权码",
  "user_type": 1,
  "token_id": 1,
  "expired_time": -1
}
```

## 安全考虑

1. **权限控制**: 只有管理员可以访问API密钥绑定功能
2. **数据脱敏**: 前端显示的API密钥进行了脱敏处理，只显示前8位和后4位
3. **状态验证**: 只有状态为"启用"的API密钥才会在选择列表中显示
4. **关联验证**: 绑定的API密钥必须属于当前用户

## 兼容性

- 向后兼容：现有的授权码不受影响，`token_id` 默认为 0 表示未绑定
- 数据库迁移：使用GORM的AutoMigrate功能自动添加新字段
- API兼容：现有API接口保持兼容，新增字段为可选参数

## 注意事项

1. API密钥绑定是可选功能，不绑定不影响授权码的正常使用
2. 删除API密钥时需要考虑已绑定的授权码的处理
3. 建议定期检查绑定关系的有效性
4. 批量操作时，所有生成的授权码将绑定到同一个API密钥

## 后续扩展

1. 支持授权码与API密钥的多对多关系
2. 添加绑定关系的使用统计
3. 支持基于绑定关系的权限控制
4. 添加绑定关系的审计日志
