package handler

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserItem struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type LoginResp struct {
	Token string    `json:"token"`
	User  *UserItem `json:"user"`
}

type GetUserListReq struct{}

type GetUserListResp struct {
	List  []*UserItem `json:"list"`
	Total int64       `json:"total"`
}

type GetUserResp struct {
	User *UserItem `json:"user"`
}
