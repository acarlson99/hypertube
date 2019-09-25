package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestUserCreation(t *testing.T) {
	name := "test1"
	bio := "hehe im a bee"
	email := "Timothy@hotmail.com"

	// make user
	body := strings.NewReader(fmt.Sprintf(`{"Name":"%s", "Bio":"%s", "Email":"%s"}`,
		name, bio, email))
	req, err := http.NewRequest("POST", "http://localhost:8800/api/user/create/", body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error("Bad response creating user:", resp)
	}
	defer resp.Body.Close()

	// get user info
	resp, err = http.Get("http://localhost:8800/api/user/info/" + name)
	if err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error("Bad response deleting user:", resp)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	got := UserData{}
	want := UserData{name, bio, email}
	err = json.Unmarshal(respBody, &got)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	if got != want {
		t.Error("Incorrect response")
	}

	// delete user
	resp, err = http.Get("http://localhost:8800/api/user/delete/" + name)
	if err != nil {
		t.Error(err)
	} else if resp.StatusCode != 200 {
		t.Error("Bad response deleting user:", resp)
	}
	defer resp.Body.Close()
}
