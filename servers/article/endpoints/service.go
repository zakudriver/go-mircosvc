package endpoints

import (
	"context"
	"database/sql"

	"github.com/kum0/go-mircosvc/shared/validator"

	articlePb "github.com/kum0/go-mircosvc/pb/article"
)

type ArticleServicer interface {
	GetCategories(context.Context) (*articlePb.GetCategoriesResponse, error)
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

func (svc *ArticleService) GetCategories(_ context.Context) (*articlePb.GetCategoriesResponse, error) {
	return &articlePb.GetCategoriesResponse{
		Count: 1,
		Data: []*articlePb.CategoryResponse{
			&articlePb.CategoryResponse{
				Id:   1,
				Name: "11",
			},
		},
	}, nil
}
