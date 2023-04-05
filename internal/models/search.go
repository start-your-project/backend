package models

type Technology struct {
	TechnologyName  string  `json:"technology_name"`
	Distance        float32 `json:"distance"`
	Professionalism float32 `json:"professionalism"`
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
}
