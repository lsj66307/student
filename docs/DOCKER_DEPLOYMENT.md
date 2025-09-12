# Docker 部署指南

本文档提供了学生管理系统的完整 Docker 部署指南，适用于 OrbStack 和 Docker Desktop。

## 📋 前置要求

- 安装 Docker 或 OrbStack
- Go 1.25+ (用于本地开发)
- 确保端口 3060 未被占用

## 🚀 快速部署

### 1. 构建 Docker 镜像

```bash
# 在项目根目录执行
docker build -t student-management-system:latest .
```

### 2. 运行 Docker 容器

```bash
# 后台运行容器，映射端口3060
docker run -d -p 3060:3060 --name student-api student-management-system:latest
```

### 3. 验证部署

```bash
# 查看容器状态
docker ps

# 查看应用日志
docker logs student-api

# 测试健康检查
curl http://localhost:3060/health

# 测试API端点
curl http://localhost:3060/api/v1/students
```

## 🔧 容器管理命令

### 基本操作

```bash
# 停止容器
docker stop student-api

# 启动已停止的容器
docker start student-api

# 重启容器
docker restart student-api

# 删除容器（需要先停止）
docker stop student-api
docker rm student-api
```

### 日志查看

```bash
# 查看容器日志
docker logs student-api

# 实时查看日志
docker logs -f student-api

# 查看最近100行日志
docker logs --tail 100 student-api
```

### 容器调试

```bash
# 进入容器内部
docker exec -it student-api sh

# 查看容器详细信息
docker inspect student-api

# 查看容器资源使用情况
docker stats student-api
```

## 🗂️ 镜像管理

```bash
# 查看本地镜像
docker images

# 删除镜像
docker rmi student-management-system:latest

# 清理未使用的镜像
docker image prune

# 重新构建镜像（无缓存）
docker build --no-cache -t student-management-system:latest .
```

## 🌐 访问地址

部署成功后，可以通过以下地址访问应用：

- **应用主页**: http://localhost:3060
- **健康检查**: http://localhost:3060/health
- **API 文档**: http://localhost:3060/

### API 端点

- **学生管理**: http://localhost:3060/api/v1/students
- **教师管理**: http://localhost:3060/api/v1/teachers
- **成绩管理**: http://localhost:3060/api/v1/grades

## 🐳 使用 Docker Compose

项目包含 `docker-compose.yml` 文件，可以使用以下命令：

```bash
# 启动所有服务（包括数据库）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 停止所有服务
docker-compose down

# 重新构建并启动
docker-compose up --build -d
```

## 🔍 故障排除

### 常见问题

1. **端口被占用**

   ```bash
   # 查看端口占用
   lsof -i :3060

   # 使用其他端口运行
   docker run -d -p 3061:3060 --name student-api student-management-system:latest
   ```

2. **容器启动失败**

   ```bash
   # 查看详细错误信息
   docker logs student-api

   # 以交互模式运行容器进行调试
   docker run -it --rm student-management-system:latest sh
   ```

3. **数据库连接问题**

   ```bash
   # 检查数据库配置
   docker exec -it student-api cat /root/configs/config.yaml

   # 查看网络连接
   docker exec -it student-api ping 192.168.31.114
   ```

### 性能监控

```bash
# 查看容器资源使用
docker stats student-api

# 查看容器进程
docker exec -it student-api ps aux

# 查看容器网络
docker network ls
docker network inspect bridge
```

## 📝 注意事项

1. **数据持久化**: 当前配置使用外部数据库，容器重启不会丢失数据
2. **环境变量**: 可以通过 `-e` 参数传递环境变量覆盖配置
3. **安全性**: 生产环境建议使用非 root 用户运行（Dockerfile 已配置）
4. **网络**: 容器默认使用 bridge 网络，如需自定义网络请参考 Docker 网络文档

## 🚀 生产环境部署建议

```bash
# 使用特定版本标签
docker build -t student-management-system:v1.0.0 .

# 设置资源限制
docker run -d \
  --name student-api \
  --memory=512m \
  --cpus=1.0 \
  -p 3060:3060 \
  --restart=unless-stopped \
  student-management-system:v1.0.0

# 使用健康检查
docker run -d \
  --name student-api \
  --health-cmd="curl -f http://localhost:3060/health || exit 1" \
  --health-interval=30s \
  --health-timeout=10s \
  --health-retries=3 \
  -p 3060:3060 \
  student-management-system:latest
```

## 📚 常用 Docker & Linux 命令参考

### Docker 基础命令 (80% 使用频率)

#### 镜像操作

```bash
# 拉取镜像
docker pull nginx:latest
docker pull postgres:13

# 查看本地镜像
docker images
docker image ls

# 删除镜像
docker rmi image_name:tag
docker rmi $(docker images -q)  # 删除所有镜像

# 构建镜像
docker build -t app_name:tag .
docker build --no-cache -t app_name:tag .

# 镜像清理
docker image prune          # 删除悬空镜像
docker image prune -a       # 删除所有未使用镜像
```

#### 容器操作

```bash
# 运行容器
docker run -d --name container_name image_name
docker run -it --rm image_name sh                    # 交互式运行
docker run -d -p 8080:80 --name web nginx           # 端口映射
docker run -d -v /host/path:/container/path image    # 挂载卷
docker run -d -e ENV_VAR=value image                # 环境变量

# 容器管理
docker ps                   # 查看运行中容器
docker ps -a               # 查看所有容器
docker start container_name
docker stop container_name
docker restart container_name
docker rm container_name
docker rm $(docker ps -aq) # 删除所有容器

# 容器交互
docker exec -it container_name bash
docker exec -it container_name sh
docker logs container_name
docker logs -f container_name              # 实时日志
docker logs --tail 100 container_name     # 最近100行

# 容器信息
docker inspect container_name
docker stats container_name
docker top container_name
```

#### Docker Compose

```bash
# 启动服务
docker-compose up
docker-compose up -d                # 后台运行
docker-compose up --build           # 重新构建

# 服务管理
docker-compose ps
docker-compose logs
docker-compose logs -f service_name
docker-compose stop
docker-compose down
docker-compose down -v              # 删除卷

# 单个服务操作
docker-compose restart service_name
docker-compose exec service_name bash
```

#### 网络和卷

```bash
# 网络管理
docker network ls
docker network create network_name
docker network inspect network_name
docker network rm network_name

# 卷管理
docker volume ls
docker volume create volume_name
docker volume inspect volume_name
docker volume rm volume_name
docker volume prune                 # 清理未使用卷
```

### Linux 基础命令 (80% 使用频率)

#### 文件和目录操作

```bash
# 导航和查看
ls -la                      # 详细列表
cd /path/to/directory
pwd                         # 当前路径
find . -name "*.go"         # 查找文件
locate filename             # 快速查找

# 文件操作
cp source destination       # 复制
mv source destination       # 移动/重命名
rm -rf directory           # 删除目录
mkdir -p path/to/dir       # 创建目录
touch filename             # 创建文件

# 文件内容
cat filename               # 查看文件
less filename              # 分页查看
head -n 20 filename        # 前20行
tail -n 20 filename        # 后20行
tail -f filename           # 实时查看
grep "pattern" filename     # 搜索内容
```

#### 进程和系统监控

```bash
# 进程管理
ps aux                     # 查看所有进程
ps aux | grep process_name # 查找特定进程
top                        # 实时进程监控
htop                       # 增强版top
kill PID                   # 终止进程
killall process_name       # 按名称终止

# 系统信息
df -h                      # 磁盘使用情况
du -sh directory           # 目录大小
free -h                    # 内存使用
uptime                     # 系统运行时间
whoami                     # 当前用户
```

#### 网络和端口

```bash
# 网络诊断
ping hostname              # 网络连通性
curl http://example.com    # HTTP请求
wget http://example.com/file # 下载文件

# 端口和连接
netstat -tulpn             # 查看端口
ss -tulpn                  # 现代版netstat
lsof -i :8080             # 查看端口占用
```

#### 文本处理

```bash
# 文本操作
grep -r "pattern" .         # 递归搜索
grep -i "pattern" file      # 忽略大小写
sed 's/old/new/g' file     # 替换文本
awk '{print $1}' file      # 提取列
sort filename              # 排序
uniq filename              # 去重
wc -l filename             # 行数统计
```

#### 权限和用户

```bash
# 权限管理
chmod 755 filename         # 修改权限
chown user:group filename  # 修改所有者
sudo command               # 以管理员运行
su - username              # 切换用户
```

#### 压缩和解压

```bash
# 压缩解压
tar -czf archive.tar.gz directory/     # 压缩
tar -xzf archive.tar.gz               # 解压
zip -r archive.zip directory/         # ZIP压缩
unzip archive.zip                     # ZIP解压
```

### 🔧 实用组合命令

```bash
# Docker 清理所有资源
docker system prune -a

# 查看容器IP地址
docker inspect container_name | grep IPAddress

# 批量停止容器
docker stop $(docker ps -q)

# 查看镜像大小排序
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | sort -k3 -h

# 实时监控容器资源
watch docker stats

# 查找大文件
find / -type f -size +100M 2>/dev/null

# 查看端口占用详情
lsof -i :3060

# 批量删除日志文件
find /var/log -name "*.log" -mtime +7 -delete

# 查看系统负载
uptime && free -h && df -h
```

### 📝 快捷键提示

```bash
# 终端快捷键
Ctrl + C        # 终止当前命令
Ctrl + Z        # 暂停当前命令
Ctrl + D        # 退出当前会话
Ctrl + L        # 清屏
Ctrl + R        # 搜索历史命令
Ctrl + A        # 光标移到行首
Ctrl + E        # 光标移到行尾

# 历史命令
history         # 查看命令历史
!!              # 重复上一条命令
!n              # 重复第n条历史命令
```

---

更多信息请参考项目文档或联系开发团队。
