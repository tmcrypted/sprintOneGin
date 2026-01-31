package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	_ = http.ListenAndServe(":8080", r)

}
