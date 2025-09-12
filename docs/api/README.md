# API 文档

## 概述

学生管理系统提供 RESTful API 接口，支持学生、教师和成绩的完整管理功能。

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`

## 认证

目前 API 不需要认证，但建议在生产环境中添加 JWT 或其他认证机制。

## 响应格式

### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "错误描述",
  "error": "详细错误信息"
}
```

## 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 分页

对于返回列表的接口，支持分页参数：

- `page`: 页码，从 1 开始，默认为 1
- `limit`: 每页数量，默认为 10，最大为 100

### 分页响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      // 数据列表
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "pages": 10
    }
  }
}
```

## API 端点

### 学生管理

详细文档: [学生 API](./students.md)

- `GET /students` - 获取学生列表
- `GET /students/{id}` - 获取学生详情
- `POST /students` - 创建学生
- `PUT /students/{id}` - 更新学生
- `DELETE /students/{id}` - 删除学生

### 教师管理

详细文档: [教师 API](./teachers.md)

- `GET /teachers` - 获取教师列表
- `GET /teachers/{id}` - 获取教师详情
- `POST /teachers` - 创建教师
- `PUT /teachers/{id}` - 更新教师
- `DELETE /teachers/{id}` - 删除教师

### 成绩管理

详细文档: [成绩 API](./grades.md)

- `GET /grades` - 获取成绩列表
- `GET /grades/{id}` - 获取成绩详情
- `POST /grades` - 创建成绩
- `PUT /grades/{id}` - 更新成绩
- `DELETE /grades/{id}` - 删除成绩
- `GET /students/{id}/grades` - 获取学生成绩

## 数据模型

### 学生 (Student)

```json
{
  "id": 1,
  "name": "张三",
  "age": 20,
  "gender": "男",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 教师 (Teacher)

```json
{
  "id": 1,
  "name": "李老师",
  "age": 35,
  "gender": "女",
  "email": "li@example.com",
  "phone": "13900139000",
  "subject": "数学",
  "department": "数学系",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 成绩 (Grade)

```json
{
  "id": 1,
  "student_id": 1,
  "math_score": 95.5,
  "english_score": 88.0,
  "chinese_score": 92.0,
  "physics_score": 87.5,
  "chemistry_score": 90.0,
  "biology_score": 85.0,
  "history_score": 89.0,
  "geography_score": 86.5,
  "politics_score": 91.0,
  "math_teacher_id": 1,
  "english_teacher_id": 2,
  "chinese_teacher_id": 3,
  "physics_teacher_id": 4,
  "chemistry_teacher_id": 5,
  "biology_teacher_id": 6,
  "history_teacher_id": 7,
  "geography_teacher_id": 8,
  "politics_teacher_id": 9,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 成绩详情 (GradeWithDetails)

```json
{
  "id": 1,
  "student_id": 1,
  "student_name": "张三",
  "math_score": 95.5,
  "math_teacher_name": "李老师",
  "english_score": 88.0,
  "english_teacher_name": "王老师",
  // ... 其他科目
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## 错误处理

### 常见错误

| 错误码 | 错误信息 | 说明 |
|--------|----------|------|
| 1001 | 参数验证失败 | 请求参数不符合要求 |
| 1002 | 资源不存在 | 请求的资源不存在 |
| 1003 | 数据库操作失败 | 数据库操作异常 |
| 1004 | 重复数据 | 尝试创建重复的数据 |

### 验证错误示例

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "name: 姓名不能为空; email: 邮箱格式不正确"
}
```

## 示例代码

### JavaScript (Fetch)

```javascript
// 获取学生列表
fetch('http://localhost:8080/api/v1/students?page=1&limit=10')
  .then(response => response.json())
  .then(data => console.log(data));

// 创建学生
fetch('http://localhost:8080/api/v1/students', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: '张三',
    age: 20,
    gender: '男',
    email: 'zhangsan@example.com',
    phone: '13800138000'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### cURL

```bash
# 获取学生列表
curl -X GET "http://localhost:8080/api/v1/students?page=1&limit=10"

# 创建学生
curl -X POST "http://localhost:8080/api/v1/students" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000"
  }'
```

## 测试

推荐使用以下工具测试 API：

- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- [HTTPie](https://httpie.io/)
- cURL

## 更新日志

- **v1.0.0** (2024-01-01): 初始版本，包含基础的学生、教师、成绩管理功能