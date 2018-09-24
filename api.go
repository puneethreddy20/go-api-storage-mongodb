package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Checks CSRF
func checkCSRF(w http.ResponseWriter, r *http.Request) (bool, error) {
	if r.Method != getMethod {
		referer := r.Referer()
		if len(referer) > 0 && len(r.Host) > 0 {
			log.Println(3, "ref =%s, host=%s", referer, r.Host)
			refererURL, err := url.Parse(referer)
			if err != nil {
				log.Println(err)
				return false, err
			}
			log.Println(3, "refHost =%s, host=%s", refererURL.Host, r.Host)
			if refererURL.Host != r.Host {
				log.Printf("CSRF detected.... rejecting with a 400")
				http.Error(w, "you are not authorized", http.StatusUnauthorized)
				err := errors.New("CSRF detected... rejecting")
				return false, err

			}
		}
	}
	return true, nil
}

func (state *RuntimeState) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the service! Please read README.md file to learn how to use the service. ")

}

//we need to have a JSON along with the request and json should have "name" key with value as user's name
func (state *RuntimeState) CreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	if r.Method != postMethod {
		log.Println("Invalid Request")
		http.Error(w, fmt.Sprintf("Invalid Request! Expected POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}
	var out map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&out)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	state.DBMutex.Lock()
	err = state.MongoDB.Insert(out)
	defer state.DBMutex.Unlock()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

}

func (state *RuntimeState) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()
	val, ok := query["username"]
	if !ok {
		log.Println("couldn't parse the URL")
		http.Error(w, "couldn't parse the URL", http.StatusBadRequest)
		return
	}
	state.DBMutex.Lock()
	err = state.MongoDB.Delete(val[0])
	defer state.DBMutex.Unlock()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

//example request: POST json to http://localhost:8080/createTags/?username=puneeth
func (state *RuntimeState) CreateandUpdateTagsforUser(w http.ResponseWriter, r *http.Request) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	if r.Method != postMethod {
		log.Println("Invalid Request")
		http.Error(w, fmt.Sprintf("Invalid Request! Expected POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	username, ok := query["username"]
	if !ok {
		log.Println("couldn't parse the URL")
		http.Error(w, "couldn't parse the URL", http.StatusBadRequest)
		return
	}

	var update map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	state.DBMutex.Lock()
	err = state.MongoDB.Upsert(username[0], update)
	defer state.DBMutex.Unlock()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

}

func (state *RuntimeState) DeleteTagsofaUser(w http.ResponseWriter, r *http.Request) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	//fmt.Println(r.URL)
	query := r.URL.Query()
	username, ok := query["username"]
	if !ok {
		log.Println("couldn't parse the URL")
		http.Error(w, "couldn't parse the URL", http.StatusBadRequest)
		return
	}
	tags, ok := query["deletetags"]
	if !ok {
		log.Println("couldn't parse the URL")
		http.Error(w, "couldn't parse the URL", http.StatusBadRequest)
		return
	}
	allTags := strings.Split(tags[0], ",")

	state.DBMutex.Lock()
	err = state.MongoDB.DeleteTags(username[0], allTags)
	defer state.DBMutex.Unlock()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
}

//example request: http://localhost:8080/listTags/?username=hello
func (state *RuntimeState) ListallTagsofaUser(w http.ResponseWriter, r *http.Request) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	query := r.URL.Query()
	val, ok := query["username"]
	if !ok {
		log.Println("couldn't parse the URL")
		http.Error(w, "couldn't parse the URL", http.StatusBadRequest)
		return
	}
	state.DBMutex.Lock()
	userinfo, err := state.MongoDB.FindByUsername(val[0])
	defer state.DBMutex.Unlock()
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(userinfo)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

}
