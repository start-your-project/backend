package models

type ProfileUserDTO struct {
	Name   string `json:"username" form:"username"`
	Email  string `json:"email" form:"email"`
	Avatar string `json:"avatar" form:"avatar"`
}

type EditProfileDTO struct {
	Name     string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type EmailUserDTO struct {
	Email string `json:"email" form:"email"`
}

type LikeDTO struct {
	Name string `json:"name" form:"name"`
}

type Favorite struct {
	ID            int64  `json:"id" form:"id"`
	Name          string `json:"name" form:"name"`
	CountAll      int64  `json:"count_all" form:"count_all"`
	CountFinished int64  `json:"count_finished" form:"count_finished"`
}

type Recommend struct {
	Learned []string `json:"learned" form:"learned"`
	ToLearn []string `json:"to_learn" form:"to_learn"`
}

type LinkDTO struct {
	Link string `json:"link" form:"link"`
}
