package dto

type Paging struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Total int `json:"total"`
}

/* Build Paging for response */
func BuildResPage(reqPage *Paging, total int) *Paging {
	return &Paging{
		Limit: reqPage.Limit,
		Page:  reqPage.Page,
		Total: total,
	}
}
