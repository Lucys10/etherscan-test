package handlers

import (
	"apietherscan/internal/store"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type Handlers struct {
	S store.Store
}

func (h *Handlers) Register(r *chi.Mux) {
	r.Get("/api", h.Api)
}

func (h *Handlers) Api(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	res := q.Get("block")

	if res == "" {
		fmt.Println("Enter block number")
	}

	block, err := strconv.Atoi(res)
	if err != nil {
		log.Println(err)
	}

	listTransInfo, err := h.S.GetTransInfo(block)
	if err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(listTransInfo)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
