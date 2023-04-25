package site_test

import (
	"bytes"
	"testing"
	"viewer/pkg/site"
)

func TestStatusPage_Works(t *testing.T) {

	param := site.StatusPageParam{
		Title: "Calculation status",
	}

    var b bytes.Buffer

	err := site.StatusPage(&b, param)
	if err != nil {
		t.Fatalf("no error should happen here: %+v", err)
	}

	if b.Len() == 0 {
		t.Errorf("length should be > 0")
	}
}
