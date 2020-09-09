package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HTTPRoutes registers endpoints and their appropriate HandlerFuncs to the provided router as well as any subrouters.
func HTTPRoutes(router *mux.Router) {
	itemSubrouter := router.PathPrefix("/items").Subrouter()
	itemSubrouter.HandleFunc("/{itemID}", item()).Methods("GET")
}

func item() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var item Item
		urlVar := mux.Vars(r)

		itemID, _ := strconv.Atoi(urlVar["itemID"]) // get rid of "_" and just handle the error
		item = GetItem(itemID)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}
