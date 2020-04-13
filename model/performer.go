package model

import (
	"database/sql"
	"fmt"
)

type Performer struct {
	ID       int
	Name     string
	Password string
}

func AddPerformer(performer Performer) {
	db := CreateDbConnection()
	defer db.Close()
	sqlStatement := `
	INSERT INTO performer (name, password)
	VALUES ($1, $2)
	RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, performer.Name, performer.Password).Scan(&id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}

func GetPerformer(id int) Performer {
	db := CreateDbConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM performer WHERE id=$1;`
	var performer Performer
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&performer.ID, &performer.Name, &performer.Password)
	switch err {
	case sql.ErrNoRows:
		return Performer{}
	case nil:
		return performer
	default:
		panic(err)
	}
}

func GetPerformerByName(name string) Performer {
	db := CreateDbConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM performer WHERE name=$1;`
	var performer Performer
	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&performer.ID, &performer.Name, &performer.Password)
	switch err {
	case sql.ErrNoRows:
		return Performer{}
	case nil:
		return performer
	default:
		panic(err)
	}
}

func ListPerformers() []Performer {
	db := CreateDbConnection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM performer ORDER BY ID")
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	performers := []Performer{}
	defer rows.Close()
	for rows.Next() {
		var performer Performer
		err = rows.Scan(&performer.ID, &performer.Name, &performer.Password)
		if err != nil {
			// handle this error
			panic(err)
		}
		performers = append(performers, performer)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return performers
}
