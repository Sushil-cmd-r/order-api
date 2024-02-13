package application

import (
	"github.com/sushil-cmd-r/order-api/order"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) loadRoutes() {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Order routes
	orderRoute := router.PathPrefix("/orders").Subrouter()
	a.registerOrderRoutes(orderRoute)

	a.router = router
}

func (a *App) registerOrderRoutes(router *mux.Router) {
	orderHandler := order.NewHandler(order.NewPostgresRepo(a.db))

	router.HandleFunc("/", orderHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/", orderHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/{id}", orderHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc("/{id}", orderHandler.UpdateById).Methods(http.MethodPut)
	router.HandleFunc("/{id}", orderHandler.DeleteById).Methods(http.MethodDelete)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %q from %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
