package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	gid := "979445890064470036"
	uid := "136856906139566081"

	url := fmt.Sprintf("/%s/%s", gid, uid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(GetUser)

	h.ServeHTTP(rr, req)
}
