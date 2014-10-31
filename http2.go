package http2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s request took %s\n", name, elapsed)
}

func TimingHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer timeTrack(time.Now(), req.URL.Path)
		fn(w, req)
	}
}

func RecoverHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("%s error %s\n", req.URL.Path, err)
				http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
				return
			}
		}()
		fn(w, req)
	}
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

func JSONResponse(rw http.ResponseWriter, r Response) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, r)
}
