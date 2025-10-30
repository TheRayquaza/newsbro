package controllers

import (
	//"net/http"

	//"github.com/gin-gonic/gin"
	"repo_feed/src/domain/services"
)

// RescoringController handles endpoints related to feed rescoring.
type RescoringController struct {
	rescoringService *services.FeedRescoringService
}

// NewRescoringController creates a new controller for feed rescoring endpoints.
func NewRescoringController(rescoringService *services.FeedRescoringService) *RescoringController {
	return &RescoringController{
		rescoringService: rescoringService,
	}
}

/*

// GetRescoringStatus godoc
// @Summary Get rescoring job status
// @Description Returns progress information about the feed rescoring job.
// @Tags Feed Rescoring
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /feed/rescoring/status [get]
// @Param Authorization header string true "Bearer <Add access token here>"
func (rc *RescoringController) GetRescoringStatus(c *gin.Context) {
	progress := rc.rescoringService.GetProgress()
	c.JSON(http.StatusOK, progress)
}

// ResetRescoringCursor godoc
// @Summary Reset rescoring cursor
// @Description Resets the rescoring job’s cursor to start rescoring all feeds from the beginning.
// @Tags Feed Rescoring
// @Produce json
// @Success 204 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /feed/rescoring/reset [post]
// @Param Authorization header string true "Bearer <Add access token here>"
func (rc *RescoringController) ResetRescoringCursor(c *gin.Context) {
	rc.rescoringService.ResetCursor()
	c.Status(http.StatusNoContent)
}
*/
