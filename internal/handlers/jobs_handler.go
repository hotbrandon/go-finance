package handlers

import (
	"net/http"

	"go-finance/internal/jobs"
	"go-finance/internal/models"

	"github.com/gin-gonic/gin"
)

type JobsHandler struct {
	Sched *jobs.Scheduler
}

func NewJobsHandler(s *jobs.Scheduler) *JobsHandler {
	return &JobsHandler{Sched: s}
}

func (h *JobsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	grp := rg.Group("/jobs")
	grp.GET("/", h.ListJobs)
}

// ListJobs godoc
// @Summary List scheduled jobs
// @Description Returns scheduled jobs with next/prev run times and metadata.
// @Tags Jobs
// @Produce json
// @Success 200 {object} models.JobsAPIResponse
// @Failure 500 {object} models.APIResponse // Error responses can still use the generic one
// @Router /jobs [get]
func (h *JobsHandler) ListJobs(c *gin.Context) {
	jobs := h.Sched.ListJobs()
	c.JSON(http.StatusOK, models.NewSuccessResponse(jobs))
}
