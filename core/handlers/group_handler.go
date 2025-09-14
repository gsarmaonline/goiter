package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsarmaonline/goiter/core/models"
	"gorm.io/gorm"
)

type (
	GroupHandler struct {
		db      *gorm.DB
		handler *Handler
	}

	GroupUpdateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		PlanID      uint   `json:"plan_id"`
	}
)

func NewGroupHandler(handler *Handler) *GroupHandler {
	return &GroupHandler{handler: handler, db: handler.db}
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	groups := []*models.Group{}
	if err := h.handler.UserScopedDB(c).Find(&groups).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Groups not found"})
		return
	}
	h.handler.WriteSuccess(c, groups)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	var group models.Group
	if err := h.handler.GetModelFromUrl(c, &group, DefaultUrlKeyName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	h.handler.WriteSuccess(c, group)
}

func (h *GroupHandler) GetGroupAncestors(c *gin.Context) {
	var group models.Group
	if err := h.handler.GetModelFromUrl(c, &group, DefaultUrlKeyName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	groups := []*models.Group{}

	// Use unrestricted DB for ancestor lookup to find all ancestors regardless of owner
	if err := group.GetGroupsAncestors(h.db, &groups); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.handler.WriteSuccess(c, groups)
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var (
		group models.Group
	)
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.handler.CreateWithUser(c, &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.handler.WriteSuccess(c, group)
}

func (h *GroupHandler) AddGroupMember(c *gin.Context) {
	var (
		groupMember models.GroupMember
		group       models.Group
	)
	if err := h.handler.GetModelFromUrl(c, &group, DefaultUrlKeyName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	if err := c.ShouldBindJSON(&groupMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Explicitly set the GroupID from the URL parameter
	groupMember.GroupID = group.ID

	if err := h.handler.CreateWithUser(c, &groupMember); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.handler.WriteSuccess(c, groupMember)
}

func (h *GroupHandler) RemoveGroupMember(c *gin.Context) {
	var (
		groupMember models.GroupMember
		group       models.Group
	)
	if err := h.handler.GetModelFromUrl(c, &group, DefaultUrlKeyName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	if err := c.ShouldBindJSON(&groupMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.handler.UserScopedDB(c).Where("group_id = ? AND member_id = ? AND member_type = ?",
		group.ID, groupMember.MemberID, groupMember.MemberType).First(&groupMember).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Group member not found"})
		return
	}
	if err := h.handler.DeleteWithUser(c, &groupMember); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.handler.WriteSuccess(c, groupMember)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {

	var group models.Group
	if err := h.handler.GetModelFromUrl(c, &group, DefaultUrlKeyName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	if err := h.handler.DeleteWithUser(c, &group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.handler.WriteSuccess(c, group)
}
