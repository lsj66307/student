# SQL 数据初始化脚本

本目录包含学生管理系统的数据初始化脚本，用于在新数据库中预装测试数据。

## 文件说明

### 1. `insert_students.sql`

- 包含 10 个学生的插入语句
- 学生信息包括：学号、姓名、年龄、性别、联系方式、专业等
- 所有学生状态为 `active`，入学时间为 2024-09-01

### 2. `insert_teachers.sql`

- 包含 5 个教师的插入语句
- 涵盖五个科目：语文、数学、英语、体育、音乐
- 教师信息包括：姓名、年龄、性别、联系方式、职称、所属系部等

### 3. `init_data.sql`

- **推荐使用的完整初始化脚本**
- 包含教师、学生、成绩和管理员数据
- 按正确的顺序插入数据（先教师，再学生，最后成绩）
- 包含数据验证查询

## 使用方法

### 方法一：使用完整初始化脚本（推荐）

```bash
# 连接到PostgreSQL数据库
psql -h localhost -U your_username -d your_database

# 执行完整初始化脚本
\i sql/init_data.sql
```

### 方法二：分别执行各个脚本

```bash
# 按顺序执行
\i sql/insert_teachers.sql
\i sql/insert_students.sql
```

### 方法三：使用命令行直接执行

```bash
# 执行完整初始化
psql -h localhost -U your_username -d your_database -f sql/init_data.sql

# 或分别执行
psql -h localhost -U your_username -d your_database -f sql/insert_teachers.sql
psql -h localhost -U your_username -d your_database -f sql/insert_students.sql
```

## 数据概览

### 学生数据（10 个）

- 学号范围：S2024001 - S2024010
- 年龄范围：19-22 岁
- 专业涵盖：计算机、软件工程、数据科学、人工智能等

### 教师数据（5 个）

- 语文教师：李明华（副教授）
- 数学教师：王晓红（教授）
- 英语教师：张建国（讲师）
- 体育教师：刘美丽（助教）
- 音乐教师：陈志强（副教授）

### 成绩数据

- 为每个学生分配了五门课程的成绩
- 成绩范围：78.5-95.0 分
- 每门课程都分配了对应的任课教师

### 管理员账户

- 账号：admin
- 密码：password（已哈希处理）
- 用于系统管理和测试

## 注意事项

1. **执行顺序**：必须先创建表结构，再执行数据插入脚本
2. **数据清理**：`init_data.sql` 中包含可选的数据清理语句（已注释）
3. **密码安全**：管理员密码已进行哈希处理，实际密码为 "password"
4. **外键约束**：成绩表依赖学生表和教师表，请按正确顺序执行
5. **数据验证**：执行完成后可运行脚本末尾的验证查询

## 数据验证

执行以下查询验证数据是否正确插入：

```sql
-- 检查各表数据量
SELECT 'students' as table_name, COUNT(*) as count FROM students
UNION ALL
SELECT 'teachers', COUNT(*) FROM teachers
UNION ALL
SELECT 'grades', COUNT(*) FROM grades
UNION ALL
SELECT 'admins', COUNT(*) FROM admins;

-- 检查学生和成绩关联
SELECT s.student_id, s.name, g.chinese_score, g.math_score, g.english_score
FROM students s
LEFT JOIN grades g ON s.id = g.student_id
ORDER BY s.id;
```

## 故障排除

如果遇到插入错误：

1. **重复键错误**：检查是否已存在相同的学号或教师数据
2. **外键约束错误**：确保先插入教师数据，再插入学生和成绩数据
3. **数据类型错误**：检查日期格式和数值范围是否正确

如需重新初始化，可以取消注释 `init_data.sql` 中的清理语句。
