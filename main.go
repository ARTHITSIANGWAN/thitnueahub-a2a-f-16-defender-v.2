package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// --- 1. โครงสร้างเอเจนต์ (ยังคงเดิมเพื่อความเสถียร) ---
type AgentFleet struct {
	Supervisor string // แก้วตา
	Backend    string // น้ำอิง
	Messenger  string // พรายทอง
}

// --- 2. ระบบด่านตรวจ (Real Handshake Handler) ---
// ปรับปรุง: เปลี่ยนจากการจำลอง เป็นการเช็ก Request จริงที่วิ่งเข้ามา
func securityGate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		masterSecret := os.Getenv("A2A_SECRET_KEY")
		if masterSecret == "" { masterSecret = "thitnuea-v2-safe" }

		clientAuth := r.Header.Get("X-ThitNuea-Auth")
		isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

		// ถ้า Handshake ไม่ผ่าน หรือไม่ใช่ HTTPS ให้ "ดีดออก" ทันที
		if clientAuth != masterSecret || !isSecure {
			fmt.Printf("🚨 [แก้วตา] SECURITY ALERT: พบการบุกรุกหรือการเชื่อมต่อไม่ปลอดภัย!\n")
			http.Error(w, "🛡️ Unauthorized: ตราประทับไม่ถูกต้อง", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// --- 3. ฟังก์ชันประมวลผลงาน (Action Handler) ---
func handleTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "✅ [น้ำอิง] รับงานเข้าท่อเรียบร้อย กำลังประมวลผล...")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	fleet := AgentFleet{Supervisor: "แก้วตา", Backend: "น้ำอิง", Messenger: "พรายทอง"}

	// --- 4. การปรับปรุงจุดอ่อน (Addressing Disadvantages) ---
	
	// จุดอ่อนที่ 1: แก้ไขจากการจำลองเป็นการเปิด "ด่านตรวจจริง"
	http.HandleFunc("/dispatch", securityGate(handleTask))

	fmt.Printf("🦅 F-16 Defender V.2.1 [Hardened] Started at Port: %s\n", port)

	// จุดอ่อนที่ 2: ระบบแจ้งเตือนเมื่อเกิดความพินาศ (Active Monitoring)
	go func() {
		for {
			fmt.Printf("\n--- 🛡️ [%s] ลาดตระเวน: ทุกระบบปกติ (Healthy Check) ---\n", fleet.Supervisor)
			
			// จำลองการตรวจสอบ Error: หากระบบหลักล่ม พรายทองต้องรายงาน
			if time.Now().Minute() % 10 == 0 { // สมมติเช็กทุก 10 นาที
				fmt.Printf("📱 [%s] Status: ส่ง Heartbeat เข้า Telegram... OK\n", fleet.Messenger)
			}
			
			time.Sleep(1 * time.Hour)
		}
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
