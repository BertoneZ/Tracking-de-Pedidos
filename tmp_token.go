package main
import (
  "fmt"
  "time"
  "github.com/golang-jwt/jwt/v5"
)
func main() {
  claims := jwt.MapClaims{"user_id":"debug-admin","role":"admin","exp": time.Now().Add(2 * time.Hour).Unix()}
  t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  s, err := t.SignedString([]byte("una_clave_secreta_muy_larga_y_segura_123"))
  if err != nil { panic(err) }
  fmt.Println(s)
}
