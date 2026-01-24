package e2e_test

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabelLive(t *testing.T) {
    t.Parallel()
    t.Run("1. Создавать шаблон этикетки товара с UUID и Наименованием организации производителя", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-create-template",
            "--id", expectedUUID,
            "--manufacturer-organization-name", expectedManufacturerOrganizationName,
        })

        require.NoError(t, err, "output:"+output)
        assert.Equal(t, "1", output)
    })
    t.Run("2. Получать данные шаблона в JSON", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-get-template",
            "--id", expectedUUID,
        })

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturer":{"organizationName":"test manufacturer organization name"}}`, output)
    })
    t.Run("5. Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        out, err := runBinary([]string{
            "labels-create-template",
            "--id", expectedUUID,
            "--manufacturer-organization-name", expectedManufacturerOrganizationName,
        })

        require.Error(t, err)
        assert.Contains(t, out, "попытка создать уже существующий шаблон")
    })
    t.Run("6. Обновлять данные шаблона", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-update-template",
            "--id", expectedUUID,
            "--manufacturer-organization-name", expectedNewManufacturerOrganizationName,
        })

        require.NoError(t, err)
        assert.Equal(t, "1", output)
        output, err = runBinary([]string{
            "labels-get-template",
            "--id", expectedUUID,
        })
        require.NoError(t, err)
        assert.Contains(t, output, expectedNewManufacturerOrganizationName)

        t.Run("и не давать это делать при ошибках из предыдущих пунктов", func(t *testing.T) {
            out, err := runBinary([]string{
                "labels-update-template",
                "--id", expectedUUID,
                "--manufacturer-organization-name", "",
            })

            require.Error(t, err)
            assert.Contains(t, out, "название организации производителя должно быть до 255 символов в длину")
        })
    })
    t.Run("4. Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-delete-template",
            "--id", expectedUUID,
        })

        require.NoError(t, err)
        assert.Equal(t, `1`, output)
    })

    t.Run("8. Указывать и получать поля Адрес, Email, сайт", func(t *testing.T) {
        t.Run("при создании", func(t *testing.T) {
            output, err := runBinary([]string{
                "labels-create-template",
                "--id", expectedUUID,
                "--manufacturer-organization-name", expectedManufacturerOrganizationName,
                "--manufacturer-organization-address", expectedManufacturerOrganizationAddress,
                "--manufacturer-email", expectedManufacturerEmail,
                "--manufacturer-site", expectedManufacturerSite,
            })

            require.NoError(t, err)
            assert.Equal(t, "1", output)
            output, err = runBinary([]string{
                "labels-get-template",
                "--id", expectedUUID,
            })

            require.NoError(t, err)
            assert.JSONEq(t, `
{
    "id":"123e4567-e89b-12d3-a456-426655440000",
    "manufacturer": {
        "organizationName": "test manufacturer organization name",
        "organizationAddress": "test manufacturer organization address",
        "email": "test@test.com",
        "site": "https://test.com"
    }
}`, output)
        })
        t.Run("и обновлении", func(t *testing.T) {
            output, err := runBinary([]string{
                "labels-update-template",
                "--id", expectedUUID,
                "--manufacturer-organization-name", expectedNewManufacturerOrganizationName,
                "--manufacturer-organization-address", expectedNewManufacturerOrganizationAddress,
                "--manufacturer-email", expectedNewManufacturerEmail,
                "--manufacturer-site", expectedNewManufacturerSite,
            })

            require.NoError(t, err)
            assert.Equal(t, "1", output)
            output, err = runBinary([]string{
                "labels-get-template",
                "--id", expectedUUID,
            })
            require.NoError(t, err)
            assert.Contains(t, output, expectedNewManufacturerOrganizationAddress)
            assert.Contains(t, output, expectedNewManufacturerEmail)
            assert.Contains(t, output, expectedNewManufacturerSite)
        })
    })
    t.Run("5. Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        _, err := runBinary([]string{
            "labels-delete-template",
            "--id", expectedUUID,
        })
        require.NoError(t, err)

        out, err := runBinary([]string{
            "labels-delete-template",
            "--id", expectedUUID,
        })

        require.Error(t, err)
        assert.Contains(t, out, "попытка удалить уже удалённый шаблон")
    })
    t.Run("6. Чтобы возвращалась уникальная ошибка при попытке создать шаблон, если длина Наименования организации производителя", func(t *testing.T) {
        t.Run("> 255", func(t *testing.T) {
            out, err := runBinary([]string{
                "labels-create-template",
                "--id", expectedUUID,
                "--manufacturer-organization-name", strings.Repeat("a", 256),
            })

            require.Error(t, err)
            assert.Contains(t, out, "название организации производителя должно быть до 255 символов в длину")
        })
        t.Run("или =0", func(t *testing.T) {
            out, err := runBinary([]string{
                "labels-create-template",
                "--id", expectedUUID,
                "--manufacturer-organization-name", "",
            })

            require.Error(t, err)
            assert.Contains(t, out, "название организации производителя должно быть до 255 символов в длину")
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
            for _, tt := range tests {
                out, err := runBinary([]string{
                    "labels-create-template",
                    "--id", expectedUUID,
                    "--manufacturer-organization-name", expectedManufacturerOrganizationName,
                    "--manufacturer-organization-address", tt.organizationAddress,
                    "--manufacturer-email", tt.email,
                    "--manufacturer-site", tt.site,
                })

                require.Error(t, err)
                assert.Contains(t, out, tt.errorMessage)
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
            for _, tt := range tests {
                out, err := runBinary([]string{
                    "labels-create-template",
                    "--id", expectedUUID,
                    "--manufacturer-organization-name", expectedManufacturerOrganizationName,
                    "--manufacturer-organization-address", tt.organizationAddress,
                    "--manufacturer-email", tt.email,
                    "--manufacturer-site", tt.site,
                })

                require.Error(t, err)
                assert.Contains(t, out, tt.errorMessage)
            }
        })
    })

    t.Run("10. Чтобы писалась история операций над шаблонами с возможностью выводить все данные в json", func(t *testing.T) {
        out, err := runBinary([]string{
            "labels-template-history",
            "--id", expectedUUID,
        })

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
    "newManufacturerEmail": "test@test.com",
    "newManufacturerSite": "https://test.com"
},
{
    "orderKey": 5,
    "action": "updated",
    "newManufacturerOrganizationName": "new test manufacturer organization name",
    "newManufacturerOrganizationAddress": "new test manufacturer organization address",
    "newManufacturerEmail": "new-test@test.com",
    "newManufacturerSite": "https://new-test.com"
},
{
    "orderKey": 6,
    "action": "deleted"
}]
`, out)
    })

    t.Run("11. Привязывать шаблон к списку категорий или категорий+типов", func(t *testing.T) {
        out, err := runBinary([]string{
            "labels-template-add-category-list",
            "--id", expectedUUID,
            "--category-id-list", "1,2-3",
        })

        require.NoError(t, err)
        assert.Equal(t, `1`, out)

        t.Run("и получать ошибку при попытке привязать уже существующую категорию", func(t *testing.T) {
            out, err := runBinary([]string{
                "labels-template-add-category-list",
                "--id", expectedUUID,
                "--category-id-list", "2-3",
            })

            require.Error(t, err)
            assert.Contains(t, out, "категория 2 с типом 3 уже привязана к шаблону")
        })
    })
}
