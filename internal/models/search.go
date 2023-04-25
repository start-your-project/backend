package models

type Technology struct {
	TechnologyName  string  `json:"technology_name" form:"technology_name"`
	Distance        float32 `json:"distance"`
	Professionalism float32 `json:"professionalism"`
	HardSkill       bool    `json:"hard_skill"`
}

//easyjson:json
type PositionData struct {
	JobName          string       `json:"job_name"`
	TechnologyNumber int          `json:"technology_number"`
	Additional       []Technology `json:"additional"`
}

//easyjson:json
type Profession struct {
	Profession string `json:"profession"`
	InBase     string `json:"in_base"`
}

//easyjson:json
type Professions struct {
	Profession []string `json:"professions"`
}

type RespProfessions struct {
	Techs      string           `json:"techs"`
	JobNumber  int              `json:"job_number"`
	Additional []RespProfession `json:"additional"`
}

type RespProfession struct {
	JobName string `json:"job_name"`
	Percent int    `json:"percent"`
}

//easyjson:json
type SearchTechs struct {
	SearchText string `json:"search_text" form:"search_text"`
}

//easyjson:json
type Techs struct {
	Techs []string `json:"techs"`
}
