package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	forceDSN := "postgresql://postgres.bahhcjwxezopjyuqmgsb:@Hmsitb1006@aws-1-ap-northeast-1.pooler.supabase.com:6543/postgres?prefer_simple_protocol=true" 
	var dsn string
	if forceDSN != "" {
		dsn = forceDSN
	} else if dbUrl := os.Getenv("DATABASE_URL"); dbUrl != "" {
		dsn = dbUrl
	} else {
		host := os.Getenv("DB_HOST")
		if host == "" { host = "localhost" }
		
		user := os.Getenv("DB_USER")
		if user == "" { user = "postgres" }
		
		password := os.Getenv("DB_PASSWORD")
		if password == "" { password = "password" }
		
		dbname := os.Getenv("DB_NAME")
		if dbname == "" { dbname = "voting_db" }
		
		port := os.Getenv("DB_PORT")
		if port == "" { port = "5432" }
		
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	}
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, 
	}), &gorm.Config{
		PrepareStmt: false, 
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}
