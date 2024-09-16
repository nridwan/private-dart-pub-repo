package appmodel

import "strconv"

type GetListRequest struct {
	Page   int `validate:"number,min=1"`
	Limit  int `validate:"number,required"`
	Search string
}

func NewGetListRequest(page string, limit string, search string) *GetListRequest {
	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	return &GetListRequest{
		Page:   pageInt,
		Limit:  limitInt,
		Search: search,
	}
}
