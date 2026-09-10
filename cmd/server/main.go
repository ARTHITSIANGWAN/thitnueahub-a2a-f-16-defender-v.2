package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// 🛡️ F16 Defender Configuration
const (
	Port = ":8080"
)

var (
	lineClient *linebot.Client
	fsClient   *firestore.Client
)

func main() {
	ctx := context.Background()
	var err error

	// 1. เริ่มต้นเชื่อมต่อ LINE Bot SDK
	lineClient, err = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Printf("⚠️ LINE Bot Warning (เช็ค Credentials): %v", err)
	}

	// 2. เริ่มต้นเชื่อมต่อ Google Firestore (Empire Ledger)
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID != "" {
		fsClient, err = firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Printf("⚠️ Firestore Warning: %v", err)
		} else {
			defer fsClient.Close()
		}
	}

	// 3. กำหนดเส้นทาง API Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/webhook/line", handleLineWebhook)
	http.HandleFunc("/api/v4/defender/status", handleDefenderStatus)

	fmt.Printf("🛩️ F16-DEFENDER-LOWCOST Active | Port %s | Ready for Action!\n", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

// 🌐 หน้าแรกสถานะเซิร์ฟเวอร์
func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>🛡️ TNH-AI-V4-F16-DEFENDER</h1><p>Status: Operational & Secure</p>")
}

// 💚 LINE Webhook Handler สำหรับรับข้อความแจ้งเตือน
func handleLineWebhook(w http.ResponseWriter, r *http.Request) {
	if lineClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := lineClient.ParseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if message, ok := event.Message.(*linebot.TextMessage); ok {
				replyText := fmt.Sprintf("🛩️ F16 Defender รับทราบข้อความ: %s", message.Text)
				lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

// ⚡ เช็คสถานะความพร้อมของ F16 Defender
func handleDefenderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"unit": "F16-DEFENDER", "status": "ONLINE", "cost": "LOW", "shield": "ACTIVE", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}
