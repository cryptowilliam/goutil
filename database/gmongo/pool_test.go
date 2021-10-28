package gmongo

import "testing"

func TestDialPool(t *testing.T) {
	invalidMongoUrl := "mongodb://192.168.1.2:3456"
	validMongoUrl := "mongodb://192.168.9.11:27717"

	_, err := DialPool(invalidMongoUrl, 20)
	if err == nil {
		t.Errorf("Dial invalid mongodb address should returns error")
		return
	}

	_, err = DialPool(validMongoUrl, 20)
	if err != nil {
		t.Errorf("Dial valid mongodb address should returns OK")
		return
	}
}
