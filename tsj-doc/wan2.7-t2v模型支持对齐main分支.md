# wan2.7-t2v 模型支持对齐 main 分支

- **日期**：2026-07-24
- **类型**：代码对齐（减少与 upstream 的冲突面）
- **关联文档**：[wan2.7-t2v模型计费支持.md](wan2.7-t2v模型计费支持.md)

## 背景

tsj 分支在 `858da443` 中自行添加了 wan2.7-t2v 模型支持；此后 main 分支通过 `52858ad1 feat: support Wan2.7 i2v media mapping (#4984)` 也加入了该模型（并完整支持 wan2.7-i2v 新协议）。模型支持部分以 main 为准，避免重复实现和合并冲突。

## 变更内容

1. **`relay/channel/task/ali/constants.go`**
   - 还原为 main 分支版本（差异仅为一处注释空格：`万相2.7 文生视频` → `万相2.7文生视频`）。

2. **`relay/channel/task/ali/adaptor.go`**（未改动）
   - 模型支持部分本就与 main 一致；保留 tsj 独有的计费配置——`ProcessAliOtherRatios` 中 wan2.7-t2v 的分辨率倍率表条目（720P=1，1080P=1/0.6）。

## 对齐后的差异面

tsj 相对 main 在阿里视频相关代码中仅剩计费配置差异：

- `setting/ratio_setting/model_ratio.go`：`defaultModelPrice` 中 `wan2.7-t2v: 0.6`（720P 每秒 0.6 元基准单价）
- `relay/channel/task/ali/adaptor.go`：`ProcessAliOtherRatios` 分辨率倍率表条目

已执行 `go build ./relay/channel/task/ali/ ./setting/ratio_setting/` 验证编译通过。
