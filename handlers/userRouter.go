package handlers

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// UsersRouter handles the request for users
func UsersRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			usersGetAll(w, r)
			return
		case http.MethodPost:
			userPostOne(w, r)
			return
		case http.MethodHead:
			usersGetAll(w, r)
			return
		case http.MethodOptions:
			postOptionsResponse(w, []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}, nil)
			return
		default:
			postError(w, http.StatusMethodNotAllowed)
			return
		}
	}

	path = strings.TrimPrefix(r.URL.Path, "/users/")
	if !bson.IsObjectIdHex(path) {
		postError(w, http.StatusNotFound)
		return
	}

	id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:
		userGetOne(w, r, id)
		return
	case http.MethodPut:
		userPutOne(w, r, id)
		return
	case http.MethodPatch:
		userPatchOne(w, r, id)
		return
	case http.MethodDelete:
		userDeleteOne(w, r, id)
		return
	case http.MethodHead:
		userGetOne(w, r, id)
		return
	case http.MethodOptions:
		postOptionsResponse(w, []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions}, nil)
		return
	default:
		postError(w, http.StatusMethodNotAllowed)
		return
	}

}
