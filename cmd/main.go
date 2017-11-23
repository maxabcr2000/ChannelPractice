package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/maxabcr2000/ChannelPractice/character"
)

var (
	pipeline = character.NewCharPipeline()
)

func main() {
	pipeline.Start()

	http.HandleFunc("/register", registerSync)
	http.HandleFunc("/registerAsync", registerAsync)
	http.HandleFunc("/delete", deleteSync)
	http.HandleFunc("/deleteAsync", deleteAsync)
	http.HandleFunc("/read", readSync)
	http.HandleFunc("/update", updateSync)
	http.HandleFunc("/updateAsync", updateAsync)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe failed: ", err)
	}
}

func registerSync(w http.ResponseWriter, req *http.Request) {
	register(w, req, true)
}

func registerAsync(w http.ResponseWriter, req *http.Request) {
	register(w, req, false)
}

func register(w http.ResponseWriter, req *http.Request, sync bool) {
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

	if !sync {
		fmt.Fprint(w, "Character registration has committed. Please check back later.")
		return
	}

	for {
		actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
		if ok {
			http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
			return
		}

		fmt.Println("main.register(): action.ID=", action.ID)

		actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
		if ok {
			fmt.Fprint(w, actionResult.Description)
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func deleteSync(w http.ResponseWriter, req *http.Request) {
	delete(w, req, true)
}

func deleteAsync(w http.ResponseWriter, req *http.Request) {
	delete(w, req, false)
}

func delete(w http.ResponseWriter, req *http.Request, sync bool) {
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

	if !sync {
		fmt.Fprint(w, "Character delete has committed. Please check back later.")
		return
	}

	for {
		actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
		if ok {
			http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
			return
		}

		fmt.Println("main.delete(): action.ID=", action.ID)

		actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
		if ok {
			fmt.Fprint(w, actionResult.Description)
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func readSync(w http.ResponseWriter, req *http.Request) {
	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "We only support application/x-www-form-urlencoded format in GET.", http.StatusUnsupportedMediaType)
		return
	}

	charID := req.FormValue("id")

	if charID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	intVal, err := strconv.Atoi(charID)
	if err != nil {
		http.Error(w, fmt.Sprintf("StatusBadRequest: %s", err.Error()), http.StatusBadRequest)
		return
	}

	char, ok := pipeline.ReadSync(intVal)

	if !ok {
		http.Error(w, fmt.Sprintf("InternalServerError: Failed to query the requested character data."), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(char)
	if err != nil {
		http.Error(w, fmt.Sprintf("InternalServerError: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(b))
}

// func readAsync(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("listJSON Endpoint: ", req.RemoteAddr)

// 	if req.Method != "GET" {
// 		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// 		return
// 	}

// 	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
// 		http.Error(w, "We only support application/x-www-form-urlencoded format in GET.", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	id := req.FormValue("id")

// 	if id == "" {
// 		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 		return
// 	}

// 	action := character.NewAction("Query character", id)
// 	isTimeout := pipeline.Read(action)
// 	if isTimeout {
// 		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
// 		return
// 	}

// 	for {
// 		actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
// 		if ok {
// 			http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
// 			return
// 		}

// 		actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
// 		if ok {
// 			char, err := actionResult.GetDataAsCharacter()
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("InternalServerError: %s", err.Error()), http.StatusInternalServerError)
// 				return
// 			}

// 			b, err := json.Marshal(char)
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("InternalServerError: %s", err.Error()), http.StatusInternalServerError)
// 				return
// 			}

// 			fmt.Fprint(w, string(b))
// 			return
// 		}

// 		time.Sleep(50 * time.Millisecond)
// 	}
// }

func updateSync(w http.ResponseWriter, req *http.Request) {
	update(w, req, true)
}

func updateAsync(w http.ResponseWriter, req *http.Request) {
	update(w, req, false)
}

func update(w http.ResponseWriter, req *http.Request, sync bool) {
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

	if !sync {
		fmt.Fprint(w, "Character update has committed. Please check back later.")
		return
	}

	for {
		actionResult, ok := pipeline.ErrorStage.CheckFailedAction(action.ID)
		if ok {
			http.Error(w, fmt.Sprintf("InternalServerError: %s", actionResult.Description), http.StatusInternalServerError)
			return
		}

		actionResult, ok = pipeline.SinkStage.CheckCompletedAction(action.ID)
		if ok {
			fmt.Fprint(w, actionResult.Description)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}
