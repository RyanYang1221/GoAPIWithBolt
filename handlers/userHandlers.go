package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"restful/user"

	"github.com/asdine/storm"

	"gopkg.in/mgo.v2/bson"
)

func bodyToUser(r *http.Request, u *user.User) error {
	if r == nil {
		return errors.New("Request is empty")
	}
	if r.Body == nil {
		return errors.New("Body is empty")
	}
	if u == nil {
		return errors.New("User is required")
	}
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bd, u)

}
func usersGetAll(w http.ResponseWriter, r *http.Request) {
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"users": users})
}

func userPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	u.ID = bson.NewObjectId()
	if err := u.Save(); err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)

}

func userGetOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}
	if r.Method == http.MethodHead {
		postBodyResponse(w, http.StatusOK, jsonResponse{})
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})
}

func userPutOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	u.ID = id
	if err := u.Save(); err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}

	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})

}

func userPatchOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}

	// Will only update visible fields an keep others the original values
	err = bodyToUser(r, u)

	if err := u.Save(); err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}

	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})

}

func userDeleteOne(w http.ResponseWriter, _ *http.Request, id bson.ObjectId) {
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
		} else {
			postError(w, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

}
