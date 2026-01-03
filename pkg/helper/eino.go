package helper

import (
	"strings"
	"vc-go/pkg/llm"

	"github.com/cloudwego/eino/schema"
)

func ParseEinoContent(message *schema.Message) string {
	if message == nil {
		return ""
	}
	return message.Content
}

func ParseEinoContentWithRemoveThink(message *schema.Message) string {
	return llm.RemoveThink(ParseEinoContent(message))
}

func ParseEinoContentWithRemoveThinkAndJSON(message *schema.Message) string {
	jsonString := strings.TrimSpace(ParseEinoContentWithRemoveThink(message))
	jsonString = strings.TrimPrefix(jsonString, "```json")
	jsonString = strings.TrimSuffix(jsonString, "```")
	return jsonString
}
