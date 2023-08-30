package router

import (
	"github.com/gorilla/mux"
	"github.com/m0rk0vka/avito-tech-backend-trainee-assigment-2023/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users/{id}", controllers.GetUserByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/createsegment", controllers.CreateSegment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/{id}", controllers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletesegment/{id}", controllers.DeleteSegment).Methods("DELETE", "OPTIONS")

	return router
}
