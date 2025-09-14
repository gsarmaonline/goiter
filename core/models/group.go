package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	UserElementType     ElementTypeT = "user_element"
	ResourceElementType ElementTypeT = "resource_element"
)

type (
	ElementTypeT string

	Group struct {
		BaseModelWithUser

		Name        string       `json:"name"`
		Description string       `json:"description"`
		MemberType  ElementTypeT `json:"member_type"`
	}

	GroupMember struct {
		BaseModelWithUser

		GroupID    uint         `json:"group_id"`
		MemberType ElementTypeT `json:"member_type"`
		MemberID   uint         `json:"member_id"`
	}

	// This struct is used to fetch groups for a given model
	GroupFetcher struct {
		tx    *gorm.DB
		model UserOwnedModel
	}
)

func (group Group) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "Group",
		ScopeType: AccountScopeType,
	}
}

func (groupMember GroupMember) GetConfig() ModelConfig {
	return ModelConfig{
		Name:      "GroupMember",
		ScopeType: AccountScopeType,
	}
}

func (group *Group) GetGroupsAncestors(tx *gorm.DB, existingGroups *[]*Group) (err error) {
	return group.getGroupsAncestorsRecursive(tx, existingGroups, make(map[uint]bool))
}

func (group *Group) getGroupsAncestorsRecursive(tx *gorm.DB, existingGroups *[]*Group, visitedGroups map[uint]bool) (err error) {
	if _, ok := visitedGroups[group.ID]; ok {
		return
	}
	visitedGroups[group.ID] = true

	// Query for group members where this group is a member of other groups
	belongsToGroups := []*GroupMember{}
	if db := tx.Where("member_type = ? AND member_id = ?", "Group", group.ID).Find(&belongsToGroups); db.Error != nil {
		err = db.Error
		return
	}

	// Extract the group IDs that this group belongs to
	ancestorGroupIDs := []uint{}
	for _, bg := range belongsToGroups {
		ancestorGroupIDs = append(ancestorGroupIDs, bg.GroupID)
	}

	// Only proceed if we found any ancestor groups
	ancestorGroups := []*Group{}
	if len(ancestorGroupIDs) == 0 {
		return
	}
	if db := tx.Where("id IN ?", ancestorGroupIDs).Find(&ancestorGroups); db.Error != nil {
		err = db.Error
		return
	}

	// Append the found ancestor groups to our result list
	*existingGroups = append(*existingGroups, ancestorGroups...)

	// Recursively find ancestors of each ancestor group
	for _, ancestorGroup := range ancestorGroups {
		if err = ancestorGroup.getGroupsAncestorsRecursive(tx, existingGroups, visitedGroups); err != nil {
			return
		}
	}

	return
}

// NewGroupFetcher returns a new GroupFetcher
func NewGroupFetcher(tx *gorm.DB, model UserOwnedModel) *GroupFetcher {
	return &GroupFetcher{
		tx:    tx,
		model: model,
	}
}

// GetGroups returns all groups for a given model
func (gf *GroupFetcher) GetGroups() (groups []*Group, err error) {
	if gf.model.GetConfig().Name == "Group" {
		err = errors.New("cannot fetch groups for Group model by GroupFetcher")
		return
	}
	groupMembers := []GroupMember{}
	if db := gf.tx.Where("group_id = ? AND group_type = ?", gf.model.GetID(), gf.model.GetConfig().Name).Find(&groupMembers); db.Error != nil {
		err = db.Error
		return
	}
	groupIds := []uint{}
	for _, gm := range groupMembers {
		groupIds = append(groupIds, gm.GroupID)
	}
	if db := gf.tx.Where("id IN ?", groupIds).Find(&groups); db.Error != nil {
		err = db.Error
		return
	}
	for _, g := range groups {
		if err = g.GetGroupsAncestors(gf.tx, &groups); err != nil {
			return
		}
	}
	return
}
