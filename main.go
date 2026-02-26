package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/genai" // ระบบสมองกลที่เจ้านายจดทะเบียนไว้
)

// โครงสร้างคำสั่งที่รับมาจาก Vercel (Elite Access)
type Command struct {
	Code   string `json:"code"`
	Action string `json:"action"`
}

func main() {
	// 1. ตรวจสอบกระสุน (Environment Variables)
	apiKey := os.Getenv("GEMINI_API_KEY")
	secretCode := "THIT1503" // รหัสปลดล็อกที่คุณปักหมุดไว้

	// 2. ตั้งค่าสมองกล Gemini V.2 (The Strategic Brain)
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.Config{APIKey: apiKey, Backend: genai.BackendVertexAI})

	// 3. ระบบ Relay & Gateway (Zero Garbage Protocol)
	http.HandleFunc("/ignite", func(w http.ResponseWriter, r *http.Request) {
		// --- CORS Management (เปิดทางให้ Vercel เข้ามา) ---
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Printf("[%s] 🚀 SIGNAL DETECTED: F-16 INITIATION SEQUENCE\n", time.Now().Format("15:04:05"))

		// --- ส่วนวิเคราะห์ภารกิจ (AI Mission Analysis) ---
		// ตัวอย่าง Logic ที่เจ้านายอยากให้ Gemini สรุปสถานะ SME
		model := client.Models.BaseModel("gemini-2.0-flash")
		resp, _ := client.Models.GenerateContent(ctx, model, genai.Text("สรุปสถานะระบบ F-16 Defender: ระบบออนไลน์ 100% พร้อมคุ้มครอง SME"), nil)
		
		fmt.Fprintf(w, "🛡️ F-16 DEFENDER V.2: %s\n", resp.Candidates[0].Content.Parts[0])
	})

	// 4. หน้าหลัก (Stealth Mode Page)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<div style="background:#020617; color:#22d3ee; padding:50px; font-family:monospace;">
				<h1>🛡️ THN_VISION_ELITE: PROTECTED</h1>
				<p>Status: CLOUD_RUN_ACTIVE | Region: asia-southeast3 (BKK)</p>
				<p>Safety Gate: %s</p>
			</div>
		`, secretCode)
	})

	// 5. รันระบบฐานทัพ
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🎯 [MISSION START] F-16 DASHBOARD PORT %s ONLINE\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ SYSTEM CRASH: %v", err)
	}
}
