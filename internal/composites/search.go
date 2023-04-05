package composites

import (
	"main/internal/microservices/search"
	"main/internal/microservices/search/repository"
	"main/internal/microservices/search/usecase"
)

type SearchComposite struct {
	Storage search.Storage
	Service *usecase.Service
}

func NewSearchComposite(postgresComposite *PostgresDBComposite) (*SearchComposite, error) {
	storage := repository.NewStorage(postgresComposite.db)
	service := usecase.NewService(storage)
	return &SearchComposite{
		Storage: storage,
		Service: service,
	}, nil
}
