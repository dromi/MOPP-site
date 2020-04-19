package model

import "fmt"

const dbsql = `

CREATE DATABASE IF NOT EXISTS mopp;

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE IF NOT EXISTS performer (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT,

	CONSTRAINT name_unique UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS  availability (
	id SERIAL PRIMARY KEY,
	date DATE NOT NULL,
	performer INT NOT NULL REFERENCES performer(id),
	report BOOLEAN NOT NULL
);`

func CreateDB() {
	db := CreateDbConnection()
	_, err := db.Exec(dbsql)
	if err != nil {
		panic(err)
	}
	db.Close()
	fmt.Println("Success!")
}
