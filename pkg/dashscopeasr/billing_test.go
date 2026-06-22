package dashscopeasr_test

import (
	"testing"

	"github.com/QuantumNous/new-api/pkg/dashscopeasr"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildUsageForBillingUsePrice(t *testing.T) {
	info := &relaycommon.RelayInfo{}
	info.PriceData = types.PriceData{
		UsePrice:   true,
		ModelPrice: 0.000045,
	}

	usage := dashscopeasr.BuildUsageForBilling(info, 5, 0)
	require.NotNil(t, usage)
	assert.Equal(t, 1, usage.TotalTokens)
	assert.Equal(t, 0, usage.PromptTokensDetails.AudioTokens)
	assert.InDelta(t, 5.0, info.PriceData.OtherRatios["seconds"], 0.001)
}

func TestBuildUsageForBillingUseRatio(t *testing.T) {
	info := &relaycommon.RelayInfo{}
	info.PriceData = types.PriceData{UsePrice: false}

	usage := dashscopeasr.BuildUsageForBilling(info, 60, 0)
	require.NotNil(t, usage)
	assert.Equal(t, 1000, usage.PromptTokensDetails.AudioTokens)
	assert.Empty(t, info.PriceData.OtherRatios)
}

func TestAdjustPreConsumeSkipsNonASR(t *testing.T) {
	info := &relaycommon.RelayInfo{OriginModelName: "whisper-1"}
	info.PriceData = types.PriceData{
		UsePrice:          true,
		QuotaToPreConsume: 100,
	}

	err := dashscopeasr.AdjustPreConsume(nil, info)
	require.NoError(t, err)
	assert.Equal(t, 100, info.PriceData.QuotaToPreConsume)
	assert.Empty(t, info.PriceData.OtherRatios)
}
