package services

import (
	"github.com/pragmaticdev85/go-microservice/internal/repositories"
)

type ExampleService struct {
	repo *repositories.ExampleRepository
}

func NewExampleService(repo *repositories.ExampleRepository) *ExampleService {
	return &ExampleService{repo: repo}
}

func (s *ExampleService) CreateExample(example *repositories.Example) (*repositories.Example, error) {
	return s.repo.Create(example)
}

func (s *ExampleService) GetExampleByID(id string) (*repositories.Example, error) {
	return s.repo.FindByID(id)
}

func (s *ExampleService) GetExamples() ([]repositories.Example, error) {
	return s.repo.GetAllExamples()
}

// Add other service methods
