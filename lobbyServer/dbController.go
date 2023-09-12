package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Highscore struct {
	Highscore int    `json:"Highscore"`
	Name      string `json:"Name"`
}

func InitDB() {
	db, err = sql.Open("mysql", "root1:root@(localhost:3306)/VerteilteSystemeDB")
	if err != nil {
		panic(err.Error())
	}
}

func CloseDB() {
	db.Close()
}

func ExecuteSQL(statement string, params ...interface{}) *sql.Rows {
	results, err := db.Query(statement, params...)
	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	return results
}

func GetFunctions() []Highscore {
	results := ExecuteSQL("SELECT Highscore, Name FROM Highscore JOIN User ON Highscore.User = User.ID ORDER BY Highscore DESC LIMIT 10")
	functions := []Highscore{}
	for results.Next() {
		var function Highscore
		results.Scan(&function.Highscore, &function.Name)
		functions = append(functions, function)
	}
	return functions
}
