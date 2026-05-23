package ai_client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const promptCacheKeyPrefix = "ai:prompt-cache:v1:"

func promptCacheKey(req completionRequest) string {
	payload := fmt.Sprintf("%s\x1f%s\x1f%s\x1f%d", req.model, req.systemPrompt, req.userPrompt, req.maxTokens)
	hash := sha256.Sum256([]byte(payload))
	return promptCacheKeyPrefix + hex.EncodeToString(hash[:])
}
