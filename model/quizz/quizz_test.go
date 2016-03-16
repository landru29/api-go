package quizz

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/landru29/api-go/helpers/mongo"

)



func init() {
	mongo.Connect("localhost", "27017", "", "","test")
}

func TestSplitFloatForTimeUnix(t *testing.T) {
	/*a, b := SplitFloatForTimeUnix(12345678.887766)
	assert.True(t, a == 12345678, "a should be equals to 12345678")
	assert.True(t, b == 887766000, "b should be equals to 887766000")*/
}