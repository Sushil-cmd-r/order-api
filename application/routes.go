package application

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sushil-cmd-r/order-api/handler"
)

func loadRoutes() *mux.Router {
  router := mux.NewRouter()
  router.Use(loggingMiddleware)


  // Order routes
  orderRoute := router.PathPrefix("/orders").Subrouter()
  registerOrderRoutes(orderRoute)

  router.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request){
    w.WriteHeader(http.StatusOK)
  })

  return router
}

func registerOrderRoutes(router *mux.Router) {
  orderHandler := &handler.Order{}

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

