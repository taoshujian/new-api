package ratio_setting

// dashscopeASRModelRatio Fun-ASR 未配置 model_price 时的 fallback 倍率（按 audio token 计量）。
// 配置了 model_price 时按「每秒单价 × usage.duration」计费，见 pkg/dashscopeasr/billing.go。
// 本文件与 upstream model_ratio.go 解耦，便于 fork 合并。
var dashscopeASRModelRatio = map[string]float64{
	"fun-asr-realtime":            15,
	"fun-asr-realtime-2026-02-28": 15,
	"fun-asr-realtime-2025-11-07": 15,
}

func init() {
	for model, ratio := range dashscopeASRModelRatio {
		defaultModelRatio[model] = ratio
	}
}
