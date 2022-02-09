package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chetanmeniyabacncy/docker_microservice3/lang"
	"github.com/chetanmeniyabacncy/docker_microservice3/models"
	"github.com/chetanmeniyabacncy/docker_microservice3/validation"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"gopkg.in/go-playground/validator.v9"
)

type LoginSuccess struct {
	Status  int64        `json:"status"`
	Message string       `json:"message"`
	Data    *models.User `json:"data"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

//Login
func (h *BaseHandlerSqlx) Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "login")
	r.ParseForm()

	w.Header().Set("content-type", "application/json")
	response := LoginSuccess{}

	var reqlogin models.ReqLogin
	reqlogin.Email = r.FormValue("email")
	reqlogin.Password = r.FormValue("password")

	v := validator.New()
	v = validation.Custom(v)

	err := v.Struct(reqlogin)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0
		response.Message = lang.Get("errors")
		response.Data = resp
		json.NewEncoder(w).Encode(response)
		return
	}

	user, errmessage := models.Login(h.db, &reqlogin)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	session.Values["authenticated"] = "1"
	session.Save(r, w)
	session, _ = store.Get(r, "login")

	response.Status = 1
	response.Message = lang.Get("login_sucess")
	response.Data = user
	json.NewEncoder(w).Encode(response)

}

func (h *BaseHandlerSqlx) Secret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "login")
		// Check if user is authenticated
		if session.Values["authenticated"] == "1" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not authorized", 401)
			return
		}
	})
}

func (h *BaseHandlerSqlx) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	session, _ := store.Get(r, "login")
	response := CommonSuccess{}
	// Revoke users authentication
	session.Values["authenticated"] = "0"
	session.Save(r, w)
	response.Status = 1
	response.Message = lang.Get("logout_success")
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandlerSqlx) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			// json.NewEncoder(w).Encode(ErrHandler(lang.Get("jwt_no_token_found")))
			// return
			http.Error(w, "Not authorized", 401)
			return
		}

		var mySigningKey = []byte(os.Getenv("SECRET_KEY"))

		token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(lang.Get("token_parsing_error"))
			}
			return mySigningKey, nil
		})

		if err != nil {
			// json.NewEncoder(w).Encode(ErrHandler(lang.Get("token_expired")))
			// return
			http.Error(w, "Not authorized", 401)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "2" {
				r.Header.Set("Role", "admin")
				next.ServeHTTP(w, r)
				return
			} else if claims["role"] == "1" {
				r.Header.Set("Role", "1")
				next.ServeHTTP(w, r)
				return
			}
		}
		if err != nil {
			json.NewEncoder(w).Encode(ErrHandler(lang.Get("not_authorized")))
			return
			http.Error(w, "Not authorized", 401)
			return
		}
	})
}
