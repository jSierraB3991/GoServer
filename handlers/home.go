package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/golang-jwt/jwt"
    "github.com/OkabeRitarou/GoServer/models"
    "github.com/OkabeRitarou/GoServer/server"
)

type HomeResponse struct {
    Message     string  `json:"message"`
    Status      bool    `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
    return func(response http.ResponseWriter, request *http.Request) {
        response.Header().Set("Content-type", "application/json")
        response.WriteHeader(http.StatusOK)
        json.NewEncoder(response).Encode(HomeResponse{
            Message: "Welcome To The Platzi Course",
            Status: true,
        })
    }
}

func GetClaimsOnJson(s server.Server, w http.ResponseWriter ,r *http.Request) *models.AppClaims {
    tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
    tokenClaims , err := jwt.ParseWithClaims(tokenString, 
                                            &models.AppClaims{},
                                            func(token *jwt.Token) (interface{}, error) {
                                                return []byte(s.Config().JwtSecret), nil
                                            })
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return nil
    }
    if claims, ok := tokenClaims.Claims.(*models.AppClaims); ok && tokenClaims.Valid {
        return claims
    }else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil
    }
}

func getParam(r *http.Request, nameParam string) string {
    return r.URL.Query().Get(nameParam)
}

func GetParamNumeric(r *http.Request, name string, defaultValue uint64) (uint64, error) {
    value := getParam(r, name)
    if value != "" {
        valueReturn, err := strconv.ParseUint(value, 10, 64)
        if err != nil {
            return uint64(0), err
        }
        return valueReturn, nil
    }
    return defaultValue, nil
}
