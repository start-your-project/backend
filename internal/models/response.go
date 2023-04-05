package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//type Technology struct {
//	Name            string  `json:"name"`
//	Distance        float32 `json:"distance"`
//	Professionalism float32 `json:"professionalism"`
//}

type ResponseTechnologies struct {
	Status       int          `json:"status"`
	PositionData PositionData `json:"position_data"`
}

type ResponseUserProfile struct {
	Status   int             `json:"status"`
	UserData *ProfileUserDTO `json:"user"`
}

type ResponseFavorites struct {
	Status        int        `json:"status"`
	FavoritesData []Favorite `json:"favorites"`
}
