package main

import (
  "fmt"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
  users []User
  hosts []Host
	creds []Credential
)

type Database struct {
	Connection *sql.DB
}

type User struct {
  Id int
  Username string
  PasswordHash string
  Admin bool
}

type Host struct {
	Id          int
	Hostname    string
	IpAddress   string
	Os          string
	Information string
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
  create table if not exists Users (
    id integer primary key autoincrement not null,
    username text not null unique, 
    passwordhash text not null,
    admin bool not null
  );
	create table if not exists Hosts (
		id integer primary key autoincrement not null,
		hostname text not null,
		ipaddress text not null,
		os text not null,
		information text not null
	);
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
		log.Print(err)
		return
	}
}

func getUsers(d *Database) ([]User, error) {
  var id int
  var username, passwordhash string
  var admin bool

  rows, err := d.Connection.Query("select * from users")
  if err != nil {
    log.Print(err)
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(&id, &username, &passwordhash, &admin)
    if err != nil {
      log.Print(err)
    }

    users = append(users, User{
      Id: id,
      Username: username,
      PasswordHash: passwordhash,
      Admin: admin,
    })
  }
  
  err = rows.Err()
  if err != nil {
    log.Print(err)
  }

  return users, nil
}

func createAdmin(d *Database) (string, error) {
  var password string

  stmt, err := d.Connection.Prepare("delete from Users where username = ?")
  if err != nil {
    log.Print(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec("admin")
  if err != nil {
    log.Print(err)
  }

  stmt_2, err := d.Connection.Prepare("insert into users(username, passwordhash, admin) values(?, ?, ?)")
  if err != nil {
    log.Print(err)
  }
  defer stmt_2.Close()

  password, err = generateRandomString(12)
  if err != nil {
    log.Print(err)
    return "", nil
  }

  hashedPassword, err := hashPassword(password)
  if err != nil {
    log.Print(err)
  }

  _, err = stmt_2.Exec("admin", hashedPassword, true)
  if err != nil {
    log.Print(err)
  }

  return password, nil
}

func insertUser(d *Database, user User) {
	tx, err := d.Connection.Begin()
	if err != nil {
		log.Print(err)
	}

	stmt, err := tx.Prepare("insert into users(username, passwordhash, admin) values(?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

  hashedPassword, err := hashPassword(user.PasswordHash)
  if err != nil {
    log.Print(err)
    return
  }

  fmt.Println(hashPassword)

	_, err = stmt.Exec(user.Username, hashedPassword, user.Admin)
	if err != nil {
		log.Print(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}
}

func getUserById(d *Database, id int) (User, error) {
  var username, password string
  var admin bool

  stmt, err := d.Connection.Prepare("select * from Users where id = ?")
  if err != nil {
    log.Print(err)
  }
  defer stmt.Close()

  err = stmt.QueryRow(id).Scan(&id, &username, &password, &admin)
  if err != nil {
    log.Print(err)
  }

  return User{
    Id: id,
    Username: username,
    PasswordHash: password,
    Admin: admin,
  }, nil
}

func getUserByUsername(d *Database, username string) (User, error) {
  var id int
  var password string
  var admin bool

  stmt, err := d.Connection.Prepare("select * from Users where username = ?")
  if err != nil {
    log.Print(err)
  }
  defer stmt.Close()

  err = stmt.QueryRow(username).Scan(&id, &username, &password, &admin)
  if err != nil {
    log.Print(err)
  }

  return User{
    Id: id,
    Username: username,
    PasswordHash: password,
    Admin: admin,
  }, nil
}

func getUserIdByUsername(d *Database, username string) (int, error) {
  var id int
  
  stmt, err := d.Connection.Prepare("select id from Users where username = ?")
  if err != nil {
    log.Print(err)
  }
  defer stmt.Close()

  err = stmt.QueryRow(username).Scan(&id)
  if err != nil {
    log.Print(err)
  }

  return id, nil
}

func getHosts(d *Database) ([]Host, error) {
	var id int
	var hostname, ipaddress, os, information string

	rows, err := d.Connection.Query("select * from hosts")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &hostname, &ipaddress, &os, &information)
		if err != nil {
			log.Print(err)
		}

		hosts = append(hosts, Host{
			Id:          id,
			Hostname:    hostname,
			IpAddress:   ipaddress,
			Os:          os,
			Information: information,
		})
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
	}

	return hosts, nil
}

func insertHost(d *Database, host Host) {
	tx, err := d.Connection.Begin()
	if err != nil {
		log.Print(err)
	}

	stmt, err := tx.Prepare("insert into hosts(hostname, ipaddress, os, information) values(?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(host.Hostname, host.IpAddress, host.Os, host.Information)
	if err != nil {
		log.Print(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}
}

func getCredentials(d *Database) ([]Credential, error) {
	var id int
	var username, password, host, information string

	rows, err := d.Connection.Query("select * from credentials")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &username, &password, &host, &information)
		if err != nil {
			log.Print(err)
		}

		creds = append(creds, Credential{
			Id:          id,
			Username:    username,
			Password:    password,
			Host:        host,
			Information: information,
		})
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
	}

	return creds, nil
}

func insertCredential(d *Database, cred Credential) {
	tx, err := d.Connection.Begin()
	if err != nil {
		log.Print(err)
	}

	stmt, err := tx.Prepare("insert into credentials(username, password, host, information) values(?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cred.Username, cred.Password, cred.Host, cred.Information)
	if err != nil {
		log.Print(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}
}
