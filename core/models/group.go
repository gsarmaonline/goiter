package models

type (
	Group struct {
		BaseModelWithUser

		Name      string `json:"name"`
		GroupType string `json:"group_type"`
	}

	GroupMember struct {
		BaseModelWithoutUser

		GroupID  uint `json:"group_id"`
		MemberID uint `json:"member_id"`
	}

	GroupParent struct {
		BaseModelWithoutUser

		GroupID       uint `json:"group_id"`
		ParentGroupID uint `json:"parent_group_id"`
	}
)
