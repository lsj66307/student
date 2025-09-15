# 数据库表查询命令

## 1. 使用 psql 命令行工具

### 连接数据库

```bash
PGPASSWORD=mm152002 psql -h 192.168.31.114 -p 5432 -U postgres -d student_management
```

### 列出所有表

```sql
\dt
```

### 查看表结构

```sql
\d table_name
```

### 列出所有表（SQL 查询）

```sql
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;
```

### 查看表的详细信息

```sql
SELECT
    table_name,
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_schema = 'public'
ORDER BY table_name, ordinal_position;
```

## 2. 使用 curl 命令（如果有数据库管理 API）

### 获取管理员列表（需要先登录获取 token）

```bash
# 先登录获取token
curl -X POST http://localhost:3060/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account": "admin", "password": "123456"}'

# 使用token获取管理员列表
curl -X GET http://localhost:3060/api/v1/admins \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 3. 当前数据库表结构总结

根据刚才的查询结果，您的数据库包含以下表：

1. **admins** - 管理员表

   - id (主键)
   - account (账号)
   - password (密码)
   - name (姓名)
   - phone (电话)
   - email (邮箱)
   - created_at, updated_at (时间戳)

2. **grades** - 成绩表

   - id (主键)
   - student_id (学生 ID)
   - subject (科目)
   - score (分数)
   - exam_date (考试日期)
   - created_at, updated_at (时间戳)

3. **students** - 学生表

   - id (主键)
   - student_id (学号)
   - name (姓名)
   - age (年龄)
   - gender (性别)
   - phone (电话)
   - email (邮箱)
   - address (地址)
   - major (专业)
   - enrollment_date (入学日期)
   - graduation_date (毕业日期)
   - status (状态)
   - created_at, updated_at (时间戳)

4. **teachers** - 教师表
   - id (主键)
   - name (姓名)
   - subject (科目)
   - email (邮箱)
   - phone (电话)
   - created_at, updated_at (时间戳)

## 4. 快速查看表数据

### 查看各表的记录数

```sql
SELECT
    schemaname,
    tablename,
    n_tup_ins - n_tup_del as row_count
FROM pg_stat_user_tables
WHERE schemaname = 'public';
```

### 查看表的大小

```sql
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public';
```
