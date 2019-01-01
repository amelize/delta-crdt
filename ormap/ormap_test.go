package ormap

import (
	"log"
	"testing"
)

func TestORMapNew(t *testing.T) {
	awsOne := NewWithAworsetStringKey("a")
	awsTwo := NewWithAworsetStringKey("b")

	awsOne.GetAsAworSet("test1").Add("testValOne")
	awsTwo.GetAsAworSet("test1").Add("testValTwo")

	awsOne.Join(awsTwo)

	values := awsOne.GetAsAworSet("test1")
	setResult := values.Value()

	log.Printf("values %+v", setResult)

	if !setResult["testValOne"] {
		t.Fatalf("fail: no testValOne")
	}

	if !setResult["testValTwo"] {
		t.Fatalf("fail: no testValTwo")
	}
}
