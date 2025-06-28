@echo off
chcp 65001 >nul
echo 开始Git提交操作...
echo.

echo 1. 检查Git状态...
git status
echo.

echo 2. 添加所有文件...
git add .
echo.

echo 3. 检查添加后的状态...
git status
echo.

echo 4. 提交更改...
git commit -m "feat: 添加Coze JWT渠道支持 - 新增Coze JWT渠道类型(50)，支持OAuth JWT认证 - 基于官方coze-go SDK实现，提供更可靠的认证机制 - 支持智能体调用、同步/异步工作流执行 - 添加Token缓存机制，自动处理过期和刷新 - 支持流式和非流式响应 - 包含完整的文档和示例代码"
echo.

echo 5. 推送到远程仓库...
git push origin main
echo.

echo Git提交操作完成！
echo.
pause
