package controllers

// ErrHandler returns error message response
func ErrHandler(errmessage string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 0;
	errresponse.Message = errmessage;
	return &errresponse
}