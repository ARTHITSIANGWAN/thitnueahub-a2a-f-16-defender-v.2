package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// --- 💎 5S Architecture & Models ---
type Mission struct {
	Platform   string // "LINE", "FB", "TELEGRAM"
	ReplyToken string
	Text       string
	UserID     string
	Timestamp  time.Time
}

type ThitNueaEmpire struct {
	bot       *linebot.Client
	db        *firestore.Client
	missionCh chan Mission
	secret    string
	wg        sync.WaitGroup
}

func main() {
	log.Println("🐅 [Tiger King]: IGNITE - Trinity Core Online...")

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	ctx := context.Background()

	// 1. เชื่อมต่อ Firestore
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dbClient, _ := firestore.NewClient(ctx, projectID)
	
	empire := &ThitNueaEmpire{
		db:        dbClient,
		missionCh: make(chan Mission, 1000),
		secret:    os.Getenv("LINE_CHANNEL_SECRET"),
	}

	// 2. เชื่อมต่อ LINE
	empire.bot, _ = linebot.New(empire.secret, os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))

	// 3. ปล่อยไอ้จอร์จ 10 คนลุยงาน
	for i := 1; i <= 10; i++ {
		empire.wg.Add(1)
		go empire.GeorgeWorker(ctx, i)
	}

	// 4. ตั้งค่า Triple Webhook (ด่านพรายทอง)
	http.HandleFunc("/webhook/line", empire.PhraiThongLine)
	http.HandleFunc("/webhook/facebook", empire.PhraiThongMeta)
	http.HandleFunc("/webhook/telegram", empire.PhraiThongTelegram)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "✅ F-16 Trinity: Stable")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// --- 🛡️ Phrai Thong Shield (Handlers) ---

func (e *ThitNueaEmpire) PhraiThongLine(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	// ตรวจสอบ Signature (5S Security)
	hash := hmac.New(sha256.New, []byte(e.secret))
	hash.Write(body)
	if base64.StdEncoding.EncodeToString(hash.Sum(nil)) != r.Header.Get("X-Line-Signature") {
		http.Error(w, "Unauthorized", 401); return
	}
	events, _ := e.bot.ParseRequest(r)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if msg, ok := event.Message.(*linebot.TextMessage); ok {
				e.missionCh <- Mission{Platform: "LINE", ReplyToken: event.ReplyToken, Text: msg.Text, UserID: event.Source.UserID}
			}
		}
	}
	w.WriteHeader(200)
}

func (e *ThitNueaEmpire) PhraiThongMeta(w http.ResponseWriter, r *http.Request) {
	// 🤫 [น้ำอิง]: ประตูเมต้าเปิดรับคำสั่งแล้วค่ะ (Logic การแกะ JSON ของ FB จะอยู่ตรงนี้)
	w.WriteHeader(200)
}

func (e *ThitNueaEmpire) PhraiThongTelegram(w http.ResponseWriter, r *http.Request) {
	// 🤖 [Optimus]: รับสัญญาณจาก Telegram แล้ว
	w.WriteHeader(200)
}

// --- 🏍️ George Worker (Processing) ---
func (e *ThitNueaEmpire) GeorgeWorker(ctx context.Context, id int) {
	defer e.wg.Done()
	for m := range e.missionCh {
		log.Printf("🛠️ [ไอ้จอร์จ-%d] รับงานจาก %s: %s", id, m.Platform, m.Text)
		// บันทึก Firestore & ตอบกลับ (ใส่ Logic น้ำอิงกวนตีนๆ ตรงนี้ได้เลย)
		e.bot.ReplyMessage(m.ReplyToken, linebot.NewTextMessage("💎 ThitNuea Vision: รับทราบครับ!")).Do()
	}
}
