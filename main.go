package main

import (
	"extenders/controller"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {

	router := httprouter.New()
	router.GET("/", controller.Index)

	router.POST("/prioritize", controller.Prioritize)

	log.Printf("start up sample-scheduler-extender!\n")
	http.ListenAndServe(":8888", router)

}
