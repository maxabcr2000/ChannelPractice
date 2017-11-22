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

	actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
	if ok {
		http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
		return
	}

	fmt.Println("main.register(): action.ID=", action.ID)

	actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
	if ok {
		fmt.Fprint(w, actionResult.Description)
	}
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
	isTimeout := pipeline.Delete(action)
	if isTimeout {
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		return
	}

	actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
	if ok {
		http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
		return
	}

	fmt.Println("main.delete(): action.ID=", action.ID)

	actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
	if ok {
		fmt.Fprint(w, actionResult.Description)
	}
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
	isTimeout := pipeline.Read(action)
	if isTimeout {
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		return
	}

	actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
	if ok {
		http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
		return
	}

	actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
	if ok {
		char, err := actionResult.GetDataAsCharacter()
		if err != nil {
			http.Error(w, fmt.Sprintf("InternalServerError: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(char)
		if err != nil {
			http.Error(w, fmt.Sprintf("InternalServerError: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(b))
	}
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
	isTimeout := pipeline.Update(action)
	if isTimeout {
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		return
	}

	actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
	if ok {
		http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
		return
	}

	actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
	if ok {
		fmt.Fprint(w, actionResult.Description)
	}
}
