package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	// URL ของ Gripen (ถ้าใช้ Cloud Run ให้ใส่ URL ที่ได้จาก Google)
	gripenURL := os.Getenv("GRIPEN_URL") 
	if gripenURL == "" { gripenURL = "http://localhost:8081/process" }

	http.HandleFunc("/scout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🚀 F-16: Signal received. Preparing Dark-Relay nudge...")

		// สร้าง Payload ลับ (Snake Nudge Protocol)
		jsonData := []byte(`{"agent": "F-16", "status": "active", "cmd": "TARGET_DETECTED"}`)
		
		// ยิง Request ไปหา Gripen พร้อมใส่รหัสลับใน Header
		req, _ := http.NewRequest("POST", gripenURL, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-ThitNuea-Auth", "DragonScale2026") // รหัสลับระหว่างเรา

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Fprintf(w, "❌ Error: Gripen not responding: %v", err)
			return
		}
		defer resp.Body.Close()

		fmt.Fprintf(w, "✅ F-16: Nudge success! Gripen status: %s", resp.Status)
	})

	http.ListenAndServe(":"+port, nil)
}
