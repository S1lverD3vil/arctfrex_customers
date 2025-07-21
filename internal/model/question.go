package model

type Survey struct {
	No       int    `json:"no"`
	Question string `json:"question"`
	Answer   bool   `json:"answer"`
}

type SurveyPayload struct {
	SurveyChecklist []Survey `json:"survey_checklist"`
	PpatkChecklist  []Survey `json:"ppatk_checklist"`
}

type SurveyResponse struct {
	Data SurveyPayload `json:"data"`
}
