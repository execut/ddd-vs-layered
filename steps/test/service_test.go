package test_test

import (
    "context"
    "strings"
    "testing"

    "effective-architecture/steps/labels"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabels_Live(t *testing.T) {
    t.Parallel()

    var service *labels.Service

    repository, err := labels.NewRepository(t.Context())
    require.NoError(t, err)
    historyRepository, err := labels.NewHistoryRepository(t.Context())
    require.NoError(t, err)

    _ = repository.Truncate(context.Background())
    _ = historyRepository.Truncate(context.Background())

    t.Cleanup(func() {
        _ = repository.Truncate(context.Background())
        _ = historyRepository.Truncate(context.Background())
    })

    t.Run("New", func(t *testing.T) {
        service = labels.NewService(repository, historyRepository)
        require.NotNil(t, service)
    })

    t.Run("CreateLabelTemplate", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturer)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturer)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyCreated)
    })

    t.Run("GetLabelTemplate", func(t *testing.T) {
        result, err := service.GetLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturer":{"organizationName":`+
            `"test manufacturer organization name"}}`, result)
    })

    t.Run("Обновлять данные шаблона", func(t *testing.T) {
        err = service.UpdateLabelTemplate(t.Context(), expectedUUID, expectedNewManufacturer)

        require.NoError(t, err)
        result, err := service.GetLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
        assert.Contains(t, result, expectedNewManufacturerOrganizationName)

        t.Run("и не давать это делать при ошибках из предыдущих пунктов", func(t *testing.T) {
            err = service.UpdateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: strings.Repeat("a", 256),
            })

            require.ErrorIs(t, labels.ErrLabelTemplateWrongManufacturerOrganizationName, err)
        })
    })

    t.Run("DeleteLabelTemplate", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), expectedUUID)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyDeleted)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать шаблон, "+
        "если длина Наименования организации производителя", func(t *testing.T) {
        t.Run("> 255", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: strings.Repeat("a", 256),
            })

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
        t.Run("= 0", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: "",
            })

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
    })

    t.Run("Указывать и получать поля Адрес, Email, сайт", func(t *testing.T) {
        t.Run("при создании", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturerWithAllFields)

            require.NoError(t, err)
            result, err := service.GetLabelTemplate(t.Context(), expectedUUID)
            require.NoError(t, err)
            assert.Contains(t, result, expectedManufacturerOrganizationAddress)
            assert.Contains(t, result, expectedManufacturerEmail)
            assert.Contains(t, result, expectedManufacturerSite)
        })
        t.Run("и обновлении", func(t *testing.T) {
            err := service.UpdateLabelTemplate(t.Context(), expectedUUID, expectedNewManufacturerWithAllFields)

            require.NoError(t, err)
            result, err := service.GetLabelTemplate(t.Context(), expectedUUID)
            require.NoError(t, err)
            assert.Contains(t, result, expectedNewManufacturerOrganizationAddress)
            assert.Contains(t, result, expectedNewManufacturerEmail)
            assert.Contains(t, result, expectedNewManufacturerSite)
        })
    })

    t.Run("При создании шаблона возвращать ошибку с понятным описанием", func(t *testing.T) {
        type args struct {
            fieldName    string
            manufacturer labels.Manufacturer
            err          error
        }

        t.Run("если длина следующих полей >255 или = 0", func(t *testing.T) {
            tests := []args{
                {
                    fieldName: "Адрес >255",
                    manufacturer: labels.Manufacturer{
                        OrganizationAddress: strings.Repeat("a", 256),
                        Email:               expectedManufacturerEmail,
                        Site:                expectedManufacturerSite,
                    },
                    err: labels.ErrLabelTemplateManufacturerOrganizationAddressWrongLen,
                },
                {
                    fieldName: "Email >255",
                    manufacturer: labels.Manufacturer{
                        OrganizationAddress: expectedManufacturerOrganizationAddress,
                        Email:               strings.Repeat("a", 256-9) + "@test.com",
                        Site:                expectedManufacturerSite,
                    },
                    err: labels.ErrLabelTemplateManufacturerEmailWrongLen,
                },
                {
                    fieldName: "Сайт >255",
                    manufacturer: labels.Manufacturer{
                        OrganizationAddress: expectedManufacturerOrganizationAddress,
                        Email:               expectedManufacturerEmail,
                        Site:                strings.Repeat("a", 256-4) + ".com",
                    },
                    err: labels.ErrLabelTemplateManufacturerSiteWrongLen,
                },
            }
            for _, testForCreateLabel := range tests {
                manufacturer := testForCreateLabel.manufacturer
                manufacturer.OrganizationName = expectedManufacturerOrganizationName

                err = service.CreateLabelTemplate(t.Context(), expectedUUID, manufacturer)

                require.ErrorIs(t, err, testForCreateLabel.err)
            }

            for _, testForUpdateLabel := range tests {
                manufacturer := testForUpdateLabel.manufacturer
                manufacturer.OrganizationName = expectedManufacturerOrganizationName

                err = service.UpdateLabelTemplate(t.Context(), expectedUUID, manufacturer)

                require.ErrorIs(t, err, testForUpdateLabel.err)
            }
        })

        t.Run("если формат не корректный", func(t *testing.T) {
            tests := []args{
                {
                    fieldName: "Email",
                    manufacturer: labels.Manufacturer{
                        OrganizationAddress: expectedManufacturerOrganizationAddress,
                        Email:               "asdasdsas",
                        Site:                expectedManufacturerSite,
                    },
                    err: labels.ErrLabelTemplateManufacturerEmailWrongFormat,
                },
                {
                    fieldName: "Сайт",
                    manufacturer: labels.Manufacturer{
                        OrganizationAddress: expectedManufacturerOrganizationAddress,
                        Email:               expectedManufacturerEmail,
                        Site:                "asdasdadasas",
                    },
                    err: labels.ErrLabelTemplateManufacturerSiteWrongFormat,
                },
            }

            for _, testForCreateLabel := range tests {
                manufacturer := testForCreateLabel.manufacturer
                manufacturer.OrganizationName = expectedManufacturerOrganizationName

                err = service.CreateLabelTemplate(t.Context(), expectedUUID, manufacturer)

                require.ErrorIs(t, err, testForCreateLabel.err)
            }

            for _, testForUpdateLabel := range tests {
                manufacturer := testForUpdateLabel.manufacturer
                manufacturer.OrganizationName = expectedManufacturerOrganizationName

                err = service.UpdateLabelTemplate(t.Context(), expectedUUID, manufacturer)

                require.ErrorIs(t, err, testForUpdateLabel.err)
            }
        })
    })

    t.Run("Чтобы писалась история операций над шаблонами с возможностью выводить все"+
        " данные в json", func(t *testing.T) {
        result, err := service.GetLabelHistory(t.Context(), expectedUUID)

        require.NoError(t, err)
        assert.JSONEq(t, `
[{
    "orderKey": 1,
    "action": "created",
    "newManufacturerOrganizationName": "test manufacturer organization name"
},
{
    "orderKey": 2,
    "action": "updated",
    "newManufacturerOrganizationName": "new test manufacturer organization name"
},
{
    "orderKey": 3,
    "action": "deleted"
},
{
    "orderKey": 4,
    "action": "created",
    "newManufacturerOrganizationName": "test manufacturer organization name",
    "newManufacturerOrganizationAddress": "test manufacturer organization address",
    "newManufacturerEmail": "email@test.com",
    "newManufacturerSite": "https://test.com"
},
{
    "orderKey": 5,
    "action": "updated",
    "newManufacturerOrganizationName": "new test manufacturer organization name",
    "newManufacturerOrganizationAddress": "new test manufacturer organization address",
    "newManufacturerEmail": "new-email@test.com",
    "newManufacturerSite": "https://new-test.com"
}]
`, result)
    })
}
