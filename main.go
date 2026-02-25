package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. ตั้งค่า Route หน้าแรก (เพื่อไล่น้องตกปลาออกไป)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🛡️ F-16 Defender V.2: System Online (Thailand Server)\n")
		fmt.Fprintf(w, "Status: Gemini AI & Communication Channels Ready.")
	})

	// 2. ฟังก์ชันตรวจสอบการทำงานของระบบ
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "✅ Dashboard is Active")
	})

	// 3. กำหนดพอร์ต (บังคับเป็น 8080 สำหรับ Cloud Run)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 F-16 Starting on Port %s...\n", port)

	// 4. สั่งรัน Server (ท่าเรือ 8080)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
