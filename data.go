package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	data []Credential
)

type Database struct {
	Connection *sql.DB
}

type Credential struct {
	Id          int
	Username    string
	Password    string
	Host        string
	Information string
}

func createDbConnection() (*Database, error) {
	db, err := sql.Open("sqlite3", "db.sql")
	if err != nil {
		log.Fatal(err)
	}

	return &Database{Connection: db}, nil
}

func (d *Database) Close() {
	d.Connection.Close()
}

func createDbSchema(d *Database) {
	const createStmt string = `
  	create table if not exists Credentials (
		id integer primary key autoincrement not null,
		username text not null,
		password text not null,
		host text not null,
		information text not null
	);
	`
	_, err := d.Connection.Exec(createStmt)
	if err != nil {
		log.Println("", err)
		return
	}
}

func insertCredential(d *Database, cred Credential) {
	tx, err := d.Connection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into credentials(username, password, host, information) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cred.Username, cred.Password, cred.Host, cred.Information)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("insert should've worked we might be fucked tho")
}

func getCredentials(d *Database) ([]Credential, error) {
	var id int
	var username, password, host, information string

	rows, err := d.Connection.Query("select * from credentials")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &username, &password, &host, &information)
		if err != nil {
			log.Fatal(err)
		}

		data = append(data, Credential{
			Id:          id,
			Username:    username,
			Password:    password,
			Host:        host,
			Information: information,
		})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

func getCredentialById(d *Database, uid int) (Credential, error) {
	var id int
	var username, password, host, information string

	stmt, err := d.Connection.Prepare("select * from credentials where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(uid).Scan(&id, &username, &password, &host, &information)
	if err != nil {
		log.Fatal(err)
	}

	return Credential{Id: id, Username: username, Password: password, Host: host, Information: information}, nil
}
