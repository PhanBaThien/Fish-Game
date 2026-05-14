package http

import (
	"net/http"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/domain"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CmsHandler xử lý các HTTP request cho Transactions, Settings, Stats, Search, Health
type CmsHandler struct {
	txUsecase      usecase.TransactionUsecase
	settingUsecase usecase.SettingUsecase
	statsUsecase   usecase.StatsUsecase
	searchUsecase  usecase.SearchUsecase
	tokenMaker     utils.TokenMaker
	startTime      time.Time
}

func NewCmsHandler(
	txUC usecase.TransactionUsecase,
	settingUC usecase.SettingUsecase,
	statsUC usecase.StatsUsecase,
	searchUC usecase.SearchUsecase,
	m utils.TokenMaker,
) *CmsHandler {
	return &CmsHandler{
		txUsecase:      txUC,
		settingUsecase: settingUC,
		statsUsecase:   statsUC,
		searchUsecase:  searchUC,
		tokenMaker:     m,
		startTime:      time.Now(),
	}
}

// RegisterRoutes gắn tất cả route CMS vào router group
func (h *CmsHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Health check — không cần auth
	router.GET("/health", h.HealthCheck)

	// Tất cả route bên dưới cần JWT
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(h.tokenMaker))
	{
		// BE-API-07: Lịch sử giao dịch
		protected.GET("/transactions", h.ListTransactions)

		// BE-API-08: Cài đặt hệ thống
		protected.GET("/settings", h.GetSettings)
		protected.PUT("/settings/:key", h.UpsertSetting)

		// BE-API-03: Dashboard stats
		protected.GET("/stats/overview", h.StatsOverview)
		protected.GET("/stats/timeseries", h.StatsTimeseries)

		// BE-API-09: Tìm kiếm tổng
		protected.GET("/search", h.Search)
	}
}

// HealthCheck godoc
// @Summary  Kiểm tra trạng thái server
// @Tags     system
// @Success  200 {object} map[string]interface{}
func (h *CmsHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime).Round(time.Second).String()
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status": "ok",
			"uptime": uptime,
		},
		"error": nil,
	})
}

// ListTransactions godoc
// @Summary  Lấy lịch sử giao dịch
// @Tags     transactions
// @Param    playerId query string false "Lọc theo UUID người chơi"
// @Param    type     query string false "shot|win|gift|deposit|withdraw|adjust"
// @Param    page     query int    false "Số trang"
// @Param    limit    query int    false "Số bản ghi mỗi trang"
func (h *CmsHandler) ListTransactions(c *gin.Context) {
	var req domain.ListTransactionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.txUsecase.ListTransactions(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// GetSettings godoc
// @Summary  Lấy tất cả cài đặt hệ thống
// @Tags     settings
func (h *CmsHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingUsecase.GetSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": settings, "error": nil})
}

// UpsertSetting godoc
// @Summary  Tạo hoặc cập nhật một cài đặt hệ thống
// @Tags     settings
// @Param    key  path string true "Tên cài đặt (key)"
func (h *CmsHandler) UpsertSetting(c *gin.Context) {
	key := c.Param("key")
	actorID, _ := c.Get("admin_id")
	actorIDStr, _ := actorID.(string)

	var req domain.UpsertSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting, err := h.settingUsecase.UpsertSetting(c.Request.Context(), key, actorIDStr, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": setting, "error": nil})
}

// StatsOverview godoc
// @Summary  Lấy các KPI tổng quan dashboard
// @Tags     stats
func (h *CmsHandler) StatsOverview(c *gin.Context) {
	resp, err := h.statsUsecase.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// StatsTimeseries godoc
// @Summary  Lấy dữ liệu biểu đồ theo thời gian
// @Tags     stats
// @Param    range query string false "7d|30d|90d (mặc định: 7d)"
func (h *CmsHandler) StatsTimeseries(c *gin.Context) {
	var req domain.TimeseriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.statsUsecase.GetTimeseries(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// Search godoc
// @Summary  Tìm kiếm tổng (players + rooms + fishes)
// @Tags     search
// @Param    q query string true "Từ khóa tìm kiếm (tối thiểu 2 ký tự)"
func (h *CmsHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if len(q) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Từ khóa tìm kiếm phải có ít nhất 2 ký tự"})
		return
	}

	result, err := h.searchUsecase.Search(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result, "error": nil})
}
