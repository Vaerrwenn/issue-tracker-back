package controllers

import (
	"issue-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IssuesCreateForm takes User's input on Issue Create form.
type IssuesCreateForm struct {
	Title    string `form:"title" binding:"required"`
	Body     string `form:"description" binding:"required"`
	Severity string `form:"severity" binding:"required"`
}

// CreateIssueHandler handles issue creation.
func CreateIssueHandler(c *gin.Context) {
	// Bind the input to variable.
	var input IssuesCreateForm
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Begin save Issue process.
	HeaderUserID := c.Request.Header.Get("userID")
	userID, err := strconv.Atoi(HeaderUserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	issue := models.Issue{
		UserID:   userID,
		Title:    input.Title,
		Body:     input.Body,
		Status:   "1",
		Severity: input.Severity,
	}
	if err := issue.SaveIssue(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    1,
		"issueID": issue.ID,
	})

	return
}