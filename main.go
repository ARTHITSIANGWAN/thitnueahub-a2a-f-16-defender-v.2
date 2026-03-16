package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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

// --- 💎 โครงสร้าง Gemini 3 Flash ---
type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type Mission struct {
	Platform   string
	ReplyToken string
	Text       string
	UserID     string
	Timestamp  time.Time
}

type ThitNueaHub struct {
	bot       *linebot.Client
	db        *firestore.Client
	missionCh chan Mission
	secret    string
	wg        sync.WaitGroup
}

type DiscordPayload struct {
	Content  string `json:"content"`
	Username string `json:"username,omitempty"`
	Avatar   string `json:"avatar_url,omitempty"`
}

// --- 🚀 ยิงรายงานเข้าห้องบัญชาการ Discord ---
func sendToDiscord(message string, agentName string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}
	payload := DiscordPayload{
		Content:  message,
		Username: agentName,
		Avatar:   "https://cdn-icons-png.flaticon.com/512/4712/4712109.png",
	}
	jsonData, _ := json.Marshal(payload)
	http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
}

// --- ⚡ หัวใจ Gemini 3: แกะเจ้าตาก ---
func askGemini(prompt string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "⚠️ กุญแจหาย! กรุณาเช็ก Secret Manager"
	}
	// ใช้โมเดล Flash ตัวแรง
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	payload := GeminiRequest{}
	payload.Contents = append(payload.Contents, struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{})
	payload.Contents[0].Parts = append(payload.Contents[0].Parts, struct {
		Text string `json:"text"`
	}{Text: "จงแกะเจ้าตากจากข้อความนี้: " + prompt})

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "❌ เชื่อมต่อ Gemini พลาด"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var geminiResp GeminiResponse
	json.Unmarshal(body, &geminiResp)

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text
	}
	return "💎 แก้วตา: กำลังประมวลผลแรงเกินไป หรือ API Key มีปัญหาค่ะ"
}

func main() {
	log.Println("🐅 [ทิศเหนือ ฮับ]: IGNITE - Full Power Online...")
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	dbClient, _ := firestore.NewClient(ctx, projectID)

	hub := &ThitNueaHub{
		db:        dbClient,
		missionCh: make(chan Mission, 10000), // 🔥 ขยายคิวเป็น 10,000 รับงานหมื่นครั้ง
		secret:    os.Getenv("LINE_CHANNEL_SECRET"),
	}

	lineToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	hub.bot, _ = linebot.New(hub.secret, lineToken)

	sendToDiscord("🚀 **[SYSTEM REIGNITE]** F-16 V.2 พร้อมถลุงงบ 9,000 แล้วเจ้านาย!", "🐅 ทิศเหนือ ฮับ (Core)")

	// ปล่อยคนงาน ไอ้จอร์จ 15 คน (เบิ้ลให้ไวขึ้น)
	for i := 1; i <= 15; i++ {
		hub.wg.Add(1)
		go hub.GeorgeWorker(ctx, i)
	}

	http.HandleFunc("/webhook/line", hub.PhraiThongLine)
	http.HandleFunc("/api/surgery", hub.NamIngSurgeryHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "✅ F-16: Active & Heavy Loaded")
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (h *ThitNueaHub) PhraiThongLine(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	hash := hmac.New(sha256.New, []byte(h.secret))
	hash.Write(body)
	sig := r.Header.Get("X-Line-Signature")
	if base64.StdEncoding.EncodeToString(hash.Sum(nil)) != sig {
		http.Error(w, "Unauthorized", 401)
		return
	}
	r.Body = io.NopCloser(strings.NewReader(string(body)))
	events, _ := h.bot.ParseRequest(r)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if msg, ok := event.Message.(*linebot.TextMessage); ok {
				h.missionCh <- Mission{
					Platform: "LINE", ReplyToken: event.ReplyToken,
					Text: msg.Text, UserID: event.Source.UserID, Timestamp: time.Now(),
				}
			}
		}
	}
	w.WriteHeader(200)
}

func (h *ThitNueaHub) NamIngSurgeryHandler(w http.ResponseWriter, r *http.Request) {
	sendToDiscord("🎨 **[น้ำอิง]** กำลังผ่าตัดระบบให้สดใสค่ะ!", "🧑‍🎨 น้ำอิง")
	fmt.Fprint(w, "น้ำอิง: ผ่าตัดเรียบร้อย!")
}

func (h *ThitNueaHub) GeorgeWorker(ctx context.Context, id int) {
	defer h.wg.Done()
	for m := range h.missionCh {
		// 🚀 เบิ้ลพลัง Gemini 3 Flash
		result := askGemini(m.Text)

		// รายงาน Discord Real-time
		report := fmt.Sprintf("📡 **[ไอ้จอร์จ-%d]**\n👤 User: `%s`\n💬 แกะได้: %s", id, m.UserID, result)
		sendToDiscord(report, "🕵️ แก้วตา")

		// ลงถัง Firestore (Disk 100GB ของพี่)
		if h.db != nil {
			h.db.Collection("missions").Add(ctx, map[string]interface{}{
				"user_id": m.UserID, "text": m.Text, "result": result,
				"platform": m.Platform, "timestamp": m.Timestamp,
			})
		}
		h.bot.ReplyMessage(m.ReplyToken, linebot.NewTextMessage(result)).Do()
	}
}
