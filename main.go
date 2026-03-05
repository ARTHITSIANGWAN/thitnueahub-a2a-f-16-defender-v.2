package autonomous

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

// 🛡️ [สมยอมสะดวก & สะอาด] แยกจุดซ่อมบำรุงออกมาเป็น Endpoint ลับ
// ให้ Cloud Scheduler ยิงมาที่ /api/v1/dark-maintenance ตอนตี 3
func HandleDarkMaintenance(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// 1. ด่านพรายทอง (Security) - กันคนนอกมารันคำสั่งทำ DB พัง
	secretKey := os.Getenv("THN_MAINTENANCE_KEY")
	if r.Header.Get("X-ThitNuea-Secret") != secretKey || secretKey == "" {
		log.Println("🚨 [พรายทอง]: ผู้บุกรุกพยายามเข้าถึงระบบซ่อมแซม!")
		http.Error(w, "🚫 Unauthorized Access", http.StatusUnauthorized)
		return
	}

	log.Println("🐅 [AI thitnueahub]: เริ่มกระบวนการล้างท่อและซ่อมแซมฐานข้อมูลรอบดึก...")

	// 2. [สมยอมสะสาง] เคลียร์คิวงานผีดิบ (Zombie Tasks) ที่ค้างคอขวด
	// ไม่ล็อก DB นาน และทำให้คิวงานเช้าวันใหม่ว่างเปล่า
	_, err := db.Exec(`
		UPDATE task_queue 
		SET status = 'Failed', error_message = 'Timeout (Cleared by Kaewta Nightly)' 
		WHERE status = 'Processing' AND updated_at < datetime('now', '-1 hour')
	`)
	if err != nil {
		log.Printf("⚠️ [แก้วตา]: เคลียร์คิวซอมบี้พลาด: %v", err)
	} else {
		log.Println("🧹 [แก้วตา]: กวาดล้างขยะใน Task Queue เรียบร้อย!")
	}

	// 3. [สมยอมสร้างนิสัย] จัดระเบียบข้อมูล (VACUUM)
	// ทำตอนตี 3 ที่มี Cloud Run รันอยู่แค่ Instance เดียว เพื่อความปลอดภัยสูงสุด (Low Cost)
	_, err = db.Exec("VACUUM;")
	if err != nil {
		log.Printf("❌ [System]: VACUUM ล้มเหลว: %v", err)
		http.Error(w, "Database Maintenance Failed", http.StatusInternalServerError)
		return
	}

	// 4. [สมยอมสุขลักษณะ] รายงานผลกลับไปยังศูนย์บัญชาการ
	log.Println("✨ [AI thitnueahub]: บำรุงรักษาเสร็จสิ้น ฐานข้อมูลคลีน 100% พร้อมรับ SME ตอนเช้า!")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "✅ Maintenance OK - Zero Garbage Achieved!")
}
