package dashscopeasr

import (
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/relay"
	"github.com/QuantumNous/new-api/relay/channel"
)

func init() {
	// 覆盖 APITypeAli 的 Adaptor 工厂：对外仍是「阿里渠道」，内部组合 dashscopeasr 能力。
	relay.RegisterAdaptorOverride(constant.APITypeAli, func() channel.Adaptor {
		return NewAdaptor()
	})
}
