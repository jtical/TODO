//Filename: internal/data/list.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	//Create a context. time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memory leaks
	defer cancel()
	//collect the data fields into a slice
	args := []interface{}{
		list.Name,
		list.Task,
		list.Status,
	}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&list.ID, &list.CreatedAt, &list.Version)
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
	//Create a context. time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memory leaks
	defer cancel()
	//execute the query using QueryRow(.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
	//Create a context. time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memory leaks
	defer cancel()

	args := []interface{}{
		list.Name,
		list.Task,
		list.Status,
		list.ID,
		list.Version,
	}
	//check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&list.Version)
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
	//Create a context. time starts when context is created
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//cleanup to prevent memory leaks
	defer cancel()

	//execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
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

// the GetAll() method returns a list of all the list sorted by id
func (m ListModel) GetAll(name string, status string, filters Filters) ([]*List, Metadata, error) {
	//construct the query to return all schools
	//make query into formated string to be able to sort by field and asc or dec dynaimicaly
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),id, created_at, name, task, status, version
		FROM lists
		WHERE (to_tsvector('simple',name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple',status) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

	//create a 3 second timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//execute the query
	args := []interface{}{name, status, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	//close the result set
	defer rows.Close()
	//store total records
	totalRecords := 0
	//intialize an empty slice to hold the list data
	lists := []*List{}
	//iterate over the rows in the result set
	for rows.Next() {
		var list List
		//scan the values from the row into the List struct
		err := rows.Scan(
			&totalRecords,
			&list.ID,
			&list.CreatedAt,
			&list.Name,
			&list.Task,
			&list.Status,
			&list.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//add the list to our slice
		lists = append(lists, &list)
	}
	//check if any errors occured while proccessing the result set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	//return the result set. the slice of lists
	return lists, metadata, nil
}
