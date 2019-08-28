package main

//go:generate go get github.com/dgrijalva/jwt-go
//go:generate go get github.com/gorilla/handlers
//go:generate go get github.com/gorilla/mux
//go:generate go get github.com/jteeuwen/go-bindata/...
//go:generate go get github.com/elazarl/go-bindata-assetfs/...
//go:generate go get github.com/googollee/go-socket.io
//go:generate go get golang.org/x/sys/...
//go:generate go get github.com/fsnotify/fsnotify

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fsnotify/fsnotify"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func _check(err error) {
	if err != nil {
		log.Println(err)
	}
}

type oServer struct {
	Secret []byte
	Ip     string
	Port   uint
	router *mux.Router
	socket *socketio.Server
}

func (s *oServer) String() string {
	return fmt.Sprintf("Served on http://%v:%d", s.Ip, s.Port)
}

func makeServer(port uint) *oServer {
	oS := new(oServer)
	ip, _ := getMyInterfaceAddr()
	oS.Ip = ip.String()
	oS.Port = port
	oS.Secret = []byte("secret")
	return oS
}

func (s *oServer) start() *oServer {
	r := mux.NewRouter()
	r.HandleFunc("/api/files", authMidlwere(filesLs)).Methods("GET")
	r.HandleFunc("/api/images/{filename}", authMidlwere(getImgFile)).Methods("GET")
	r.HandleFunc("/api/get", GetTokenHandler).Methods("POST")
	r.HandleFunc("/api/upload", authMidlwere(uploadImgFile)).Methods("POST")
	r.HandleFunc("/api/remove", authMidlwere(removeImgFile)).Methods("POST")
	r.HandleFunc("/api/edit", authMidlwere(editImgFile)).Methods("POST")
	if s.socket != nil {
		r.Handle("/api/socket/", sockMidlware(s.socket))
	}
	r.PathPrefix("/").Handler(http.FileServer(assetFS()))
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), handlers.LoggingHandler(os.Stdout, r))
	_check(err)
	return s
}

func (s *oServer) Close() {
	s.socket.Close()
}

func (srv *oServer) createSocketIOserv() *oServer {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join("fs")
		return nil
	})
	server.OnEvent("/", "APIREQ", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("HALLO")
		fmt.Println("HALLOO!!:", last)
		return last
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
	go server.Serve()
	srv.socket = server
	return srv
}

var mySigningKey = []byte("secret")

func fileWatcher(s *socketio.Server) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add("./")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			s.BroadcastToRoom("fs", "CHANGEFS")
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
	//
}

func main() {
	mServer := makeServer(3000)
	defer mServer.Close()
	log.Println(mServer)
	mServer.createSocketIOserv()
	go fileWatcher(mServer.socket)
	mServer.start()
}

func sockMidlware(res http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Debug for socket
		// w.Header().Set("Access-Control-Allow-Origin", strings.Replace(r.Host, "3000", "3001", 1))
		w.Header().Set("Access-Control-Allow-Origin", r.Host)
		res.ServeHTTP(w, r)
	}
}

func authMidlwere(res http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cook, _ := r.Cookie("auth")
		if cook != nil {
			token, err := jwt.Parse(cook.Value, func(token *jwt.Token) (interface{}, error) {
				return mySigningKey, nil
			})
			if err != nil || !token.Valid {
				log.Println("Authentication failed " + err.Error())
				w.WriteHeader(http.StatusForbidden)
			} else {
				res.ServeHTTP(w, r)
			}
		} else {
			// http.Redirect(w, r, "/", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

type FileI struct {
	Name    string
	Size    int64
	ModTime int64
	IsDir   bool
}

func getMyInterfaceAddr() (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses[0], nil
}

type authParams struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

var editImgFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"type":"err",message:"Error in editImgFile"}`))
		return
	}
	var ar map[string]string
	json.Unmarshal([]byte(reqBody), &ar)
	os.Rename("./"+ar["oldname"], "./"+ar["newname"])
	w.Write([]byte(`{"type":"ok","message":"Renamed:` + ar["oldname"] + " to " + ar["newname"] + `"}`))
})

var removeImgFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"type":"err",message:"Error in removeImgFile"}`))
		return
	}
	var ar map[string]string
	json.Unmarshal([]byte(reqBody), &ar)

	os.Remove("./" + ar["image"])
	w.Write([]byte(`{"type":"ok","message":"Deleted:` + ar["image"] + `"}`))
})

var uploadImgFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		w.Write([]byte(`{"type":"err",message:"Error in uploadImgFile"}`))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.Write([]byte(`{"type":"ok","message":"Uploaded:"}`))
})

var getImgFile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	dat, err := ioutil.ReadFile("./" + filename)
	if err != nil {
		panic(err)
	}
	w.Write(dat)
})

var filesLs = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	files, err := ioutil.ReadDir("./")
	s := make([]FileI, 0)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		s = append(s, FileI{Name: f.Name(), Size: f.Size(), ModTime: f.ModTime().Unix(), IsDir: f.IsDir()})
	}
	if err != nil {
		log.Println(err)
	}
	payload, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	w.Write([]byte(payload))

})

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	var ap authParams
	json.Unmarshal([]byte(reqBody), &ap)

	if ap.Login == "admin" && ap.Pass == "123" {
		token := jwt.New(jwt.SigningMethodHS256)
		expire := time.Now().AddDate(0, 0, 1)
		if err != nil {
			log.Fatal(err)
		}

		token.Claims = jwt.MapClaims{
			"admin": true,
			"name":  ap.Login,
			"exp":   expire.Unix(),
		}
		tokenString, _ := token.SignedString(mySigningKey)
		var cook = http.Cookie{Name: "auth", Value: tokenString, Expires: expire, Path: "/"}
		w.Header().Set("Authorization", tokenString)
		http.SetCookie(w, &cook)
		http.Redirect(w, r, "/", http.StatusAccepted)
	} else {
		http.Redirect(w, r, "/", http.StatusForbidden)
	}
})
