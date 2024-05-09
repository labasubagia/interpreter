package object

import "testing"

func TestHash(t *testing.T) {
	key1 := &String{Value: "key"}
	key2 := &String{Value: "key"}
	value1 := &String{Value: "value"}

	pairs := map[HashKey]Object{}
	pairs[key1.HashKey()] = value1

	if pairs[key2.HashKey()] != pairs[key1.HashKey()] {
		t.Errorf("same key should have same value")
	}
}

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff2.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
