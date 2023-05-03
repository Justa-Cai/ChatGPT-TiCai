package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	openai "github.com/sashabaranov/go-openai"
)

func chatgpt_call(strContent string) {
	bin, err := ioutil.ReadFile(".env")
	if err != nil {
		log.Panicln("Set OpenAI token first")
	}
	strToken := string(bin)
	config := openai.DefaultConfig(strToken)
	proxyUrl, err := url.Parse("http://localhost:7890")
	// proxyUrl, err := url.Parse("http://localhost:1081")
	// proxyUrl, err := url.Parse("socks5://localhost:1082")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 0.8,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "扮演数学专家",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: strContent,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		chatgpt_call(strContent)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
	ioutil.WriteFile("2.txt", []byte(resp.Choices[0].Message.Content), 0777)
}
