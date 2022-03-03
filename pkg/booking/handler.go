package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type Handler struct {
	service        *Service
	launchpadPat   *regexp.Regexp
	destinationPat *regexp.Regexp
	bookingPat     *regexp.Regexp
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:        service,
		launchpadPat:   regexp.MustCompile("^/launchpad/?$"),
		destinationPat: regexp.MustCompile("^/destination/?$"),
		bookingPat:     regexp.MustCompile("^/booking/$"),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGET(w, r)
	case http.MethodPost:
		h.handlePOST(w, r)
	default:
		badRequest(w)
	}
}

func (h *Handler) handleGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s", r.RequestURI)
	switch {
	case h.launchpadPat.MatchString(r.RequestURI):
		all, err := h.service.launchpadRepository.GetAllActive()
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeJsonResponse(w, http.StatusOK, all)
	case h.destinationPat.MatchString(r.RequestURI):
		all, err := h.service.destinationRepository.GetAll()
		if err != nil {
			internalServerError(w, err)
			break
		}
		writeJsonResponse(w, http.StatusOK, all)
	default:
		notFound(w)
	}
}

func (h *Handler) handlePOST(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST %s", r.RequestURI)
	switch {
	case h.bookingPat.MatchString(r.RequestURI):
		request := Request{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			badRequest(w)
			break
		}
		response, err := h.service.AddBooking(request)
		if err != nil {
			internalServerError(w, err)
			break
		}
		if errorResponse, ok := response.(*ErrorResponse); ok {
			writeJsonResponse(w, http.StatusBadRequest, errorResponse)
		} else {
			writeJsonResponse(w, http.StatusOK, response)
		}
	default:
		notFound(w)
	}
}

func badRequest(w http.ResponseWriter) {
	writeString(w, http.StatusBadRequest, "bad request")
}

func notFound(w http.ResponseWriter) {
	writeString(w, http.StatusNotFound, "not found")
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Printf("internal server error: %s", err)
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
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		log.Println("can't write response")
	}
}
