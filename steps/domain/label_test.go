package domain_test

import (
    "testing"

    "effective-architecture/steps/contract/external"
    "effective-architecture/steps/domain"

    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

const (
    expectedSKU int64 = 1
)

var (
    expectedLabelID, _   = domain.NewLabelID(expectedUUIDValue)
    expectedCategoryList = []domain.Category{
        expectedCategory1,
        expectedCategory2WithoutType,
    }
)

func TestLabel(t *testing.T) {
    ctrl := gomock.NewController(t)
    ozonServiceMock := domain.NewMockIServiceOzon(ctrl)
    labelRepository := domain.NewMockIRepository(ctrl)
    label := domain.NewLabel(expectedLabelID, ozonServiceMock, labelRepository)

    t.Run("и получать ошибку, если SKU отсутствует", func(t *testing.T) {
        ozonServiceMock.EXPECT().CategoryList(gomock.Any(), gomock.Any()).Return(nil, external.ErrSkuNotFound)

        err := label.StartGeneration(t.Context(), expectedSKU)

        require.ErrorIs(t, err, external.ErrSkuNotFound)
    })

    t.Run("или для категории SKU нет шаблона", func(t *testing.T) {
        ozonServiceMock.EXPECT().CategoryList(gomock.Any(), gomock.Any()).Return(expectedCategoryList, nil)
        labelRepository.EXPECT().LoadByCategories(t.Context(), expectedCategoryList).Return(nil, nil)

        err := label.StartGeneration(t.Context(), expectedSKU)

        require.ErrorIs(t, err, domain.ErrLabelTemplateForCategoryNotFound)
    })

    t.Run("14. Начинать генерацию этикетки по SKU", func(t *testing.T) {
        ozonServiceMock.EXPECT().CategoryList(t.Context(), expectedSKU).Return(expectedCategoryList, nil)
        labelRepository.EXPECT().LoadByCategories(t.Context(), []domain.Category{
            expectedCategory1,
            expectedCategory2WithoutType,
        }).Return(&domain.LabelTemplate{
            ID: expectedLabelTemplateID,
        }, nil)

        err := label.StartGeneration(t.Context(), expectedSKU)

        require.NoError(t, err)

        t.Run("или такая генерация уже была запущена", func(t *testing.T) {
            err := label.StartGeneration(t.Context(), expectedSKU)

            require.ErrorContains(t, err, "генерация этикетки с таким идентификатором уже существует")
        })
    })
}
