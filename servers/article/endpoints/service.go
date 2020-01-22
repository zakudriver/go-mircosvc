package endpoints

import (
	"context"
	"database/sql"

	"github.com/kum0/go-mircosvc/shared/validator"
)

type ArticleServicer interface {
	GetCategories(context.Context) (error, error)
}

func NewArticleService(db *sql.DB) ArticleServicer {
	return &ArticleService{
		db,
		validator.NewValidator(),
	}
}

type ArticleService struct {
	mysql     *sql.DB
	validator *validator.Validator
}

func (svc *ArticleService) GetCategories(_ context.Context) (error, error) {
	return nil, nil
}
