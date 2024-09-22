package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DELETE /records/{id}/version/{ver}
// DeleteRecord deletes the specified version for the record.
func (a *API) DeleteRecord(w http.ResponseWriter, r *http.Request) {
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
		err := writeError(w, "invalid ver; version must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	err = a.records.DeleteRecordForVersion(
		ctx,
		int(idNumber),
		int(verNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v with version %v does not exist", idNumber, verNumber), http.StatusBadRequest)
		logError(err)
		return
	}
	err = writeJSON(w, fmt.Sprintf(""), http.StatusOK)
	logError(err)
}

