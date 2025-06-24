package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
	"github.com/lenarlenar/mygokeeper/internal/server/service"
)

type RecordHandler struct {
	Service *service.RecordService
}

func (h *RecordHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	var rec entity.Record
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "bad input", http.StatusBadRequest)
		return
	}
	rec.UserID = userID

	if err := h.Service.Save(r.Context(), &rec); err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RecordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	records, err := h.Service.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func (h *RecordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	idStr := strings.TrimPrefix(r.URL.Path, "/record/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid record id", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(r.Context(), userID, id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RecordHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int64)

	idStr := strings.TrimPrefix(r.URL.Path, "/record/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid record id", http.StatusBadRequest)
		return
	}

	var rec entity.Record
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	rec.ID = id
	rec.UserID = userID

	if err := h.Service.Update(r.Context(), &rec); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
