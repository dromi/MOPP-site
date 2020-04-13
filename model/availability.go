package model

import (
	"time"

	"github.com/lib/pq/hstore"
)

type Availability struct {
	ID      int
	Date    time.Time
	Reports hstore.Hstore
}

// func AddAvailability(availability Availability) {
// 	db := CreateDbConnection()
// 	defer db.Close()
// 	sqlStatement := `
// 	INSERT INTO performer (name)
// 	VALUES ($1)
// 	RETURNING id` 
// 	id := 0
// 	err := db.QueryRow(sqlStatement, performer.Name).Scan(&id)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("New record ID is:", id)
// }

// func GetPerformer(id int) Performer {
// 	db := CreateDbConnection()
// 	defer db.Close()
// 	sqlStatement := `SELECT * FROM performer WHERE id=$1;`
// 	var performer Performer
// 	row := db.QueryRow(sqlStatement, id)
// 	err := row.Scan(&performer.ID, &performer.Name)
// 	switch err {
// 	case sql.ErrNoRows:
// 		return Performer{}
// 	case nil:
// 		return performer
// 	default:
// 		panic(err)
// 	}
//}

func ListAvailability() []Availability {
	db := CreateDbConnection()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM availability")
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	avails := []Availability{}
	defer rows.Close()
	for rows.Next() {
		var avail Availability
		err = rows.Scan(&avail.ID, &avail.Date, &avail.Reports)
		if err != nil {
			// handle this error
			panic(err)
		}
		avails = append(avails, avail)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return avails
}
