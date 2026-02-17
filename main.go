package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

// เพิ่มระบบ Logging และ Error Handling ที่เข้มข้นขึ้น
func dispatchWithSmartRetry(target string, msg AIDispatch) error {
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		err := sendToBot(target, msg)
		if err == nil {
			fmt.Printf("✅ [Attempt %d] ส่งคำสั่ง %s สำเร็จ!\n", i+1, msg.Action)
			return nil
		}

		// คำนวณเวลาการรอแบบ Exponential (2^i) เพื่อไม่ให้ยิงรัวเกินไป
		// ครั้งที่ 1 รอ 2 วิ, ครั้งที่ 2 รอ 4 วิ, ครั้งที่ 3 รอ 8 วิ
		waitTime := time.Duration(math.Pow(2, float64(i+1))) * time.Second
		
		fmt.Printf("⚠️ [Attempt %d] ล้มเหลว: %v. จะลองใหม่ในอีก %v...\n", i+1, err, waitTime)
		time.Sleep(waitTime)
	}

	return fmt.Errorf("❌ ระบบยอมแพ้: ไม่สามารถเชื่อมต่อกับ %s ได้หลังจากพยายาม %d ครั้ง", target, maxRetries)
}

// ปรับปรุงส่วนการส่งให้รองรับ Timeout
func sendToBot(target string, msg AIDispatch) error {
	jsonData, _ := json.Marshal(msg)
	
	// ตั้งค่า Timeout ป้องกันค้าง (สำคัญมากสำหรับ Security)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(target, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	return nil
}
