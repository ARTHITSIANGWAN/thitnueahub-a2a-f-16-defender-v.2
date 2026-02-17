package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

// 1. โครงสร้างข้อมูลกลาง (เน้นประหยัด Memory)
type AIDispatch struct {
	Sender    string                 `json:"sender"`
	Action    string                 `json:"action"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp"`
}

// 2. Safety Gate: ป้องกัน thitnueahub จากการโดนแบน (Compliance Check)
func validateSecurity(msg AIDispatch) error {
	content, _ := msg.Payload["target"].(string)
	// ดักจับ Keyword เสี่ยง (Prohibited Use Policy)
	prohibited := []string{"malware", "exploit", "bypass"}
	for _, word := range prohibited {
		if strings.Contains(strings.ToLower(content), word) {
			return fmt.Errorf("🛡️ SAFETY_BLOCK: ตรวจพบคำต้องห้าม [%s]", word)
		}
	}
	return nil
}

// 3. Smart Dispatcher: ส่งงานไป Gripen พร้อมระบบคืนชีพ (Smart Retry)
func dispatchToAgent(target string, msg AIDispatch, secret string) error {
	maxRetries := 3
	client := &http.Client{Timeout: 10 * time.Second}

	for i := 0; i < maxRetries; i++ {
		jsonData, _ := json.Marshal(msg)
		req, _ := http.NewRequest("POST", target, bytes.NewBuffer(jsonData))
		
		// A2A Handshake Security
		req.Header.Set("X-ThitNuea-Auth", secret)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("✅ [Attempt %d] F-16 ส่งงานสำเร็จ!\n", i+1)
			resp.Body.Close()
			return nil
		}

		// Exponential Backoff: พักเครื่องก่อนลองใหม่ (2, 4, 8 วินาที)
		waitTime := time.Duration(math.Pow(2, float64(i+1))) * time.Second
		fmt.Printf("⚠️ [Attempt %d] ส่งพลาด (อาจเพราะเน็ตมือถือสวิง) ลองใหม่ใน %v...\n", i+1, waitTime)
		time.Sleep(waitTime)
	}
	return fmt.Errorf("❌ F-16 ยอมแพ้: ติดต่อเอเจนต์ปลายทางไม่ได้")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	secret := os.Getenv("A2A_SECRET_KEY")
	if secret == "" { secret = "thitnuea-v2-safe" }

	fmt.Printf("🦅 F-16 Defender V.2 เริ่มปฏิบัติการที่พอร์ต %s [Low Cost Mode]\n", port)

	// จำลองการทำงาน: ตรวจสอบความปลอดภัยอุปกรณ์ทุกๆ 1 นาที
	go func() {
		for {
			task := AIDispatch{
				Sender: "F16-Core",
				Action: "HEALTH_CHECK",
				Payload: map[string]interface{}{"target": "Mobile-System"},
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// 1. ตรวจความปลอดภัยก่อนส่ง
			if err := validateSecurity(task); err != nil {
				fmt.Println(err)
			} else {
				// 2. ส่งต่อให้ Gripen (หรือ Agent อื่นๆ)
				targetAgent := "http://localhost:8081/process"
				dispatchToAgent(targetAgent, task, secret)
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	// รัน Server นิ่งๆ กิน RAM น้อยๆ
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
