# wan2.7-t2v 模型计费支持

- **日期**：2026-06-18
- **提交**：`858da443` perf: 添加wan2.7-t2v模型计费支持
- **类型**：功能增强（计费）

## 背景

阿里云万相 2.7 文生视频模型 `wan2.7-t2v` 需要接入计费体系，按分辨率区分单价。

## 计费规则

- **720P**：0.6 元/秒（基准单价）
- **1080P**：1 元/秒（通过分辨率倍率 `1/0.6 ≈ 1.667` 换算）

## 变更内容

1. **`setting/ratio_setting/model_ratio.go`**
   - `defaultModelPrice` 新增 `"wan2.7-t2v": 0.6`（基础单价按 720P 每秒 0.6 元计）。

2. **`relay/channel/task/ali/adaptor.go`**
   - `ProcessAliOtherRatios` 分辨率倍率表新增 `wan2.7-t2v`：`720P → 1`、`1080P → 1/0.6`。

3. **`relay/channel/task/ali/constants.go`**
   - 支持模型列表新增 `wan2.7-t2v`（万相2.7 文生视频）。

4. **`.env`（新增，14 行）**
   - 本地 `go run main.go` 使用的环境配置：连接 `tsj-mysql`（宿主机 3307）与 `tsj-redis`（宿主机 6380），开启 `ERROR_LOG_ENABLED`、`BATCH_UPDATE_ENABLED`。
   - 注意：`.env` 包含本地密码，且本次同时从 `.gitignore` 中移除了 `.env` 条目使其可被提交，合并回上游前需评估。

5. **`web/default/src/i18n/locales/zh.json`**
   - 一条页脚文案末尾多了 `1`，疑似误改，建议后续清理。
