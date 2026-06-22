package dashscopeasr

import "strings"

// modelPrefix 是 Fun-ASR Realtime 系列模型的统一前缀。
const modelPrefix = "fun-asr-realtime"

// DefaultModels 返回内置支持的 DashScope Fun-ASR Realtime 模型名列表。
// 渠道配置中也可手动添加同前缀的其他快照版本。
func DefaultModels() []string {
	return []string{
		"fun-asr-realtime",
		"fun-asr-realtime-2026-02-28",
		"fun-asr-realtime-2025-11-07",
	}
}

// Supports 判断模型名是否应由 dashscopeasr 模块处理。
func Supports(modelName string) bool {
	normalized := strings.ToLower(strings.TrimSpace(modelName))
	if normalized == "" {
		return false
	}
	return strings.HasPrefix(normalized, modelPrefix)
}
