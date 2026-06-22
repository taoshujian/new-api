package dashscopeasr

import (
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
)

// ShouldHandle 判断当前请求是否应由 dashscopeasr 模块接管。
func ShouldHandle(info *relaycommon.RelayInfo) bool {
	if info == nil {
		return false
	}
	if info.RelayMode != relayconstant.RelayModeAudioTranscription {
		return false
	}
	return Supports(info.OriginModelName)
}
