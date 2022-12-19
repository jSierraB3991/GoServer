package repository

import (
    "context"
    "github.com/OkabeRitarou/GoServer/models"
)

type PostRepository interface {
    InsertPost(ctx context.Context, post *models.Post) error
    GetPostById(ctx context.Context, id string) (*models.Post, error)
    UpdatePost(ctx context.Context, post *models.Post) error
    ListPost(ctx context.Context, page uint64, size uint64) ([]*models.Post, error)
    Close() error
}

var postRepositoryImpl PostRepository

func SetPostRepository(repository PostRepository) {
    postRepositoryImpl = repository
}

func InsertPost(ctx context.Context, post *models.Post) error {
    return postRepositoryImpl.InsertPost(ctx, post)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
    return postRepositoryImpl.GetPostById(ctx, id)
}

func UpdatePost(ctx context.Context, post *models.Post) error {
    return postRepositoryImpl.UpdatePost(ctx, post)
}

func Close() error {
    return postRepositoryImpl.Close()
}

func ListPost(ctx context.Context, page uint64, size uint64) ([]*models.Post, error){
    return postRepositoryImpl.ListPost(ctx, page, size)
}
