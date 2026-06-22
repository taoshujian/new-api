package dashscopeasr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/logger"
	"github.com/gin-gonic/gin"
)

// HandleSyncResponse 解析百炼同步响应并写回 OpenAI 兼容格式。
func HandleSyncResponse(c *gin.Context, statusCode int, body []byte, responseFormat string) (*dto.Usage, error) {
	logger.LogDebug(c, "[dashscopeasr] 上游响应 status=%d bodyBytes=%d", statusCode, len(body))

	if statusCode != http.StatusOK {
		return nil, parseUpstreamError(c, statusCode, body)
	}

	var syncResp SyncResponse
	if err := common.Unmarshal(body, &syncResp); err != nil {
		return nil, fmt.Errorf("解析 DashScope ASR 响应失败: %w, body=%s", err, truncateForLog(body, 512))
	}

	if syncResp.Code != "" || syncResp.Message != "" {
		logger.LogError(c, fmt.Sprintf("[dashscopeasr] 上游业务错误 code=%s message=%s requestId=%s",
			syncResp.Code, syncResp.Message, syncResp.RequestID))
		return nil, fmt.Errorf("DashScope ASR error: %s - %s (request_id=%s)", syncResp.Code, syncResp.Message, syncResp.RequestID)
	}

	text := strings.TrimSpace(syncResp.Output.Text)
	if text == "" && syncResp.Output.Sentence != nil {
		text = strings.TrimSpace(syncResp.Output.Sentence.Text)
	}

	logger.LogInfo(c, fmt.Sprintf("[dashscopeasr] 识别完成 requestId=%s durationSec=%d textLen=%d preview=%q",
		syncResp.RequestID, syncResp.Usage.Duration, len(text), truncateForLog([]byte(text), 120)))

	responseBytes, contentType, err := buildClientResponse(text, syncResp, responseFormat)
	if err != nil {
		return nil, err
	}

	c.Data(http.StatusOK, contentType, responseBytes)
	usage := BuildUsageFromDuration(syncResp.Usage.Duration)
	return usage, nil
}

func buildClientResponse(text string, syncResp SyncResponse, responseFormat string) ([]byte, string, error) {
	switch strings.ToLower(strings.TrimSpace(responseFormat)) {
	case "text":
		return []byte(text), "text/plain; charset=utf-8", nil
	case "verbose_json":
		resp := dto.WhisperVerboseJSONResponse{
			Task: "transcribe",
			Text: text,
		}
		if syncResp.Output.Sentence != nil {
			s := syncResp.Output.Sentence
			resp.Duration = float64(s.EndTime) / 1000.0
			resp.Segments = []dto.Segment{
				{
					Id:    s.SentenceID,
					Start: float64(s.BeginTime) / 1000.0,
					End:   float64(s.EndTime) / 1000.0,
					Text:  s.Text,
				},
			}
		}
		payload, err := common.Marshal(resp)
		return payload, "application/json", err
	default:
		payload, err := common.Marshal(dto.AudioResponse{Text: text})
		return payload, "application/json", err
	}
}

func parseUpstreamError(c *gin.Context, statusCode int, body []byte) error {
	var errResp ErrorResponse
	if err := common.Unmarshal(body, &errResp); err == nil && (errResp.Message != "" || errResp.Code != "") {
		logger.LogError(c, fmt.Sprintf("[dashscopeasr] 上游 HTTP 错误 status=%d code=%s message=%s requestId=%s",
			statusCode, errResp.Code, errResp.Message, errResp.RequestID))
		return fmt.Errorf("DashScope ASR upstream error (%d): %s - %s (request_id=%s)",
			statusCode, errResp.Code, errResp.Message, errResp.RequestID)
	}
	logger.LogError(c, fmt.Sprintf("[dashscopeasr] 上游 HTTP 错误 status=%d body=%s", statusCode, truncateForLog(body, 1024)))
	return fmt.Errorf("DashScope ASR upstream error (%d): %s", statusCode, truncateForLog(body, 256))
}

func truncateForLog(data []byte, max int) string {
	if len(data) <= max {
		return string(data)
	}
	return string(data[:max]) + "...(truncated)"
}
