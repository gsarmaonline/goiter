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

		Name      string       `json:"name"`
		GroupType ElementTypeT `json:"group_type"`
	}

	GroupMember struct {
		BaseModelWithoutUser

		GroupID   uint         `json:"group_id"`
		GroupType ElementTypeT `json:"group_type"`
		MemberID  uint         `json:"member_id"`
	}

	GroupParent struct {
		BaseModelWithoutUser

		GroupID       uint `json:"group_id"`
		ParentGroupID uint `json:"parent_group_id"`
	}

	// This struct is used to fetch groups for a given model
	GroupFetcher struct {
		tx    *gorm.DB
		model UserOwnedModel
	}
)

func (group *Group) GetGroupsTillRoot(tx *gorm.DB, existingGroups []*Group) (err error) {
	tmpGroups := []*Group{}
	if db := tx.Where("parent_group_id = ?", group.ID).Find(&tmpGroups); db.Error != nil {
		err = db.Error
		return
	}
	existingGroups = append(existingGroups, tmpGroups...)
	for _, g := range tmpGroups {
		if err = g.GetGroupsTillRoot(tx, existingGroups); err != nil {
			return
		}
	}
	return
}

func NewGroupFetcher(tx *gorm.DB, model UserOwnedModel) *GroupFetcher {
	return &GroupFetcher{
		tx:    tx,
		model: model,
	}
}

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
	if len(groupIds) > 0 {
		if db := gf.tx.Where("id IN ?", groupIds).Find(&groups); db.Error != nil {
			err = db.Error
			return
		}
	}
	for _, g := range groups {
		if err = g.GetGroupsTillRoot(gf.tx, groups); err != nil {
			return
		}
	}
	return
}
