package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	fmt.Fprintf(w, "Not implemented yet\n\ninfo of `%s`\n%s", vars["username"], "{user info goes here}")
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

	fmt.Println("Starting server on port 8800")
	log.Fatal(http.ListenAndServe(":8800", router))
}
