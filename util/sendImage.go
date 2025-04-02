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
	"strings"
)

type MessageContentSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type MessageContent struct {
	Type   string                `json:"type"`
	Source *MessageContentSource `json:"source,omitempty"`
	Text   string                `json:"text,omitempty"`
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type RequestBody struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
}

// ContentItem は content 配列内の各要素を表す構造体
type ContentItem struct {
	Type string `json:"type"` // コンテンツのタイプ（例: text, image）
	Text string `json:"text"` // コンテンツのテキスト
}

// ResponseBody はレスポンス全体を表す構造体
type ResponseBody struct {
	ID      string        `json:"id"`      // メッセージID
	Type    string        `json:"type"`    // メッセージタイプ（例: message）
	Role    string        `json:"role"`    // ロール（例: assistant, user）
	Model   string        `json:"model"`   // 使用したモデル名
	Content []ContentItem `json:"content"` // content 配列
}

func removeCodeTags(input string) string {
	// Remove the opening and closing ``````
	result := strings.ReplaceAll(input, "```html", "")
	result = strings.ReplaceAll(result, "```", "")
	return result
}

func SendImage(imageData []byte) string {
	// OpenAI APIキーを設定
	apiKey := os.Getenv("API_KEY")
	endpoint := os.Getenv("END_POINT")

	// モデル、画像メディアタイプ、Base64エンコードされた画像データを設定
	model := os.Getenv("MODEL") // 適切なモデルを選択
	imageMediaType := "image/png"

	// 画像ファイルを読み込み、Base64エンコード
	base64Encoded := base64.StdEncoding.EncodeToString(imageData)

	system_message := "You are an Web Designer and image editor. You can edit the image. You can return image which improved. no Explanation just image."
	//user_message := "Create improve image  from the freehands image of this and return only improved image by html. Draw basic clean image. No Explanation return just html."
	user_message := "This freehand image should be neatly formatted and converted with basic Shapes like straight line, circle or square,etc  " +
		" back to HTML and CSS. No explanation is required, just return the HTML."
	//user_message := "Describe this image."
	// リクエストボディを構築
	requestBody := RequestBody{
		Model:     model,
		System:    system_message,
		MaxTokens: 1024,
		Messages: []Message{
			{
				Role: "user",
				Content: []MessageContent{
					{
						Type: "image",
						Source: &MessageContentSource{
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

	// endpoint にリクエストを送信
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatalf("リクエストの作成に失敗しました: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("x-api-key", apiKey)

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
	if len(response.Content) > 0 {
		return removeCodeTags(response.Content[0].Text)
	} else {
		fmt.Println("No response content found.")
		return ""
	}
}
