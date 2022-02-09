package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/chetanmeniyabacncy/docker_microservice5/generallib"
	"github.com/chetanmeniyabacncy/docker_microservice5/lang"
	"github.com/chetanmeniyabacncy/docker_microservice5/models"
	"github.com/chetanmeniyabacncy/docker_microservice5/repositories"
	"github.com/chetanmeniyabacncy/docker_microservice5/validation"
	"github.com/streadway/amqp"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	db *sql.DB
}

// BaseHandler will hold everything that controller needs
type BaseHandlerGorm struct {
	db *gorm.DB
}

// BaseHandler will hold everything that controller needs
type BaseHandlerSqlx struct {
	db *sqlx.DB
}

// BaseHandler will hold everything that controller needs
type BaseHandlerMongo struct {
	db *mongo.Client
}

// swagger:model GetCompanies
type GetCompanies struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string            `json:"message"`
	Data    *models.Companies `json:"data"`
}

type GetCompaniesGorm struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string           `json:"message"`
	Data    *models.Company1 `json:"data"`
}

type GetCompaniesDataTables struct {
	Status  int64                       `json:"status"`
	Message string                      `json:"message"`
	Data    *models.DataTablesCompanies `json:"data"`
}

// swagger:model GetCompany
type GetCompany struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Companies for this user
	Data *models.Company `json:"data"`
}

type GetCompanyGorm struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Companies for this user
	Data *models.Company1 `json:"data"`
}

type CustomValidationMessages struct {
	messages map[string]string
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerGorm(db *gorm.DB) *BaseHandlerGorm {
	return &BaseHandlerGorm{
		db: db,
	}
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerSqlx(db *sqlx.DB) *BaseHandlerSqlx {
	return &BaseHandlerSqlx{
		db: db,
	}
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerMongo(db *mongo.Client) *BaseHandlerMongo {
	return &BaseHandlerMongo{
		db: db,
	}
}

// HelloWorld returns Hello, World
func (h *BaseHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies := models.GetCompanies(h.db)

	if err := h.db.Ping(); err != nil {
		fmt.Println("DB Error")
	}

	for _, elem := range *companies {
		w.Write([]byte(elem.Name))
	}

}

type Trainer struct {
	Name string
	Age  int
	City string
}

type CompanyBaseHandler struct {
	companyRepo repositories.CompanyRepository
}

// NewCompanyRepository returns a new CompanyBaseHandler
func NewCompanyRepository(companyRepo repositories.CompanyRepository) *CompanyBaseHandler {
	return &CompanyBaseHandler{
		companyRepo: companyRepo,
	}
}

// HelloWorld returns Hello, World
func (h *CompanyBaseHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	response := GetCompany{}

	company, err := h.companyRepo.FindByID(id)
	if err != nil {
		fmt.Println("Error", err)
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = company

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HelloWorld returns Hello, World
func (h *BaseHandlerMongo) GetCompaniesMongo(w http.ResponseWriter, r *http.Request) {
	response := GetCompanies{}

	companieslist := models.Companies{}
	searchtext := r.URL.Query().Get("search")

	companies, errmessage := models.GetCompaniesMongo(h.db, searchtext)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	// collection := h.db.Database("company").Collection("companies")
	// ash := Trainer{"Ash", 10, "Pallet Town"}

	// var result bson.M
	// err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "name", Value: "Ash"}}).Decode(&result)

	// // collection := h.db.Database("company").Collection("trainers")
	// // ash := Trainer{"Ash", 10, "Pallet Town"}
	// // misty := Trainer{"Misty", 10, "Cerulean City"}
	// // brock := Trainer{"Brock", 15, "Pewter City"}

	// // insertResult, err := collection.InsertOne(context.TODO(), ash)
	// if err != nil {
	// 	// log.Fatal(result)
	// }
	// fmt.Println(err)
	// fmt.Println("test")
	// log.Fatal(collection)
	// fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	for _, elem := range *companies {
		company := models.Company{}
		company.Name = elem["name"].(string)
		company.Status = elem["status"].(int64)
		company.Id = elem["id"].(int64)
		companieslist = append(companieslist, company)
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = &companieslist

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// HelloWorld returns Hello, World
func (h *BaseHandlerMongo) RabbitMQ(w http.ResponseWriter, r *http.Request) {
	response := GetCompanies{}

	companieslist := models.Companies{}
	searchtext := r.URL.Query().Get("search")

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler("Failed to open a channel1"))
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler("Failed to open a channel2"))
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"companyadd", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler("Failed to open a channel3"))
		return
	}

	body := "Company Get Successfully!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler("Failed to send message in channel"))
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)

	companies, errmessage := models.GetCompaniesMongo(h.db, searchtext)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	for _, elem := range *companies {
		company := models.Company{}
		company.Name = elem["name"].(string)
		company.Status = elem["status"].(int64)
		company.Id = elem["id"].(int64)
		companieslist = append(companieslist, company)
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = &companieslist

}

// HelloWorld returns Hello, World
func (h *BaseHandlerMongo) RPCget(w http.ResponseWriter, r *http.Request) {
	response := GetCompanies{}

	companieslist := models.Companies{}
	searchtext := r.URL.Query().Get("search")

	var reply *[]bson.M

	client, err := rpc.DialHTTP("tcp", "microservicestest.microservice_4:8092")

	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	err = client.Call("API.GetCompaniesMongoRpc", searchtext, &reply)
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(err.Error()))
		return
	}
	for _, elem := range *reply {
		company := models.Company{}
		company.Name = elem["name"].(string)
		company.Status = elem["status"].(int64)
		company.Id = elem["id"].(int64)
		companieslist = append(companieslist, company)
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = &companieslist

}

func (h *BaseHandlerSqlx) GetCompaniesSqlxDataTables(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.DataTablesRequest
	err := decoder.Decode(&reqcompany)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	response := GetCompaniesDataTables{}
	data := models.GetCompaniesSqlxDataTables(h.db, reqcompany)

	// for _, elem := range *companies {
	// 	w.Write([]byte(elem.Name))
	// }
	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = data

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route GET /admin/company/list admin listCompany
// Get companies list
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetCompanies
func (h *BaseHandlerSqlx) GetCompaniesSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetCompanies{}
	companies := models.GetCompaniesSqlx(h.db)

	// for _, elem := range *companies {
	// 	w.Write([]byte(elem.Name))
	// }
	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = companies

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerGorm) GetCompaniesGorm(w http.ResponseWriter, r *http.Request) {
	response := GetCompaniesGorm{}
	companies := models.GetCompaniesGorm(h.db)

	// for _, elem := range *companies {
	// 	w.Write([]byte(elem.Name))
	// }
	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = companies

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /admin/company/ admin addCompany
// Create a new company
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetCompany
func (h *BaseHandlerSqlx) PostCompanySqlx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := GetCompany{}

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.ReqCompany
	err := decoder.Decode(&reqcompany)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	v := validator.New()
	v = validation.Custom(v)

	err = v.Struct(reqcompany)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0
		response.Message = lang.Get("errors")
		response.Data = resp
		json.NewEncoder(w).Encode(response)
		return
	}

	company, errmessage := models.PostCompanySqlx(h.db, reqcompany)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	generallib.Measure()
	generallib.GoChannleExample()
	go generallib.SendMail()
	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerGorm) PostCompanyGorm(w http.ResponseWriter, r *http.Request) {
	req := r.Context()

	w.Header().Set("content-type", "application/json")
	response := GetCompanyGorm{}

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.ReqCompany
	err := decoder.Decode(&reqcompany)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	v := validator.New()
	v = validation.Custom(v)

	err = v.Struct(reqcompany)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0
		response.Message = lang.Get("errors")
		response.Data = resp
		json.NewEncoder(w).Encode(response)
		return
	}

	company, errmessage := models.PostCompanyGorm(req, h.db, reqcompany)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	// generallib.Measure()
	// generallib.GoChannleExample()
	// go generallib.SendMail()
	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerMongo) PostCompanyMongo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := GetCompany{}

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.ReqCompany
	err := decoder.Decode(&reqcompany)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	v := validator.New()
	v = validation.Custom(v)

	err = v.Struct(reqcompany)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0
		response.Message = lang.Get("errors")
		response.Data = resp
		json.NewEncoder(w).Encode(response)
		return
	}

	company, errmessage := models.PostCompanyMongo(h.db, reqcompany)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	generallib.Measure()
	generallib.GoChannleExample()
	go generallib.SendMail()
	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerMongo) EditCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	response := GetCompany{}

	var reqcompany models.ReqCompany
	reqcompany.Status, _ = strconv.ParseInt(r.FormValue("status"), 10, 64)
	reqcompany.Name = r.FormValue("name")

	company, errmessage := models.EditCompanyMongo(h.db, &reqcompany, vars["id"])
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("update_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

// GetCompany returns company
func (h *BaseHandlerSqlx) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := GetCompany{}

	company, errmessage := models.GetCompany(h.db, vars["id"])

	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = company

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerGorm) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := GetCompanyGorm{}
	fmt.Println(vars["id"])
	company, errmessage := models.GetCompanyGorm(h.db, vars["id"])

	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = company

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route  PUT /admin/company/{id} admin editCompany
// Edit a company
//
// consumes:
//         - application/x-www-form-urlencoded
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetCompany
func (h *BaseHandlerSqlx) EditCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	response := GetCompany{}
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	var reqcompany models.ReqCompany
	reqcompany.Status, err = strconv.ParseInt(r.FormValue("status"), 10, 64)
	reqcompany.Name = r.FormValue("name")

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	company, errmessage := models.EditCompany(h.db, &reqcompany, id)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("update_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerGorm) EditCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	response := GetCompanyGorm{}
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	fmt.Println(vars["id"])
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	var reqcompany models.ReqCompany
	reqcompany.Status, err = strconv.ParseInt(r.FormValue("status"), 10, 64)
	reqcompany.Name = r.FormValue("name")

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	company, errmessage := models.EditCompanyGorm(h.db, &reqcompany, id)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("update_success")
	response.Data = company
	json.NewEncoder(w).Encode(response)
}

// swagger:route DELETE /admin/company/{id} admin deleteCompany
// Delete company
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: CommonSuccess
// Create handles Delete get company
func (h *BaseHandlerSqlx) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	errmessage := models.DeleteCompany(h.db, vars["id"])

	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	successresponse := CommonSuccess{}
	successresponse.Status = 1
	successresponse.Message = lang.Get("delete_success")

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(successresponse)
}

func (h *BaseHandlerGorm) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	errmessage := models.DeleteCompanyGorm(h.db, vars["id"])

	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	successresponse := CommonSuccess{}
	successresponse.Status = 1
	successresponse.Message = lang.Get("delete_success")

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(successresponse)
}

func (h *BaseHandlerMongo) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	errmessage := models.DeleteCompanyMongo(h.db, vars["id"])

	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	successresponse := CommonSuccess{}
	successresponse.Status = 1
	successresponse.Message = lang.Get("delete_success")

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(successresponse)
}
