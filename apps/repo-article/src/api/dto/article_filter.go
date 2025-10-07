package dto

type ArticleFilters struct {
	Category    string `form:"category"`
	Subcategory string `form:"subcategory"`
	Search      string `form:"search"`
	Limit       uint   `form:"limit,default=10"`
	Offset      uint   `form:"offset,default=0"`
}
