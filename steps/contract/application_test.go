package contract_test

import (
	"strings"
	"testing"

	"effective-architecture/steps/contract"
	"effective-architecture/steps/contract/external"
	"effective-architecture/steps/presentation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	expectedTemplateID                               = "123e4567-e89b-12d3-a456-426655440001"
	expectedLabelGenerationID                        = "223e4567-e89b-12d3-a456-426655440001"
	expectedManufacturerOrganizationName             = "test manufacturer organization name"
	expectedManufacturerOrganizationAddress          = "test manufacturer organization address"
	expectedManufacturerEmail                        = "test@test.com"
	expectedManufacturerSite                         = "https://test.com"
	expectedNewManufacturerOrganizationName          = "new test manufacturer organization name"
	expectedNewManufacturerOrganizationAddress       = "new test manufacturer organization address"
	expectedNewManufacturerEmail                     = "new-test@test.com"
	expectedNewManufacturerSite                      = "https://new-test.com"
	expectedSKU                                int64 = 555
)

var (
	expectedCategory2TypeID = "3"
)

func TestLabelTemplate_Live(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	externalOzonMock := external.NewMockIExternalServiceOzon(ctrl)
	app, err := presentation.NewApplication(t.Context(), externalOzonMock)
	require.NoError(t, err)
	_ = app.Cleanup(t.Context(), expectedTemplateID)
	_ = app.Cleanup(t.Context(), expectedLabelGenerationID)
	t.Cleanup(func() {
		_ = app.Cleanup(t.Context(), expectedTemplateID)
		_ = app.Cleanup(t.Context(), expectedLabelGenerationID)
	})

	t.Run("1. Создавать шаблон этикетки товара с UUID и Наименованием организации производителя", func(t *testing.T) {
		err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
			OrganizationName: expectedManufacturerOrganizationName,
		})

		require.NoError(t, err)
	})
	t.Run("2. Получать данные шаблона в JSON", func(t *testing.T) {
		result, err := app.Get(t.Context(), expectedTemplateID)

		require.NoError(t, err)
		assert.Equal(t, contract.LabelTemplate{
			ID: expectedTemplateID,
			Manufacturer: contract.Manufacturer{
				OrganizationName: expectedManufacturerOrganizationName,
			},
		}, result)
	})
	t.Run("4. Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
		err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
			OrganizationName: expectedManufacturerOrganizationName,
		})

		require.Error(t, err)
		assert.ErrorContains(t, err, "попытка создать уже существующий шаблон")
	})
	t.Run("7. Обновлять данные шаблона", func(t *testing.T) {
		err := app.Update(t.Context(), expectedTemplateID, contract.Manufacturer{
			OrganizationName: expectedNewManufacturerOrganizationName,
		})

		require.NoError(t, err)
		result, err := app.Get(t.Context(), expectedTemplateID)
		require.NoError(t, err)
		assert.Contains(t, result.Manufacturer.OrganizationName, expectedNewManufacturerOrganizationName)

		t.Run("и не давать это делать при ошибках из предыдущих пунктов", func(t *testing.T) {
			err := app.Update(t.Context(), expectedTemplateID, contract.Manufacturer{
				OrganizationName: "",
			})

			require.Error(t, err)
			assert.ErrorContains(t, err, "название организации производителя должно быть до 255 символов в длину")
		})
	})
	t.Run("3. Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
		err := app.Delete(t.Context(), expectedTemplateID)

		require.NoError(t, err)
	})

	t.Run("8. Указывать и получать поля Адрес, Email, сайт", func(t *testing.T) {
		t.Run("при создании", func(t *testing.T) {
			err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
				OrganizationName:    expectedManufacturerOrganizationName,
				OrganizationAddress: expectedManufacturerOrganizationAddress,
				Email:               expectedManufacturerEmail,
				Site:                expectedManufacturerSite,
			})

			require.NoError(t, err)
			result, err := app.Get(t.Context(), expectedTemplateID)
			require.NoError(t, err)
			assert.Equal(t, expectedManufacturerOrganizationName, result.Manufacturer.OrganizationName)
			assert.Equal(t, expectedManufacturerOrganizationAddress, result.Manufacturer.OrganizationAddress)
			assert.Equal(t, expectedManufacturerEmail, result.Manufacturer.Email)
			assert.Equal(t, expectedManufacturerSite, result.Manufacturer.Site)
		})
		t.Run("и обновлении", func(t *testing.T) {
			err := app.Update(t.Context(), expectedTemplateID, contract.Manufacturer{
				OrganizationName:    expectedNewManufacturerOrganizationName,
				OrganizationAddress: expectedNewManufacturerOrganizationAddress,
				Email:               expectedNewManufacturerEmail,
				Site:                expectedNewManufacturerSite,
			})

			require.NoError(t, err)
			result, err := app.Get(t.Context(), expectedTemplateID)
			require.NoError(t, err)
			assert.Equal(t, expectedNewManufacturerOrganizationName, result.Manufacturer.OrganizationName)
			assert.Equal(t, expectedNewManufacturerOrganizationAddress, result.Manufacturer.OrganizationAddress)
			assert.Equal(t, expectedNewManufacturerEmail, result.Manufacturer.Email)
			assert.Equal(t, expectedNewManufacturerSite, result.Manufacturer.Site)
		})
	})
	t.Run("6. Чтобы возвращалась уникальная ошибка при попытке создать шаблон, если длина Наименования "+
		"организации производителя", func(t *testing.T) {
		t.Run("> 255", func(t *testing.T) {
			err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
				OrganizationName: strings.Repeat("a", 256),
			})

			require.Error(t, err)
			require.ErrorContains(t, err, "название организации производителя должно быть до 255 символов в длину")
		})
		t.Run("или =0", func(t *testing.T) {
			err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
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
				err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
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
				err := app.Create(t.Context(), expectedTemplateID, contract.Manufacturer{
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

	var (
		expectedCategory1IDAsInt64 int64 = 1
		expectedCategory1ID              = "1"
		expectedCategory2ID              = "2"
		expectedCategory1                = contract.Category{
			CategoryID: expectedCategory1ID,
		}
		expectedCategory2 = contract.Category{
			CategoryID: expectedCategory2ID,
			TypeID:     &expectedCategory2TypeID,
		}
	)

	t.Run("11. Привязывать шаблон к списку категорий или категорий+типов", func(t *testing.T) {
		err = app.AddCategoryList(t.Context(), expectedTemplateID, []contract.Category{
			expectedCategory1,
			expectedCategory2,
		})

		require.NoError(t, err)

		t.Run("и получать ошибку при попытке привязать уже существующую категорию", func(t *testing.T) {
			err = app.AddCategoryList(t.Context(), expectedTemplateID, []contract.Category{
				expectedCategory2,
			})

			require.Error(t, err)
			assert.ErrorContains(t, err, "категория уже привязана к шаблону (категория 2, тип 3)")
		})
	})

	t.Run("12. Отвязывать шаблон от списка категорий или категорий+типов", func(t *testing.T) {
		err = app.UnlinkCategoryList(t.Context(), expectedTemplateID, []contract.Category{
			expectedCategory1,
			expectedCategory2,
		})

		require.NoError(t, err)

		t.Run("и получать ошибку при попытке отвязать уже отвязанную категорию", func(t *testing.T) {
			err = app.UnlinkCategoryList(t.Context(), expectedTemplateID, []contract.Category{
				expectedCategory2,
			})

			require.Error(t, err)
			assert.ErrorContains(t, err, "категория уже отвязана от шаблона (категория 2, тип 3)")
		})
	})

	t.Run("10. Чтобы писалась история операций над шаблонами с возможностью"+
		" выводить все данные в json"+
		"12. Смотреть историю добавления и удаления категорий в шаблоне", func(t *testing.T) {
		result, err := app.HistoryList(t.Context(), expectedTemplateID)

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
				Action:   "category_list_added",
				CategoryList: []contract.Category{
					expectedCategory1,
					expectedCategory2,
				},
			},
			{
				OrderKey: 7,
				Action:   "category_list_unlinked",
				CategoryList: []contract.Category{
					expectedCategory1,
					expectedCategory2,
				},
			},
		}, result)
	})

	t.Run("14. Начинать генерацию этикетки по SKU", func(t *testing.T) {
		t.Run("и получать ошибку, если SKU отсутствует", func(t *testing.T) {
			externalOzonMock.EXPECT().Product(gomock.Any(), gomock.Any()).
				Return(external.Product{}, external.ErrSkuNotFound)

			err = app.StartLabelGeneration(t.Context(), expectedLabelGenerationID, expectedSKU)

			require.ErrorContains(t, err, "sku не найден")
		})

		t.Run("или для категории SKU нет шаблона", func(t *testing.T) {
			externalOzonMock.EXPECT().Product(gomock.Any(), gomock.Any()).Return(external.Product{
				Category: external.CategoryWithType{
					Category: external.Category{
						ID: 11,
					},
					TypeID: 11,
				},
			}, nil)

			err = app.StartLabelGeneration(t.Context(), expectedLabelGenerationID, expectedSKU)

			require.ErrorContains(t, err, "шаблон этикетки для SKU не найден")
		})

		err = app.AddCategoryList(t.Context(), expectedTemplateID, []contract.Category{
			expectedCategory1,
		})
		require.NoError(t, err)
		externalOzonMock.EXPECT().Product(t.Context(), expectedSKU).Return(external.Product{
			Category: external.CategoryWithType{
				Category: external.Category{
					ID: 5,
					ParentCategory: &external.Category{
						ParentCategory: &external.Category{
							ID: expectedCategory1IDAsInt64,
						},
					},
				},
				TypeID: 2,
			},
		}, nil)

		err := app.StartLabelGeneration(t.Context(), expectedLabelGenerationID, expectedSKU)

		require.NoError(t, err)

		t.Run("или такая генерация уже была запущена", func(t *testing.T) {
			err = app.StartLabelGeneration(t.Context(), expectedLabelGenerationID, expectedSKU)

			require.ErrorContains(t, err, "генерация этикетки с таким идентификатором уже существует")
		})
	})

	t.Run("15. Получать статус генерации этикетки по SKU", func(t *testing.T) {
		response, err := app.LabelGeneration(t.Context(), expectedLabelGenerationID)

		require.NoError(t, err)

		assert.Equal(t, contract.LabelGenerationStatusCreated, response.Status)
	})

	t.Run("5. Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
		err := app.Delete(t.Context(), expectedTemplateID)
		require.NoError(t, err)

		err = app.Delete(t.Context(), expectedTemplateID)

		require.Error(t, err)
		require.ErrorContains(t, err, "попытка удалить уже удалённый шаблон")
	})
}
