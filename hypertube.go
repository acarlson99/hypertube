package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func yonp(predicate string) bool {
	fmt.Print(predicate + " [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	if input == 'y' {
		fmt.Println("OK!")
		return true
	}
	return false
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<h1>Hello from index</h1>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "\"%s\" not found. OH NO.", r.URL)
}

func userUsageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Routes: /create/, /login/, /logout/, /update/, /delete/, /info/")
}

func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	var user UserData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	err = DBUserCreate(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Not implemented yet\n\nLogging in `%s`", vars["username"])
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Not implemented yet\n\nLogging out `%s` with token: `%s`", vars["username"], r.Body)
}

func userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Not implemented yet\n\nUpdating `%s`\nbody:%s", vars["username"], r.Body)
}

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	err := DBUserDelete(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	return
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := DBUserInfo(vars["username"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "user `%s` not found", vars["username"])
		return
	}

	info, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(info))
}

func videoInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	video, err := DBVideoInfo(vars["uuid"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "video `%s` not found", vars["uuid"])
		return
	}

	info, err := json.Marshal(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(info))
}

func videoCreateHandler(w http.ResponseWriter, r *http.Request) {
	var video VideoData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	err = json.Unmarshal(body, &video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	err = DBVideoCreate(video)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
}

func videoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	UUID := mux.Vars(r)["uuid"]
	err := DBVideoDelete(UUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	return
}

func videoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, "`%s` updated with body `%s`", vars["videoID"], r.Body)
}

func videoWatchHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("static/lemon-demon.mp4")
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	reader := bufio.NewReader(file)
	data := make([]byte, 64)
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		n, err := reader.Read(data)
		if n == 0 || err != nil {
			return
		}

		err = conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			return
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/api/user/", userUsageHandler)
	router.HandleFunc("/api/user/create/", userCreateHandler)
	router.HandleFunc("/api/user/login/{username}", userLoginHandler)
	router.HandleFunc("/api/user/logout/{username}", userLogoutHandler)
	router.HandleFunc("/api/user/update/{username}", userUpdateHandler)
	router.HandleFunc("/api/user/delete/{username}", userDeleteHandler)
	router.HandleFunc("/api/user/info/{username}", userInfoHandler)
	router.HandleFunc("/api/video/info/{uuid}", videoInfoHandler)
	router.HandleFunc("/api/video/create/", videoCreateHandler)
	router.HandleFunc("/api/video/delete/{uuid}", videoDeleteHandler)
	router.HandleFunc("/api/video/update/{uuid}", videoUpdateHandler)
	router.HandleFunc("/ws/video/watch/", videoWatchHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	fmt.Println("Connecting to database..")
	err := DBInit()
	if err != nil {
		log.Fatal(err)
	}
	if DBHasUserTable() == false {
		if yonp("User table does not exist, create one?") {
			DBCreateUserTable()
		}
		if yonp("Add four test users?") {
			DBGenerateTrash()
		}
	}
	if DBHasVideoTable() == false {
		if yonp("Video table does not exist, create one?") {
			DBCreateVideoTable()
		}
	}

	fmt.Println("Listening on port :8800")
	log.Fatal(http.ListenAndServe(":8800", router))
}
