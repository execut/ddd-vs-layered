package domain_test

import (
	"effective-architecture/steps/contract/external"
	"testing"

	"effective-architecture/steps/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	expectedSKU           int64 = 1
	expectedProductName         = "test product name"
	expectedLabelFilePath       = "test label file path"
)

var (
	expectedLabelID, _   = domain.NewLabelID(expectedUUIDValue)
	expectedCategoryList = []domain.Category{
		expectedCategory1,
		expectedCategory2WithoutType,
	}
	expectedProduct, _   = domain.NewProduct(expectedProductName)
	expectedLabelFile, _ = domain.NewLabelFile(expectedLabelFilePath)
)

func TestLabel(t *testing.T) {
	ctrl := gomock.NewController(t)
	ozonServiceMock := domain.NewMockIServiceOzon(ctrl)
	labelRepository := domain.NewMockIRepository(ctrl)
	labelGeneratorMock := domain.NewMockILabelGenerator(ctrl)
	label := domain.NewLabel(expectedLabelID, ozonServiceMock, labelRepository, labelGeneratorMock)

	t.Run("14. Начинать генерацию этикетки по SKU", func(t *testing.T) {
		err := label.StartGeneration(expectedUserID, expectedSKU)

		require.NoError(t, err)

		t.Run("или такая генерация уже была запущена", func(t *testing.T) {
			err := label.StartGeneration(expectedUserID, expectedSKU)

			require.ErrorContains(t, err, "генерация этикетки с таким идентификатором уже существует")
		})
	})

	t.Run("20. Чтобы нельзя было работать с чужими записями", func(t *testing.T) {
		err := label.FillData(t.Context(), expectedOtherUserID)

		require.ErrorIs(t, err, domain.ErrLabelWrongUser)
	})

	t.Run("16. Наполнять этикетку данными из внешнего API и вычислять по ним шаблон", func(t *testing.T) {
		t.Run("и получать ошибку, если SKU отсутствует", func(t *testing.T) {
			ozonServiceMock.EXPECT().ProductData(gomock.Any(), gomock.Any()).Return(nil, domain.Product{},
				external.ErrSkuNotFound)

			err := label.FillData(t.Context(), expectedUserID)

			require.ErrorIs(t, err, external.ErrSkuNotFound)
		})

		t.Run("или для категории SKU нет шаблона", func(t *testing.T) {
			ozonServiceMock.EXPECT().ProductData(gomock.Any(), gomock.Any()).Return(expectedCategoryList, domain.Product{}, nil)
			labelRepository.EXPECT().LoadByCategoryList(t.Context(), expectedCategoryList).Return(nil, nil)

			err := label.FillData(t.Context(), expectedUserID)

			require.ErrorIs(t, err, domain.ErrLabelTemplateForCategoryNotFound)
		})

		ozonServiceMock.EXPECT().ProductData(t.Context(), expectedSKU).Return(expectedCategoryList, expectedProduct, nil)
		labelRepository.EXPECT().LoadByCategoryList(t.Context(), []domain.Category{
			expectedCategory1,
			expectedCategory2WithoutType,
		}).Return(&domain.LabelTemplate{
			ID:           expectedLabelTemplateID,
			Manufacturer: expectedManufacturer,
		}, nil)

		err := label.FillData(t.Context(), expectedUserID)

		require.NoError(t, err)
		assert.Equal(t, expectedLabelTemplateID, label.TemplateID)
		assert.Equal(t, expectedProduct, label.Product)
		assert.Equal(t, domain.LabelStatusDataFilled, label.Status)
	})

	t.Run("18. Генерировать этикетку через внешний сервис, передавая ему все нужные данные", func(t *testing.T) {
		labelGeneratorMock.EXPECT().Generate(t.Context(), expectedProduct, expectedManufacturer, expectedSKU).
			Return(expectedLabelFile, nil)

		err := label.Generate(t.Context(), expectedUserID)

		require.NoError(t, err)
		assert.Equal(t, domain.LabelStatusGenerated, label.Status)
		require.NotNil(t, label.File)
		assert.Equal(t, expectedLabelFile, *label.File)
	})
}
