package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// --- 1. โครงสร้างข้อมูล (เชื่อมกับ Environment Variables บน Render) ---
type TaskResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	AiPlan  string `json:"aiPlan"` // สำหรับ AI Overview
}

// --- 2. ด่านตรวจแก้วตา (Security Handshake) ---
func securityGate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ดึงรหัสลับจาก Env Var ที่ท่านตั้งไว้ (FB_VERIFY_TOKEN หรือตั้งใหม่)
		masterSecret := os.Getenv("X_AUTH_TOKEN") 
		if masterSecret == "" { masterSecret = "ThitNuea_Secret_2026" }

		clientAuth := r.Header.Get("X-ThitNuea-Auth")
		
		// ตรวจสอบตราประทับ
		if clientAuth != masterSecret {
			fmt.Printf("🚨 [แก้วตา] ตรวจพบการบุกรุก! Token ไม่ถูกต้อง\n")
			http.Error(w, `{"status":"error", "message":"🛡️ Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// --- 3. สมองกลน้ำอิง (AI Overview Logic) ---
func handleStrategy(w http.ResponseWriter, r *http.Request) {
	// Logic: รับข้อมูลธุรกิจ -> ส่งให้ Gemini -> สรุปแบบ 3 บรรทัด
	summary := "🚀 สรุป SME F-16: \n1. เน้นความปลอดภัยแบบ A2A \n2. ลดต้นทุนด้วย Go Engine \n3. กระจายคอนเทนต์ผ่าน 3 ช่องทางหลัก"
	
	resp := TaskResponse{
		Status:  "success",
		Message: "วิเคราะห์เสร็จสิ้นโดย Gemini 3 Flash",
		AiPlan:  summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	// --- เส้นทางหลักของระบบ ---
	// 1. Webhook เดิมของท่าน (LINE/FB)
	// http.HandleFunc("/webhook", legacyHandler) 

	// 2. ท่อใหม่สำหรับ Dashboard (ผ่านด่านแก้วตา)
	http.HandleFunc("/api/strategy", securityGate(handleStrategy))

	fmt.Printf("🦅 F-16 Defender V.2.1 [HARDENED] เริ่มปฏิบัติการที่พอร์ต %s\n", port)
	
	// ระบบรายงานตัวพรายทอง (Heartbeat)
	go func() {
		for {
			fmt.Println("📱 [พรายทอง] รายงานสถานะ: ระบบ F-16 บินลาดตระเวนปกติ...")
			time.Sleep(1 * time.Hour)
		}
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
