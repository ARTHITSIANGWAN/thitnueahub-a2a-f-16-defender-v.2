package main

import (
	"fmt"
	"strings"
	"time"
)

// รายการคำต้องห้ามเบื้องต้นเพื่อป้องกันการโดน Flag จากระบบ Google (Prohibited Use)
var prohibitedKeywords = []string{
	"malware", "exploit", "bypass security", "hacking tool",
	"generate fake identity", "phishing", "hate speech",
}

func validatePayload(action string, payload map[string]interface{}) error {
	// 1. ตรวจสอบขนาด Payload (ป้องกัน Denial of Service)
	content, ok := payload["content"].(string)
	if ok && len(content) > 15000 { // ปรับเกณฑ์ตามความเหมาะสมของ Gemini API
		return fmt.Errorf("SAFETY_ERR: Payload size limit exceeded (Max 15,000 chars)")
	}

	// 2. ตรวจสอบเนื้อหาที่ผิดนโยบาย (Generative AI Prohibited Use Check)
	contentLower := strings.ToLower(content)
	for _, word := range prohibitedKeywords {
		if strings.Contains(contentLower, word) {
			// แจ้งเตือนระดับวิกฤต (Critical Alert)
			return fmt.Errorf("POLICY_VIOLATION: Detect prohibited keyword [%s]. Action Blocked to protect thitnueahub.", word)
		}
	}

	// 3. ตรวจสอบ Action ที่มีความเสี่ยงสูง
	if action == "SENSITIVE_ACCESS" {
		if _, authorized := payload["auth_token"]; !authorized {
			return fmt.Errorf("SECURITY_ERR: Unauthorized sensitive access attempt")
		}
	}

	return nil
}

// ระบบเตือนภัยที่ชัดเจน (Unified Alert System)
func logSecurityAlert(err error) {
	fmt.Printf("\n[🚨 SECURITY ALERT - %s]\n", time.Now().Format("15:04:05"))
	fmt.Printf("Message: %v\n", err)
	fmt.Println("Status: COMMAND_REJECTED")
	fmt.Println("Protection: thitnueahub account status remains SAFE")
	fmt.Println("--------------------------------------------")
}

func dispatchWithSmartRetry(target string, msg AIDispatch) error {
	// 1. ตรวจสอบนโยบายก่อนส่ง (Compliance Check)
	if err := validatePayload(msg.Action, msg.Payload); err != nil {
		logSecurityAlert(err) // แสดงระบบเตือนที่ชัดเจน
		return err
	}

	// ส่วนที่เหลือของ Retry Logic...
    // (โค้ด Smart Retry เดิมของเจ้านายจะทำงานต่อเมื่อผ่าน Gate นี้เท่านั้น)
    return nil 
}
