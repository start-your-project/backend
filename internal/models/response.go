package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseTechnologies struct {
	Status       int          `json:"status"`
	PositionData PositionData `json:"position_data"`
	InBase       string       `json:"in_base"`
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
	Status      int          `json:"status"`
	Top         []Profession `json:"professions"`
	TipsToLearn string       `json:"tips_to_learn"`
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

type ResponseTechs struct {
	Status int      `json:"status"`
	Techs  []string `json:"techs"`
}

type ResponseCheck struct {
	Status int  `json:"status"`
	IsWork bool `json:"is_work"`
}
