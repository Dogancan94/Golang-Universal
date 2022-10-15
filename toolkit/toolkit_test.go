package toolkit

import "testing"

func TestToolkit_CreateRandomString(t *testing.T) {
	var testTools Tools
	randomString := testTools.CreateRandomString(10)
	if len(randomString) != 10 {
		t.Error("wrong length random string return")
	}
}
