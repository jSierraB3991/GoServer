package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/OkabeRitarou/GoServer/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgreRepository(urlDatabase string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", urlDatabase)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repository PostgresRepository) getUserByEmailValidate(ctx context.Context, email string) bool {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, email FROM users WHERE email = $1", email)
	user, _ := continueSearchSql(rows, err)
	return user.Email != ""
}

func (repository *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	exists := repository.getUserByEmailValidate(ctx, user.Email)
	if exists {
		return errors.New("Duplicate Email")
	}
	_, err := repository.db.ExecContext(ctx,
		"INSERT INTO users (id, email, password) VALUES($1, $2, $3)",
		user.Id, user.Email, user.Password)
	return err
}

func (repository *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	return continueSearchSql(rows, err)
}

func (repository *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return &user, nil
}

func continueSearchSql(rows *sql.Rows, err error) (*models.User, error) {
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return &user, nil
}

func (repository *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repository.db.ExecContext(ctx,
		"INSERT INTO posts (id, content, user_id) VALUES($1, $2, $3)",
		post.Id, post.Content, post.UserId)
	return err
}

func (repository *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repository.GetPostById(ctx, post.Id)
	if err != nil {
		return errors.New("Post not exists")
	}
	_, err = repository.db.ExecContext(ctx,
		"UPDATE posts SET content = $1 WHERE id = $2 AND user_id = $3",
		post.Content, post.Id, post.UserId)
	return err
}

func (repository *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, content, user_id, create_at FROM posts WHERE id = $1", id)

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}

	for rows.Next() {
		if err = rows.Scan(&post.Id, &post.Content, &post.UserId, &post.CreateAt); err == nil {
			return &post, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return &post, nil
}

func (repository *PostgresRepository) ListPost(ctx context.Context, page uint64, size uint64) ([]*models.Post, error) {
	fmt.Printf("SELECT id, content, user_id, create_at FROM posts LIMIT %d OFFSET %d\n", size, page*size)
	rows, err := repository.db.QueryContext(ctx, "SELECT id, content, user_id, create_at FROM posts LIMIT $1 OFFSET $2", size, page*size)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var posts = []*models.Post{}

	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.Content, &post.UserId, &post.CreateAt); err == nil {
			posts = append(posts, &post)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, rows.Err()
	}
	return posts, nil
}

func (repository *PostgresRepository) Close() error {
	return repository.db.Close()
}
