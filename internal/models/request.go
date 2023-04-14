package models

type ResumeRequest struct {
	CvText string `json:"cv_text" form:"cv_text"`
	NTech  int    `json:"n_tech" form:"n_tech"`
	NProf  int    `json:"n_prof" form:"n_prof"`
}

type LetterRequest struct {
	Resume  string `json:"resume" form:"resume"`
	Vacancy string `json:"vacancy" form:"vacancy"`
}
