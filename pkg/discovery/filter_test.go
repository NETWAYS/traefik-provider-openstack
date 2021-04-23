package discovery

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterOptions_Matches(t *testing.T) {
	f := FilterOptions{}
	assert.True(t, f.Matches("test"))

	f.Excludes = []string{"te*"}
	assert.False(t, f.Matches("test"))

	f.Includes = []string{"serv*"}
	assert.False(t, f.Matches("test"))
	assert.True(t, f.Matches("server1"))
	assert.False(t, f.Matches("other"))

}
