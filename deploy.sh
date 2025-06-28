#!/bin/bash

# New API Docker 快速部署脚本
# 使用方法: ./deploy.sh [选项]
# 选项:
#   --with-db    使用 MySQL + Redis 完整部署
#   --simple     仅使用 SQLite 简单部署
#   --build      从源码构建部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    print_message "Docker 环境检查通过"
}

# 创建部署目录
create_directories() {
    print_step "创建部署目录..."
    mkdir -p data logs
    print_message "目录创建完成"
}

# 生成简单部署配置
generate_simple_compose() {
    print_step "生成简单部署配置..."
    cat > docker-compose.yml << 'EOF'
version: '3.4'

services:
  new-api:
    build: .
    image: new-api:custom
    container_name: new-api
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
      - ERROR_LOG_ENABLED=true
    healthcheck:
      test: ["CMD-SHELL", "wget -q -O - http://localhost:3000/api/status | grep -o '\"success\":\\s*true' | awk -F: '{print $$2}'"]
      interval: 30s
      timeout: 10s
      retries: 3
EOF
    print_message "简单部署配置生成完成"
}

# 生成完整部署配置
generate_full_compose() {
    print_step "生成完整部署配置..."
    
    # 生成随机密码
    MYSQL_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-16)
    
    cat > docker-compose.yml << EOF
version: '3.4'

services:
  new-api:
    build: .
    image: new-api:custom
    container_name: new-api
    restart: always
    command: --log-dir /app/logs
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - SQL_DSN=root:${MYSQL_PASSWORD}@tcp(mysql:3306)/new-api
      - REDIS_CONN_STRING=redis://redis
      - TZ=Asia/Shanghai
      - ERROR_LOG_ENABLED=true
    depends_on:
      - redis
      - mysql
    healthcheck:
      test: ["CMD-SHELL", "wget -q -O - http://localhost:3000/api/status | grep -o '\"success\":\\s*true' | awk -F: '{print \$\$2}'"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:latest
    container_name: redis
    restart: always
    volumes:
      - redis_data:/data

  mysql:
    image: mysql:8.2
    container_name: mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_PASSWORD}
      MYSQL_DATABASE: new-api
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
  redis_data:
EOF

    # 保存密码到文件
    echo "MySQL Root Password: ${MYSQL_PASSWORD}" > .env
    print_message "完整部署配置生成完成"
    print_warning "MySQL 密码已保存到 .env 文件中，请妥善保管"
}

# 生成构建部署配置
generate_build_compose() {
    print_step "生成构建部署配置..."
    cat > docker-compose.yml << 'EOF'
version: '3.4'

services:
  new-api:
    build: .
    image: new-api:custom
    container_name: new-api
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
      - ERROR_LOG_ENABLED=true
    healthcheck:
      test: ["CMD-SHELL", "wget -q -O - http://localhost:3000/api/status | grep -o '\"success\":\\s*true' | awk -F: '{print $$2}'"]
      interval: 30s
      timeout: 10s
      retries: 3
EOF
    print_message "构建部署配置生成完成"
}

# 启动服务
start_services() {
    print_step "启动服务..."
    docker-compose up -d
    
    print_step "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    if docker-compose ps | grep -q "Up"; then
        print_message "服务启动成功！"
        print_message "访问地址: http://localhost:3000"
        print_message "默认管理员账号: root"
        print_message "默认管理员密码: 123456"
        print_warning "请及时修改默认密码！"
    else
        print_error "服务启动失败，请检查日志"
        docker-compose logs
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "New API Docker 快速部署脚本"
    echo ""
    echo "使用方法:"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  --simple     仅使用 SQLite 简单部署（默认）"
    echo "  --with-db    使用 MySQL + Redis 完整部署"
    echo "  --build      从源码构建部署"
    echo "  --help       显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0                # 简单部署"
    echo "  $0 --with-db      # 完整部署"
    echo "  $0 --build        # 构建部署"
}

# 主函数
main() {
    echo "========================================"
    echo "    New API Docker 快速部署脚本"
    echo "========================================"
    echo ""
    
    # 解析参数
    case "${1:-}" in
        --simple)
            DEPLOY_TYPE="simple"
            ;;
        --with-db)
            DEPLOY_TYPE="full"
            ;;
        --build)
            DEPLOY_TYPE="build"
            ;;
        --help)
            show_help
            exit 0
            ;;
        "")
            DEPLOY_TYPE="simple"
            ;;
        *)
            print_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
    
    print_message "部署类型: $DEPLOY_TYPE"
    echo ""
    
    # 执行部署步骤
    check_docker
    create_directories
    
    case "$DEPLOY_TYPE" in
        "simple")
            generate_simple_compose
            ;;
        "full")
            generate_full_compose
            ;;
        "build")
            if [ ! -f "Dockerfile" ]; then
                print_error "未找到 Dockerfile，请在项目根目录运行此脚本"
                exit 1
            fi
            generate_build_compose
            ;;
    esac
    
    start_services
    
    echo ""
    echo "========================================"
    print_message "部署完成！"
    echo "========================================"
    echo ""
    echo "常用命令:"
    echo "  查看服务状态: docker-compose ps"
    echo "  查看日志:     docker-compose logs -f new-api"
    echo "  停止服务:     docker-compose down"
    echo "  重启服务:     docker-compose restart"
    echo "  更新服务:     docker-compose pull && docker-compose up -d"
    echo ""
}

# 运行主函数
main "$@"
