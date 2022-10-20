//Filename: internal/data/list.go

package data

import (
	"time"

	"todo.joelical.net/internal/validator"
)

type List struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	Version   int32     `json:"version"`
}

func ValidateList(v *validator.Validator, list *List) {
	// use the check() method to execute our validation checks
	v.Check(list.Name != "", "name", "must be provied")
	v.Check(len(list.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(list.Task != "", "task", "must be provied")
	v.Check(len(list.Task) <= 800, "name", "must not be more than 200 bytes long")

	v.Check(list.Status != "", "status", "must be provied")
	v.Check(len(list.Status) <= 300, "status", "must not be more than 200 bytes long")

}
