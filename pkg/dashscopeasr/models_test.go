package dashscopeasr_test

import (
	"testing"

	"github.com/QuantumNous/new-api/pkg/dashscopeasr"
	"github.com/stretchr/testify/assert"
)

func TestSupports(t *testing.T) {
	tests := []struct {
		model    string
		expected bool
	}{
		{"fun-asr-realtime-2026-02-28", true},
		{"fun-asr-realtime", true},
		{"FUN-ASR-REALTIME-2026-02-28", true},
		{"fun-asr", false},
		{"whisper-1", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			assert.Equal(t, tt.expected, dashscopeasr.Supports(tt.model))
		})
	}
}

func TestDetectFormat(t *testing.T) {
	assert.Equal(t, "wav", dashscopeasr.DetectFormat("sample.wav"))
	assert.Equal(t, "mp3", dashscopeasr.DetectFormat("sample.mp3"))
	assert.Equal(t, "wav", dashscopeasr.DetectFormat("unknown.bin"))
}

func TestBuildUsageFromDuration(t *testing.T) {
	usage := dashscopeasr.BuildUsageFromDuration(60)
	assert.Equal(t, 1000, usage.PromptTokensDetails.AudioTokens)
	assert.Equal(t, 1000, usage.TotalTokens)
}
