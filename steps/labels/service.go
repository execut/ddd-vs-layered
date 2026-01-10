package labels

import (
    "context"
    "encoding/json"
    "errors"
    "strings"
)

var (
    ErrLabelTemplateAlreadyCreated                    = errors.New("попытка создать уже существующий шаблон")
    ErrLabelTemplateAlreadyDeleted                    = errors.New("попытка удалить уже удалённый шаблон")
    ErrLabelTemplateWrongManufacturerOrganizationName = errors.New("название организации производителя должно " +
        "быть до 255 символов в длину")
)

type IRepository interface {
    Insert(ctx context.Context, model LabelTemplate) error
    Find(ctx context.Context, id string) (LabelTemplate, error)
    Update(ctx context.Context, model LabelTemplate) error
    Truncate(ctx context.Context) error
    Delete(ctx context.Context, id string) error
}

type Service struct {
    repository IRepository
}

func NewService(repository IRepository) *Service {
    return &Service{
        repository: repository,
    }
}

func (s Service) CreateLabelTemplate(ctx context.Context, labelTemplateID string,
    manufacturerOrganizationName string) error {
    err := s.validateVarchar(manufacturerOrganizationName)
    if err != nil {
        return err
    }

    model := LabelTemplate{
        ID:                           labelTemplateID,
        ManufacturerOrganizationName: manufacturerOrganizationName,
    }

    err = s.repository.Insert(ctx, model)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"label_templates_pkey\"") {
            return ErrLabelTemplateAlreadyCreated
        }

        return err
    }

    return nil
}

func (s Service) UpdateLabelTemplate(ctx context.Context, labelTemplateID string,
    manufacturerOrganizationName string) error {
    err := s.validateVarchar(manufacturerOrganizationName)
    if err != nil {
        return err
    }

    model := LabelTemplate{
        ID:                           labelTemplateID,
        ManufacturerOrganizationName: manufacturerOrganizationName,
    }

    err = s.repository.Update(ctx, model)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) GetLabelTemplate(ctx context.Context, id string) (string, error) {
    model, err := s.repository.Find(ctx, id)
    if err != nil {
        return "", err
    }

    resultMarshaled, err := json.Marshal(model)
    if err != nil {
        return "", err
    }

    return string(resultMarshaled), nil
}

func (s Service) DeleteLabelTemplate(ctx context.Context, id string) error {
    err := s.repository.Delete(ctx, id)
    if err != nil {
        if errors.Is(err, ErrCouldNotDelete) {
            return ErrLabelTemplateAlreadyDeleted
        }

        return err
    }

    return nil
}

func (s Service) validateVarchar(manufacturerOrganizationName string) error {
    if len(manufacturerOrganizationName) > 255 || len(manufacturerOrganizationName) == 0 {
        return ErrLabelTemplateWrongManufacturerOrganizationName
    }

    return nil
}
