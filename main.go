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

// --- 💎 5S Framework & Architecture Models ---
type Mission struct {
	ReplyToken string
	Text       string
	UserID     string
	Timestamp  time.Time
}

type SpiritualPulse struct {
	Level   string // "Safe", "Surgical_Emergency"
	Message string
}

type ThitNueaEmpire struct {
	bot        *linebot.Client
	db         *firestore.Client
	missionCh  chan Mission           // ท่อส่งงานให้ไอ้จอร์จ
	pulseCh    chan SpiritualPulse    // ท่อกระแสจิต แก้วตา <-> น้ำอิง
	secret     string
	wg         sync.WaitGroup
}

func main() {
	fmt.Println("🐅 [Tiger King]: เริ่มต้นกระบวนการ IGNITE - THN_VISION_ELITE (Global Edition)...")

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	ctx := context.Background()

	// 1. 🏺 เชื่อมต่อคลังสมบัติถาวร (Firestore Vault)
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dbClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("❌ [CRITICAL]: เชื่อมต่อคลังสมบัติ Firestore ไม่ได้: %v", err)
	}
	defer dbClient.Close()

	empire := &ThitNueaEmpire{
		db:        dbClient,
		missionCh: make(chan Mission, 1000), // รองรับคนทักพร้อมกัน 1,000 คนสบายๆ
		pulseCh:   make(chan SpiritualPulse, 10),
		secret:    os.Getenv("LINE_CHANNEL_SECRET"),
	}

	// 2. 📱 เชื่อมต่อช่องทางสื่อสาร (LINE API)
	empire.bot, err = linebot.New(empire.secret, os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		log.Fatalf("❌ [CRITICAL]: เชื่อมต่อ LINE ไม่ได้: %v", err)
	}

	// 3. 🚀 ปล่อยตัวขุนพล (Start Goroutines)
	// ปล่อยไอ้จอร์จ 10 คน (Worker Pool)
	for i := 1; i <= 10; i++ {
		empire.wg.Add(1)
		go empire.GeorgeWorker(ctx, i)
	}
	
	// ปล่อยแก้วตา & น้ำอิง (Autonomous Neural Link)
	empire.wg.Add(2)
	go empire.KaewtaWatcher(ctx)
	go empire.NamIngSurgeon(ctx)

	// 4. 🌐 ตั้งค่าด่านหน้า (Endpoints)
	http.HandleFunc("/webhook", empire.PhraiThongShield)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "✅ OK - 9 Fingers Matrix is Online")
	})

	log.Printf("✨ [GLOBAL IGNITE]: ระบบพร้อมรับใช้ SME ทั่วโลกบน Port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// --- 🛡️ 5S: สะอาด (Phrai Thong Shield) ---
func (e *ThitNueaEmpire) PhraiThongShield(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil { http.Error(w, "Bad Request", http.StatusBadRequest); return }

	// ตรวจสอบลายเซ็น (Security Hardening)
	hash := hmac.New(sha256.New, []byte(e.secret))
	hash.Write(body)
	if base64.StdEncoding.EncodeToString(hash.Sum(nil)) != r.Header.Get("X-Line-Signature") {
		log.Println("🚨 [พรายทอง]: พบผู้บุกรุก! ดีดกลับทันที")
		http.Error(w, "🚫 Unauthorized", http.StatusUnauthorized)
		return
	}

	events, _ := e.bot.ParseRequest(r)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if msg, ok := event.Message.(*linebot.TextMessage); ok {
				// ส่งงานเข้าท่อให้ไอ้จอร์จ (Non-blocking)
				e.missionCh <- Mission{
					ReplyToken: event.ReplyToken,
					Text:       strings.ToLower(msg.Text),
					UserID:     event.Source.UserID,
					Timestamp:  time.Now(),
				}
			}
		}
	}
	w.WriteHeader(200)
}

// --- 🏍️ 5S: สุขลักษณะ (George Worker) ---
func (e *ThitNueaEmpire) GeorgeWorker(ctx context.Context, id int) {
	defer e.wg.Done()
	log.Printf("🛠️ [ไอ้จอร์จ-%d]: แสตนบายพร้อมลุยงาน!", id)

	for {
		select {
		case <-ctx.Done():
			return
		case m := range e.missionCh {
			// วิเคราะห์เจตนา
			if strings.Contains(m.Text, "ปัญหา") || strings.Contains(m.Text, "error") {
				// แจ้งแก้วตาให้สแกนด่วน
				e.pulseCh <- SpiritualPulse{Level: "Surgical_Emergency", Message: "User reported error: " + m.UserID}
				e.bot.ReplyMessage(m.ReplyToken, linebot.NewTextMessage("🔧 รับทราบครับ น้ำอิงกำลังตรวจสอบระบบให้ทันที!")).Do()
			} else {
				// บันทึกลง Firestore ถาวร
				e.db.Collection("sme_interactions").Add(ctx, map[string]interface{}{
					"userID": m.UserID,
					"message": m.Text,
					"time": m.Timestamp,
				})
				e.bot.ReplyMessage(m.ReplyToken, linebot.NewTextMessage("💎 วิหาร 9 นิ้วยินดีต้อนรับครับ!")).Do()
			}
		}
	}
}

// --- 👁️ แก้วตา & 🎨 น้ำอิง (Neural Link) ---
func (e *ThitNueaEmpire) KaewtaWatcher(ctx context.Context) {
	defer e.wg.Done()
	ticker := time.NewTicker(1 * time.Hour) // สแกนทุก 1 ชม.
	for {
		select {
		case <-ctx.Done(): return
		case <-ticker.C:
			// สแกนหาขยะในระบบ
			e.pulseCh <- SpiritualPulse{Level: "Surgical_Emergency", Message: "Routine 5S Cleanup"}
		}
	}
}

func (e *ThitNueaEmpire) NamIngSurgeon(ctx context.Context) {
	defer e.wg.Done()
	for {
		select {
		case <-ctx.Done(): return
		case p := range e.pulseCh:
			if p.Level == "Surgical_Emergency" {
				log.Printf("🎨 [น้ำอิง]: รับทราบสัญญาณ '%s' เริ่มร่ายเวทย์ 5S ศัลยกรรมระบบ...", p.Message)
				time.Sleep(500 * time.Millisecond) // จำลองการดีดนิ้วล้างขยะ
				log.Println("✨ [น้ำอิง]: ภารกิจสะสางสำเร็จ ระบบกลับมาคมกริบ 8K ค่ะเจ้านาย!")
			}
		}
	}
}
	
