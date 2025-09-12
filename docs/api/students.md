# 学生 API

## 概述

学生 API 提供学生信息的完整管理功能，包括创建、查询、更新和删除操作。

## 端点列表

| 方法   | 端点             | 描述         |
| ------ | ---------------- | ------------ |
| GET    | `/students`      | 获取学生列表 |
| GET    | `/students/{id}` | 获取学生详情 |
| POST   | `/students`      | 创建学生     |
| PUT    | `/students/{id}` | 更新学生     |
| DELETE | `/students/{id}` | 删除学生     |

## 详细接口

### 1. 获取学生列表

获取所有学生的分页列表。

**请求**

```
GET /api/v1/students
```

**查询参数**

| 参数    | 类型    | 必需 | 默认值 | 描述                |
| ------- | ------- | ---- | ------ | ------------------- |
| page    | integer | 否   | 1      | 页码                |
| limit   | integer | 否   | 10     | 每页数量 (最大 100) |
| name    | string  | 否   | -      | 按姓名模糊搜索      |
| gender  | string  | 否   | -      | 按性别筛选 (男/女)  |
| age_min | integer | 否   | -      | 最小年龄            |
| age_max | integer | 否   | -      | 最大年龄            |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "张三",
        "age": 20,
        "gender": "男",
        "email": "zhangsan@example.com",
        "phone": "13800138000",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "李四",
        "age": 19,
        "gender": "女",
        "email": "lisi@example.com",
        "phone": "13900139000",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "pages": 1
    }
  }
}
```

### 2. 获取学生详情

根据学生 ID 获取详细信息。

**请求**

```
GET /api/v1/students/{id}
```

**路径参数**

| 参数 | 类型    | 必需 | 描述    |
| ---- | ------- | ---- | ------- |
| id   | integer | 是   | 学生 ID |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "张三",
    "age": 20,
    "gender": "男",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**错误响应**

```json
{
  "code": 404,
  "message": "学生不存在",
  "error": "student not found"
}
```

### 3. 创建学生

创建新的学生记录。

**请求**

```
POST /api/v1/students
```

**请求体**

```json
{
  "name": "王五",
  "age": 18,
  "gender": "男",
  "email": "wangwu@example.com",
  "phone": "13700137000"
}
```

**字段说明**

| 字段   | 类型    | 必需 | 验证规则    | 描述     |
| ------ | ------- | ---- | ----------- | -------- |
| name   | string  | 是   | 长度 1-50   | 学生姓名 |
| age    | integer | 是   | 范围 1-150  | 学生年龄 |
| gender | string  | 是   | 枚举: 男/女 | 学生性别 |
| email  | string  | 是   | 邮箱格式    | 学生邮箱 |
| phone  | string  | 是   | 手机号格式  | 学生电话 |

**响应示例**

```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": 3,
    "name": "王五",
    "age": 18,
    "gender": "男",
    "email": "wangwu@example.com",
    "phone": "13700137000",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**验证错误响应**

```json
{
  "code": 400,
  "message": "参数验证失败",
  "error": "name: 姓名不能为空; email: 邮箱格式不正确"
}
```

### 4. 更新学生

更新现有学生的信息。

**请求**

```
PUT /api/v1/students/{id}
```

**路径参数**

| 参数 | 类型    | 必需 | 描述    |
| ---- | ------- | ---- | ------- |
| id   | integer | 是   | 学生 ID |

**请求体**

```json
{
  "name": "张三三",
  "age": 21,
  "gender": "男",
  "email": "zhangsan_new@example.com",
  "phone": "13800138001"
}
```

**字段说明**

所有字段都是可选的，只更新提供的字段。验证规则与创建接口相同。

**响应示例**

```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "name": "张三三",
    "age": 21,
    "gender": "男",
    "email": "zhangsan_new@example.com",
    "phone": "13800138001",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:30:00Z"
  }
}
```

### 5. 删除学生

删除指定的学生记录。

**请求**

```
DELETE /api/v1/students/{id}
```

**路径参数**

| 参数 | 类型    | 必需 | 描述    |
| ---- | ------- | ---- | ------- |
| id   | integer | 是   | 学生 ID |

**响应示例**

```json
{
  "code": 200,
  "message": "删除成功"
}
```

**错误响应**

```json
{
  "code": 404,
  "message": "学生不存在",
  "error": "student not found"
}
```

## 使用示例

### cURL 示例

```bash
# 获取学生列表
curl -X GET "http://localhost:8080/api/v1/students?page=1&limit=5"

# 获取学生详情
curl -X GET "http://localhost:8080/api/v1/students/1"

# 创建学生
curl -X POST "http://localhost:8080/api/v1/students" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "赵六",
    "age": 19,
    "gender": "女",
    "email": "zhaoliu@example.com",
    "phone": "13600136000"
  }'

# 更新学生
curl -X PUT "http://localhost:8080/api/v1/students/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三丰",
    "age": 22
  }'

# 删除学生
curl -X DELETE "http://localhost:8080/api/v1/students/1"
```

### JavaScript 示例

```javascript
const API_BASE = "http://localhost:8080/api/v1";

// 获取学生列表
async function getStudents(page = 1, limit = 10) {
  const response = await fetch(
    `${API_BASE}/students?page=${page}&limit=${limit}`
  );
  return await response.json();
}

// 获取学生详情
async function getStudent(id) {
  const response = await fetch(`${API_BASE}/students/${id}`);
  return await response.json();
}

// 创建学生
async function createStudent(studentData) {
  const response = await fetch(`${API_BASE}/students`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(studentData),
  });
  return await response.json();
}

// 更新学生
async function updateStudent(id, studentData) {
  const response = await fetch(`${API_BASE}/students/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(studentData),
  });
  return await response.json();
}

// 删除学生
async function deleteStudent(id) {
  const response = await fetch(`${API_BASE}/students/${id}`, {
    method: "DELETE",
  });
  return await response.json();
}

// 使用示例
(async () => {
  // 创建学生
  const newStudent = await createStudent({
    name: "孙七",
    age: 20,
    gender: "男",
    email: "sunqi@example.com",
    phone: "13500135000",
  });
  console.log("创建的学生:", newStudent);

  // 获取学生列表
  const students = await getStudents(1, 10);
  console.log("学生列表:", students);
})();
```

## 注意事项

1. **数据验证**: 所有输入数据都会进行严格验证，确保数据的完整性和正确性。

2. **唯一性约束**: 邮箱和电话号码必须唯一，不能重复。

3. **软删除**: 删除操作是物理删除，请谨慎操作。

4. **分页限制**: 每页最多返回 100 条记录，建议使用合适的分页大小。

5. **搜索功能**: 姓名搜索支持模糊匹配，不区分大小写。

6. **错误处理**: 所有错误都会返回统一格式的错误信息，便于客户端处理。
