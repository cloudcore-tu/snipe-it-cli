package snipeit_test

import (
	"testing"

	"github.com/cloudcore-tu/snipe-it-cli/internal/snipeit"
	"github.com/stretchr/testify/assert"
)

// --- MaskAuthHeader (セキュリティテスト) ---

func TestMaskAuthHeader_Bearer(t *testing.T) {
	// Bearer トークンは *** に置き換えられること
	masked := snipeit.MaskAuthHeader("Bearer secret-api-token-123")
	assert.Equal(t, "Bearer ***", masked)
	assert.NotContains(t, masked, "secret")
}

func TestMaskAuthHeader_Empty(t *testing.T) {
	assert.Equal(t, "", snipeit.MaskAuthHeader(""))
}

func TestMaskAuthHeader_NonBearer(t *testing.T) {
	// Bearer 以外も REDACTED になること
	masked := snipeit.MaskAuthHeader("Token abc123")
	assert.Equal(t, "***REDACTED***", masked)
	assert.NotContains(t, masked, "abc123")
}

func TestMaskAuthHeader_BearerOnly(t *testing.T) {
	// "Bearer " のみでトークンなしの場合
	masked := snipeit.MaskAuthHeader("Bearer ")
	assert.Equal(t, "Bearer ***", masked)
}
