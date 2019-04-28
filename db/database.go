package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PlayerInfo struct {
	ID    string
	PosX  float64
	PosY  float64
	Class string
}

var Database *sql.DB

type DbConfig struct {
	DbHost string
	DbPort string
	DbName string
	DbUser string
	DbPwd  string
}

var DbCfg DbConfig = DbConfig{
	DbHost: "0.0.0.0",
	DbPort: "5432",
	DbName: "postgres",
	DbUser: "postgres",
	DbPwd:  "1Qh0RjfU7T!",
}

func InitDB() (err error) {

	connInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DbCfg.DbHost, DbCfg.DbPort, DbCfg.DbUser, DbCfg.DbPwd, DbCfg.DbName)

	Database, err = sql.Open("postgres", connInfo)
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)
		if err = Database.Ping(); err == nil {
			break
		}
		return err
	}

	log.Println("Successfully connected to database")

	return nil
}

func CreateDbTable(DB *sql.DB) error {
	DB.Exec(`CREATE SCHEMA public;`)
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS "public"."users"(
		"username" character varying(50) NOT NULL UNIQUE,
		"password" TEXT NOT NULL,
		"id" character varying(50) UNIQUE,
		"playername" character varying(50),
		"posx" float DEFAULT 250,
		"posy" float DEFAULT 250,
		"class" TEXT NOT NULL,
		"world" character varying(50)
	)`)
	if err != nil {

		return err

	}
	log.Println("Players Table Created")

	return nil
}

func DeleteDbTable(DB *sql.DB, ID []string) error {
	_, err := DB.Exec(`DROP TABLE users`)
	if err != nil {
		log.Println(err)
		return (err)
	}
	log.Println("Deleted Users")

	for _, v := range ID {

		fmt.Sprintf(`DROP TABLE %s`, v)
		_, err = DB.Exec(fmt.Sprintf(`DROP TABLE "%s"`, v))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	log.Println("Deleted inventories")
	return nil
}

// TODO need update
func DeleteInventory(DB *sql.DB, ID string) error {

	query := `DELETE FROM "public"."inventory" WHERE playerid = $1`
	_, err := DB.Exec(query, ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inventory of playerid " + ID + " deleted")
	return nil
}

func RegisterPlayer(DB *sql.DB, username string, password string, player PlayerInfo) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Println("Problem with hashing password")
		return err
	}

	query := `INSERT INTO users(username,password,id,class) VALUES ($1,$2,$3,$4)`

	_, err = DB.Exec(query, username, hashedPassword, player.ID, player.Class)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func LoginPlayer(DB *sql.DB, username string, password string) (authorized bool, player PlayerInfo) {

	var hashed_pass string

	err := DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashed_pass)
	if err != nil {
		log.Println("Could not parse password")
		return false, PlayerInfo{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed_pass), []byte(password))
	if err != nil {
		log.Println("Unauthorized")
		return false, PlayerInfo{}
	}

	err = DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&player.ID)
	if err != nil {
		log.Println("Could not parse id ", err)
		return false, PlayerInfo{}
	}
	err = DB.QueryRow("SELECT posx FROM users WHERE username = $1", username).Scan(&player.PosX)
	if err != nil {
		log.Println("Could not parse posx", err)
		return false, PlayerInfo{}
	}
	err = DB.QueryRow("SELECT posy FROM users WHERE username = $1", username).Scan(&player.PosY)
	if err != nil {
		log.Println("Could not parse posy ", err)
		return false, PlayerInfo{}
	}
	err = DB.QueryRow("SELECT class FROM users WHERE username = $1", username).Scan(&player.Class)
	if err != nil {
		log.Println("Could not parse class ", err)
		return false, PlayerInfo{}
	}

	return true, player

}

//CreateInventoryTable create new Inventory Table shared accross all players, every slot has ID of player
func CreateInventoryTable(DB *sql.DB, ID string) error {

	_, err := DB.Exec(fmt.Sprintf(`CREATE SEQUENCE IF NOT EXISTS "public".%s INCREMENT 1 START 1`, ID))
	if err != nil {
		return err
	}

	_, err = DB.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "public"."%s"(
		"id" integer DEFAULT nextval(('"public".%s'::text)) NOT NULL,
		"slot" character varying(20)
		
	)`, ID, ID))
	if err != nil {
		return err
	}
	log.Println("Inventory Created")

	return nil
}

//AddToInventory add new item to inventory , need ID of player, and name of Item
func AddToInventory(DB *sql.DB, ID string, Item string) error {
	query := fmt.Sprintf(`INSERT INTO "%s"(slot) VALUES ($1)`, ID)

	_, err := DB.Exec(query, Item)
	if err != nil {
		return err
	}
	log.Println("Item: " + Item + " added to inventory of player with ID: " + ID)
	return nil
}

func GetInventory(DB *sql.DB, ID string) (err error, items []string) {

	query, err := DB.Prepare(fmt.Sprintf(`SELECT slot FROM "%s"`, ID))
	if err != nil {
		return err, nil
	}

	rows, err := query.Query()
	defer rows.Close()

	for rows.Next() {
		var item string

		if err := rows.Scan(&item); err != nil {
			log.Println(err)
		}
		items = append(items, item)
	}

	fmt.Println(items)
	log.Println("Player with ID: " + ID + " got inventory from DB")
	return nil, nil

}
