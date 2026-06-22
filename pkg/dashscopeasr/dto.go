package dashscopeasr

// SyncGenerationPath 百炼 Fun-ASR-Realtime HTTP 同步识别路径（与万相同步图像共用 generation 端点）。
const SyncGenerationPath = "/api/v1/services/aigc/multimodal-generation/generation"

// MaxEncodedAudioBytes Base64 编码后音频上限（百炼文档建议控制编码后体积）。
const MaxEncodedAudioBytes = 10 * 1024 * 1024

// SyncRequest Fun-ASR-Realtime 同步 HTTP 请求体。
// 文档：https://help.aliyun.com/zh/model-studio/fun-asr-recorded-speech-recognition-http-api
type SyncRequest struct {
	Model      string           `json:"model"`
	Input      SyncInput        `json:"input"`
	Parameters SyncParameters   `json:"parameters"`
	Resources  []map[string]any `json:"resources,omitempty"`
}

type SyncInput struct {
	Messages []SyncMessage `json:"messages"`
}

type SyncMessage struct {
	Role    string        `json:"role"`
	Content []SyncContent `json:"content"`
}

type SyncContent struct {
	Audio string `json:"audio,omitempty"`
}

// SyncParameters 识别参数；format 必填，其余可选。
type SyncParameters struct {
	Format       string `json:"format"`
	SampleRate   string `json:"sample_rate,omitempty"`
	VADEnabled   *bool  `json:"vad_enabled,omitempty"`
	AudioAddress string `json:"audio_address,omitempty"`
}

// SyncResponse 百炼同步识别响应。
type SyncResponse struct {
	Output    SyncOutput `json:"output"`
	Usage     SyncUsage  `json:"usage"`
	RequestID string     `json:"request_id"`
	Code      string     `json:"code,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type SyncOutput struct {
	Text     string        `json:"text"`
	Sentence *SyncSentence `json:"sentence,omitempty"`
}

type SyncSentence struct {
	SentenceID  int        `json:"sentence_id"`
	SentenceEnd bool       `json:"sentence_end"`
	BeginTime   int        `json:"begin_time"`
	EndTime     int        `json:"end_time"`
	Text        string     `json:"text"`
	ChannelID   int        `json:"channel_id"`
	Words       []SyncWord `json:"words,omitempty"`
}

type SyncWord struct {
	Text        string `json:"text"`
	BeginTime   int    `json:"begin_time"`
	EndTime     int    `json:"end_time"`
	Punctuation string `json:"punctuation"`
	Fixed       bool   `json:"fixed"`
}

type SyncUsage struct {
	Duration int `json:"duration"`
}

// ErrorResponse 百炼通用错误响应。
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}
