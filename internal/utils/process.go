package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)


type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ResponseBody struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CallAPIWithRetry(apiKey, apiURL, content string) (string, error) {
	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := CallAPI(apiKey, apiURL, content)
		if err != nil {
			// check if the error is due to rate limit
			if strings.Contains(err.Error(), "rate_limit_exceeded") {
				wait := time.Duration(attempt*10) * time.Second
				fmt.Printf("⚠️  Rate limit exceeded. Waiting %v before retrying...\n", wait)
				time.Sleep(wait)
				continue
			}
			return "", err
		}
		return result, nil
	}
	return "", fmt.Errorf("failed after %d attempts (rate limit)", maxRetries)
}

func CallAPI(apiKey, apiURL, content string) (string, error) {
	body := RequestBody{
		Model: "openai/gpt-oss-20b",
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to encode JSON: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result ResponseBody
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w\nResponse: %s", err, string(respBody))
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("response does not contain any content. Full response: %s", string(respBody))
}

func SplitContent(content string, chunkSize int) []string {
	var chunks []string
	for start := 0; start < len(content); start += chunkSize {
		end := start + chunkSize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[start:end])
	}
	return chunks
}
