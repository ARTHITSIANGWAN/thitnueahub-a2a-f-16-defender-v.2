package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// --- 1. โครงสร้างข้อมูล (Models) ---
type StrategyRequest struct {
	BusinessInfo string `json:"business_info"`
}

type StrategyResponse struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

// --- 2. ระบบรักษาความปลอดภัย โดย "แก้วตา" ---
func validateSecurity(content string) bool {
	prohibited := []string{"malware", "exploit", "bypass", "hack", "script"}
	for _, word := range prohibited {
		if strings.Contains(strings.ToLower(content), word) {
			return false
		}
	}
	return true
}

func main() {
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	// --- 3. Endpoints ---

	// [ท่อที่ 1]: ส่งหน้า Dashboard ให้เจ้านาย (Embedded HTML)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, dashboardHTML)
	})

	// [ท่อที่ 2]: API ประมวลผลกลยุทธ์ (The Intelligence Hub)
	http.HandleFunc("/api/strategy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req StrategyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return
		}

		// ด่านตรวจของแก้วตา
		if !validateSecurity(req.BusinessInfo) {
			json.NewEncoder(w).Encode(StrategyResponse{
				Status:  "Fail",
				Content: "🚨 [Security Alert] แก้วตาตรวจพบเนื้อหาไม่ปลอดภัยในคำสั่งของคุณค่ะ",
			})
			return
		}

		// จำลองการประมวลผลของ "น้ำอิง" (ในอนาคตจะต่อ Gemini API ตรงนี้)
		result := fmt.Sprintf("✨ แผนกลยุทธ์สำหรับ '%s': \n1. เน้นการทำ Micro-Targeting ในพื้นที่ \n2. ใช้ Content แบบ Short-Video เพื่อดึงดูด SME \n3. ปรับโครงสร้างต้นทุนด้วยระบบ Zero Garbage ของ F-16", req.BusinessInfo)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StrategyResponse{
			Status:  "Success",
			Content: result,
		})
	})

	fmt.Printf("🚀 [แก้วตา] F-16 Defender V.2 Online! Port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// --- 4. หน้า Dashboard (Embedded HTML) ---
const dashboardHTML = `
<!DOCTYPE html>
<html lang="th">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>F-16 Defender V.2: SME Dashboard</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Sarabun:wght@300;400;700&display=swap');
        body { font-family: 'Sarabun', sans-serif; background-color: #f1f5f0; }
        .glass-panel { background: rgba(255, 255, 255, 0.8); backdrop-filter: blur(10px); border: 1px solid rgba(255, 255, 255, 0.5); }
    </style>
</head>
<body class="antialiased">
    <div class="max-w-4xl mx-auto p-6 space-y-8 pt-10">
        <header class="flex items-center gap-4">
            <div class="w-12 h-12 bg-emerald-600 rounded-2xl flex items-center justify-center text-white font-black">F16</div>
            <div>
                <h1 class="text-2xl font-black text-slate-900">DEFENDER V.2</h1>
                <p class="text-xs font-bold text-emerald-600 uppercase tracking-widest">Supervisor: KAEWTA 🛡️</p>
            </div>
        </header>

        <section class="glass-panel rounded-[2.5rem] p-8 shadow-xl">
            <h3 class="text-xl font-bold mb-4">ศูนย์บัญชาการกลยุทธ์ ✨</h3>
            <textarea id="bizInput" class="w-full h-32 p-4 rounded-2xl border border-slate-200 focus:ring-2 focus:ring-emerald-500 outline-none mb-4" placeholder="ป้อนข้อมูลธุรกิจที่นี่..."></textarea>
            <button onclick="sendStrategy()" id="btn" class="w-full py-4 bg-emerald-600 text-white rounded-2xl font-bold hover:bg-emerald-700 transition-all">วิเคราะห์แผนการบิน</button>
            
            <div id="result" class="mt-6 p-6 bg-slate-900 text-emerald-400 rounded-2xl hidden font-mono text-sm whitespace-pre-line"></div>
        </section>
    </div>

    <script>
        async function sendStrategy() {
            const input = document.getElementById('bizInput').value;
            const btn = document.getElementById('btn');
            const resDiv = document.getElementById('result');
            
            if(!input) return alert('กรุณาพิมพ์ข้อมูลก่อนค่ะเจ้านาย');

            btn.disabled = true;
            btn.innerText = 'กำลังคำนวณ...';
            
            try {
                const response = await fetch('/api/strategy', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({business_info: input})
                });
                const data = await response.json();
                resDiv.classList.remove('hidden');
                resDiv.innerText = data.content;
            } catch (e) {
                alert('เกิดข้อผิดพลาดในการเชื่อมต่อท่อข้อมูลค่ะ');
            } finally {
                btn.disabled = false;
                btn.innerText = 'วิเคราะห์แผนการบิน';
            }
        }
    </script>
</body>
</html>
`
