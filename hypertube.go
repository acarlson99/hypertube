package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	WSWriteSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
		goto handleErr
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		goto handleErr
	}
	err = DBUserCreate(user)
	if err != nil {
		goto handleErr
	}
	return

handleErr:
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

// request should look something like:
// https://github.com/login/oauth/authorize?client_id=CLIENT_ID&redirect_uri=http://localhost:8800/oauth/redirect
func oauthRedirectGithubHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err) // TODO: address error
	}
	code := r.FormValue("code")
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID_github, clientSecret_github, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		panic(err) // TODO: address error
	}
	req.Header.Set("accept", "application/json")

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		panic(err) // TODO: address error
	}
	defer res.Body.Close()

	var t tokenJSON
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		panic(err)
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
	io.WriteString(w, t.AccessToken)
	// TODO: do something with AccessToken
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
		fmt.Fprint(w, err)
	}
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
		fmt.Fprint(w, err)
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
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(info))
}

func videoCreateHandler(w http.ResponseWriter, r *http.Request) {
	var video VideoData
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		goto handleErr
	}
	err = json.Unmarshal(body, &video)
	if err != nil {
		goto handleErr
	}
	err = DBVideoCreate(video)
	if err != nil {
		goto handleErr
	}
	return

handleErr:
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func videoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	UUID := mux.Vars(r)["uuid"]
	err := DBVideoDelete(UUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
}

func videoUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, `Not implemented yet`)
	fmt.Fprintf(w, "`%s` updated with body `%s`", vars["videoID"], r.Body)
}

func videoDownloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	index, err := strconv.Atoi(mux.Vars(r)["index"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	}
	_, err = TClientAdd("magnet:?xt=urn:btih:"+string(hash), index)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err)
	}
}

func videoWatchHandlerTest(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./static/bunny.webm")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Println(stat.Size())

	reader := bufio.NewReader(file)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	bytes, _ := ioutil.ReadAll(reader)
	chunks := splitBytes(bytes, WSWriteSize)
	for _, chunk := range chunks {
		err = conn.WriteMessage(websocket.BinaryMessage, chunk)
		if err != nil {
			return
		}
	}
}

// split `buf` into byte slices of at most `size` size
func splitBytes(buf []byte, size int) [][]byte {
	chunks := [][]byte{}
	for len(buf) >= size {
		var chunk []byte
		chunk, buf = buf[:size], buf[size:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

// NOTE: this is horribly broken.  TeeReader does not fill buf with info
// need way to read from reader without modifying reader
func videoWatchHandlerHash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	reader, ok := openTorrents[hash]
	if !ok {
		fmt.Fprint(w, "Video not found")
		return
	}
	var buf bytes.Buffer
	io.TeeReader(reader, &buf)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	chunks := splitBytes(buf.Bytes(), WSWriteSize)
	for _, chunk := range chunks {
		err = conn.WriteMessage(websocket.BinaryMessage, chunk)
		if err != nil {
			return
		}
	}
}

func createRoutes(router *mux.Router) {
	router.HandleFunc("/", indexHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/oauth/redirect/github", oauthRedirectGithubHandler)

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
	router.HandleFunc("/api/video/download/{hash}/{index}", videoDownloadHandler)
	//router.HandleFunc("/api/video/watch/local/{uuid}", videoWatchLocalHandler)
	router.HandleFunc("/ws/video/watch/", videoWatchHandlerTest)
	router.HandleFunc("/ws/video/watch/{hash}", videoWatchHandlerHash)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
}

var clientID_github = ""
var clientSecret_github = ""

func main() {
	router := mux.NewRouter()
	createRoutes(router)

	clientID_github = os.Getenv("clientID_github")
	clientSecret_github = os.Getenv("clientSecret_github")

	DBInit("postgres")
	DBGenerateTablesPrompt()

	TClientStart()

	log.Print("Listening on port :8800")
	log.Fatal(http.ListenAndServe(":8800", router))
}
