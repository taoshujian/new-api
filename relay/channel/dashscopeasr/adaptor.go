package dashscopeasr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/QuantumNous/new-api/dto"
	asr "github.com/QuantumNous/new-api/pkg/dashscopeasr"
	"github.com/QuantumNous/new-api/relay/channel"
	"github.com/QuantumNous/new-api/relay/channel/ali"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

const channelName = "ali+dashscopeasr"

// Adaptor 组合适配器：Fun-ASR 走 pkg/dashscopeasr，其余能力委托给原生 ali.Adaptor。
// 设计目标：fork 扩展与 upstream ali 包解耦，合并 upstream 时冲突最小。
type Adaptor struct {
	inner ali.Adaptor
}

func NewAdaptor() *Adaptor {
	return &Adaptor{}
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
	a.inner.Init(info)
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	if asr.ShouldHandle(info) {
		return asr.BuildRequestURL(info.ChannelBaseUrl), nil
	}
	return a.inner.GetRequestURL(info)
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	if asr.ShouldHandle(info) {
		channel.SetupApiRequestHeader(info, c, req)
		req.Set("Authorization", "Bearer "+info.ApiKey)
		req.Set("Content-Type", "application/json")
		req.Set("Accept", "application/json")
		req.Set("X-DashScope-SSE", "disable")
		return nil
	}
	return a.inner.SetupRequestHeader(c, req, info)
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	return a.inner.ConvertOpenAIRequest(c, info, request)
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return a.inner.ConvertRerankRequest(c, relayMode, request)
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return a.inner.ConvertEmbeddingRequest(c, info, request)
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	if !asr.ShouldHandle(info) {
		return a.inner.ConvertAudioRequest(c, info, request)
	}

	payload, err := asr.BuildSyncRequest(c, info.UpstreamModelName, request)
	if err != nil {
		return nil, fmt.Errorf("dashscope asr convert request failed: %w", err)
	}
	info.UpstreamRequestBodySize = int64(len(payload))
	return bytes.NewReader(payload), nil
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	return a.inner.ConvertImageRequest(c, info, request)
}

func (a *Adaptor) ConvertOpenAIResponsesRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.OpenAIResponsesRequest) (any, error) {
	return a.inner.ConvertOpenAIResponsesRequest(c, info, request)
}

func (a *Adaptor) ConvertClaudeRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.ClaudeRequest) (any, error) {
	return a.inner.ConvertClaudeRequest(c, info, request)
}

func (a *Adaptor) ConvertGeminiRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeminiChatRequest) (any, error) {
	return a.inner.ConvertGeminiRequest(c, info, request)
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *types.NewAPIError) {
	if !asr.ShouldHandle(info) {
		return a.inner.DoResponse(c, resp, info)
	}

	defer service.CloseResponseBodyGracefully(resp)
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, types.NewErrorWithStatusCode(
			fmt.Errorf("读取 DashScope ASR 响应失败: %w", readErr),
			types.ErrorCodeReadResponseBodyFailed,
			http.StatusInternalServerError,
		)
	}

	responseFormat := c.GetString(asr.ContextKeyResponseFormat())
	if responseFormat == "" {
		responseFormat = "json"
	}

	usageDTO, handleErr := asr.HandleSyncResponse(c, info, resp.StatusCode, body, responseFormat, info.UpstreamRequestBodySize)
	if handleErr != nil {
		statusCode := http.StatusBadGateway
		if resp.StatusCode >= 400 {
			statusCode = resp.StatusCode
		}
		return nil, types.NewErrorWithStatusCode(handleErr, types.ErrorCodeBadResponse, statusCode)
	}
	return usageDTO, nil
}

func (a *Adaptor) GetModelList() []string {
	base := a.inner.GetModelList()
	return append(base, asr.DefaultModels()...)
}

func (a *Adaptor) GetChannelName() string {
	return channelName
}
