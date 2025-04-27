package model

type MovieRequest struct {
	Hero      string `json:"hero"`
	Genre     string `json:"genre"`
	Language  string `json:"language"`
	Reception string `json:"reception"`
}
type ParsedMovie struct {
	Title    string
	Year     string
	Language string
}
