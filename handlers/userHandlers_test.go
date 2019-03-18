package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"restful/user"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestBodyToUser(t *testing.T) {
	// define test suite
	valid := &user.User{
		ID:   bson.NewObjectId(),
		Name: "Ryan",
		Role: "Engineer",
	}
	valid2 := &user.User{
		ID:   valid.ID,
		Name: "Ryan",
		Role: "tester",
	}
	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error marshalling a valid user: %s", err)
		t.FailNow()
	}
	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "empty body request",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "empty user request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		},
		{
			txt: "empty user request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id": 123}`)),
			},
			u:   &user.User{},
			err: true,
		},
		{
			txt: "valid user request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			err: false,
			exp: valid,
		},
		{
			txt: "valid partial request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"role": "tester", "age": 30}`)),
			},
			u:   valid,
			err: false,
			exp: valid2,
		},
	}

	for _, tc := range ts {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expect get error, get none")
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
