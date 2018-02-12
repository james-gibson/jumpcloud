package main

import "fmt"
import "time"
import "crypto/sha512"
import "encoding/base64"
import "net/http"
import "log"

func hashEncode(password string) string {
    hash := sha512.New()
    hash.Write([]byte(password))
    byteSlice := hash.Sum(nil)

    result := base64.StdEncoding.EncodeToString([]byte(byteSlice))
    return result
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    formValues := r.Form
    password := formValues.Get("password")

    // Delaying to simulate external work taking time
    artificialDelay := time.NewTimer(5 * time.Second)
    <-artificialDelay.C

    fmt.Fprintf(w, hashEncode(password))
}

func main() {
    fmt.Println("Server launching")
    http.HandleFunc("/hash", hashHandler)

    err := http.ListenAndServe(":8000", nil)

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    fmt.Println("Server listening on port 8000")

}
