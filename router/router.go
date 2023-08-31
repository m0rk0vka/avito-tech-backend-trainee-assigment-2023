package router

import (
	"github.com/gorilla/mux"
	"github.com/m0rk0vka/avito-tech-backend-trainee-assigment-2023/controllers"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/getusersegments", controllers.GetUserSegments).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/createsegment", controllers.CreateSegment).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/updateusersegments", controllers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletesegment", controllers.DeleteSegment).Methods("DELETE", "OPTIONS")

	return router
}
