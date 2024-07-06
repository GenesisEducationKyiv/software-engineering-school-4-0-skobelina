package rest

import "github.com/gorilla/mux"

type Registrable interface {
	Register(router *mux.Router)
}
