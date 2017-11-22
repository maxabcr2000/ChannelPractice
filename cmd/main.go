package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maxabcr2000/ChannelPractice/character"
)

var (
	pipeline = character.NewCharPipeline()
)

func main() {
	pipeline.Start()

	http.HandleFunc("/register", register)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe failed: ", err)
	}
}

func register(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "We only support application/json format in POST.", http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var char *character.Character
	err = json.Unmarshal(body, &char)
	if err != nil || char.ID <= 0 || char.Name == "" {
		http.Error(w, fmt.Sprintf("Bad Request: %s", string(body)), http.StatusBadRequest)
		return
	}

	action := character.NewAction("Register new character", *char)
	isTimeout := pipeline.Register(action)

	if isTimeout {
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		return
	}

	if pipeline.ErrorStage.CheckFailedAction(action.ID) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Println("main.register(): action.ID=", action.ID)

	actionDesc := pipeline.SinkStage.CheckCompletedAction(action.ID)
	fmt.Println("main.register(): actionDesc=", actionDesc)
	fmt.Fprint(w, actionDesc)
}

func delete(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "We only support application/x-www-form-urlencoded format in DELETE.", http.StatusUnsupportedMediaType)
		return
	}

	id := req.FormValue("id")

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	action := character.NewAction("Delete character", id)
	pipeline.Delete(action)

	//#Todo: Waiting for sinkStage , errorStage or timeout after 5 secs
}

func read(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "We only support application/x-www-form-urlencoded format in GET.", http.StatusUnsupportedMediaType)
		return
	}

	id := req.FormValue("id")

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	action := character.NewAction("Query character", id)
	pipeline.Read(action)

	//#Todo: Waiting for sinkStage , errorStage or timeout after 5 secs

	// b, err := json.Marshal(char)
	// if err != nil {
	// 	log.Println("Serialize Error")
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Fprint(w, string(b))
}

func update(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "We only support application/json format in PUT.", http.StatusUnsupportedMediaType)
		return
	}

	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var char *character.Character
	err = json.Unmarshal(body, &char)
	if err != nil || char.ID == 0 || char.Name == "" {
		http.Error(w, fmt.Sprintf("Bad Request: %s", string(body)), http.StatusBadRequest)
		return
	}

	action := character.NewAction("Update character", *char)
	pipeline.Update(action)
	//#Todo: Waiting for sinkStage , errorStage or timeout after 5 secs
}
