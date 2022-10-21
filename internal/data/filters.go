//Filename: internal/data/filters.go

package data

import (
	"todo.joelical.net/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	SortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	//check page and page size parameters
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 1000, "page", "must be a maxinum of 1000")

	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maxinum of 100")
	//check that the sort parameter matches a value in the acceptable sort list
	v.Check(validator.In(f.Sort, f.SortList...), "sort", "invalid sort value")

}
