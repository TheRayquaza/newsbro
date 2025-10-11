package dto

type ArticleHistoryFilters struct {
	Limit  uint `form:"limit,default=10"`
	Offset uint `form:"offset,default=0"`
}
