package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type fn func(map[string]any) string

type ChatClient struct {
	client *openai.Client
	model  string
	funcs  map[string]fn
}

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("OPENAI_API_BASE")
	model := os.Getenv("OPENAI_API_MODEL")

	if apiKey == "" || baseURL == "" || model == "" {
		fmt.Println("检查环境变量设置")
		return
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	// 构建注册函数
	chatClient := &ChatClient{
		client: client,
		model:  model,
	}

	// 注册函数到chatClient
	chatClient.funcs = map[string]fn{
		"getWeather": getWeather,
		"getTime":    getTime,
	}

	chatClient.ChatLoop()
}

func (c *ChatClient) ChatLoop() {
	fmt.Print("Type your queries or 'quit' to exit.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nUser: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if strings.ToLower(userInput) == "quit" {
			break
		}
		if userInput == "" {
			continue
		}

		response, err := c.ProcessQuery(userInput)
		if err != nil {
			fmt.Printf("请求失败: %v\n", err)
			continue
		}

		fmt.Printf("Assistant: %s\n", response)
	}
}

func (c *ChatClient) ProcessQuery(userInput string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: []openai.ChatCompletionMessage{{Role: "user", Content: userInput}},
		Functions: []openai.FunctionDefinition{
			{
				Name:        "getWeather",
				Description: "Get weather for a given city",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"city": { "type": "string" }
					},
					"required": ["city"]
				}`),
			},
			{
				Name:        "getTime",
				Description: "Get current time for a given city",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"city": { "type": "string" }
					},
					"required": ["city"]
				}`),
			},
		},
		FunctionCall: "auto",
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	choice := resp.Choices[0].Message
	if choice.FunctionCall != nil {
		fnName := choice.FunctionCall.Name
		fn, ok := c.funcs[fnName]
		if !ok {
			return "", fmt.Errorf("函数未注册: %s", fnName)
		}
		var args map[string]any
		if err := json.Unmarshal([]byte(choice.FunctionCall.Arguments), &args); err != nil {
			return "", fmt.Errorf("参数解析失败: %v", err)
		}
		result := fn(args)
		return result, nil
	}
	return choice.Content, nil
}

func getWeather(args map[string]any) string {
	city, _ := args["city"].(string)
	weatherData := map[string]string{
		"New York":      "Sunny, 25°C",
		"Tokyo":         "Cloudy, 22°C",
		"San Francisco": "Foggy, 18°C",
	}
	if val, ok := weatherData[city]; ok {
		return val
	}
	return "未知城市的天气"
}

func getTime(args map[string]any) string {
	city, _ := args["city"].(string)
	timeData := map[string]string{
		"New York":      "14:30 PM",
		"Tokyo":         "03:30 AM",
		"San Francisco": "11:30 AM",
	}
	if val, ok := timeData[city]; ok {
		return val
	}
	return "未知城市的时间"
}
