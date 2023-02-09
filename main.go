package main

import (
	"log"
	"net/http"
	"photo-media/pkg/controllers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fileServer := http.FileServer(http.Dir("./pkg/view"))
	//Acc Auth
	http.Handle("/", fileServer)
	http.HandleFunc("/signup", controllers.SignupHandler)
	http.HandleFunc("/login", controllers.LoginHandler)

	//Content
	http.HandleFunc("/home", controllers.HomeHandler)

	//Subcontent: Cards
	http.HandleFunc("/home/card", controllers.CardHandler)
	http.HandleFunc("/home/card/comment", controllers.CommentHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
