package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// หน้าแรก (Home) - ปราบตุ๊กตาตกปลา
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🛡️ F-16 Defender V.2: System Online (Bangkok Server)\n")
		fmt.Fprintf(w, "Status: Ready to Serve")
	})

	// กำหนดพอร์ตให้ตรงกับที่ Google Cloud Run ต้องการ (8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	fmt.Printf("🚀 F-16 Starting on Port %s...\n", port)

	// สั่งรัน Server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
