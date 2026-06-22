package ratio_setting

// dashscopeASRModelRatio Fun-ASR Realtime 默认倍率（按音频 token 计量，约对齐 whisper 量级）。
// 管理员可在后台模型倍率中覆盖；本文件与 upstream model_ratio.go 解耦，便于 fork 合并。
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
