package handlers

import (
	"apietherscan/internal/store"
	"apietherscan/pkg/logger"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handlers struct {
	S    store.Store
	Logs *logger.Log
}

func (h *Handlers) Register(r *chi.Mux) {
	r.Get("/api", h.Api)
}

func (h *Handlers) Api(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	res := q.Get("block")

	if res == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	block, err := strconv.Atoi(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"function": "api",
			"error":    err,
		}).Error("failed convert string to number")
	}

	listTransInfo, err := h.S.GetTransInfo(block)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"function": "api",
			"error":    err,
		}).Error("failed insert to mongodb")
	}

	data, err := json.Marshal(listTransInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"function": "api",
			"error":    err,
		}).Error("failed marshal to json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Logs.WithFields(logrus.Fields{
			"package":  "handlers",
			"function": "api",
			"error":    err,
		}).Error("failed send response")
	}
}
