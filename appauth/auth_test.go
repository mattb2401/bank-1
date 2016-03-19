package appauth

import "testing"

func TestProcessAppAuthParams(t *testing.T) {
	data := []string{}
	_, err := ProcessAppAuth(data)
	if err == nil {
		t.Errorf("TestProcessAppAuthParams does not pass. Looking for %v, got %v", "Too few parameters", nil)
	}

	data = []string{"", "1"}
	_, err = ProcessAppAuth(data)
	if err == nil {
		t.Errorf("TestProcessAppAuthParams does not pass. Looking for %v, got %v", "Too few parameters", nil)
	}

	data = []string{"", "", "2", ""}
	_, err = ProcessAppAuth(data)
	if err == nil {
		t.Errorf("TestProcessAppAuthParams does not pass. Looking for %v, got %v", "Too few parameters", nil)
	}

	data = []string{"", "", "3", ""}
	_, err = ProcessAppAuth(data)
	if err == nil {
		t.Errorf("TestProcessAppAuthParams does not pass. Looking for %v, got %v", "Too few parameters", nil)
	}
}
