-- 数据库初始化脚本
-- 包含学生、教师和成绩数据的完整初始化

-- 清理现有数据（可选，谨慎使用）
-- DELETE FROM grades;
-- DELETE FROM students;
-- DELETE FROM teachers;
-- DELETE FROM admins;

-- 重置序列（可选，谨慎使用）
-- ALTER SEQUENCE students_id_seq RESTART WITH 1;
-- ALTER SEQUENCE teachers_id_seq RESTART WITH 1;
-- ALTER SEQUENCE grades_id_seq RESTART WITH 1;
-- ALTER SEQUENCE admins_id_seq RESTART WITH 1;

-- 插入教师数据（必须先插入教师，因为成绩表有外键引用）
INSERT INTO teachers (name, age, gender, email, phone, subject, title, department) VALUES
('李明华', 35, '男', 'liminghua@school.edu.cn', '13900139001', '语文', '副教授', '中文系'),
('王晓红', 42, '女', 'wangxiaohong@school.edu.cn', '13900139002', '数学', '教授', '数学系'),
('张建国', 38, '男', 'zhangjianguo@school.edu.cn', '13900139003', '英语', '讲师', '外语系'),
('刘美丽', 29, '女', 'liumeili@school.edu.cn', '13900139004', '体育', '助教', '体育系'),
('陈志强', 45, '男', 'chenzhiqiang@school.edu.cn', '13900139005', '音乐', '副教授', '艺术系');

-- 插入学生数据
INSERT INTO students (student_id, name, age, gender, phone, email, address, major, enrollment_date, graduation_date, status) VALUES
('S2024001', '张三', 20, '男', '13800138001', 'zhangsan@example.com', '北京市朝阳区学院路1号', '计算机科学与技术', '2024-09-01', '2028-06-30', 'active'),
('S2024002', '李四', 19, '女', '13800138002', 'lisi@example.com', '上海市浦东新区张江路2号', '软件工程', '2024-09-01', '2028-06-30', 'active'),
('S2024003', '王五', 21, '男', '13800138003', 'wangwu@example.com', '广州市天河区科技路3号', '数据科学与大数据技术', '2024-09-01', '2028-06-30', 'active'),
('S2024004', '赵六', 20, '女', '13800138004', 'zhaoliu@example.com', '深圳市南山区高新路4号', '人工智能', '2024-09-01', '2028-06-30', 'active'),
('S2024005', '钱七', 22, '男', '13800138005', 'qianqi@example.com', '杭州市西湖区文三路5号', '网络工程', '2024-09-01', '2028-06-30', 'active'),
('S2024006', '孙八', 19, '女', '13800138006', 'sunba@example.com', '南京市鼓楼区中山路6号', '信息安全', '2024-09-01', '2028-06-30', 'active'),
('S2024007', '周九', 20, '男', '13800138007', 'zhoujiu@example.com', '武汉市洪山区珞喻路7号', '物联网工程', '2024-09-01', '2028-06-30', 'active'),
('S2024008', '吴十', 21, '女', '13800138008', 'wushi@example.com', '成都市高新区天府大道8号', '电子信息工程', '2024-09-01', '2028-06-30', 'active'),
('S2024009', '郑一', 19, '男', '13800138009', 'zhengyi@example.com', '西安市雁塔区科技路9号', '通信工程', '2024-09-01', '2028-06-30', 'active'),
('S2024010', '陈二', 20, '女', '13800138010', 'chener@example.com', '重庆市渝北区龙溪路10号', '自动化', '2024-09-01', '2028-06-30', 'active');

-- 插入成绩数据（为每个学生分配随机成绩和对应的教师）
INSERT INTO grades (student_id, chinese_score, math_score, english_score, sports_score, music_score, chinese_teacher_id, math_teacher_id, english_teacher_id, sports_teacher_id, music_teacher_id) VALUES
(1, 85.5, 92.0, 78.5, 88.0, 82.5, 1, 2, 3, 4, 5),
(2, 90.0, 87.5, 85.0, 90.5, 88.0, 1, 2, 3, 4, 5),
(3, 78.5, 95.0, 82.0, 85.5, 79.0, 1, 2, 3, 4, 5),
(4, 88.0, 89.5, 91.0, 87.0, 85.5, 1, 2, 3, 4, 5),
(5, 82.5, 84.0, 79.5, 92.0, 86.0, 1, 2, 3, 4, 5),
(6, 91.5, 88.0, 87.5, 89.0, 90.5, 1, 2, 3, 4, 5),
(7, 86.0, 91.5, 83.0, 86.5, 84.0, 1, 2, 3, 4, 5),
(8, 89.5, 86.5, 88.0, 91.5, 87.5, 1, 2, 3, 4, 5),
(9, 84.0, 90.0, 86.5, 88.5, 83.0, 1, 2, 3, 4, 5),
(10, 87.5, 85.5, 89.0, 87.5, 89.0, 1, 2, 3, 4, 5);

-- 插入默认管理员账户（密码需要在应用中进行哈希处理）
-- 注意：这里的密码是明文，实际使用时应该通过应用程序的注册接口创建管理员账户
INSERT INTO admins (account, password, name, phone, email) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统管理员', '13900000000', 'admin@school.edu.cn');

-- 数据初始化完成
-- 可以通过以下查询验证数据：
-- SELECT COUNT(*) FROM students; -- 应该返回 10
-- SELECT COUNT(*) FROM teachers; -- 应该返回 5
-- SELECT COUNT(*) FROM grades;   -- 应该返回 10
-- SELECT COUNT(*) FROM admins;   -- 应该返回 1