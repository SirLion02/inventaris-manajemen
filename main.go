package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Selamat datang di Backend Golang!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Halo! Kamu sedang mengakses halaman Hello.")
}

func main() {
    http.HandleFunc("/", homeHandler)

    http.HandleFunc("/hello", helloHandler)

    fmt.Println("Server berjalan di http://localhost:3000")

    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        fmt.Println("Error saat menjalankan server:", err)
    }
}
