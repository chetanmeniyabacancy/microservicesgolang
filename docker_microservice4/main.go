package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chetanmeniyabacncy/docker_microservice4/config"
	"github.com/chetanmeniyabacncy/docker_microservice4/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type API int

// GetCompanies returns companies list
func (a *API) GetCompaniesMongoRpc(searchtext string, reply *models.Companies) error {

	companieslist := models.Companies{}

	db := config.ConnectDBmongo()
	opt := options.FindOptions{}

	collection := db.Database(os.Getenv("DBNAMEMONGO")).Collection("companies1")

	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "name", Value: "test"}}, &opt)
	var companies []bson.M

	if err != nil {
		fmt.Println("Err1", err)
		return err

	}
	if err = cursor.All(context.TODO(), &companies); err != nil {
		fmt.Println("Err2", err)
		return err
	}
	if len(companies) <= 0 {
		returnerr := errors.New("no_result")
		return returnerr
	}

	for _, elem := range companies {
		company := models.Company{}
		company.Name = elem["name"].(string)
		company.Status, _ = strconv.ParseInt(elem["status"].(string), 10, 64)
		company.Id, _ = strconv.ParseInt(elem["id"].(string), 10, 64)
		companieslist = append(companieslist, company)
	}

	fmt.Println("Err3", err)
	*reply = companieslist
	return nil
}

func main() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)
	// err := godotenv.Load(os.ExpandEnv("$GOPATH/src/golang-master/.env"))

	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	api := new(API)
	err = rpc.Register(api)
	fmt.Println("Error1", err)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+os.Getenv("RPCPORT"))
	fmt.Println("Error2", err)

	http.Serve(listener, nil)
}
