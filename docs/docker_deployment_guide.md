# Docker 部署指南

本指南将帮助您在服务器上使用 Docker 部署 New API 项目。

## 部署方式选择

### 方式一：使用预构建镜像（推荐）

这是最简单的部署方式，直接使用官方预构建的镜像。

#### 1. 准备环境

确保服务器已安装：
- Docker
- Docker Compose

```bash
# 安装 Docker（Ubuntu/Debian）
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. 创建部署目录

```bash
mkdir -p /opt/new-api
cd /opt/new-api
```

#### 3. 创建 docker-compose.yml

```yaml
version: '3.4'

services:
  new-api:
    image: calciumion/new-api:latest
    container_name: new-api
    restart: always
    command: --log-dir /app/logs
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - SQL_DSN=root:123456@tcp(mysql:3306)/new-api
      - REDIS_CONN_STRING=redis://redis
      - TZ=Asia/Shanghai
      - ERROR_LOG_ENABLED=true
    depends_on:
      - redis
      - mysql
    healthcheck:
      test: ["CMD-SHELL", "wget -q -O - http://localhost:3000/api/status | grep -o '\"success\":\\s*true' | awk -F: '{print $$2}'"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:latest
    container_name: redis
    restart: always

  mysql:
    image: mysql:8.2
    container_name: mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: new-api
    volumes:
      - mysql_data:/var/lib/mysql

volumes:
  mysql_data:
```

#### 4. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f new-api
```

### 方式二：仅使用 SQLite（轻量级部署）

如果不需要 MySQL 和 Redis，可以使用更简单的配置：

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

### 方式三：从源码构建

如果需要自定义修改或使用最新代码：

#### 1. 克隆代码

```bash
git clone https://github.com/Calcium-Ion/new-api.git
cd new-api
```

#### 2. 构建镜像

```bash
# 构建镜像
docker build -t new-api:custom .

# 或使用 docker-compose 构建
docker-compose -f docker-compose.build.yml up -d
```

#### 3. 创建构建用的 docker-compose.build.yml

```yaml
version: '3.4'

services:
  new-api:
    build: .
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

## 环境变量配置

### 基础配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `TZ` | 时区设置 | `Asia/Shanghai` |
| `ERROR_LOG_ENABLED` | 启用错误日志 | `true` |

### 数据库配置

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `SQL_DSN` | 数据库连接字符串 | `root:password@tcp(mysql:3306)/new-api` |

### Redis 配置

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `REDIS_CONN_STRING` | Redis连接字符串 | `redis://redis:6379` |

### 多节点部署配置

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `SESSION_SECRET` | 会话密钥（多机部署必须） | `your_random_string` |
| `NODE_TYPE` | 节点类型 | `slave` |
| `SYNC_FREQUENCY` | 同步频率（秒） | `60` |
| `FRONTEND_BASE_URL` | 前端基础URL | `https://your-domain.com` |

## 常用操作命令

### 服务管理

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f new-api

# 更新镜像
docker-compose pull
docker-compose up -d
```

### 数据备份

```bash
# 备份 SQLite 数据库
docker cp new-api:/data/one-api.db ./backup/

# 备份 MySQL 数据库
docker exec mysql mysqldump -u root -p123456 new-api > backup.sql

# 恢复 MySQL 数据库
docker exec -i mysql mysql -u root -p123456 new-api < backup.sql
```

### 日志管理

```bash
# 查看应用日志
docker-compose logs new-api

# 清理日志
docker-compose logs --tail=0 new-api

# 查看容器资源使用
docker stats new-api
```

## 反向代理配置

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Traefik 配置示例

```yaml
version: '3.4'

services:
  new-api:
    image: calciumion/new-api:latest
    container_name: new-api
    restart: always
    volumes:
      - ./data:/data
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.new-api.rule=Host(`your-domain.com`)"
      - "traefik.http.routers.new-api.entrypoints=websecure"
      - "traefik.http.routers.new-api.tls.certresolver=myresolver"
      - "traefik.http.services.new-api.loadbalancer.server.port=3000"
```

## 故障排除

### 常见问题

1. **容器启动失败**
   ```bash
   # 查看详细错误信息
   docker-compose logs new-api
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库容器状态
   docker-compose ps mysql
   
   # 查看数据库日志
   docker-compose logs mysql
   ```

3. **端口冲突**
   ```bash
   # 修改 docker-compose.yml 中的端口映射
   ports:
     - "8080:3000"  # 改为其他端口
   ```

4. **权限问题**
   ```bash
   # 确保数据目录权限正确
   sudo chown -R 1000:1000 ./data ./logs
   ```

### 性能优化

1. **资源限制**
   ```yaml
   services:
     new-api:
       deploy:
         resources:
           limits:
             memory: 1G
             cpus: '0.5'
   ```

2. **日志轮转**
   ```yaml
   services:
     new-api:
       logging:
         driver: "json-file"
         options:
           max-size: "10m"
           max-file: "3"
   ```

## 安全建议

1. **修改默认密码**
   - 修改 MySQL root 密码
   - 设置强密码策略

2. **网络安全**
   - 使用防火墙限制端口访问
   - 配置 HTTPS

3. **定期更新**
   - 定期更新镜像版本
   - 关注安全公告

4. **数据备份**
   - 定期备份数据库
   - 测试恢复流程

## 监控和维护

### 健康检查

项目已内置健康检查，可通过以下方式监控：

```bash
# 检查健康状态
curl http://localhost:3000/api/status

# 查看健康检查日志
docker inspect new-api | grep -A 10 Health
```

### 资源监控

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
df -h

# 查看内存使用
free -h
```

通过以上配置，您就可以在服务器上成功部署 New API 项目了。建议先在测试环境验证配置，然后再部署到生产环境。
