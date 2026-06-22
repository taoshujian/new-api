package relay

import (
	"sync"

	"github.com/QuantumNous/new-api/relay/channel"
)

// AdaptorFactory 用于按需创建渠道 Adaptor 实例。
// 每次请求应返回新实例，避免可变状态在并发请求间泄漏。
type AdaptorFactory func() channel.Adaptor

var (
	adaptorOverrideMu sync.RWMutex
	adaptorOverrides  = map[int]AdaptorFactory{}
)

// RegisterAdaptorOverride 注册渠道 Adaptor 覆盖工厂。
// 典型用途：fork 扩展在不动原渠道包内部逻辑的情况下，通过组合适配器注入新能力。
func RegisterAdaptorOverride(apiType int, factory AdaptorFactory) {
	if factory == nil {
		return
	}
	adaptorOverrideMu.Lock()
	defer adaptorOverrideMu.Unlock()
	adaptorOverrides[apiType] = factory
}

func getAdaptorOverride(apiType int) channel.Adaptor {
	adaptorOverrideMu.RLock()
	factory, ok := adaptorOverrides[apiType]
	adaptorOverrideMu.RUnlock()
	if !ok || factory == nil {
		return nil
	}
	return factory()
}
