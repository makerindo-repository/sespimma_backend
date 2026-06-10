package handler

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sespima_api/config"
	"sespima_api/internal/service"
	"sespima_api/models"
	"sespima_api/pkg/response"
)

// DashboardHandler aggregates a Serdik's live picture (scores, approved
// rewards/punishments, recommendations, activity history) and exposes the
// unified Korsis approval inbox.
type DashboardHandler struct {
	kegiatan   *service.KegiatanService
	reward     *service.RewardService
	punishment *service.PunishmentService
}

func NewDashboardHandler(
	kegiatan *service.KegiatanService,
	reward *service.RewardService,
	punishment *service.PunishmentService,
) *DashboardHandler {
	return &DashboardHandler{kegiatan: kegiatan, reward: reward, punishment: punishment}
}

// ---- response shapes ----

type scoreBlock struct {
	NilaiAkademik   float64 `json:"nilai_akademik"`
	NilaiJasmani    float64 `json:"nilai_jasmani"`
	NilaiMental     float64 `json:"nilai_mental"`
	RewardTotal     float64 `json:"reward_total"`
	PunishmentTotal float64 `json:"punishment_total"`
	// NetAdjustment is reward_total - punishment_total: the live delta applied
	// to the serdik's accumulated score from approved maker-checker records.
	NetAdjustment float64 `json:"net_adjustment"`
}

type rpItem struct {
	ID    int64     `json:"id"`
	Name  string    `json:"name"`
	Point float64   `json:"point"`
	Qty   int       `json:"qty"`
	Date  time.Time `json:"date"`
	Notes *string   `json:"notes,omitempty"`
}

type activityItem struct {
	Type        string    `json:"type"` // reward | punishment | attendance | izin
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Points      float64   `json:"points"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetSerdikDashboard returns the aggregated dashboard for the authenticated serdik.
func (h *DashboardHandler) GetSerdikDashboard(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		response.Unauthorized(c, "sesi tidak valid")
		return
	}

	var serdik models.Serdik
	if err := config.DB.Where("user_id = ?", userID).First(&serdik).Error; err != nil {
		response.NotFound(c, "data serdik tidak ditemukan untuk akun ini")
		return
	}
	db := config.DB

	// ── Approved rewards / punishments (the only ones that count) ─────────
	var rewards []rpItem
	db.Raw(`
		SELECT ur.id, ri.kegiatan AS name, ur.point, ur.qty, ur.reward_date AS date, ur.notes
		FROM user_rewards ur
		JOIN reward_items ri ON ri.id = ur.reward_item_id
		WHERE ur.user_id = ? AND ur.status = 'approved'
		ORDER BY ur.reward_date DESC`, userID).Scan(&rewards)

	var punishments []rpItem
	db.Raw(`
		SELECT pl.id, pi.activity AS name, pl.point, pl.qty, pl.violation_date AS date, pl.notes
		FROM punishment_logs pl
		JOIN punishment_items pi ON pi.id = pl.punishment_item_id
		WHERE pl.user_id = ? AND pl.status = 'approved'
		ORDER BY pl.violation_date DESC`, userID).Scan(&punishments)

	var rewardTotal, punishmentTotal float64
	for _, r := range rewards {
		rewardTotal += r.Point * float64(maxInt(r.Qty, 1))
	}
	for _, p := range punishments {
		punishmentTotal += p.Point * float64(maxInt(p.Qty, 1))
	}

	// ── Assessment scores (avg per category, by serdik_id) ────────────────
	scores := scoreBlock{
		NilaiAkademik:   avgNilai(db, "penilaian_akademik", serdik.ID),
		NilaiJasmani:    avgNilai(db, "penilaian_jasmani", serdik.ID),
		NilaiMental:     avgNilai(db, "penilaian_mental", serdik.ID),
		RewardTotal:     round2(rewardTotal),
		PunishmentTotal: round2(punishmentTotal),
		NetAdjustment:   round2(rewardTotal - punishmentTotal),
	}

	// ── Attendance summary + records ──────────────────────────────────────
	absList, summary, err := h.kegiatan.GetAttendanceBySerdik(uint64(serdik.ID))
	if err != nil {
		absList = nil
		summary = nil
	}

	// ── Activity history (merged, newest first) ───────────────────────────
	history := make([]activityItem, 0, len(rewards)+len(punishments)+len(absList))
	for _, r := range rewards {
		history = append(history, activityItem{
			Type: "reward", Title: r.Name, Description: derefNote(r.Notes),
			Points: r.Point, Status: "approved", Timestamp: r.Date,
		})
	}
	for _, p := range punishments {
		history = append(history, activityItem{
			Type: "punishment", Title: p.Name, Description: derefNote(p.Notes),
			Points: -p.Point, Status: "approved", Timestamp: p.Date,
		})
	}
	for _, a := range absList {
		history = append(history, activityItem{
			Type: "attendance", Title: "Presensi", Description: string(a.Status),
			Points: 0, Status: string(a.Status), Timestamp: a.Datetime,
		})
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.After(history[j].Timestamp)
	})
	if len(history) > 50 {
		history = history[:50]
	}

	// ── Recommendations (rule-based, dynamic) ─────────────────────────────
	recommendations := buildRecommendations(scores, summary)

	response.OK(c, "success", gin.H{
		"serdik": gin.H{
			"id":        serdik.ID,
			"no_serdik": serdik.NoSerdik,
			"nama":      serdik.NamaLengkap,
			"pokjar_id": serdik.PokjarID,
		},
		"scores":           scores,
		"attendance":       summary,
		"rewards":          rewards,
		"punishments":      punishments,
		"recommendations":  recommendations,
		"activity_history": history,
	})
}

// GetKorsisInbox returns everything awaiting Korsis approval: pending rewards,
// pending punishments, and pending leave (izin) requests.
func (h *DashboardHandler) GetKorsisInbox(c *gin.Context) {
	rewards, err := h.reward.GetPending()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	punishments, err := h.punishment.GetPending()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	izin, err := h.kegiatan.GetPendingIzin()
	if err != nil {
		izin = nil
	}

	response.OK(c, "success", gin.H{
		"rewards":       rewards,
		"punishments":   punishments,
		"izin":          izin,
		"total_pending": len(rewards) + len(punishments) + len(izin),
	})
}

// ---- helpers ----

func avgNilai(db *gorm.DB, table string, serdikID uint) float64 {
	var v *float64
	db.Raw("SELECT AVG(nilai) FROM "+table+" WHERE serdik_id = ?", serdikID).Scan(&v)
	if v == nil {
		return 0
	}
	return round2(*v)
}

func buildRecommendations(s scoreBlock, summary interface{}) []string {
	recs := []string{}

	if s.PunishmentTotal > 0 && s.PunishmentTotal >= s.RewardTotal {
		recs = append(recs, "Tingkatkan kedisiplinan: akumulasi pelanggaran lebih tinggi dari penghargaan.")
	}
	if s.RewardTotal > s.PunishmentTotal && s.RewardTotal > 0 {
		recs = append(recs, "Pertahankan prestasi: akumulasi penghargaan Anda sangat baik.")
	}
	if s.NilaiAkademik > 0 && s.NilaiAkademik < 70 {
		recs = append(recs, "Fokus pada peningkatan nilai akademik yang masih di bawah standar.")
	}
	if s.NilaiJasmani > 0 && s.NilaiJasmani < 70 {
		recs = append(recs, "Tingkatkan kesemaptaan jasmani melalui latihan rutin.")
	}
	if s.NilaiMental > 0 && s.NilaiMental < 70 {
		recs = append(recs, "Perkuat aspek mental kepribadian dan kepemimpinan.")
	}
	if len(recs) == 0 {
		recs = append(recs, "Pertahankan performa Anda secara keseluruhan.")
	}
	return recs
}

func derefNote(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func round2(f float64) float64 {
	return float64(int64(f*100+0.5)) / 100
}
