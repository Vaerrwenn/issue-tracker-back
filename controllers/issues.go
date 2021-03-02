package controllers

import (
	"issue-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IssueCreateForm takes User's input on Issue Create form.
type IssueCreateForm struct {
	Title    string `form:"title" binding:"required"`
	Body     string `form:"description" binding:"required"`
	Severity string `form:"severity" binding:"required"`
}

// IssueUpdateForm is for Updating.
type IssueUpdateForm struct {
	Title    string `form:"title" binding:"required"`
	Body     string `form:"description" binding:"required"`
	Status   string `form:"status" binding:"required"`
	Severity string `form:"severity" binding:"required"`
}

// returnErrorAndAbort returns a JSON with "error": errorText in it. After that,
// it aborts and stop the running function.
//
// Takes Gin's context, the HTTP Code, and error text.
func returnErrorAndAbort(ctx *gin.Context, code int, errorText string) {
	ctx.JSON(code, gin.H{
		"error": errorText,
	})
	ctx.Abort()
}

// CreateIssueHandler handles issue creation.
func CreateIssueHandler(c *gin.Context) {
	// Bind the input to variable.
	var input IssueCreateForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Begin save Issue process.
	HeaderUserID := c.Request.Header.Get("userID")
	userID, err := strconv.Atoi(HeaderUserID)

	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	issue := models.Issue{
		UserID:   userID,
		Title:    input.Title,
		Body:     input.Body,
		Status:   "1",
		Severity: input.Severity,
	}
	if err := issue.ValidateIssue(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := issue.SaveIssue(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": issue.ID,
		"msg":  "Data succesfully created.",
	})

	return
}

// IndexIssueHandler shows ALL issues. 😢😢😢😢😢
func IndexIssueHandler(c *gin.Context) {
	var issue models.Issue
	result, err := issue.IndexIssues()

	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qty":  len(*result),
		"data": result,
	})
	return
}

// ShowIssueHandler fetch ONE issue by ID.
//
// /v1/protected/issue/:id
//
// For example: /v1/protected/issue/1
func ShowIssueHandler(c *gin.Context) {
	var issue models.Issue

	id := c.Param("id")
	result, err := issue.FindIssueByID(id)

	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// UpdateIssueHandler is used for updating. Works similarly to CreateIssueHandler.
//
// Only the poster and developer can Update an Issue.
func UpdateIssueHandler(c *gin.Context) {
	var issue models.Issue
	var user models.User

	// Get Issue ID from param
	// For example, update/3
	// get that 3.
	id := c.Param("id")
	source, err := issue.FindFirstIssueByID(id)

	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	var input IssueUpdateForm
	if err := c.ShouldBind(&input); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	HeaderUserID := c.Request.Header.Get("userID")
	userID, err := strconv.Atoi(HeaderUserID)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "couldn't parse UserID into Int")
		return
	}

	userRole, err := user.GetUserRoleByID(userID)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	// Checks whether User with the same ID can do this request.
	if userID != source.UserID {
		// Checks whether a User with different ID as the poster/author is a Developer.
		if userRole != "2" {
			returnErrorAndAbort(c, http.StatusBadRequest, "User is unauthorized for this request.")
			return
		}
	}

	issue = models.Issue{
		Title:    input.Title,
		Body:     input.Body,
		Status:   input.Status,
		Severity: input.Severity,
	}

	if err := issue.ValidateIssue(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := issue.UpdateIssue(source); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": idNum,
		"msg":  "Data has been updated succesfully.",
	})

	return
}

// DeleteIssueHandler deletes an Issue by ID.
func DeleteIssueHandler(c *gin.Context) {
	var issue models.Issue
	var user models.User

	id := c.Param("id")
	source, err := issue.FindFirstIssueByID(id)

	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	HeaderUserID := c.Request.Header.Get("userID")
	userID, err := strconv.Atoi(HeaderUserID)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "couldn't parse UserID into Int")
		return
	}

	userRole, err := user.GetUserRoleByID(userID)
	if err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, err.Error())
		return
	}

	if userID != source.UserID {
		if userRole != "2" {
			returnErrorAndAbort(c, http.StatusBadRequest, "User is unauthorized for this request.")
			return
		}
	}

	if err := source.DeleteIssue(); err != nil {
		returnErrorAndAbort(c, http.StatusBadRequest, "Unable to delete issue")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "deleted",
		"msg":  "issue is deleted successfully",
	})
	return
}
