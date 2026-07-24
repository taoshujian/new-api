# 阿里云 DashScope 语音（ASR）模型支持

- **日期**：2026-06-22
- **提交**：`d9b2cac3` feat: 新增 pkg/dashscopeasr/*。支持阿里云语音模型
- **类型**：核心功能新增

## 背景

为支持阿里云百炼 DashScope 的 Fun-ASR Realtime 系列语音识别模型，以 OpenAI 兼容的音频转写（audio transcription）接口对外提供服务。整体设计原则是与上游代码解耦，便于 fork 后续合并 upstream。

## 支持模型

- `fun-asr-realtime`
- `fun-asr-realtime-2026-02-28`
- `fun-asr-realtime-2025-11-07`
- 同前缀的其他快照版本可通过渠道配置手动添加（`Supports()` 按 `fun-asr-realtime` 前缀匹配）。

## 变更内容（13 个文件，+686 行）

### 新增 `pkg/dashscopeasr/` 包

| 文件 | 说明 |
| --- | --- |
| `models.go` | 内置模型列表 `DefaultModels()`、模型匹配 `Supports()`（前缀匹配） |
| `router.go` | `ShouldHandle()` 判断是否接管请求：`RelayModeAudioTranscription` + 模型前缀命中 |
| `dto.go` | 请求/响应 DTO 定义（80 行） |
| `request.go` | 构造 DashScope 上游请求（166 行） |
| `response.go` | 解析百炼同步响应并写回 OpenAI 兼容格式（97 行） |
| `url.go` | 上游 URL 构建（54 行） |
| `models_test.go` | 模型匹配测试（39 行） |

### 新增渠道适配器 `relay/channel/dashscopeasr/`

- `adaptor.go`（138 行）：实现 channel.Adaptor 接口，处理请求转换与响应。
- `register.go`（14 行）：`init()` 中注册渠道。

### 新增 Adaptor 覆盖注册机制 `relay/adaptor_registry.go`（37 行）

- `RegisterAdaptorOverride(apiType, factory)`：fork 扩展可在不修改原渠道包内部逻辑的情况下，通过组合适配器注入新能力。
- 每次请求通过工厂返回新实例，避免可变状态在并发请求间泄漏。

### 默认倍率配置 `setting/ratio_setting/dashscope_asr.go`（15 行）

- `dashscopeASRModelRatio`：Fun-ASR Realtime 默认倍率 15（按音频 token 计量，约对齐 whisper 量级），管理员可在后台模型倍率中覆盖。
- 通过 `init()` 注入 `defaultModelRatio`，与 upstream `model_ratio.go` 文件级解耦。

### 其他

- `main.go`：导入注册 dashscopeasr 渠道包。
- `relay/relay_adaptor.go`：接入适配器选择逻辑（+3 行）。
