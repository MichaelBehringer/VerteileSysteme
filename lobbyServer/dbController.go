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
	db, err = sql.Open("mysql", "gogo:gogo@tcp(db:3306)/gogoGameDB")
	if err != nil {
		fmt.Println(err)
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

func ExecuteSQLRow(statement string, params ...interface{}) *sql.Row {
	return db.QueryRow(statement, params...)
}

func ExecuteDDL(statement string, params ...interface{}) sql.Result {
	result, _ := db.Exec(statement, params...)
	return result
}

func GetFunctions() []Highscore {
	results := ExecuteSQL("SELECT Highscore, Username FROM Highscore JOIN Player ON Highscore.Player_ID = Player.ID ORDER BY Highscore DESC LIMIT 10")
	functions := []Highscore{}
	for results.Next() {
		var function Highscore
		results.Scan(&function.Highscore, &function.Name)
		functions = append(functions, function)
	}
	return functions
}

func GetGameServers() []Server {
	results := ExecuteSQL("select ID, Servername, Servernumber, PlayerCounter from ActiveGameServer")
	servers := []Server{}
	for results.Next() {
		var server Server
		results.Scan(&server.ID, &server.PetName, &server.Address, &server.PlayerCount)
		servers = append(servers, server)
	}
	return servers
}
