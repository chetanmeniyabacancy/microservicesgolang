package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chetanmeniyabacncy/docker_microservice1/controllers"
	"github.com/streadway/amqp"

	"github.com/chetanmeniyabacncy/docker_microservice1/config"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"net/http/httputil"
	"net/url"
	// "context"
	// "crypto/tls"
	// "flag"
	// "fmt"
	// "io"
	// "time"
	// "golang.org/x/crypto/acme/autocert"
)

// const (
// 	htmlIndex = `<html><body>Welcome!</body></html>`
// 	httpPort  = "127.0.0.1:5001"
// )

// var (
// 	flgProduction          = true
// 	flgRedirectHTTPToHTTPS = true
// )

// func handleIndex(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, htmlIndex)
// }

// func makeServerFromMux(mux *http.ServeMux) *http.Server {
// 	// set timeouts so that a slow or malicious client doesn't
// 	// hold resources forever
// 	return &http.Server{
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 5 * time.Second,
// 		IdleTimeout:  120 * time.Second,
// 		Handler:      mux,
// 	}
// }

// func makeHTTPServer() *http.Server {
// 	mux := &http.ServeMux{}
// 	mux.HandleFunc("/", handleIndex)
// 	return makeServerFromMux(mux)

// }

// func makeHTTPToHTTPSRedirectServer() *http.Server {
// 	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
// 		newURI := "https://" + r.Host + r.URL.String()
// 		http.Redirect(w, r, newURI, http.StatusFound)
// 	}
// 	mux := &http.ServeMux{}
// 	mux.HandleFunc("/", handleRedirect)
// 	return makeServerFromMux(mux)
// }

// func parseFlags() {
// 	flag.BoolVar(&flgProduction, "production", true, "if true, we start HTTPS server")
// 	flag.BoolVar(&flgRedirectHTTPToHTTPS, "redirect-to-https", true, "if true, we redirect HTTP to HTTPS")
// 	flag.Parse()
// }

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

// // serves index file
// func admin(w http.ResponseWriter, r *http.Request) {
// 	p := path.Dir("../backend/src/index.html")
// 	// set header
// 	w.Header().Set("Content-type", "text/html")
// 	http.ServeFile(w, r, p)
// }

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

	r := mux.NewRouter()
	db := config.ConnectDB()
	// gormdb := config.ConnectDBGorm()
	// dbsqlx := config.ConnectDBSqlx()
	dbmongo := config.ConnectDBmongo()
	h := controllers.NewBaseHandler(db)
	// hgorm := controllers.NewBaseHandlerGorm(gormdb)
	// hsqlx := controllers.NewBaseHandlerSqlx(dbsqlx)
	hmongo := controllers.NewBaseHandlerMongo(dbmongo)

	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	// r.Handle("/backend1", http.FileServer(http.Dir("../backend/dist")))
	fs := http.FileServer(http.Dir("/backend/src/"))
	http.Handle("/backend2/", http.StripPrefix("/backend2", fs))

	// handler for documentation
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs1"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs1", sh1)

	opts2 := middleware.RapiDocOpts{SpecURL: "/swagger.yaml", Path: "docs2"}
	sh2 := middleware.RapiDoc(opts2, nil)
	r.Handle("/docs2", sh2)

	// user := r.PathPrefix("/admin").Subrouter()
	// user.HandleFunc("/backend", admin).Methods("GET")
	// user.HandleFunc("/login", hsqlx.Login).Methods("POST")
	// user.HandleFunc("/logout", hsqlx.Logout).Methods("GET")

	company := r.PathPrefix("/admin/company").Subrouter()
	// company.HandleFunc("/listfordatatables", hsqlx.GetCompaniesSqlxDataTables).Methods("POST")
	// company.HandleFunc("/list", hsqlx.GetCompaniesSqlx).Methods("GET")
	// company.HandleFunc("/", hsqlx.PostCompanySqlx).Methods("POST")
	// company.HandleFunc("/", hsqlx.GetCompany).Methods("GET")
	// company.HandleFunc("/{id}", hsqlx.EditCompany).Methods("PUT")
	// company.HandleFunc("/{id}", hsqlx.DeleteCompany).Methods("DELETE")
	// company.Use(hsqlx.IsAuthorized)

	// company.HandleFunc("/listgorm", hgorm.GetCompaniesGorm).Methods("GET")
	// company.HandleFunc("/gorm/", hgorm.PostCompanyGorm).Methods("POST")
	// company.HandleFunc("/gorm/{id}", hgorm.GetCompany).Methods("GET")
	// company.HandleFunc("/gorm/{id}", hgorm.EditCompany).Methods("PUT")
	// company.HandleFunc("/gorm/{id}", hgorm.DeleteCompany).Methods("DELETE")
	// company.Use(hsqlx.IsAuthorized)

	company.HandleFunc("/listcompaniesmongo", hmongo.GetCompaniesMongo).Methods("GET")
	company.HandleFunc("/rabbitmq", hmongo.RabbitMQ).Methods("GET")
	// company.HandleFunc("/postcompaniesmongo", hmongo.PostCompanyMongo).Methods("POST")
	// company.HandleFunc("/postcompaniesmongo/{id}", hmongo.EditCompany).Methods("PUT")
	// company.HandleFunc("/deletecompanymongo/{id}", hmongo.DeleteCompany).Methods("DELETE")
	r.HandleFunc("/", h.GetCompanies)

	// api := new(models.API)
	// err = rpc.Register(api)
	// fmt.Println("Error1", err)
	// rpc.HandleHTTP()

	// listener, err := net.Listen("tcp", ":8092")
	// fmt.Println("Error2", err)

	// http.Serve(listener, nil)

	// r.HandleFunc("/sqlx", hsqlx.GetCompaniesSqlx)

	// Create repos
	// userRepo := repositories.NewUserRepo(dbsqlx)

	// hinterface := controllers.NewCompanyRepository(userRepo)

	// company.HandleFunc("/interface/{id}", hinterface.GetCompany).Methods("GET")

	// parseFlags()
	// var m *autocert.Manager

	// var httpsSrv *http.Server
	// fmt.Println(flgProduction)

	// if flgProduction {

	// 	hostPolicy := func(ctx context.Context, host string) error {
	// 		// Note: change to your real host
	// 		allowedHost := "localhost"
	// 		if host == allowedHost {
	// 			return nil
	// 		}
	// 		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	// 	}

	// 	dataDir := "."
	// 	m = &autocert.Manager{
	// 		Prompt:     autocert.AcceptTOS,
	// 		HostPolicy: hostPolicy,
	// 		Cache:      autocert.DirCache(dataDir),
	// 	}

	// 	httpsSrv = makeHTTPServer()
	// 	httpsSrv.Addr = ":443"
	// 	httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	// 	go func() {
	// 		fmt.Printf("Starting HTTPS server on %s\n", httpsSrv.Addr)
	// 		err := httpsSrv.ListenAndServeTLS("", "")
	// 		if err != nil {
	// 			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
	// 		}
	// 	}()
	// }

	// var httpSrv *http.Server
	// if flgRedirectHTTPToHTTPS {
	// 	httpSrv = makeHTTPToHTTPSRedirectServer()
	// } else {
	// 	httpSrv = makeHTTPServer()
	// }
	// // allow autocert handle Let's Encrypt callbacks over http
	// if m != nil {
	// 	httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)
	// }

	// httpSrv.Addr = httpPort
	// fmt.Printf("Starting HTTP server on %s\n", httpPort)
	// err = httpSrv.ListenAndServe()
	// if err != nil {
	// 	log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	// }

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT"},
	})

	s := c.Handler(r)
	// http.ListenAndServe(":5000", s)

	origin, _ := url.Parse("http://localhost:5001/")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	conn, err := amqp.Dial(os.Getenv("DBURLRABBITMQ"))
	fmt.Println("Error", err)
	defer conn.Close()

	ch, err := conn.Channel()
	fmt.Println("Error", err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"companyadd", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	fmt.Println("Error", err)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	fmt.Println("Error", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	http.ListenAndServe(":5000", s)

	<-forever

}
