package dashscopeasr

import (
	"fmt"
	"path/filepath"
	"strings"
)

var formatByExtension = map[string]string{
	".wav":  "wav",
	".mp3":  "mp3",
	".mpeg": "mp3",
	".opus": "opus",
	".ogg":  "opus",
	".flac": "flac",
	".m4a":  "mp3",
	".aac":  "mp3",
}

var mimeByFormat = map[string]string{
	"wav":  "audio/wav",
	"mp3":  "audio/mpeg",
	"opus": "audio/opus",
	"flac": "audio/flac",
}

// DetectFormat 根据文件名推断百炼 format 参数。
func DetectFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if format, ok := formatByExtension[ext]; ok {
		return format
	}
	return "wav"
}

// MimeType 返回 format 对应的 Data URI MIME 类型。
func MimeType(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	if mime, ok := mimeByFormat[format]; ok {
		return mime
	}
	return "audio/wav"
}

// BuildDataURI 构造百炼 Base64 Data URI：data:{mime};base64,{data}
func BuildDataURI(format string, base64Data string) string {
	return fmt.Sprintf("data:%s;base64,%s", MimeType(format), base64Data)
}

// BuildRequestURL 拼接完整上游 URL。
func BuildRequestURL(baseURL string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	return base + SyncGenerationPath
}
