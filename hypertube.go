package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func yonp(predicate string) bool {
	fmt.Print(predicate + " [y/n] ")
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

// User handlers
func userHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Routes: /create/, /login/, /logout/, /update/, /delete/, /info/")
}

func usernameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"username":"beanboy","bio":"I'm a very beany boy!"}`)
}

func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Not implemented yet\n\nCreate `%s` with body: %s", vars["username"], r.Body)
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
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Not implemented yet\n\nDeleting `%s`", vars["username"])
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
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, `{"%s": "bullshit"}`, vars["videoID"])
}

func videoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, "`%s` deleted", vars["videoID"])
}

func videoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, "`%s` updated with body `%s`", vars["videoID"], r.Body)
}

func videoWatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, "watching `%s`", vars["videoID"])
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler)
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/api/user/", userHandler)
	router.HandleFunc("/api/user/create/{username}", userCreateHandler)
	router.HandleFunc("/api/user/login/{username}", userLoginHandler)
	router.HandleFunc("/api/user/logout/{username}", userLogoutHandler)
	router.HandleFunc("/api/user/update/{username}", userUpdateHandler)
	router.HandleFunc("/api/user/delete/{username}", userDeleteHandler)
	router.HandleFunc("/api/user/info/{username}", userInfoHandler)

	router.HandleFunc("/api/video/info/{videoID}", videoInfoHandler)
	router.HandleFunc("/api/video/delete/{videoID}", videoDeleteHandler)
	router.HandleFunc("/api/video/update/{videoID}", videoUpdateHandler)
	router.HandleFunc("/api/video/watch/{videoID}", videoWatchHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// set up database
	fmt.Println("Connecting to database..")
	err := DBInit()
	if err != nil {
		log.Fatal(err)
	}
	if DBHasUserTable() == false {
		if yonp("User table does not exist, create one?") {
			DBCreateUserTable()
		}
	}
	if DBHasVideoTable() == false {
		if yonp("Video table does not exist, create one?") {
			DBCreateVideoTable()
		}
	}

	fmt.Println("Listening on port 8800")
	log.Fatal(http.ListenAndServe(":8800", router))
}
