package middleware

import (
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt"
    "github.com/OkabeRitarou/GoServer/server"
    "github.com/OkabeRitarou/GoServer/models"
)

var (
    NO_AUTH_NEED = []string {
        "login",
        "signup",
    }
)

func shouldCheckToken(route string) bool {
    for _, value := range NO_AUTH_NEED {
        if strings.Contains(route, value) {
            return false
        }
    }
    return true
}

func CheckoutMiddleware(s server.Server) (func (http.Handler) http.Handler) {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
            if !shouldCheckToken(r.URL.Path) {
                next.ServeHTTP(w, r)
                return
            }
            tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
            _, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
                return []byte(s.Config().JwtSecret), nil
            })

            if err != nil {
                http.Error(w, err.Error(), http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
