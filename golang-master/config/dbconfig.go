package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chetanmeniyabacancy/golang-master/constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// ConnectDB opens a connection to the database
func ConnectDB() *sql.DB {
	db, err := sql.Open(os.Getenv("DBTYPE"), os.Getenv("DBUSERNAME")+":"+os.Getenv("DBPASSWORD")+"@/"+os.Getenv("DBNAME"))
	if err != nil {
		panic(err.Error())
	}
	return db
}

// ConnectDB opens a connection to the database
func ConnectDBGorm() *gorm.DB {
	db, err := sql.Open(os.Getenv("DBTYPE"), os.Getenv("DBUSERNAME")+":"+os.Getenv("DBPASSWORD")+"@/"+os.Getenv("DBNAME")+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	return gormDB
}

// ConnectDB opens a connection to the database
func ConnectDBSqlx() *sqlx.DB {
	db, err := sqlx.Open(constants.DBTYPE, constants.DBUSERNAME+":"+constants.DBPASSWORD+"@/"+constants.DBNAME)
	if err != nil {
		panic(err.Error())
	}

	return db
}

// ConnectDB opens a connection to the database
func ConnectDBmongo() *mongo.Client {
	// Set client options

	clientOptions := options.Client().ApplyURI(os.Getenv("DBURLMONGO"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
