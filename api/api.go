package api

import (
	"context"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/service"
)

type API struct {
	records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records}
}

const requestSourceKey string = "src"

func addRequestSource(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), requestSourceKey, "v2")
        r = r.WithContext(ctx)

        // Call the next handler
        next.ServeHTTP(w, r)
    }
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}
func (a *API) CreateRoutes2(routes *mux.Router) {
	routes.Path("/records/{id}/version/{ver}").HandlerFunc(addRequestSource(a.GetRecords)).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(addRequestSource(a.GetRecords)).Methods("GET")
	routes.Path("/records/{id}/latest").HandlerFunc(addRequestSource(a.GetLatestVersion)).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(addRequestSource(a.PostRecords)).Methods("POST")
}
