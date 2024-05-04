package data

import (
	"goproject/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"
	"fmt"

)


type Group struct{
	Id 					int		`json:"id"`
	Name 				string	`json:"name"`
	Num_of_members 		int		`json:"num_of_members"`		
} 

func ValidateGroup(v *validator.Validator, group *Group){
	v.Check(group.Name != "", "name", "must be provided")
	v.Check(group.Num_of_members != 0, "num_of_members", "must be greater than 0")
}


type GroupModel struct {
	DB *sql.DB
}


func (g GroupModel) Insert(group *Group) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system-generated data.
	query := `
		INSERT INTO groups(name, num_of_members)
		VALUES ($1, $2)
		RETURNING group_id, name, num_of_members;`

	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{group.Name, group.Num_of_members}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the movie struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return g.DB.QueryRowContext(ctx, query, args...).Scan(&group.Id, &group.Name, &group.Num_of_members)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (g GroupModel) Get(id int64) (*Group, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT group_id, name, num_of_members
		FROM groups
		WHERE group_id = $1;`

	var group Group
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := g.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&group.Id, &group.Name, &group.Num_of_members)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &group, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (g GroupModel) Update(group *Group) error {
	query := `
		UPDATE groups
		SET name = $1, num_of_members = $2
		WHERE group_id = $3
		RETURNING group_id, name, num_of_members;
		`

	args := []interface{}{group.Name, group.Num_of_members, group.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return g.DB.QueryRowContext(ctx, query, args...).Scan(&group.Id, &group.Name, &group.Num_of_members)
}


func (g GroupModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the movie ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM groups
		WHERE group_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := g.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}


func (g GroupModel) GetAll(name string, num_of_members int, filters Filters) ([]*Group, Metadata, error) {
	// Construct the SQL query to retrieve all movie records.
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), group_id, name, num_of_members
		FROM groups
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (num_of_members = $2 OR $2 = 1)
		ORDER BY %s %s, group_id
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{name, num_of_members, filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := g.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	// Initialize an empty slice to hold the movie data.
	groups := []*Group{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var group Group
		// Scan the values from the row into the Movie struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&group.Id,
			&group.Name,
			&group.Num_of_members,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Movie struct to the slice.
		groups = append(groups, &group)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// If everything went OK, then return the slice of movies.
	return groups, metadata, nil
}
