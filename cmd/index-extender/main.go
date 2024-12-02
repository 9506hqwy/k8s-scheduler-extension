package main

import (
	"encoding/json"
	"net/http"
	"os"

	v1 "k8s.io/kube-scheduler/extender/v1"

	"github.com/9506hqwy/k8s-scheduler-extension/pkg/indexscheduling"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var args v1.ExtenderArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if args.Pod == nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ret := indexscheduling.Filter(&args)

	if body, err := json.Marshal(ret); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		w.Write(body)
	}
}

func main() {
	http.HandleFunc("/api/scheduler/filter", handler)

	if err := http.ListenAndServe(":10261", nil); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
