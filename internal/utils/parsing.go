package utils

import (
	"encoding/json"
	"fmt"
)
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func ParsingContent(jsonData []byte) (string, error) {
	var resp ChatResponse
	if err := json.Unmarshal(jsonData, &resp); err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices found")
	}
	return resp.Choices[0].Message.Content, nil
}