package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/chetanmeniyabacncy/docker_microservice2/lang"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type Company1 struct {
	gorm.Model
	Name   string
	Status int64
}

// swagger:model Company
type Company struct {
	// Id of the company
	// in: int64
	Id int64 `json:"id"`
	// Name of the company
	// in: string
	Name string `json:"name"`
	// Status of the company
	// in: int64
	Status int64 `json:"status"`
}

type Companies []Company

type DataTablesCompanies struct {
	RecordsTotal int64     `json:"recordsTotal"`
	Companies    Companies `json:"companies"`
}

// swagger:parameters admin deleteCompany
type ReqDeleteCompany struct {
	// name: id
	// in: path
	// type: integer
	// required: true
	Id string `json:"id"`
}

type ReqAddCompany struct {
	// Name of the company
	// in: string
	Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
	// Status of the company
	// in: int64
	Status int64 `json:"status" validate:"required"`
}

// swagger:parameters admin editCompany
type ReqCompany struct {
	// name: id
	// in: path
	// type: integer
	// required: true
	Id int64 `json:"id"`
	// name: Name
	// in: formData
	// type: string
	// required: true
	Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
	// name: Status
	// in: formData
	// type: int64
	// required: true
	Status int64 `json:"status" validate:"required"`
}

// swagger:parameters admin addCompany
type ReqCompanyBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddCompany"
	//  required: true
	Body ReqAddCompany `json:"body"`
}

// GetCompanies returns companies list
func GetCompaniesMongo(db *mongo.Client, searchtext string) (*[]bson.M, string) {
	// companies := Companies{}

	opt := options.FindOptions{}

	collection := db.Database(os.Getenv("DBNAMEMONGO")).Collection("companies1")

	cursor, err := collection.Find(context.TODO(), bson.D{{Key: "name", Value: searchtext}}, &opt)
	var companies []bson.M

	if err != nil {
		return &companies, ErrHandler(err)

	}
	if err = cursor.All(context.TODO(), &companies); err != nil {
		return &companies, ErrHandler(err)

	}
	if len(companies) <= 0 {
		return &companies, lang.Get("no_result")
	}
	fmt.Println("Test")
	// fmt.Println(err)

	// fmt.Println(companies)
	//collection.find({name: 'Harry'})

	return &companies, ""
}

// GetCompanies returns companies list
func GetCompanies(db *sql.DB) *Companies {
	companies := Companies{}

	res, err := db.Query("SELECT id,name,status FROM companies")

	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var company Company
		err := res.Scan(&company.Id, &company.Name, &company.Status)
		if err != nil {
			log.Fatal(err)
		}
		companies = append(companies, company)

	}
	return &companies

}

func GetCompaniesSqlxDataTables(db *sqlx.DB, reqcompany *DataTablesRequest) *DataTablesCompanies {
	dataTablesFields := map[int]string{0: "id", 1: "name", 2: "status"}
	datatablecompanies := DataTablesCompanies{}
	orderByQuery := ``
	for _, service := range reqcompany.Order {
		fmt.Println(service.Column)
		if len(orderByQuery) > 0 {
			orderByQuery = orderByQuery + `, ` + dataTablesFields[service.Column] + ` ` + service.Dir
		} else {
			orderByQuery = ` Order By ` + dataTablesFields[service.Column] + ` ` + service.Dir
		}
	}
	fmt.Println(orderByQuery)
	searchQuery := ``
	if reqcompany.Search.Value != "" {
		searchQuery = ` AND name like '%` + reqcompany.Search.Value + `%'`
	}
	err := db.Select(&datatablecompanies.Companies, "SELECT SQL_CALC_FOUND_ROWS id,name,status FROM companies WHERE 1=1 "+searchQuery+" "+orderByQuery+" LIMIT ? OFFSET ? ", reqcompany.Length, reqcompany.Start)
	err = db.Get(&datatablecompanies.RecordsTotal, "SELECT count(0) FROM companies WHERE  1=1 "+searchQuery+"")

	if err != nil {
		log.Fatal(err)
	}
	return &datatablecompanies

}

func GetCompaniesSqlx(db *sqlx.DB) *Companies {
	companies := Companies{}

	err := db.Select(&companies, "SELECT id,name,status FROM companies")

	if err != nil {
		log.Fatal(err)
	}
	return &companies

}

func GetCompaniesGorm(gormDB *gorm.DB) *Company1 {

	gormDB.AutoMigrate(&Company1{})
	companies := Company1{}

	result := gormDB.Find(&companies)

	if result.RowsAffected == 0 {
		fmt.Println(lang.Get("no_result"))
	}

	return &companies

}

// PostCompanySqlx insert company
func PostCompanySqlx(db *sqlx.DB, reqcompany *ReqCompany) (*Company, string) {
	name := reqcompany.Name
	status := reqcompany.Status

	var company Company

	stmt, err := db.Prepare("INSERT INTO companies(name,status) VALUES(?,?)")
	if err != nil {
		return &company, ErrHandler(err)
	}
	result, err := stmt.Exec(name, status)
	if err != nil {
		return &company, ErrHandler(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return &company, ErrHandler(err)
	}

	err = db.Get(&company, "SELECT id,name,status FROM companies where id = ?", id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// PostCompanyGorm insert company
func PostCompanyGorm(ctx context.Context, db *gorm.DB, reqcompany *ReqCompany) (*Company1, string) {

	name := reqcompany.Name
	status := reqcompany.Status
	var wg sync.WaitGroup

	var company Company1
	ch := make(chan uint, 1)

	// for i := 0; i < 5000; i++ {
	// 	company1 := Company1{Name: name + " test" + strconv.Itoa(i), Status: status}
	// 	db.Create(&company1)

	// }
	for i := 0; i < 1; i++ {
		select {
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			fmt.Fprint(os.Stderr, "request cancelled\n")
		case <-time.After(1 * time.Millisecond):
			wg.Add(1)
			go func() {
				defer wg.Done()
				company1 := Company1{Name: name + strconv.Itoa(i), Status: status}
				db.Create(&company1) // pass pointer of data to Create
				// id := company1.ID
				ch <- company1.ID
			}()
		}
	}
	go func() {
		for i := 0; i < 1; i++ {
			var id uint
			id = <-ch
			fmt.Println(id)
		}
	}()
	wg.Wait()
	// db.First(&company, id)
	return &company, ""
}

// PostCompanyMongo insert company
func PostCompanyMongo(db *mongo.Client, reqcompany *ReqCompany) (*Company, string) {
	// name := reqcompany.Name
	// status := reqcompany.Status

	var company Company

	collection := db.Database(os.Getenv("DBNAMEMONGO")).Collection("companies1")
	result, err := collection.InsertOne(context.TODO(), reqcompany)

	if err != nil {
		return &company, ErrHandler(err)
	}
	id := result.InsertedID

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&company)

	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// GetCompany get company
func GetCompany(db *sqlx.DB, id string) (*Company, string) {
	var company Company
	err := db.Get(&company, "SELECT id,name,status FROM companies where id = ?", id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

func GetCompanyGorm(db *gorm.DB, id string) (*Company1, string) {
	var company Company1
	db.First(&company, id)
	return &company, ""
}

// PostCompanySqlx insert company
func EditCompany(db *sqlx.DB, reqcompany *ReqCompany, id int64) (*Company, string) {
	name := reqcompany.Name
	status := reqcompany.Status

	var company Company

	stmt, err := db.Prepare("Update companies set name=?, status=? where id = ?")
	if err != nil {
		return &company, ErrHandler(err)
	}
	_, err = stmt.Exec(name, status, id)
	if err != nil {
		return &company, ErrHandler(err)
	}

	err = db.Get(&company, "SELECT id,name,status FROM companies where id = ?", id)
	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// PostCompanySqlx insert company
func EditCompanyGorm(db *gorm.DB, reqcompany *ReqCompany, id int64) (*Company1, string) {
	name := reqcompany.Name
	status := reqcompany.Status

	var company Company1
	db.First(&company, id)
	company.Name = name
	company.Status = status
	db.Save(&company)

	return &company, ""
}

// PostCompanySqlx insert company
func EditCompanyMongo(db *mongo.Client, reqcompany *ReqCompany, id string) (*Company, string) {
	var company Company

	collection := db.Database(os.Getenv("DBNAMEMONGO")).Collection("companies1")

	name := reqcompany.Name
	status := reqcompany.Status

	_id, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{primitive.E{Key: "_id", Value: _id}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "name", Value: name}, {Key: "status", Value: status}}},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &company, ErrHandler(err)
	}

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: _id}}).Decode(&company)

	if err != nil {
		return &company, lang.Get("no_result")
	}

	if err != nil {
		return &company, lang.Get("no_result")
	}
	return &company, ""
}

// DeleteCompany get company
func DeleteCompany(db *sqlx.DB, id string) string {
	stmt, err := db.Prepare("DELETE FROM companies where id = ?")
	if err != nil {
		return ErrHandler(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return ErrHandler(err)
	}
	return ""
}

// DeleteCompany get company
func DeleteCompanyMongo(db *mongo.Client, id string) string {
	collection := db.Database(os.Getenv("DBNAMEMONGO")).Collection("companies1")

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{primitive.E{Key: "_id", Value: _id}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
	if err != nil {
		return ErrHandler(err)
	}
	return ""
}

// DeleteCompany get company
func DeleteCompanyGorm(db *gorm.DB, id string) string {
	var company Company1
	db.Delete(&company, id)
	return ""
}
