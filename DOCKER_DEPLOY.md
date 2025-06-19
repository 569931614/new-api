# Docker 部署 New API 项目

## 快速开始

### 方法一：使用预构建镜像（推荐）

这是最简单的部署方式，适合大多数用户。

#### 1. 简单部署（仅使用 SQLite）

```bash
# 创建部署目录
mkdir new-api && cd new-api

# 创建数据和日志目录
mkdir -p data logs

# 下载简单配置文件
curl -o docker-compose.yml https://raw.githubusercontent.com/Calcium-Ion/new-api/main/docker-compose.simple.yml

# 启动服务
docker-compose up -d
```

或者直接使用以下配置创建 `docker-compose.yml`：

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

#### 2. 完整部署（MySQL + Redis）

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
```

### 方法二：从源码构建

如果您需要自定义修改：

```bash
# 克隆代码
git clone https://github.com/Calcium-Ion/new-api.git
cd new-api

# 构建并启动
docker-compose up -d --build
```

## 部署步骤

### 1. 准备服务器环境

确保服务器已安装 Docker 和 Docker Compose：

```bash
# 安装 Docker（Ubuntu/Debian）
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. 创建部署目录

```bash
mkdir -p /opt/new-api
cd /opt/new-api
mkdir -p data logs
```

### 3. 创建配置文件

选择上面的配置之一，创建 `docker-compose.yml` 文件。

### 4. 启动服务

```bash
# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f new-api
```

### 5. 访问应用

打开浏览器访问：`http://your-server-ip:3000`

默认管理员账号：
- 用户名：`root`
- 密码：`123456`

**重要：请立即修改默认密码！**

## 常用管理命令

```bash
# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f new-api

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 更新镜像
docker-compose pull
docker-compose up -d

# 备份数据
docker cp new-api:/data ./backup/

# 进入容器
docker exec -it new-api sh
```

## 环境变量配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `TZ` | 时区 | `Asia/Shanghai` | `Asia/Shanghai` |
| `SQL_DSN` | 数据库连接 | SQLite | `root:pass@tcp(mysql:3306)/new-api` |
| `REDIS_CONN_STRING` | Redis连接 | 无 | `redis://redis:6379` |
| `ERROR_LOG_ENABLED` | 错误日志 | `true` | `true/false` |
| `SESSION_SECRET` | 会话密钥 | 随机生成 | `your_secret_key` |

## 反向代理配置

### Nginx 配置

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

### 使用 SSL

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 修改端口映射
   ports:
     - "8080:3000"  # 改为其他端口
   ```

2. **权限问题**
   ```bash
   sudo chown -R 1000:1000 ./data ./logs
   ```

3. **容器无法启动**
   ```bash
   # 查看详细日志
   docker-compose logs new-api
   ```

4. **数据库连接失败**
   ```bash
   # 检查数据库容器
   docker-compose logs mysql
   ```

### 性能优化

1. **限制资源使用**
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

## 数据备份

### SQLite 备份

```bash
# 备份数据库文件
docker cp new-api:/data/one-api.db ./backup/

# 恢复数据库
docker cp ./backup/one-api.db new-api:/data/
```

### MySQL 备份

```bash
# 备份数据库
docker exec mysql mysqldump -u root -p new-api > backup.sql

# 恢复数据库
docker exec -i mysql mysql -u root -p new-api < backup.sql
```

## 安全建议

1. **修改默认密码**
2. **使用强密码**
3. **启用 HTTPS**
4. **定期更新镜像**
5. **配置防火墙**
6. **定期备份数据**

## 监控

### 健康检查

```bash
# 检查服务健康状态
curl http://localhost:3000/api/status
```

### 资源监控

```bash
# 查看容器资源使用
docker stats new-api
```

通过以上步骤，您就可以在服务器上成功部署 New API 项目了！
