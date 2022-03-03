package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

type Handler struct {
	service          *Service
	launchpadPat     *regexp.Regexp
	destinationPat   *regexp.Regexp
	bookingPat       *regexp.Regexp
	deleteBookingPat *regexp.Regexp
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:          service,
		launchpadPat:     regexp.MustCompile("^/launchpad/?$"),
		destinationPat:   regexp.MustCompile("^/destination/?$"),
		bookingPat:       regexp.MustCompile("^/booking/$"),
		deleteBookingPat: regexp.MustCompile("^/booking/([0-9a-f-]{36})$"),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGET(w, r)
	case http.MethodPost:
		h.handlePOST(w, r)
	case http.MethodDelete:
		h.handleDELETE(w, r)
	default:
		badRequest(w)
	}
}

func (h *Handler) handleGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s", r.RequestURI)
	switch {
	case h.launchpadPat.MatchString(r.RequestURI):
		all, err := h.service.GetAllLaunchpads()
		if err != nil {
			internalServerError(w, err)
			return
		}
		writeJsonResponse(w, http.StatusOK, all)
	case h.destinationPat.MatchString(r.RequestURI):
		all, err := h.service.GetAllDestinations()
		if err != nil {
			internalServerError(w, err)
			break
		}
		writeJsonResponse(w, http.StatusOK, all)
	case h.bookingPat.MatchString(r.RequestURI):
		all, err := h.service.GetAllBookings()
		if err != nil {
			internalServerError(w, err)
			break
		}
		writeJsonResponse(w, http.StatusOK, all)
	default:
		http.NotFound(w, r)
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
		http.NotFound(w, r)
	}
}

func (h *Handler) handleDELETE(w http.ResponseWriter, r *http.Request) {
	log.Printf("DELETE %s", r.RequestURI)
	switch {
	case h.deleteBookingPat.MatchString(r.RequestURI):
		str := h.deleteBookingPat.FindStringSubmatch(r.RequestURI)
		id := str[1]
		_, err := uuid.Parse(id)
		if err != nil {
			badRequest(w)
			break
		}
		err = h.service.DeleteBooking(id)
		if err != nil {
			internalServerError(w, err)
		}
		writeString(w, http.StatusNoContent, "")
	default:
		http.NotFound(w, r)
	}
}

func badRequest(w http.ResponseWriter) {
	writeString(w, http.StatusBadRequest, "bad request")
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
