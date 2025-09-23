package dto

type ArticleFilters struct {
	Category    string `form:"category"`
	Subcategory string `form:"subcategory"`
	Search      string `form:"search"`
	Limit       int    `form:"limit,default=10"`
	Offset      int    `form:"offset,default=0"`
}
