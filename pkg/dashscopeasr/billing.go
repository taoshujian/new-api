package dashscopeasr

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

const otherRatioKeySeconds = "seconds"

// AdjustPreConsume 按音频时长估算预扣费倍率（model_price 视为每秒单价，与视频 Task 的 seconds 倍率一致）。
func AdjustPreConsume(c *gin.Context, info *relaycommon.RelayInfo) error {
	if info == nil || !ShouldHandle(info) || !info.PriceData.UsePrice {
		return nil
	}

	seconds, err := estimateAudioSeconds(c)
	if err != nil {
		return fmt.Errorf("dashscope asr 预扣费估算时长失败: %w", err)
	}
	seconds = normalizeBillingSeconds(seconds)

	info.PriceData.AddOtherRatio(otherRatioKeySeconds, float64(seconds))
	applyOtherRatiosToPreConsume(&info.PriceData)

	logger.LogDebug(c, "[dashscopeasr] 预扣费估算 seconds=%d quotaToPreConsume=%d modelPrice=%.8f",
		seconds, info.PriceData.QuotaToPreConsume, info.PriceData.ModelPrice)
	return nil
}

// ApplySettleDuration 用上游 usage.duration 更新结算倍率（覆盖预扣时的估算值）。
func ApplySettleDuration(info *relaycommon.RelayInfo, durationSec int) {
	if info == nil || !info.PriceData.UsePrice {
		return
	}
	seconds := normalizeBillingSeconds(durationSec)
	info.PriceData.AddOtherRatio(otherRatioKeySeconds, float64(seconds))
}

// BuildUsageForBilling 按计费模式构造 Usage：按价走 Text 结算链路，按倍率走 Audio 结算链路。
func BuildUsageForBilling(info *relaycommon.RelayInfo, durationSec int, fileSize int64) *dto.Usage {
	if info != nil && info.PriceData.UsePrice {
		ApplySettleDuration(info, durationSec)
		return buildSettleUsage()
	}
	return BuildUsageFromDurationPrecise(durationSec, fileSize)
}

func buildSettleUsage() *dto.Usage {
	// 不含 AudioTokens，使 AudioHelper 走 PostTextConsumeQuota（已支持 OtherRatios × 按次单价）。
	return &dto.Usage{
		PromptTokens: 1,
		TotalTokens:  1,
	}
}

func normalizeBillingSeconds(seconds int) int {
	if seconds <= 0 {
		return 1
	}
	return seconds
}

func applyOtherRatiosToPreConsume(priceData *types.PriceData) {
	if priceData == nil || len(priceData.OtherRatios) == 0 {
		return
	}
	for _, ratio := range priceData.OtherRatios {
		if ratio != 1.0 && ratio > 0 {
			priceData.QuotaToPreConsume = int(float64(priceData.QuotaToPreConsume) * ratio)
		}
	}
}

func estimateAudioSeconds(c *gin.Context) (int, error) {
	multiForm, err := common.ParseMultipartFormReusable(c)
	if err != nil {
		return 0, fmt.Errorf("解析 multipart 表单失败: %w", err)
	}
	fileHeaders := multiForm.File["file"]
	if len(fileHeaders) == 0 {
		return 1, nil
	}

	totalSeconds := 0
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			return 0, fmt.Errorf("打开音频文件失败: %w", err)
		}
		ext := filepath.Ext(fileHeader.Filename)
		duration, err := common.GetAudioDuration(c.Request.Context(), file, ext)
		_ = file.Close()
		if err != nil {
			return 0, fmt.Errorf("读取音频时长失败: %w", err)
		}
		totalSeconds += int(math.Ceil(duration))
	}
	return totalSeconds, nil
}
