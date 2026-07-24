# 替换 favicon 图标与微调部署配置

- **日期**：2026-06-17
- **提交**：`dd443fca` chore: 修改.ico
- **类型**：资源与配置微调

## 变更内容

1. **`web/classic/public/favicon.ico`**
   - 替换站点图标（15 KB → 2.9 MB）。

2. **`docker-compose.alone.yml`**
   - 镜像名由 `new-api:alone-local` 改为 `new-api:alone`。

3. **`docker-compose.prod.yml`**
   - 生产环境连接宿主机 MySQL/Redis 的端口由 3307/6380 改为 3306/6379。

4. **`web/classic/src/i18n/locales/zh-CN.json`**
   - 两条中文文案微调：`版权所有` → `版权-所有`、`设计与开发由` → `设计与开发由：`。
