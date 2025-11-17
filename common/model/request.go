package model

type UserInfoUpdate struct {
	NickName    string `json:"nickname"`
	Description string `json:"description"`
	Gender      string `json:"gender"`
}

type UserAvatatUpdate struct {
	Avatar string `json:"avatar"`
}

type CommentReq struct {
	VideoID  string `json:"video_id"`
	Content  string `json:"content"`
	ParentID string `json:"parent_id"`
}

type CommentLikeReq struct {
	CommentID uint64 `json:"comment_id"`
	Action    string `josn:"action"`
}
