package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func reloadable() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)
	go func() {
		for {
			<-s
			log.Println("Reloaded")
		}
	}()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./assets/index.html")
}

func contact(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Form)
}

func register() *http.ServeMux {
	mux := http.NewServeMux()
	staticHandler := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", staticHandler))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/contact", contact)

	return mux
}

func listen(mux *http.ServeMux) {

	port := os.Getenv("PORT")

	if port == "" {
		log.Println("$PORT not set")
		os.Exit(1)
	}
	log.Println("Server starting on port: " + port)

	if os.Getenv("ENV") != "development" {

		//cert := os.Getenv("NEST_CERT")
		//if cert == "" {
		//	log.Println("$NEST_CERT not set")
		//	os.Exit(1)
		//}
		//
		//key := os.Getenv("NEST_KEY")
		//if key == "" {
		//	log.Println("$NEST_KEY not set")
		//	os.Exit(1)
		//}
		// Start server on $PORT
		log.Println("using TLS")
		log.Fatal(http.ListenAndServeTLS(":"+port, "./server.crt", "./server.key", mux))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}

}

func main() {
	mux := register()
	reloadable()
	listen(mux)
}
