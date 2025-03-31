// main_test.go
package main

import (
	"goWhiteBoard/util"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestSend(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	imageData, err := os.ReadFile("whiteboard.png")
	if err != nil {
		t.Fatalf("画像の読み込みに失敗しました: %v", err)
	}
	util.SendImage(imageData)
}
