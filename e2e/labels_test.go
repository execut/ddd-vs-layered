package e2e_test

import (
    "encoding/json"
    "fmt"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type MainStructure struct {
    Text  string      `json:"text,omitempty"`
    Array []TestArray `json:"test_array,omitempty"`
}

type TestArray struct {
    ArrayText string `json:"array_text,omitempty"`
}

func TestLabelLive(t *testing.T) {
    Test := MainStructure{
        Text: "test",
        Array: []TestArray{
            {ArrayText: "test1"},
            {ArrayText: "test2"},
        },
    }
    json, _ := json.Marshal(Test)
    fmt.Println(string(json))
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
