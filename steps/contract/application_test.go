package contract_test

import (
    "strings"
    "testing"

    "effective-architecture/steps/contract"
    "effective-architecture/steps/presentation"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const (
    expectedID                                 = "123e4567-e89b-12d3-a456-426655440000"
    expectedManufacturerOrganizationName       = "test manufacturer organization name"
    expectedManufacturerOrganizationAddress    = "test manufacturer organization address"
    expectedManufacturerEmail                  = "test@test.com"
    expectedManufacturerSite                   = "https://test.com"
    expectedNewManufacturerOrganizationName    = "new test manufacturer organization name"
    expectedNewManufacturerOrganizationAddress = "new test manufacturer organization address"
    expectedNewManufacturerEmail               = "new-test@test.com"
    expectedNewManufacturerSite                = "https://new-test.com"
)

var (
    expectedCategory2TypeID = "3"
)

func TestLabelTemplate_Live(t *testing.T) {
    t.Parallel()

    app, err := presentation.NewApplication(t.Context())
    require.NoError(t, err)
    _ = app.Cleanup(t.Context(), expectedID)
    t.Cleanup(func() {
        _ = app.Cleanup(t.Context(), expectedID)
    })

    t.Run("1. Создавать шаблон этикетки товара с UUID и Наименованием организации производителя", func(t *testing.T) {
        err := app.Create(t.Context(), expectedID, contract.Manufacturer{
            OrganizationName: expectedManufacturerOrganizationName,
        })

        require.NoError(t, err)
    })
    t.Run("2. Получать данные шаблона в JSON", func(t *testing.T) {
        result, err := app.Get(t.Context(), expectedID)

        require.NoError(t, err)
        assert.Equal(t, contract.LabelTemplate{
            ID: expectedID,
            Manufacturer: contract.Manufacturer{
                OrganizationName: expectedManufacturerOrganizationName,
            },
        }, result)
    })
    t.Run("4. Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        err := app.Create(t.Context(), expectedID, contract.Manufacturer{
            OrganizationName: expectedManufacturerOrganizationName,
        })

        require.Error(t, err)
        assert.ErrorContains(t, err, "попытка создать уже существующий шаблон")
    })
    t.Run("7. Обновлять данные шаблона", func(t *testing.T) {
        err := app.Update(t.Context(), expectedID, contract.Manufacturer{
            OrganizationName: expectedNewManufacturerOrganizationName,
        })

        require.NoError(t, err)
        result, err := app.Get(t.Context(), expectedID)
        require.NoError(t, err)
        assert.Contains(t, result.Manufacturer.OrganizationName, expectedNewManufacturerOrganizationName)

        t.Run("и не давать это делать при ошибках из предыдущих пунктов", func(t *testing.T) {
            err := app.Update(t.Context(), expectedID, contract.Manufacturer{
                OrganizationName: "",
            })

            require.Error(t, err)
            assert.ErrorContains(t, err, "название организации производителя должно быть до 255 символов в длину")
        })
    })
    t.Run("3. Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
        err := app.Delete(t.Context(), expectedID)

        require.NoError(t, err)
    })

    t.Run("8. Указывать и получать поля Адрес, Email, сайт", func(t *testing.T) {
        t.Run("при создании", func(t *testing.T) {
            err := app.Create(t.Context(), expectedID, contract.Manufacturer{
                OrganizationName:    expectedManufacturerOrganizationName,
                OrganizationAddress: expectedManufacturerOrganizationAddress,
                Email:               expectedManufacturerEmail,
                Site:                expectedManufacturerSite,
            })

            require.NoError(t, err)
            result, err := app.Get(t.Context(), expectedID)
            require.NoError(t, err)
            assert.Equal(t, expectedManufacturerOrganizationName, result.Manufacturer.OrganizationName)
            assert.Equal(t, expectedManufacturerOrganizationAddress, result.Manufacturer.OrganizationAddress)
            assert.Equal(t, expectedManufacturerEmail, result.Manufacturer.Email)
            assert.Equal(t, expectedManufacturerSite, result.Manufacturer.Site)
        })
        t.Run("и обновлении", func(t *testing.T) {
            err := app.Update(t.Context(), expectedID, contract.Manufacturer{
                OrganizationName:    expectedNewManufacturerOrganizationName,
                OrganizationAddress: expectedNewManufacturerOrganizationAddress,
                Email:               expectedNewManufacturerEmail,
                Site:                expectedNewManufacturerSite,
            })

            require.NoError(t, err)
            result, err := app.Get(t.Context(), expectedID)
            require.NoError(t, err)
            assert.Equal(t, expectedNewManufacturerOrganizationName, result.Manufacturer.OrganizationName)
            assert.Equal(t, expectedNewManufacturerOrganizationAddress, result.Manufacturer.OrganizationAddress)
            assert.Equal(t, expectedNewManufacturerEmail, result.Manufacturer.Email)
            assert.Equal(t, expectedNewManufacturerSite, result.Manufacturer.Site)
        })
    })
    t.Run("5. Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := app.Delete(t.Context(), expectedID)
        require.NoError(t, err)

        err = app.Delete(t.Context(), expectedID)

        require.Error(t, err)
        require.ErrorContains(t, err, "попытка удалить уже удалённый шаблон")
    })
    t.Run("6. Чтобы возвращалась уникальная ошибка при попытке создать шаблон, если длина Наименования "+
        "организации производителя", func(t *testing.T) {
        t.Run("> 255", func(t *testing.T) {
            err := app.Create(t.Context(), expectedID, contract.Manufacturer{
                OrganizationName: strings.Repeat("a", 256),
            })

            require.Error(t, err)
            require.ErrorContains(t, err, "название организации производителя должно быть до 255 символов в длину")
        })
        t.Run("или =0", func(t *testing.T) {
            err := app.Create(t.Context(), expectedID, contract.Manufacturer{
                OrganizationName: "",
            })

            require.Error(t, err)
            require.ErrorContains(t, err, "название организации производителя должно быть до 255 символов в длину")
        })
    })

    t.Run("9. При создании шаблона возвращать ошибку с понятным описанием", func(t *testing.T) {
        type args struct {
            fieldName           string
            organizationAddress string
            email               string
            site                string
            errorMessage        string
        }

        t.Run("если длина следующих полей >255 или = 0", func(t *testing.T) {
            tests := []args{
                {
                    fieldName:           "Адрес >255",
                    organizationAddress: strings.Repeat("a", 256),
                    email:               expectedManufacturerEmail,
                    site:                expectedManufacturerSite,
                    errorMessage:        "адрес должен быть до 255 символов в длину",
                },
                {
                    fieldName:           "Email >255",
                    organizationAddress: expectedManufacturerOrganizationAddress,
                    email:               strings.Repeat("a", 256-9) + "@test.com",
                    site:                expectedManufacturerSite,
                    errorMessage:        "email должен быть до 255 символов в длину",
                },
                {
                    fieldName:           "Сайт >255",
                    organizationAddress: expectedManufacturerOrganizationAddress,
                    email:               expectedManufacturerEmail,
                    site:                strings.Repeat("a", 256-4) + ".com",
                    errorMessage:        "сайт должен быть до 255 символов в длину",
                },
            }
            for _, ttForLengthErrors := range tests {
                err := app.Create(t.Context(), expectedID, contract.Manufacturer{
                    OrganizationName:    expectedManufacturerOrganizationName,
                    OrganizationAddress: ttForLengthErrors.organizationAddress,
                    Email:               ttForLengthErrors.email,
                    Site:                ttForLengthErrors.site,
                })

                require.Error(t, err)
                assert.ErrorContains(t, err, ttForLengthErrors.errorMessage)
            }
        })

        t.Run("если формат не корректный", func(t *testing.T) {
            tests := []args{
                {
                    fieldName:           "Email",
                    organizationAddress: expectedManufacturerOrganizationAddress,
                    email:               "asdasdsas",
                    site:                expectedManufacturerSite,
                    errorMessage:        "email имеет не корректный формат",
                },
                {
                    fieldName:           "Сайт",
                    organizationAddress: expectedManufacturerOrganizationAddress,
                    email:               expectedManufacturerEmail,
                    site:                "asdasdadasas",
                    errorMessage:        "сайт имеет не корректный формат",
                },
            }
            for _, ttForWrongFormat := range tests {
                err := app.Create(t.Context(), expectedID, contract.Manufacturer{
                    OrganizationName:    expectedManufacturerOrganizationName,
                    OrganizationAddress: ttForWrongFormat.organizationAddress,
                    Email:               ttForWrongFormat.email,
                    Site:                ttForWrongFormat.site,
                })

                require.Error(t, err)
                assert.ErrorContains(t, err, ttForWrongFormat.errorMessage)
            }
        })
    })

    t.Run("10. Чтобы писалась история операций над шаблонами с возможностью"+
        " выводить все данные в json", func(t *testing.T) {
        result, err := app.HistoryList(t.Context(), expectedID)

        require.NoError(t, err)
        assert.Equal(t, []contract.LabelTemplateHistoryRow{
            {
                OrderKey:                        1,
                Action:                          "created",
                NewManufacturerOrganizationName: expectedManufacturerOrganizationName,
            },
            {
                OrderKey:                        2,
                Action:                          "updated",
                NewManufacturerOrganizationName: expectedNewManufacturerOrganizationName,
            },
            {
                OrderKey: 3,
                Action:   "deleted",
            },
            {
                OrderKey:                           4,
                Action:                             "created",
                NewManufacturerOrganizationName:    expectedManufacturerOrganizationName,
                NewManufacturerOrganizationAddress: expectedManufacturerOrganizationAddress,
                NewManufacturerEmail:               expectedManufacturerEmail,
                NewManufacturerSite:                expectedManufacturerSite,
            },
            {
                OrderKey:                           5,
                Action:                             "updated",
                NewManufacturerOrganizationName:    expectedNewManufacturerOrganizationName,
                NewManufacturerOrganizationAddress: expectedNewManufacturerOrganizationAddress,
                NewManufacturerEmail:               expectedNewManufacturerEmail,
                NewManufacturerSite:                expectedNewManufacturerSite,
            },
            {
                OrderKey: 6,
                Action:   "deleted",
            },
        }, result)
    })

    t.Run("11. Привязывать шаблон к списку категорий или категорий+типов", func(t *testing.T) {
        err = app.AddCategoryList(t.Context(), expectedID, []contract.Category{
            {
                CategoryID: "1",
            },
            {
                CategoryID: "2",
                TypeID:     &expectedCategory2TypeID,
            },
        })

        require.NoError(t, err)

        t.Run("и получать ошибку при попытке привязать уже существующую категорию", func(t *testing.T) {
            err = app.AddCategoryList(t.Context(), expectedID, []contract.Category{
                {
                    CategoryID: "2",
                    TypeID:     &expectedCategory2TypeID,
                },
            })

            require.NoError(t, err)

            require.Error(t, err)
            assert.ErrorContains(t, err, "категория уже привязана к шаблону (категория 2, тип 3)")
        })
    })
}
