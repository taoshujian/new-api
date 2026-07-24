# 使用 MySQL 替换 PostgreSQL

- **日期**：2026-06-16
- **提交**：`88b38f8c` chore: 使用MySQL
- **类型**：部署配置变更

## 背景

项目默认 `docker-compose.yml` 使用 PostgreSQL 作为主数据库，本分支将开发/默认部署环境切换为 MySQL。

## 变更内容

仅修改 `docker-compose.yml`（+26 / -26）：

1. **数据库服务切换**
   - 注释掉 `postgres` 服务（postgres:15，`pg_data` 数据卷）
   - 启用 `mysql` 服务（mysql:8.2，`mysql_data` 数据卷，root 密码 123456，数据库 `new-api`）

2. **应用连接串切换**
   - `SQL_DSN` 由 `postgresql://root:123456@postgres:5432/new-api` 改为 `root:123456@tcp(mysql:3306)/new-api`

3. **端口映射调整**
   - new-api 服务端口由 `3000:3000` 改为 `3000:8080`

4. **数据卷切换**
   - `pg_data` 卷注释，启用 `mysql_data` 卷
