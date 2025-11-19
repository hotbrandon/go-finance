package handlers

import (
	"net/http"

	"go-finance/internal/jobs"

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
// @Success 200 {array} jobs.JobStatus
// @Failure 500 {object} map[string]string
// @Router /jobs [get]
func (h *JobsHandler) ListJobs(c *gin.Context) {
	jobs := h.Sched.ListJobs()
	c.JSON(http.StatusOK, jobs)
}
