package handlers

import (
    "encoding/json"
    "net/http"
    "time"
   
    "github.com/golang-jwt/jwt"
    "github.com/OkabeRitarou/GoServer/models"
    "github.com/OkabeRitarou/GoServer/repository"
    "github.com/OkabeRitarou/GoServer/server"
    "github.com/segmentio/ksuid"

    "golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
    Token   string `json:"token"`
}

type LoginRequest struct {
    Email string    `json:"email"`
    Password string `json:"password"`
}

type SignUpRequest struct {
    Email string    `json:"email"`
    Password string `json:"password"`
}

type SignUpResponse struct {
    Id string       `json:"id"`
    Email string    `json:"email"`
}

func LoginHandler(s server.Server) http.HandlerFunc {
    return func(response http.ResponseWriter, request *http.Request) {
        var loginRequest = LoginRequest{}
        err := json.NewDecoder(request.Body).Decode(&loginRequest)
        if err != nil {
            http.Error(response, err.Error(), http.StatusBadRequest)
            return
        }
        user, err := repository.GetUserByEmail(request.Context(), loginRequest.Email)
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }
        if user == nil {
            http.Error(response, "Invalid Credentials", http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
            http.Error(response, "Invalid Credentials", http.StatusUnauthorized)
            return
        }

        claims := models.AppClaims{
            UserId: user.Id,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
            },
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString([]byte(s.Config().JwtSecret))

        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }
        response.Header().Set("Content-type", "application/json")
        json.NewEncoder(response).Encode(LoginResponse{
            Token: tokenString,
        })
    }
}

func SignUpHandler(s server.Server) http.HandlerFunc {
    return func(response http.ResponseWriter, request *http.Request) {
        var userRequest = SignUpRequest{}
        err := json.NewDecoder(request.Body).Decode(&userRequest)
        if err != nil {
            http.Error(response, err.Error(), http.StatusBadRequest)
            return
        }
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
        }
        id, err := ksuid.NewRandom()
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }
        var user = models.User {
            Email: userRequest.Email,
            Password: string(hashedPassword),
            Id: id.String(),
        }

        err = repository.InsertUser(request.Context(), &user)
        if err != nil {
            http.Error(response, err.Error(), http.StatusInternalServerError)
            return
        }

        response.Header().Set("Content-type", "application/json")
        json.NewEncoder(response).Encode(SignUpResponse{
            Id: user.Id,
            Email: user.Email,
        })
    }
}

func MyHandler(s server.Server) (http.HandlerFunc) {
    return func (w http.ResponseWriter, r *http.Request){

        claims := GetClaimsOnJson(s, w, r)
        if claims == nil {
            return
        }
        user, err := repository.GetUserById(r.Context(), claims.UserId)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return 
        }
        w.Header().Set("Content-type", "application/json")
        json.NewEncoder(w).Encode(*user)
    }
}
