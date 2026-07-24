# Fun-ASR Realtime 计费修正：单次价格 × 秒数

- **日期**：2026-06-22
- **提交**：`2b696d36` fix: 修改 DashScope Fun-ASR Realtime 系列 模型计费，改为 单次价格x秒数
- **类型**：计费逻辑修正

## 背景

上一版本中 Fun-ASR 按音频 token 倍率计费。修正为：当模型配置了 `model_price` 时，按 **每秒单价 × usage.duration（秒）** 计费，与视频 Task 的 `seconds` 倍率机制保持一致；未配置 `model_price` 时仍 fallback 到 audio token 倍率。

## 变更内容（10 个文件，+208 / -7）

### 新增 `pkg/dashscopeasr/billing.go`（108 行）

- `AdjustPreConsume(c, info)`：预扣费调整器。仅当 `ShouldHandle(info)` 且 `UsePrice` 时生效；解析 multipart 表单中的音频文件，用 `common.GetAudioDuration` 估算时长，`AddOtherRatio("seconds", seconds)` 并将倍率乘入 `QuotaToPreConsume`。
- `ApplySettleDuration(info, durationSec)`：结算时用上游返回的 `usage.duration` 覆盖预扣时的估算值。
- `BuildUsageForBilling(info, durationSec, fileSize)`：按计费模式构造 Usage——按价（UsePrice）走 Text 结算链路（构造不含 AudioTokens 的 Usage，使 `AudioHelper` 走 `PostTextConsumeQuota`，支持 OtherRatios × 按次单价）；按倍率走 Audio 结算链路。
- `normalizeBillingSeconds`：时长 <= 0 时归一为 1 秒。

### 新增预扣费调整器注册机制 `relay/billing_adjust.go`（31 行）

- `RegisterPreConsumeAdjuster(adjuster)` / `ApplyPreConsumeAdjusters(c, info)`：在 `ModelPriceHelper` 之后、`PreConsumeBilling` 之前执行已注册的调整器。fork 扩展通过注册注入，无需修改 controller 业务逻辑。

### 接入点

- `controller/relay.go`：在定价之后、预扣费之前调用 `relay.ApplyPreConsumeAdjusters`，出错返回 400 `ErrorCodeModelPriceError`。
- `relay/channel/dashscopeasr/register.go`：`init()` 中注册 `relay.RegisterPreConsumeAdjuster(asr.AdjustPreConsume)`。
- `pkg/dashscopeasr/response.go`：`HandleSyncResponse` 签名增加 `info` 与 `fileSize` 参数，改用 `BuildUsageForBilling` 构造结算 Usage，按秒结算时输出日志。
- `relay/channel/dashscopeasr/adaptor.go`：适配新签名，传入 `info.UpstreamRequestBodySize`。

### 配置与文档注释

- `setting/ratio_setting/dashscope_asr.go`：注释更新，说明 `dashscopeASRModelRatio` 仅为未配置 `model_price` 时的 fallback。
- `main.go`：渠道包导入路径修正为 `relay/channel/dashscopeasr`。
- `docker-compose.alone.yml`：镜像 tag 更新为 `new-api:alone-0.0.1`。

### 测试 `pkg/dashscopeasr/billing_test.go`（48 行）

- `TestBuildUsageForBillingUsePrice`：UsePrice 时 TotalTokens=1、无 AudioTokens、`OtherRatios["seconds"]` 为时长。
- `TestBuildUsageForBillingUseRatio`：倍率模式 60 秒 → 1000 AudioTokens，不写 OtherRatios。
- `TestAdjustPreConsumeSkipsNonASR`：非 ASR 模型（whisper-1）跳过调整，预扣额度不变。
