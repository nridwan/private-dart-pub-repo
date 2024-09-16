package appmodel

// Response : Base response
type Response struct {
	ResponseSchema *ResponseSchema `json:"response_schema"`
	ResponseOutput interface{}     `json:"response_output"`
}

// Error : error detail
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ResponseSchema : ResponseSchema response
type ResponseSchema struct {
	ResponseCode    *string `json:"response_code"`
	ResponseMessage *string `json:"response_message"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

type PaginationResponsePagination struct {
	Page  *int `json:"page"`
	Total *int `json:"total"`
	Size  *int `json:"size"`
}

type PaginationResponseList struct {
	Pagination *PaginationResponsePagination `json:"pagination"`
	Content    interface{}                   `json:"content"`
}

type PaginationResponse struct {
	List *PaginationResponseList `json:"list"`
}
