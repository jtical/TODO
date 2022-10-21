//Filename: internal/data/list.go

package data

import (
	"database/sql"
	"errors"
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

// define a ListModel which wraps a sql.db connection pool
type ListModel struct {
	DB *sql.DB
}

// Insert() allows us to creat a new list
func (m ListModel) Insert(list *List) error {
	query := `
		INSERT INTO lists (name, task, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
	`
	//collect the data fields into a slice
	args := []interface{}{
		list.Name,
		list.Task,
		list.Status,
	}
	return m.DB.QueryRow(query, args...).Scan(&list.ID, &list.CreatedAt, &list.Version)
}

// Get() allows us to retrieve a specfic list
func (m ListModel) Get(id int64) (*List, error) {
	//ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	//create the query
	query := `
		SELECT id, created_at, name, task, status,version
		FROM lists
		WHERE id = $1
	`
	//declare a list variable to hold the returned data
	var list List
	//execute the query using QueryRow(
	err := m.DB.QueryRow(query, id).Scan(
		&list.ID,
		&list.CreatedAt,
		&list.Name,
		&list.Task,
		&list.Status,
		&list.Version,
	)
	//handle any errors
	if err != nil {
		//check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//success
	return &list, nil
}

// Update() allows us edit a specific list
// optimistic locking on the version # enssure version has not changed from when i first read it to when will write it back with new changes
func (m ListModel) Update(list *List) error {
	//create a query using the newly updated data
	query := `
		UPDATE lists
		SET name = $1,
			task = $2,
			status = $3,
			version = version + 1
		WHERE id = $4
		AND version = $5 
		RETURNING version
	`
	args := []interface{}{
		list.Name,
		list.Task,
		list.Status,
		list.ID,
		list.Version,
	}
	//check for edit conflicts
	err := m.DB.QueryRow(query, args...).Scan(&list.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

// Delete() allows us to remove a specific list
func (m ListModel) Delete(id int64) error {
	//check if the id exist
	if id < 1 {
		return ErrRecordNotFound
	}
	//create the delete query
	query := `
		DELETE FROM lists
		WHERE id = $1
	`
	//execute the query
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	//check how many rows affected by the delte operation. we will use the RowsAffected() on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	//check to see if zero rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
