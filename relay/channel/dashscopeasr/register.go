package dashscopeasr

import (
	"github.com/QuantumNous/new-api/constant"
	asr "github.com/QuantumNous/new-api/pkg/dashscopeasr"
	"github.com/QuantumNous/new-api/relay"
	"github.com/QuantumNous/new-api/relay/channel"
)

func init() {
	// 覆盖 APITypeAli 的 Adaptor 工厂：对外仍是「阿里渠道」，内部组合 dashscopeasr 能力。
	relay.RegisterAdaptorOverride(constant.APITypeAli, func() channel.Adaptor {
		return NewAdaptor()
	})
	// 按秒预扣费：model_price 为每秒单价，OtherRatios["seconds"] 与视频 Task 一致。
	relay.RegisterPreConsumeAdjuster(asr.AdjustPreConsume)
}
