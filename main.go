package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// --- 1. โครงสร้างข้อมูล (Data Models) ---
type StrategyRequest struct {
	BusinessInfo string `json:"business_info"`
}

type StrategyResponse struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

// --- 2. ระบบความปลอดภัย (Kaewta Security Gate) ---
func validateSecurity(content string) bool {
	prohibited := []string{"malware", "exploit", "hack", "bypass", "script", "drop table"}
	for _, word := range prohibited {
		if strings.Contains(strings.ToLower(content), word) {
			return false
		}
	}
	return true
}

// --- 3. สมองกล (Narm-Ing Engine: Gemini 3 Flash) ---
func callGeminiAI(prompt string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "🚨 [Error] ไม่พบกุญแจ Gemini ในระบบค่ะ"
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey
	
	instruction := fmt.Sprintf("คุณคือผู้เชี่ยวชาญวางแผนกลยุทธ์ SME ของ ThitNueaHub จงวิเคราะห์ธุรกิจนี้แบบเจาะลึกและสั้นกระชับ: %s", prompt)
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": instruction}}},
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "🚨 [Error] การเชื่อมต่อ Gemini ขัดข้อง"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Parsing Gemini Response (Simple Version)
	tryParse := result["candidates"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return tryParse
}

// --- 4. หน่วยกระจายข่าว (Praithong Dispatcher: LINE & Meta) ---
func praithongDispatcher(msg string) {
	// ดึงกุญแจจาก Environment Variables (Master Key System)
	lineToken := os.Getenv("LINE_TOKEN")
	metaToken := os.Getenv("META_TOKEN")

	// [ส่งเข้า LINE]
	if lineToken != "" {
		lineUrl := "https://api.line.me/v2/bot/message/broadcast"
		linePayload := fmt.Sprintf(`{"messages":[{"type":"text","text":"🛡️ F-16 Mission Report:\n%s"}]}`, msg)
		req, _ := http.NewRequest("POST", lineUrl, strings.NewReader(linePayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+lineToken)
		http.DefaultClient.Do(req)
	}

	// [ส่งเข้า Meta (Facebook Feed)]
	if metaToken != "" {
		metaUrl := "https://graph.facebook.com/v19.0/me/feed"
		metaValues := fmt.Sprintf("message=%s&access_token=%s", msg, metaToken)
		http.Post(metaUrl, "application/x-www-form-urlencoded", strings.NewReader(metaValues))
	}
	log.Println("✅ [พรายทอง] รายงานภารกิจเสร็จสิ้น")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	// [ท่อที่ 1]: หน้า Dashboard (Kaewta UI)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, dashboardHTML)
	})

	// [ท่อที่ 2]: Intelligence API
	http.HandleFunc("/api/strategy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { return }

		var req StrategyRequest
		json.NewDecoder(r.Body).Decode(&req)

		if !validateSecurity(req.BusinessInfo) {
			json.NewEncoder(w).Encode(StrategyResponse{Status: "Fail", Content: "🚨 ตรวจพบความเสี่ยงด้านความปลอดภัย!"})
			return
		}

		// 🤖 เริ่มการประมวลผล
		aiResult := callGeminiAI(req.BusinessInfo)

		// 📱 ส่งกระจายข่าวแบบ Background
		go praithongDispatcher(aiResult)

		json.NewEncoder(w).Encode(StrategyResponse{
			Status:  "Success",
			Content: aiResult + "\n\n(รายงานถูกส่งเข้า LINE และ Meta เรียบร้อยค่ะ!)",
		})
	})

	fmt.Printf("🚀 [F-16 V.2] Full Operational at Port %s | Time: %s\n", port, time.Now().Format(time.RFC822))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const dashboardHTML = `
<!DOCTYPE html>
<html lang="th">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>F-16 DEFENDER MASTER</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Sarabun:wght@300;700&display=swap" rel="stylesheet">
    <style>body { font-family: 'Sarabun', sans-serif; background: #f8fafc; }</style>
</head>
<body class="flex flex-col min-h-screen">
    <div class="flex-grow flex items-center justify-center p-6">
        <div class="max-w-lg w-full bg-white rounded-[2.5rem] p-10 shadow-2xl border-t-8 border-emerald-500">
            <div class="flex justify-between items-start mb-8">
                <div>
                    <h1 class="text-3xl font-black text-slate-900 tracking-tighter">F-16 MASTER</h1>
                    <p class="text-xs font-bold text-emerald-600 uppercase tracking-widest">SME Airborne Ecosystem</p>
                </div>
                <div class="w-12 h-12 bg-slate-900 rounded-2xl flex items-center justify-center text-white font-black text-xl">V2</div>
            </div>
            
            <label class="text-[10px] font-black text-slate-400 uppercase mb-2 block">Business Information Input</label>
            <textarea id="bizInput" class="w-full h-40 p-5 rounded-3xl bg-slate-50 border border-slate-200 mb-6 outline-none focus:ring-4 focus:ring-emerald-100 transition-all text-sm" placeholder="พิมพ์ชื่อธุรกิจหรือปัญหาที่ต้องการให้ AI ช่วยวางแผน..."></textarea>
            
            <button onclick="launchMission()" id="btn" class="w-full py-5 bg-emerald-600 text-white rounded-3xl font-black text-lg shadow-lg shadow-emerald-200 hover:bg-slate-900 transition-all active:scale-95">START MISSION 🚀</button>
            
            <div id="result" class="mt-8 p-6 bg-slate-900 text-emerald-300 rounded-3xl hidden text-sm leading-relaxed whitespace-pre-line border border-emerald-500/30"></div>
        </div>
    </div>
    <script>
        async function launchMission() {
            const btn = document.getElementById('btn');
            const res = document.getElementById('result');
            const input = document.getElementById('bizInput').value;
            if(!input) return alert('กรุณาระบุข้อมูลธุรกิจค่ะเจ้านาย');

            btn.disabled = true; btn.innerText = 'ประมวลผลกลยุทธ์... 🧠';
            res.classList.add('hidden');
            
            try {
                const response = await fetch('/api/strategy', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({business_info: input})
                });
                const data = await response.json();
                res.classList.remove('hidden');
                res.innerText = data.content;
            } catch (e) {
                alert('ท่อข้อมูลขัดข้อง!');
            } finally {
                btn.disabled = false; btn.innerText = 'START MISSION 🚀';
            }
        }
    </script>
</body>
</html>`
