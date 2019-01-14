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
	awsTwo.Join(awsOne)
	awsTwo.Join(awsTwo)

	values := awsOne.GetAsAworSet("test1")
	setResult := values.Value()

	log.Printf("values %+v", setResult)

	if !setResult["testValOne"] {
		t.Fatalf("fail: no testValOne")
	}

	if !setResult["testValTwo"] {
		t.Fatalf("fail: no testValTwo")
	}

	values = awsTwo.GetAsAworSet("test1")
	setResult = values.Value()

	if !setResult["testValOne"] {
		t.Fatalf("fail: no testValOne")
	}

	if !setResult["testValTwo"] {
		t.Fatalf("fail: no testValTwo")
	}

}

func TestORMapNewCCouner(t *testing.T) {
	awsOne := NewWithStingKey("a", IntCounter)
	awsTwo := NewWithStingKey("b", IntCounter)

	awsOne.GetAsIntCounter("a").Inc(2)
	awsTwo.GetAsIntCounter("a").Inc(3)

	awsTwo.Join(awsOne)
	awsOne.Join(awsTwo)

	val := awsTwo.GetAsIntCounter("a").Value()
	if val != 5 {
		t.Fatalf("expect %d, but has %d", 5, val)
	}

	val = awsOne.GetAsIntCounter("a").Value()
	if val != 5 {
		t.Fatalf("expect %d, but has %d", 5, val)
	}
}
