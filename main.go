package main

import (
	"log"
	"net/http"
	"spbkluapp/config"
	"spbkluapp/controllers/authcontroller"
	"spbkluapp/controllers/bsscontroller"
	"spbkluapp/controllers/homecontroller"
	"spbkluapp/middlewares/authmiddleware"
)

func main() {
	config.Connect()

	//static for css, js, images
	staticDir := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	// Serve the static files under the "/static/" URL path
	http.Handle("/static/", staticDir)

	//Routes

	//auth
	http.HandleFunc("/", authcontroller.Index)
	http.HandleFunc("/login", authcontroller.Login)
	http.HandleFunc("/logout", authcontroller.Logout)

	//Home
	http.HandleFunc("/dashboard", authmiddleware.AuthMiddleware(homecontroller.Index))

	//Bss
	http.HandleFunc("/bss", authmiddleware.AuthMiddleware(bsscontroller.Index))
	http.HandleFunc("/bss_ajax", authmiddleware.AuthMiddleware(bsscontroller.Get))
	http.HandleFunc("/bss_add", authmiddleware.AuthMiddleware(bsscontroller.Add))
	http.HandleFunc("/bss_edit", authmiddleware.AuthMiddleware(bsscontroller.Edit))
	http.HandleFunc("/bss_delete", authmiddleware.AuthMiddleware(bsscontroller.Delete))

	log.Println("Server running on port 9000")
	http.ListenAndServe(":9000", nil)
}
