package http

import (
	"csm/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (h *Http) getContactsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling get all contacts")

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	limit := 20
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Cannot convert string to int when getting limit for contacts. Error %v. Limit passed: %s", err, limitStr)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		limit = l
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Printf("Cannot convert string to int when getting page for contacts. Error %v. Page passed: %s", err, pageStr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if limit <= 0 || page <= 0 {
		log.Printf("Limit or page passed as url query param is not valid. Error %v. Limit: %d. Page: %d", err, limit, page)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	contacts, hasMore, err := h.db.GetContacts(page, limit)
	if err != nil {
		log.Printf("Failed to get contacts. Error %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Contacts []database.Contact `json:"contacts"`
		HasMore  bool               `json:"has_more"`
	}{
		Contacts: contacts,
		HasMore:  hasMore,
	}

	encodedResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal when getting contacts. Error %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedResp)
}

func (h *Http) getDeltaContactsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling get delta contacts")

	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")
	timestampStr := r.URL.Query().Get("timestamp")

	limit := 20
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Cannot convert string to int when getting limit for contacts. Error %v. Limit passed: %s", err, limitStr)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		limit = l
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Printf("Cannot convert string to int when getting page for contacts. Error %v. Page passed: %s", err, pageStr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	timestamp, err := strconv.Atoi(timestampStr)
	if err != nil {
		log.Printf("Cannot convert string to int when getting last sync time for contacts. Error %v. Timestamp passed: %s", err, timestampStr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if limit <= 0 || page <= 0 || timestamp <= 0 {
		log.Printf("Limit or page or timestamp passed as url query param is not valid. Error %v. Limit: %d. Page: %d. Timestamp: %d", err, limit, page, timestamp)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	updatedContacts, deletedContactIDs, hasMore, err := h.db.GetUpdatedOrDeletedContactsSinceLastSync(page, limit, timestamp)
	if err != nil {
		log.Printf("Failed to get updated or deleted contacts. Error %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		UpdatedContacts   []database.Contact `json:"updated"`
		DeletedContactIDs []string           `json:"deleted"`
		HasMore           bool               `json:"has_more"`
	}{
		UpdatedContacts:   updatedContacts,
		DeletedContactIDs: deletedContactIDs,
		HasMore:           hasMore,
	}

	encodedResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal when getting updated or deleted contacts. Error %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encodedResp)
}
