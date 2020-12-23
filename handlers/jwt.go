package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetJwks(key jwk.Key) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		keys := make(map[string][]jwk.Key)
		keys["keys"] = []jwk.Key{key}
		jsonbuf, err := json.Marshal(keys)
		if err != nil {
			log.Printf("failed to generate JSON: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonbuf)
	}

}

func GetTokenForSub(privateKey jwk.Key) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		t := jwt.New()
		t.Set(jwt.SubjectKey, vars["sub"])
		t.Set(jwt.IssuedAtKey, time.Now().Unix())

		signed, err := jwt.Sign(t, jwa.RS256, privateKey)
		if err != nil {
			log.Printf("failed to sign token: %s", err)
			return
		}
		w.Write(signed)
	}
}

func GenerateTokenWithPayload(privateKey jwk.Key) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		var m map[string]interface{}
		err = json.Unmarshal(body, &m)
		if err != nil {
			panic(err)
		}

		t := jwt.New()
		for k, v := range m {
			t.Set(k, v)
		}
		t.Set(jwt.IssuedAtKey, time.Now().Unix())

		signed, err := jwt.Sign(t, jwa.RS256, privateKey)
		if err != nil {
			log.Printf("failed to sign token: %s", err)
			return
		}
		w.Write(signed)
	}
}
