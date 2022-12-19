package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/OkabeRitarou/GoServer/models"
	"github.com/OkabeRitarou/GoServer/repository"
	"github.com/OkabeRitarou/GoServer/server"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetClaimsOnJson(s, w, r)
		if claims == nil {
			return
		}
		var postRequest = InsertPostRequest{}
		if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post := models.Post{
			Id:      id.String(),
			Content: postRequest.Content,
			UserId:  claims.UserId,
		}
		err = repository.InsertPost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var postMessage = models.WebSocketMessage{
			Type:    "Post_Create",
			Payload: post,
		}
		s.Hub().Broadcast(postMessage, nil)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PostResponse{
			Id:      post.Id,
			Content: post.Content,
		})
	}
}

func GetPostByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		claims := GetClaimsOnJson(s, w, r)
		if claims == nil {
			return
		}
		var postRequest = InsertPostRequest{}
		if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		post := models.Post{
			Id:      params["id"],
			Content: postRequest.Content,
			UserId:  claims.UserId,
		}

		err := repository.UpdatePost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PostResponse{
			Id:      post.Id,
			Content: post.Content,
		})
	}
}

func ListPostsHandlers(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := GetParamNumeric(r, "page", uint64(0))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		size, err := GetParamNumeric(r, "size", uint64(2))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		posts, err := repository.ListPost(r.Context(), page, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
