package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/i18n"
	"github.com/songquanpeng/one-api/common/logger"
)

func mapPublicError(message string) (i18nKey string, code string, ok bool) {
	normalized := strings.ToLower(strings.TrimSpace(message))
	if normalized == "" {
		return "", "", false
	}

	if strings.Contains(normalized, "该令牌状态不可用") ||
		strings.Contains(normalized, "token status is unavailable") ||
		strings.Contains(normalized, "该令牌额度已用尽") ||
		strings.Contains(normalized, "额度已用尽") ||
		strings.Contains(normalized, "token quota exhausted") {
		return "billing_credits_exhausted", "billing_credits_exhausted", true
	}

	return "", "", false
}

func abortWithMessage(c *gin.Context, statusCode int, message string) {
	displayMessage := message
	errorPayload := gin.H{
		"type": "one_api_error",
	}
	if i18nKey, code, ok := mapPublicError(message); ok {
		displayMessage = i18n.Translate(c, i18nKey)
		errorPayload["code"] = code
	}
	errorPayload["message"] = helper.MessageWithRequestId(displayMessage, c.GetString(helper.RequestIdKey))

	c.JSON(statusCode, gin.H{
		"error": errorPayload,
	})
	c.Abort()
	logger.Error(c.Request.Context(), message)
}

func getRequestModel(c *gin.Context) (string, error) {
	var modelRequest ModelRequest
	err := common.UnmarshalBodyReusable(c, &modelRequest)
	if err != nil {
		return "", fmt.Errorf("common.UnmarshalBodyReusable failed: %w", err)
	}
	if strings.HasPrefix(c.Request.URL.Path, "/v1/moderations") {
		if modelRequest.Model == "" {
			modelRequest.Model = "text-moderation-stable"
		}
	}
	if strings.HasSuffix(c.Request.URL.Path, "embeddings") {
		if modelRequest.Model == "" {
			modelRequest.Model = c.Param("model")
		}
	}
	if strings.HasPrefix(c.Request.URL.Path, "/v1/images/generations") {
		if modelRequest.Model == "" {
			modelRequest.Model = "dall-e-2"
		}
	}
	if strings.HasPrefix(c.Request.URL.Path, "/v1/audio/transcriptions") || strings.HasPrefix(c.Request.URL.Path, "/v1/audio/translations") {
		if modelRequest.Model == "" {
			modelRequest.Model = "whisper-1"
		}
	}
	return modelRequest.Model, nil
}

func isModelInList(modelName string, models string) bool {
	modelList := strings.Split(models, ",")
	for _, model := range modelList {
		if modelName == model {
			return true
		}
	}
	return false
}
