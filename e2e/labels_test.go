package e2e_test

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabelLive(t *testing.T) {
    t.Parallel()
    bins := []string{
        "./bin/1-layered",
        //"./bin/2-ddd-event-sourcing",
    }
    for _, bin := range bins {
        t.Run(bin+": Создавать шаблон этикетки товара с UUID и Наименованием организации производителя", func(t *testing.T) {
            output, err := runBinary(bin, []string{
                "labels-create-template",
                "--id", testUUID,
                "--manufacturer-organization-name", testManufacturerOrganizationName,
            })
            output = strings.ReplaceAll(output, "\n", "")

            require.NoError(t, err, "output:"+output)
            assert.Equal(t, "1", output)
        })
        //        t.Run(bin+": Получать данные шаблона в JSON", func(t *testing.T) {
        //            output, err := runBinary(bin, []string{
        //                "labels-get-template",
        //                "--id", testUUID,
        //            })
        //
        //            require.NoError(t, err)
        //            assert.Equal(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturer_organization_name":"test manufacturer organization name"}}
        //`, output)
        //        })
    }
}
