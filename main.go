package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"precisely/handler"
	"precisely/model"
	"time"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	db := model.DocumentRepository.Init(
		viper.GetString("DRIVER"),
		viper.GetString("MS_USERNAME"),
		viper.GetString("MS_PASSWORD"),
		viper.GetString("MS_PORT"),
		viper.GetString("MS_HOST"),
		viper.GetString("MS_DB"),
	)
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/documents", handler.CreateHandler).Methods("POST")
	r.HandleFunc("/documents/{id:[0-9]+}", handler.UpdateHandler).Methods("PUT")
	r.HandleFunc("/documents/{id:[0-9]+}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/documents/{id:[0-9]+}", handler.GetByIdHandler)
	r.HandleFunc("/documents", handler.GetAllHandler)

	r.Use(commonMiddleware)

	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:" + viper.GetString("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server started on port " + viper.GetString("PORT"))
	log.Fatal(srv.ListenAndServe())
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
