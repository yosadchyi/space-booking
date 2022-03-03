package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type Handler struct {
	service         *Service
	launchpadsPat   *regexp.Regexp
	destinationsPat *regexp.Regexp
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:         service,
		launchpadsPat:   regexp.MustCompile("^/launch/?$"),
		destinationsPat: regexp.MustCompile("^/destination/?$"),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGET(w, r)
	default:
		writeString(w, http.StatusBadRequest, "bad request")
	}
}

func (h *Handler) handleGET(w http.ResponseWriter, r *http.Request) {
	println(r.RequestURI)
	switch {
	case h.launchpadsPat.MatchString(r.RequestURI):
		all, err := h.service.launchpadRepository.GetAllActive()
		if err != nil {
			internalServerError(w)
			return
		}
		writeJsonResponse(w, http.StatusOK, all)
		break
	case h.destinationsPat.MatchString(r.RequestURI):
		all, err := h.service.destinationRepository.GetAll()
		if err != nil {
			internalServerError(w)
			return
		}
		writeJsonResponse(w, http.StatusOK, all)
		break
	default:
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	writeString(w, http.StatusNotFound, "not found")
}

func internalServerError(w http.ResponseWriter) {
	writeString(w, http.StatusInternalServerError, "internal server error")
}

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("can't serialize response")
		return
	}
	writeResponse(w, status, bytes)
}

func writeString(w http.ResponseWriter, status int, data string) {
	writeResponse(w, status, []byte(data))
}

func writeResponse(w http.ResponseWriter, status int, data []byte) {
	_, err := w.Write([]byte(data))
	if err != nil {
		log.Println("can't write response")
	}
	w.WriteHeader(status)
}
