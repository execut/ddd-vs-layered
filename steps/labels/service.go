package labels

import (
    "context"
)

type IRepository interface {
    Insert(ctx context.Context, model LabelTemplate) error
    Truncate(ctx context.Context) error
}

type Service struct {
    repository IRepository
}

func NewService(repository IRepository) *Service {
    return &Service{
        repository: repository,
    }
}

func (s Service) CreateLabelTemplate(ctx context.Context, id string, manufacturerOrganizationName string) error {
    model := LabelTemplate{
        ID:                           id,
        ManufacturerOrganizationName: manufacturerOrganizationName,
    }

    err := s.repository.Insert(ctx, model)
    if err != nil {
        return err
    }

    return nil
}
