package main

import (
	"fmt"
	"net/http"
	"os"
)

// f-16 scout: low cost, zero garbage, 24/7 monitoring
func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	http.HandleFunc("/scout", func(w http.ResponseWriter, r *http.Request) {
		// a2a: รับข้อมูลแล้ว "สะกิด" (nudge) ไปที่ gripen
		fmt.Fprintf(w, "f-16: signal received. nudging gripen engine...")
	})

	http.ListenAndServe(":"+port, nil)
}

