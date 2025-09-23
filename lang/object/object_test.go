package object

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("string with same content has different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("string with same content has different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("string with different content has same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	trueKey1 := &Boolean{Value: true}
	trueKey2 := &Boolean{Value: true}
	falseKey1 := &Boolean{Value: false}
	falseKey2 := &Boolean{Value: false}

	if trueKey1.HashKey() != trueKey2.HashKey() {
		t.Errorf("same booleans has different hash keys")
	}

	if falseKey1.HashKey() != falseKey2.HashKey() {
		t.Errorf("same booleans has different hash keys")
	}

	if trueKey1.HashKey() == falseKey1.HashKey() {
		t.Errorf("different booleans has same hash key")
	}
}

func TestIntegerHashKey(t *testing.T) {
	number1 := &Integer{Value: 1}
	number2 := &Integer{Value: 1}
	number3 := &Integer{Value: 11}
	number4 := &Integer{Value: 11}

	if number1.HashKey() != number2.HashKey() {
		t.Errorf("same integers has different hash keys")
	}

	if number3.HashKey() != number4.HashKey() {
		t.Errorf("same integers has different hash keys")
	}

	if number1.HashKey() == number3.HashKey() {
		t.Errorf("different integers has same hash key")
	}
}
