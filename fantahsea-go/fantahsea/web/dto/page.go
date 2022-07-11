package dto

type Paging struct {

	Limit int `json:"limit"`
	Page int `json:"page"`
	Total int `json:"total"`

}