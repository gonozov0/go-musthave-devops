package handlers

import "net/http"

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.Repo.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
