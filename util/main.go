package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := "https://router.huggingface.co/together/v1/chat/completions"
	bearerToken := os.Getenv("API_KEY")

	// リクエストボディの構造体
	requestBody := map[string]interface{}{
		"model": "deepseek-ai/DeepSeek-R1",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "What is the capital of France?",
			},
		},
		"max_tokens": 500,
		"stream":     false,
	}

	// JSONにエンコード
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("JSONのエンコードに失敗: %v", err)
		return
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("リクエストの作成に失敗: %v", err)
		return
	}

	// ヘッダーの設定
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントの作成とリクエストの送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("リクエストの送信に失敗: %v", err)
		return
	}
	defer resp.Body.Close()

	// レスポンスの読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンスの読み込みに失敗: %v", err)
		return
	}

	// レスポンスの出力
	fmt.Println("Response Body:", string(body))
}
