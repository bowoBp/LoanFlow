package dto

type GetListQuery struct {
	PerPage int    `json:"perPage"`
	Page    int    `json:"page"`
	Search  string `json:"search"`
}

type Sorting struct {
	CreatedAt string `json:"createdAt" validate:"enum=ASC DESC"`
}
