package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/buger/jsonparser"
    "fmt"
    "log"
    "net/http"
    "os"
)

func getConnectionString() string {
	vcapservices := []byte(os.Getenv("VCAP_SERVICES"))
	hostname, _ := jsonparser.GetString(vcapservices, "p.mysql", "[0]", "credentials", "hostname")
	database, _ := jsonparser.GetString(vcapservices, "p.mysql", "[0]", "credentials", "name")
	username, _ := jsonparser.GetString(vcapservices, "p.mysql", "[0]", "credentials", "username")
	password, _ := jsonparser.GetString(vcapservices, "p.mysql", "[0]", "credentials", "password")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, password, hostname, database)
	fmt.Println("connectionString: " + connectionString)
	return connectionString
}

func writeMysql(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        petname := r.FormValue("petname")
		connectionString := getConnectionString()
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			fmt.Println("sql.Open error")
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		err = db.Ping()
		if err != nil {
		fmt.Println("db.Ping error")
		panic(err.Error()) // proper error handling instead of panic in your app
		}

		
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS pet ( id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, name VARCHAR(30) NOT NULL)");
		if err != nil {
			log.Fatal(err)
		}
		
		sql_statement := fmt.Sprintf("insert into pet (name) values ('%s')", petname)
		_, err = db.Exec(sql_statement)
		if err != nil {
			log.Fatal(err)
		}
	}
	
}

func readMysql(w http.ResponseWriter, r *http.Request) {
	connectionString := getConnectionString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("sql.Open error")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	err = db.Ping()
	if err != nil {
	  fmt.Println("db.Ping error")
	  panic(err.Error()) // proper error handling instead of panic in your app
	}
	
	rows, err := db.Query("select * from pet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id int
		name string
	)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w,name + "\n")

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
}