package api

import (
	"bytes"
	"effective-mobile/store"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func executeRequest(req *http.Request, router *http.ServeMux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(rr, req)

	return rr
}

func TestUsersHandler_CreateUser(t *testing.T) {
	db := store.SetupTestDB(t)
	usersStore := store.NewPostgresUsersStore(db)
	usersHandler := NewUsersHandler(usersStore)

	router := &http.ServeMux{}
	router.HandleFunc("POST /users", usersHandler.HandleCreateUser)

	t.Run("Should not be able to create a user if body is missing", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", nil)

		rr := executeRequest(req, router)

		if rr.Code != http.StatusBadRequest {
			t.Fatal("Should have returned a BadRequest status code")
		}
	})

	t.Run("Should be able to create a user", func(t *testing.T) {
		user := store.User{
			FirstName: "John",
			LastName:  "Doe",
		}

		b, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(b))

		rr := executeRequest(req, router)

		if rr.Code != http.StatusCreated {
			t.Fatal("Should have returned a StatusCreated code")
		}
	})
}
