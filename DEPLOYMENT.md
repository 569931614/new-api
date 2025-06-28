# New API Linux 服务器部署指南

## 快速部署（推荐）

本项目提供了一键部署脚本，支持多种部署方式。

### 1. 准备服务器环境

#### Ubuntu/Debian 系统：
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Docker
curl -fsSL https://get.docker.com | bash
sudo systemctl start docker
sudo systemctl enable docker

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 将当前用户添加到docker组
sudo usermod -aG docker $USER
# 重新登录或执行：newgrp docker
```

#### CentOS/RHEL 系统：
```bash
# 更新系统
sudo yum update -y

# 安装Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 将当前用户添加到docker组
sudo usermod -aG docker $USER
```

### 2. 一键部署

```bash
# 下载项目
git clone https://github.com/your-repo/new-api.git
cd new-api

# 给部署脚本执行权限
chmod +x deploy.sh

# 选择部署方式：

# 方式1：简单部署（使用SQLite，推荐个人使用）
./deploy.sh --simple

# 方式2：完整部署（使用MySQL+Redis，推荐生产环境）
./deploy.sh --with-db

# 方式3：从源码构建部署
./deploy.sh --build
```

### 3. 部署完成后

部署成功后，您将看到以下信息：

```
========================================
           New API 部署完成
========================================

访问地址: http://your-server-ip:3000
默认管理员账号: root
默认密码: 123456

常用命令:
  查看服务状态: docker-compose ps
  查看日志:     docker-compose logs -f new-api
  停止服务:     docker-compose down
  重启服务:     docker-compose restart
  更新服务:     docker-compose pull && docker-compose up -d
```

### 4. 首次登录设置

1. 打开浏览器访问：`http://your-server-ip:3000`
2. 使用默认账号密码登录：`root` / `123456`
3. **立即修改默认密码**
4. 配置系统设置和渠道信息

## 手动部署方式

如果您需要更多控制，可以手动部署：

### 1. 创建项目目录

```bash
mkdir -p /opt/new-api
cd /opt/new-api
mkdir -p data logs
```

### 2. 创建 docker-compose.yml

#### 简单部署（SQLite）：
```yaml
version: '3.4'

services:
  new-api:
    image: calciumion/new-api:latest
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
```

#### 完整部署（MySQL + Redis）：
```yaml
version: '3.4'

services:
  new-api:
    image: calciumion/new-api:latest
    container_name: new-api
    restart: always
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - SQL_DSN=root:your_password@tcp(mysql:3306)/new-api
      - REDIS_CONN_STRING=redis://redis
      - TZ=Asia/Shanghai
      - ERROR_LOG_ENABLED=true
    depends_on:
      - redis
      - mysql

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
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: new-api
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
  redis_data:
```

### 3. 启动服务

```bash
# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

## 方案二：自定义Docker镜像部署

### 1. 构建自定义镜像

```bash
# 在项目根目录执行
docker build -t new-api:custom .

# 运行容器
docker run -d \
  --name new-api \
  -p 3000:3000 \
  -v $(pwd)/data:/data \
  -v $(pwd)/logs:/app/logs \
  -e TZ=Asia/Shanghai \
  new-api:custom
```

## 方案三：直接二进制部署

### 1. 编译项目

```bash
# 在开发机器上编译
cd web
npm install
npm run build

cd ..
go mod download
go build -ldflags "-s -w" -o new-api

# 上传二进制文件到服务器
scp new-api user@server:/opt/new-api/
scp -r web/dist user@server:/opt/new-api/web/
```

### 2. 服务器配置

```bash
# 创建服务用户
sudo useradd -r -s /bin/false new-api

# 创建目录
sudo mkdir -p /opt/new-api/{data,logs}
sudo chown -R new-api:new-api /opt/new-api

# 创建systemd服务
sudo nano /etc/systemd/system/new-api.service
```

### 3. systemd 服务配置

```ini
[Unit]
Description=New API Service
After=network.target

[Service]
Type=simple
User=new-api
Group=new-api
WorkingDirectory=/opt/new-api
ExecStart=/opt/new-api/new-api
Restart=always
RestartSec=5
Environment=PORT=3000
Environment=TZ=Asia/Shanghai

[Install]
WantedBy=multi-user.target
```

### 4. 启动服务

```bash
# 重载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start new-api
sudo systemctl enable new-api

# 查看状态
sudo systemctl status new-api

# 查看日志
sudo journalctl -u new-api -f
```

## 环境变量配置

### 常用环境变量

```bash
# 数据库配置
SQL_DSN=root:password@tcp(localhost:3306)/new-api

# Redis配置
REDIS_CONN_STRING=redis://localhost:6379

# 端口配置
PORT=3000

# 日志配置
LOG_DIR=/app/logs
ERROR_LOG_ENABLED=true

# 会话配置（多机部署必须设置）
SESSION_SECRET=your-random-secret-string

# 多节点部署
NODE_TYPE=slave
SYNC_FREQUENCY=60
FRONTEND_BASE_URL=https://your-domain.com

# 时区设置
TZ=Asia/Shanghai
```

## 反向代理配置

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    # SSL证书配置
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## 数据备份

### 1. 数据库备份

```bash
# MySQL备份
docker exec mysql mysqldump -u root -p123456 new-api > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复
docker exec -i mysql mysql -u root -p123456 new-api < backup.sql
```

### 2. 文件备份

```bash
# 备份数据目录
tar -czf new-api-data-$(date +%Y%m%d_%H%M%S).tar.gz data/

# 备份日志
tar -czf new-api-logs-$(date +%Y%m%d_%H%M%S).tar.gz logs/
```

## 监控和维护

### 1. 健康检查

```bash
# 检查服务状态
curl http://localhost:3000/api/status

# 检查Docker容器
docker-compose ps
docker stats
```

### 2. 日志管理

```bash
# 查看应用日志
docker-compose logs -f new-api

# 清理旧日志
find logs/ -name "*.log" -mtime +30 -delete
```

### 3. 更新部署

```bash
# Docker Compose方式更新
docker-compose pull
docker-compose up -d

# 自定义镜像更新
docker build -t new-api:custom .
docker-compose up -d
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   sudo netstat -tlnp | grep :3000
   sudo lsof -i :3000
   ```

2. **权限问题**
   ```bash
   sudo chown -R new-api:new-api /opt/new-api
   sudo chmod +x /opt/new-api/new-api
   ```

3. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接字符串
   - 检查防火墙设置

4. **内存不足**
   ```bash
   free -h
   docker stats
   ```

## 安全建议

1. **修改默认密码**
2. **配置防火墙**
3. **使用HTTPS**
4. **定期备份数据**
5. **监控系统资源**
6. **及时更新系统和应用**

## 性能优化

1. **数据库优化**
   - 配置合适的连接池
   - 定期清理日志表
   - 添加必要的索引

2. **缓存配置**
   - 启用Redis缓存
   - 配置适当的缓存策略

3. **负载均衡**
   - 多实例部署
   - 使用Nginx负载均衡
