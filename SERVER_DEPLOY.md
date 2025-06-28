# New API 服务器部署快速指南

## 🚀 5分钟快速部署

### 第一步：准备服务器

确保您的Linux服务器满足以下要求：
- **系统**：Ubuntu 18.04+、CentOS 7+、Debian 9+
- **内存**：至少 1GB RAM
- **存储**：至少 10GB 可用空间
- **网络**：可访问互联网

### 第二步：安装Docker环境

#### 自动安装脚本（推荐）：
```bash
# 一键安装Docker和Docker Compose
curl -fsSL https://get.docker.com | bash
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 重新登录或执行以下命令使docker组生效
newgrp docker
```

#### 验证安装：
```bash
docker --version
docker-compose --version
```

### 第三步：部署应用

```bash
# 1. 下载项目
git clone https://github.com/your-repo/new-api.git
cd new-api

# 2. 一键部署（选择其中一种）

# 简单部署（个人使用，使用SQLite）
chmod +x deploy.sh
./deploy.sh --simple

# 或者完整部署（生产环境，使用MySQL+Redis）
./deploy.sh --with-db
```

### 第四步：访问应用

部署完成后：
1. 打开浏览器访问：`http://your-server-ip:3000`
2. 默认账号：`root`，密码：`123456`
3. **立即修改默认密码！**

## 🔧 常用管理命令

```bash
# 查看服务状态
docker-compose ps

# 查看实时日志
docker-compose logs -f new-api

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 更新到最新版本
docker-compose pull
docker-compose up -d

# 备份数据
tar -czf backup-$(date +%Y%m%d).tar.gz data/

# 进入容器调试
docker exec -it new-api sh
```

## 🌐 配置域名和HTTPS

### 1. 安装Nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
```

### 2. 配置Nginx反向代理

创建配置文件：
```bash
sudo nano /etc/nginx/sites-available/new-api
```

添加以下内容：
```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为您的域名
    
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

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/new-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. 配置SSL证书（使用Let's Encrypt）

```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加以下行：
# 0 12 * * * /usr/bin/certbot renew --quiet
```

## 🔒 安全配置

### 1. 防火墙设置

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 2. 修改默认端口（可选）

如果您想修改默认端口3000，编辑docker-compose.yml：
```yaml
ports:
  - "8080:3000"  # 将3000改为8080或其他端口
```

然后重启服务：
```bash
docker-compose down
docker-compose up -d
```

## 📊 监控和维护

### 1. 设置日志轮转

```bash
sudo nano /etc/logrotate.d/new-api
```

添加内容：
```
/opt/new-api/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
}
```

### 2. 系统监控脚本

创建监控脚本：
```bash
nano monitor.sh
```

```bash
#!/bin/bash
# New API 服务监控脚本

check_service() {
    if docker-compose ps | grep -q "Up"; then
        echo "✅ 服务运行正常"
        return 0
    else
        echo "❌ 服务异常，尝试重启..."
        docker-compose restart
        return 1
    fi
}

check_disk() {
    USAGE=$(df /opt/new-api | awk 'NR==2 {print $5}' | sed 's/%//')
    if [ $USAGE -gt 80 ]; then
        echo "⚠️  磁盘使用率过高: ${USAGE}%"
    else
        echo "✅ 磁盘使用正常: ${USAGE}%"
    fi
}

check_memory() {
    USAGE=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
    if [ $USAGE -gt 80 ]; then
        echo "⚠️  内存使用率过高: ${USAGE}%"
    else
        echo "✅ 内存使用正常: ${USAGE}%"
    fi
}

echo "=== New API 服务监控 $(date) ==="
cd /opt/new-api
check_service
check_disk
check_memory
echo "=================================="
```

设置定时检查：
```bash
chmod +x monitor.sh
crontab -e
# 添加：每5分钟检查一次
# */5 * * * * /opt/new-api/monitor.sh >> /var/log/new-api-monitor.log 2>&1
```

## 🔄 备份和恢复

### 自动备份脚本

```bash
nano backup.sh
```

```bash
#!/bin/bash
# 自动备份脚本

BACKUP_DIR="/opt/backups/new-api"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# 备份数据目录
tar -czf $BACKUP_DIR/data_$DATE.tar.gz -C /opt/new-api data/

# 备份配置文件
cp /opt/new-api/docker-compose.yml $BACKUP_DIR/docker-compose_$DATE.yml

# 如果使用MySQL，备份数据库
if docker-compose ps | grep -q mysql; then
    docker exec mysql mysqldump -u root -p$(grep MYSQL_ROOT_PASSWORD /opt/new-api/.env | cut -d'=' -f2) new-api > $BACKUP_DIR/database_$DATE.sql
fi

# 清理30天前的备份
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete

echo "备份完成: $BACKUP_DIR"
```

设置自动备份：
```bash
chmod +x backup.sh
crontab -e
# 添加：每天凌晨2点备份
# 0 2 * * * /opt/new-api/backup.sh
```

## 🆘 故障排除

### 常见问题

1. **服务无法启动**
   ```bash
   # 查看详细日志
   docker-compose logs new-api
   
   # 检查端口占用
   sudo netstat -tlnp | grep :3000
   ```

2. **内存不足**
   ```bash
   # 查看内存使用
   free -h
   docker stats
   
   # 清理Docker缓存
   docker system prune -f
   ```

3. **磁盘空间不足**
   ```bash
   # 查看磁盘使用
   df -h
   
   # 清理日志文件
   docker-compose exec new-api sh -c "find /app/logs -name '*.log' -mtime +7 -delete"
   ```

4. **数据库连接失败**
   ```bash
   # 检查数据库容器状态
   docker-compose ps mysql
   
   # 查看数据库日志
   docker-compose logs mysql
   ```

### 紧急恢复

如果服务完全无法访问：
```bash
# 1. 停止所有服务
docker-compose down

# 2. 清理并重新启动
docker system prune -f
docker-compose up -d

# 3. 如果仍有问题，重新部署
./deploy.sh --simple
```

## 📞 获取帮助

- 查看项目文档：`README.md`
- 查看详细部署指南：`DEPLOYMENT.md`
- 检查GitHub Issues
- 社区支持论坛

---

**重要提醒**：
1. 定期备份数据
2. 及时更新系统和应用
3. 监控服务状态
4. 保护好管理员密码
5. 配置HTTPS加密
