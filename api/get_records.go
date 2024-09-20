package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GET /records/{id}/latest
// GetLatestVersion retrieves the latest version for the record.
func (a *API) GetLatestVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}
	ver, err := a.records.GetLatestVersion(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}
	err = writeJSON(w, fmt.Sprintf("{version:%v}", ver), http.StatusOK)
	logError(err)
}

// GET /records/{id}
// GetRecord retrieves the record.
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	ver := mux.Vars(r)["ver"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	verNumber, err := strconv.ParseInt(ver, 10, 32)
	if err != nil {
		verNumber = -1
	}
	record, err := a.records.GetRecord(
		ctx,
		int(idNumber),
		int(verNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}
