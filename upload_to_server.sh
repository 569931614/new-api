#!/bin/bash

# 上传到服务器的部署脚本
# 使用方法: ./upload_to_server.sh user@server:/path/to/deploy

set -e

# 检查参数
if [ $# -eq 0 ]; then
    echo "使用方法: $0 user@server:/path/to/deploy"
    echo "示例: $0 root@192.168.1.100:/opt/new-api"
    exit 1
fi

SERVER_PATH=$1

echo "🚀 开始部署到服务器..."

# 1. 确保前端已构建
echo "📦 检查前端构建文件..."
if [ ! -d "web/dist" ]; then
    echo "❌ 前端文件未构建，正在构建..."
    cd web
    npm run build
    cd ..
    echo "✅ 前端构建完成"
else
    echo "✅ 前端文件已存在"
fi

# 2. 创建临时目录，排除不需要的文件
echo "📁 准备上传文件..."
TEMP_DIR=$(mktemp -d)
rsync -av --exclude='node_modules' \
          --exclude='.git' \
          --exclude='*.log' \
          --exclude='logs/*' \
          --exclude='data/*' \
          --exclude='.env' \
          --exclude='*.exe' \
          . $TEMP_DIR/

# 3. 上传到服务器
echo "⬆️  上传文件到服务器..."
rsync -avz --delete $TEMP_DIR/ $SERVER_PATH/

# 4. 清理临时目录
rm -rf $TEMP_DIR

echo "✅ 文件上传完成！"

# 5. 在服务器上执行部署命令
SERVER_USER=$(echo $SERVER_PATH | cut -d: -f1)
SERVER_DIR=$(echo $SERVER_PATH | cut -d: -f2)

echo "🐳 在服务器上构建和启动服务..."
ssh $SERVER_USER << EOF
cd $SERVER_DIR
echo "停止现有服务..."
docker-compose down 2>/dev/null || true

echo "构建新镜像..."
docker build -t new-api:latest .

echo "启动服务..."
docker-compose up -d

echo "查看服务状态..."
docker-compose ps

echo "🎉 部署完成！"
echo "访问地址: http://\$(curl -s ifconfig.me):3000"
echo "默认账号: root / 123456"
EOF

echo "🎉 部署脚本执行完成！"
