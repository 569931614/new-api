# 🚀 Git 提交指南 - Coze JWT 渠道

## 📋 需要提交的文件

### 1. 核心实现文件
```bash
# 添加 Coze JWT 渠道实现
git add relay/channel/coze_jwt/

# 添加常量和配置更新
git add common/constants.go
git add relay/constant/api_type.go
git add relay/relay_adaptor.go
git add controller/channel-test.go

# 添加前端配置
git add web/src/constants/channel.constants.js

# 添加依赖更新
git add go.mod
git add go.sum

# 添加文档和示例
git add docs/coze_jwt_channel.md
git add coze_exmplate/
```

### 2. 一键添加所有文件
```bash
git add .
```

## 📝 提交信息

```bash
git commit -m "feat: 添加Coze JWT渠道支持

- 新增Coze JWT渠道类型(50)，支持OAuth JWT认证
- 基于官方coze-go SDK实现，提供更可靠的认证机制
- 支持智能体调用、同步/异步工作流执行
- 添加Token缓存机制，自动处理过期和刷新
- 支持流式和非流式响应
- 包含完整的文档和示例代码

主要文件:
- relay/channel/coze_jwt/: Coze JWT渠道实现
- docs/coze_jwt_channel.md: 使用文档
- coze_exmplate/: 官方示例代码
- 前后端常量和配置更新"
```

## 🌐 推送到远程仓库

```bash
git push origin main
```

## 🔍 检查状态

```bash
# 查看当前状态
git status

# 查看提交历史
git log --oneline -5
```

## 📊 完整的提交流程

```bash
# 1. 检查当前状态
git status

# 2. 添加所有文件
git add .

# 3. 提交更改
git commit -m "feat: 添加Coze JWT渠道支持

- 新增Coze JWT渠道类型(50)，支持OAuth JWT认证
- 基于官方coze-go SDK实现，提供更可靠的认证机制
- 支持智能体调用、同步/异步工作流执行
- 添加Token缓存机制，自动处理过期和刷新
- 支持流式和非流式响应
- 包含完整的文档和示例代码

主要文件:
- relay/channel/coze_jwt/: Coze JWT渠道实现
- docs/coze_jwt_channel.md: 使用文档
- coze_exmplate/: 官方示例代码
- 前后端常量和配置更新"

# 4. 推送到远程仓库
git push origin main

# 5. 验证推送成功
git status
```

## ✅ 提交内容总结

### 🆕 新增文件
- `relay/channel/coze_jwt/` - Coze JWT渠道完整实现
- `docs/coze_jwt_channel.md` - 使用文档
- `coze_exmplate/` - 官方示例代码

### 🔧 修改文件
- `common/constants.go` - 添加渠道类型常量
- `relay/constant/api_type.go` - 添加API类型映射
- `relay/relay_adaptor.go` - 注册新渠道适配器
- `controller/channel-test.go` - 添加渠道测试支持
- `web/src/constants/channel.constants.js` - 前端渠道选项
- `go.mod` & `go.sum` - 添加coze-go SDK依赖

### 🎯 功能特性
- ✅ OAuth JWT 认证
- ✅ 智能体调用
- ✅ 同步/异步工作流
- ✅ Token 缓存机制
- ✅ 流式响应支持
- ✅ 完整文档和示例

## 🎉 提交完成后

提交成功后，您的GitHub仓库将包含完整的Coze JWT渠道实现，其他开发者可以：

1. 克隆仓库获取最新代码
2. 参考文档配置Coze JWT渠道
3. 使用示例代码进行测试
4. 基于实现进行二次开发

**现在请在终端中执行上述命令来完成Git提交！** 🚀
