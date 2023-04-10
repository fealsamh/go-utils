package textutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakecasedFromCamelcased(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(SnakecasedFromCamelcased("NewYork"), "new_york")
	assert.Equal(SnakecasedFromCamelcased("SQLQuery"), "sql_query")
	assert.Equal(SnakecasedFromCamelcased("AnotherSQLQuery"), "another_sql_query")
	assert.Equal(SnakecasedFromCamelcased("OpenAPI"), "open_api")
	assert.Equal(SnakecasedFromCamelcased("NEXTSTEP"), "nextstep")
	assert.Equal(SnakecasedFromCamelcased("SomeAPIsXyz", "api"), "some_apis_xyz")
	assert.Equal(SnakecasedFromCamelcased("SomeIDsXyz", "id"), "some_ids_xyz")
}
