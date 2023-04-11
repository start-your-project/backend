package models

type ResumeRequest struct {
	CvText string `json:"cv_text"`
	NTech  int    `json:"n_tech"`
	NProf  int    `json:"n_prof"`
}
