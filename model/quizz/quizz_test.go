package quizz

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/landru29/api-go/helpers/mongo"
	"github.com/landru29/api-go/helpers/json_load"
)

func init() {
	enreg := []Model {}
	mongo.Connect("127.0.0.1", "27017", "", "", "test")
	json_load.LoadJson("./fixtures.json", &enreg)
	index := 1
	for ; index < len(enreg); {
		enreg[index].Save()
	}
}

func TestInsert(t *testing.T) {
	a:=10
	assert.True(t, a == 10, "")
	/*a, b := SplitFloatForTimeUnix(12345678.887766)
	assert.True(t, a == 12345678, "a should be equals to 12345678")
	assert.True(t, b == 887766000, "b should be equals to 887766000")*/
}
