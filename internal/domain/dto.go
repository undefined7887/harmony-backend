package domain

type IdParam struct {
	ID string `uri:"id" binding:"id"`
}

type PaginationQuery struct {
	Offset int64 `form:"offset" binding:"min=0"`
	Limit  int64 `form:"limit" binding:"min=0"`
}
