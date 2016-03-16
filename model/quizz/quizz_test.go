package quizz

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"github.com/landru29/api-go/helpers/mongo"

)



func init() {
	mongo.Connect("localhost", "27017", "", "","test")
}

func TestInsert(t *testing.T) {
	quizz := Model{
		explaination : "Explaination",
		image        : "imgUrl",
		level        : 0,
		published    : true,
		tags         : "tag1",
		text         : "Question",
		//choices      : []Choice{},
	}
	quizz.Save()
	/*a, b := SplitFloatForTimeUnix(12345678.887766)
	assert.True(t, a == 12345678, "a should be equals to 12345678")
	assert.True(t, b == 887766000, "b should be equals to 887766000")*/
}