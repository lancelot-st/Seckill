package datamoudles

type User struct {
	ID           int64  `json:"id " from:"ID" sql:"ID"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
	UserName     string `json:"userName" from:"userName" sql:"userName"`
	HashPassword string `json:"-" from:"passWord" sql:"passWord"`
}
