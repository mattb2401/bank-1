package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	// Auth
	// Extend token
	Route{
		"AuthIndex",
		"POST",
		"/auth",
		AuthIndex,
	},
	// Get token
	Route{
		"AuthLogin",
		"POST",
		"/auth/login",
		AuthLogin,
	},
	// Create auth account
	Route{
		"AuthCreate",
		"POST",
		"/auth/account",
		AuthCreate,
	},
	// Remove auth account
	Route{
		"AuthRemove",
		"DELETE",
		"/auth/account",
		AuthRemove,
	},
	// Accounts
	// Get account details
	Route{
		"AccountIndex",
		"GET",
		"/account",
		AccountIndex,
	},
	// Create account
	Route{
		"AccountCreate",
		"POST",
		"/account",
		AccountCreate,
	},
	// Get single account
	Route{
		"AccountGet",
		"GET",
		"/account/{accountId}",
		AccountGet,
	},
	// Get all accounts
	Route{
		"AccountGetAll",
		"GET",
		"/account/all",
		AccountGetAll,
	},
	// Payments
	// Credit initiation
	Route{
		"PaymentCreditInitiation",
		"POST",
		"/payment/credit",
		PaymentCreditInitiation,
	},
	// Deposit initiation
	Route{
		"PaymentDepositInitiation",
		"POST",
		"/payment/deposit",
		PaymentDepositInitiation,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
