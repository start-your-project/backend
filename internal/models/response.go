package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

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

type ResponseTop struct {
	Status int          `json:"status"`
	Top    []Profession `json:"professions"`
}

type ResponseProfessions struct {
	Status      int      `json:"status"`
	Professions []string `json:"professions"`
}

type ResponseProfessionsWithTechnology struct {
	Status      int              `json:"status"`
	Professions *RespProfessions `json:"professions"`
}
type ResponseResume struct {
	Status    int         `json:"status"`
	Recommend []Recommend `json:"recommend"`
}

type ResponseFinished struct {
	Status   int      `json:"status"`
	Finished []string `json:"finished"`
}

type ResponseLetter struct {
	Status      int    `json:"status"`
	CoverLetter string `json:"cover_letter"`
}
