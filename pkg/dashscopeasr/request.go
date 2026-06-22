package dashscopeasr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/gin-gonic/gin"
)

const (
	contextKeyResponseFormat = "dashscope_asr_response_format"
)

// ContextKeyResponseFormat 供组合 Adaptor 读取 response_format。
func ContextKeyResponseFormat() string {
	return contextKeyResponseFormat
}

// BuildSyncRequest 将 OpenAI /v1/audio/transcriptions 表单请求转换为百炼同步 JSON 请求体。
func BuildSyncRequest(c *gin.Context, modelName string, request dto.AudioRequest) ([]byte, error) {
	formData, err := common.ParseMultipartFormReusable(c)
	if err != nil {
		return nil, fmt.Errorf("解析 multipart 表单失败: %w", err)
	}

	fileHeaders := formData.File["file"]
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("缺少必填字段 file")
	}
	fileHeader := fileHeaders[0]

	logger.LogInfo(c, fmt.Sprintf("[dashscopeasr] 收到识别请求 model=%s file=%s size=%d contentType=%s",
		modelName, fileHeader.Filename, fileHeader.Size, fileHeader.Header.Get("Content-Type")))

	if fileHeader.Size > MaxEncodedAudioBytes {
		return nil, fmt.Errorf("音频文件过大（原始 %d bytes），请控制在约 7MB 以内以避免 Base64 超限", fileHeader.Size)
	}

	audioBytes, err := readUploadFile(fileHeader)
	if err != nil {
		return nil, err
	}

	format := DetectFormat(fileHeader.Filename)
	sampleRate := "16000"
	var vadEnabled *bool

	// 允许通过 metadata 覆盖百炼原生 parameters
	if len(request.Metadata) > 0 {
		var override SyncParameters
		if err := json.Unmarshal(request.Metadata, &override); err != nil {
			return nil, fmt.Errorf("metadata 无法解析为 DashScope ASR parameters: %w", err)
		}
		if strings.TrimSpace(override.Format) != "" {
			format = override.Format
		}
		if strings.TrimSpace(override.SampleRate) != "" {
			sampleRate = override.SampleRate
		}
		if override.VADEnabled != nil {
			vadEnabled = override.VADEnabled
		}
	}

	encoded := base64.StdEncoding.EncodeToString(audioBytes)
	if len(encoded) > MaxEncodedAudioBytes {
		return nil, fmt.Errorf("Base64 编码后音频超过 %d bytes 上限", MaxEncodedAudioBytes)
	}

	syncReq := SyncRequest{
		Model: modelName,
		Input: SyncInput{
			Messages: []SyncMessage{
				{
					Role: "user",
					Content: []SyncContent{
						{Audio: BuildDataURI(format, encoded)},
					},
				},
			},
		},
		Parameters: SyncParameters{
			Format:     format,
			SampleRate: sampleRate,
			VADEnabled: vadEnabled,
		},
		Resources: []map[string]any{},
	}

	payload, err := common.Marshal(syncReq)
	if err != nil {
		return nil, fmt.Errorf("序列化 DashScope ASR 请求失败: %w", err)
	}

	responseFormat := strings.TrimSpace(request.ResponseFormat)
	if responseFormat == "" {
		responseFormat = "json"
	}
	c.Set(contextKeyResponseFormat, responseFormat)

	logger.LogDebug(c, "[dashscopeasr] 上游请求已构建 model=%s format=%s sampleRate=%s payloadBytes=%d responseFormat=%s",
		modelName, format, sampleRate, len(payload), responseFormat)

	return payload, nil
}

func readUploadFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer file.Close()

	audioBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件失败: %w", err)
	}
	if len(audioBytes) == 0 {
		return nil, fmt.Errorf("上传文件为空")
	}
	return audioBytes, nil
}

// BuildUsageFromDuration 将百炼 usage.duration（秒）映射为 OpenAI Usage，便于走音频计费链路。
func BuildUsageFromDuration(durationSec int) *dto.Usage {
	audioTokens := durationSec * 1000 / 60
	if audioTokens <= 0 && durationSec > 0 {
		audioTokens = 1
	}
	if durationSec <= 0 {
		audioTokens = 1
	}

	usage := &dto.Usage{
		PromptTokens:     audioTokens,
		CompletionTokens: 0,
		TotalTokens:      audioTokens,
	}
	usage.PromptTokensDetails.AudioTokens = audioTokens
	return usage
}

// BuildUsageFromDurationPrecise 在 duration 为 0 时回退到文件大小估算（极端容错）。
func BuildUsageFromDurationPrecise(durationSec int, fileSize int64) *dto.Usage {
	if durationSec > 0 {
		return BuildUsageFromDuration(durationSec)
	}
	if fileSize <= 0 {
		return BuildUsageFromDuration(1)
	}
	// 16kHz mono 16bit ≈ 32000 bytes/s
	estimatedSec := int(math.Ceil(float64(fileSize) / 32000.0))
	if estimatedSec <= 0 {
		estimatedSec = 1
	}
	logger.LogDebug(nil, "[dashscopeasr] 上游未返回 duration，按文件大小估算计费秒数=%d fileSize=%d", estimatedSec, fileSize)
	return BuildUsageFromDuration(estimatedSec)
}
