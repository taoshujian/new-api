# 新增生产环境 Docker 部署配置

- **日期**：2026-06-17
- **提交**：`eec70b0d` chore: 新增生产环境yml
- **类型**：部署配置新增

## 背景

上游仓库只有统一的 `docker-compose.yml`，缺少独立部署（不依赖编排内数据库）和生产环境的部署模板。

## 变更内容

新增 3 个文件，微调 `docker-compose.yml`，共 +174 行：

1. **`Dockerfile.alone`**（64 行，新增）
   - 独立部署镜像构建文件，产出 `new-api:alone` 系列镜像。

2. **`docker-compose.alone.yml`**（54 行，新增）
   - 独立部署编排：仅 new-api 单容器，连接宿主机（`host.docker.internal`）上的 MySQL（3307）/ Redis（6380）。

3. **`docker-compose.prod.yml`**（55 行，新增）
   - 生产环境编排模板，含 `SQL_DSN`、`REDIS_CONN_STRING`、`TZ`、`ERROR_LOG_ENABLED`、`BATCH_UPDATE_ENABLED`、`NODE_NAME` 等环境变量配置。

4. **`docker-compose.yml`**
   - 1 行微调。
