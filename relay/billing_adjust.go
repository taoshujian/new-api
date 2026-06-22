package relay

import (
	relaycommon "github.com/QuantumNous/new-api/relay/common"

	"github.com/gin-gonic/gin"
)

// PreConsumeAdjuster 在 ModelPriceHelper 之后、PreConsumeBilling 之前调整预扣额度。
// fork 扩展可通过 RegisterPreConsumeAdjuster 注入，无需修改 controller 业务逻辑。
type PreConsumeAdjuster func(c *gin.Context, info *relaycommon.RelayInfo) error

var preConsumeAdjusters []PreConsumeAdjuster

// RegisterPreConsumeAdjuster 注册预扣费调整器（可多次注册）。
func RegisterPreConsumeAdjuster(adjuster PreConsumeAdjuster) {
	if adjuster == nil {
		return
	}
	preConsumeAdjusters = append(preConsumeAdjusters, adjuster)
}

// ApplyPreConsumeAdjusters 依次执行已注册的预扣费调整器。
func ApplyPreConsumeAdjusters(c *gin.Context, info *relaycommon.RelayInfo) error {
	for _, adjuster := range preConsumeAdjusters {
		if err := adjuster(c, info); err != nil {
			return err
		}
	}
	return nil
}
