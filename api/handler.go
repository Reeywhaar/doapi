package api

import (
	"doctl/internal"
	"encoding/json"
	"net/http"
)

type CreateDNSRequest struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Data   string `json:"data"`
	TTL    int    `json:"ttl"`
}

type DeleteDNSRequest struct {
	Domain string `json:"domain"`
	ID     int    `json:"id"`
}

type Handler struct {
	client *internal.Client
}

func NewHandler(client *internal.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /dns", h.handleList)
	mux.HandleFunc("POST /dns", h.handleCreate)
	mux.HandleFunc("DELETE /dns", h.handleDelete)
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	result, err := h.client.ListDomains()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateDNSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Domain == "" || req.Type == "" || req.Name == "" || req.Data == "" {
		http.Error(w, "domain, type, name, and data are required", http.StatusBadRequest)
		return
	}
	data, status, err := h.client.CreateRecord(req.Domain, req.Type, req.Name, req.Data, req.TTL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var req DeleteDNSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Domain == "" || req.ID == 0 {
		http.Error(w, "domain and id are required", http.StatusBadRequest)
		return
	}
	status, err := h.client.DeleteRecord(req.Domain, req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.WriteHeader(status)
}
