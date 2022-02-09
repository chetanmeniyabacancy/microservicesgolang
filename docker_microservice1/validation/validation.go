package validation

import (
	"fmt"
	"regexp"

	"github.com/chetanmeniyabacncy/docker_microservice1/lang"
	"gopkg.in/go-playground/validator.v9"
)

const (
	alphaSpaceRegexString string = "^[a-zA-Z ]+$"
	dateRegexString       string = "^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[469]|11)[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[13578]|1[02])[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

type FinalErrResponse struct {
	Status  int64        `json:"status"`
	Message string       `json:"message"`
	Data    *ErrResponse `json:"data"`
}

func Custom(validate *validator.Validate) *validator.Validate {
	validate.RegisterValidation("alpha_space", isAlphaSpace)
	validate.RegisterValidation("date", isDate)
	return validate
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}
func isDate(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(dateRegexString)
	return reg.MatchString(fl.Field().String())
}

func ToErrResponse(err error) *ErrResponse {
	resp := ErrResponse{
		Errors: make([]string, len(err.(validator.ValidationErrors))),
	}
	for i, e := range err.(validator.ValidationErrors) {
		// json.NewEncoder(w).Encode(ErrHandler(e))
		switch e.Tag() {
		case "required":
			resp.Errors[i] = fmt.Sprintf("%s "+lang.Get("field_is_a_required_field")+" ", lang.Get(e.Field()))
		case "max":
			resp.Errors[i] = fmt.Sprintf("%s "+lang.Get("must_be_a_maximum_of")+" %s "+lang.Get("in_length"), lang.Get(e.Field()), e.Param())
		case "min":
			resp.Errors[i] = fmt.Sprintf("%s "+lang.Get("must_be_a_minimum_of")+" %s "+lang.Get("in_length"), lang.Get(e.Field()), e.Param())
		case "url":
			resp.Errors[i] = fmt.Sprintf("%s "+lang.Get("must_be_a_valid_url")+" ", lang.Get(e.Field()))
		case "alpha_space":
			resp.Errors[i] = fmt.Sprintf("%s "+lang.Get("can_only_contain_alphabetic_and_space_characters")+" ", lang.Get(e.Field()))
		default:
			resp.Errors[i] = fmt.Sprintf(" "+lang.Get("something_wrong_on")+" %s; %s", lang.Get(e.Field()), lang.Get(e.Tag()))
		}
	}
	return &resp
}
