package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raulsilva-tech/DVRStreamAdapter/configs"
	"github.com/raulsilva-tech/DVRStreamAdapter/internal/webserver/handlers"
)

func main() {

	//getting parameters that set port and directory
	// port := flag.String("p", "8888", "Porta para servir os arquivos")
	// directory := flag.String("d", "./videos", "O diretório que deve ser servido")
	// flag.Parse()

	//loading configuration
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	vh := handlers.NewVideoHandler()
	r.Get("/stream/{file_name}", vh.Stream)

	log.Printf("Serving on HTTP port: %s \n", cfg.Port)
	//starting the server
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))

}
