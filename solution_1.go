package main

import "fmt"
import "crypto/sha512"
import "encoding/base64"


func hashEncode(password string) string {
    hash := sha512.New()
    hash.Write([]byte(password))
    byteSlice := hash.Sum(nil)

    result := base64.StdEncoding.EncodeToString([]byte(byteSlice))
    return result
}

func main() {
    data := "angryMonkey"

    fmt.Println(string(hashEncode(data)))
}
