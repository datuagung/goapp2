// app.go -> models2

package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(user, password, host, port, dbname string) {
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) Run(addr string) {
    log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() { 
    a.Router.HandleFunc("/models2", a.getModels2).Methods("GET")
    a.Router.HandleFunc("/models", a.createModels).Methods("POST")
    a.Router.HandleFunc("/models/{models2id:[0-9]+}", a.getModels).Methods("GET")
    a.Router.HandleFunc("/models/{models2id:[0-9]+}", a.updateModels).Methods("PUT")
    a.Router.HandleFunc("/models/{models2id:[0-9]+}", a.deleteModels).Methods("DELETE")
}

func (a *App) getModels2(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    getModels2Listing, err := getModels2(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, getModels2Listing)
}

func (a *App) createModels(w http.ResponseWriter, r *http.Request) {
    var u models
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := u.createModels(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, "SQL error"+err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getModels(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    models2id, err := strconv.Atoi(vars["models2id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid models ID")
        return
    }

    u := models{Models2id: models2id}
    if err := u.getModels(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Models not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateModels(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    models2id, err := strconv.Atoi(vars["models2id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Models ID")
        return
    }

    var u models
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    u.Models2id = models2id

    if err := u.updateModels(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteModels(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    models2id, err := strconv.Atoi(vars["models2id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid User ID")
        return
    }

    u := models{Models2id: models2id}
    if err := u.deleteModels(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
