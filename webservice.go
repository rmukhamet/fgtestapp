package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type dupeResult struct {
	Dupes bool `json:"dupes"`
}
type webService struct {
	logService *logService
	router     *chi.Mux
}

func newWebService(ls *logService) *webService {
	ws := new(webService)
	router := chi.NewRouter()
	router.Get("/{userID1}/{userID2}", ws.checkDupesHandler)

	ws.router = router
	ws.logService = ls
	return ws
}

func (ws *webService) checkDupesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID1, userID2, err := getUserIDs(r)
	if err != nil {
		renderError(w, "INVALID USERID", http.StatusBadRequest)
		return
	}

	isDupe, err := ws.logService.checkDupes(userID1, userID2)
	if err != nil {
		log.Println(err)
		renderError(w, "ERROR_CHECKING_DUPES", http.StatusBadRequest)
		return
	}
	err = renderResult(w, &dupeResult{Dupes: isDupe})
	if err != nil {
		renderError(w, "ERROR_GENERATING_RESULT", http.StatusBadRequest)
	}
}
func (ws *webService) run() {
	log.Fatal(http.ListenAndServe(":4000", ws.router))

}

func getUserIDs(r *http.Request) (userID1 int64, userID2 int64, err error) {
	if userID1Str := chi.URLParam(r, "userID1"); userID1Str != "" {
		//convert to int
		userID1, err = strconv.ParseInt(userID1Str, 10, 64)
		if err != nil {
			return
		}
	}
	if userID2Str := chi.URLParam(r, "userID2"); userID2Str != "" {
		//convert to int
		userID2, err = strconv.ParseInt(userID2Str, 10, 64)
		if err != nil {
			return
		}
	}
	return
}
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}
func renderResult(w http.ResponseWriter, result interface{}) (err error) {
	resultB, err := json.Marshal(result)
	if err != nil {
		return
	}
	_, err = w.Write(resultB)
	if err != nil {
		log.Println(err)
	}
	return
}
