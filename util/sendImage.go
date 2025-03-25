package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MessageContent struct {
	Type   string `json:"type"`
	Source struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      string `json:"data"`
	} `json:"source,omitempty"`
	Text string `json:"text,omitempty"`
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type RequestBody struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
}

type ResponseBody struct {
	Choices []struct {
		Message struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func SendImage(imageData []byte) {
	// OpenAI APIキーを設定
	apiKey := os.Getenv("API_KEY")
	endpoint := os.Getenv("END_POINT")

	// モデル、画像メディアタイプ、Base64エンコードされた画像データを設定
	model := os.Getenv("MODEL") // 適切なモデルを選択
	imageMediaType := "image/png"

	// 画像ファイルを読み込み、Base64エンコード
	base64Encoded := base64.StdEncoding.EncodeToString(imageData)

	system_message := "You are an Web Designer and image editor. You can edit the image. You can return image which improved. no Explanation just image."
	user_message := "Create improve image from the freehands image of this and return only image. Draw basic clean image. No Explanation just image."
	// リクエストボディを構築
	requestBody := RequestBody{
		Model:     model,
		System:    system_message,
		MaxTokens: 4096,
		Messages: []Message{
			{
				Role: "user",
				Content: []MessageContent{
					{
						Type: "image",
						Source: struct {
							Type      string `json:"type"`
							MediaType string `json:"media_type"`
							Data      string `json:"data"`
						}{
							Type:      "base64",
							MediaType: imageMediaType,
							Data:      base64Encoded,
						},
					},
					{
						Type: "text",
						Text: user_message,
					},
				},
			},
		},
	}

	// JSONにエンコード
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("JSONエンコードに失敗しました: %v", err)
	}

	// OpenAI APIにリクエストを送信
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatalf("リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("リクエストの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディを解析
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンスの読み込みに失敗しました: %v", err)
	}

	var response ResponseBody
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Fatalf("JSONデコードに失敗しました: %v", err)
	}

	// 結果を出力
	if len(response.Choices) > 0 && len(response.Choices[0].Message.Content) > 0 {
		fmt.Println(response.Choices[0].Message.Content[0].Text)
	} else {
		fmt.Println("No response content found.")
	}
}
