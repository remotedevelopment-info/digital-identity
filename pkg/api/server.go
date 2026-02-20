package api

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/nicholasalexander/digital-identity/pkg/auth"
	"github.com/nicholasalexander/digital-identity/pkg/chain"
	"github.com/nicholasalexander/digital-identity/pkg/store"
)

type Server struct {
	store store.ChainStore
	mux   *http.ServeMux
}

func NewServer(s store.ChainStore) *Server {
	srv := &Server{store: s, mux: http.NewServeMux()}
	srv.routes()
	return srv
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealth)
	s.mux.HandleFunc("/chains", s.handleChains)
	s.mux.HandleFunc("/chains/", s.handleChainByOwner)
}

type createChainRequest struct {
	OwnerID       string `json:"owner_id"`
	RootPublicKey string `json:"root_public_key"`
}

type createChainResponse struct {
	Chain          *chain.IdentityChain `json:"chain"`
	RootPrivateKey string               `json:"root_private_key,omitempty"`
}

type appendEventRequest struct {
	Event            chain.IdentityEvent `json:"event"`
	SignerPrivateKey string              `json:"signer_private_key"`
	Auth             auth.AuthContext    `json:"auth"`
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleChains(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createChain(w, r)
	case http.MethodGet:
		s.listChains(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleChainByOwner(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/chains/")
	parts := strings.Split(path, "/")
	ownerID := strings.TrimSpace(parts[0])
	if ownerID == "" {
		writeError(w, http.StatusBadRequest, "owner_id is required")
		return
	}

	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.getChain(w, r, ownerID)
		return
	}

	if len(parts) == 2 && parts[1] == "events" {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.appendEvent(w, r, ownerID)
		return
	}

	if len(parts) == 2 && parts[1] == "verify" {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		s.verifyChain(w, r, ownerID)
		return
	}

	writeError(w, http.StatusNotFound, "not found")
}

func (s *Server) createChain(w http.ResponseWriter, r *http.Request) {
	var req createChainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.OwnerID == "" {
		writeError(w, http.StatusBadRequest, "owner_id is required")
		return
	}

	var pub ed25519.PublicKey
	var priv ed25519.PrivateKey
	var err error

	if req.RootPublicKey == "" {
		pub, priv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to generate key")
			return
		}
	} else {
		decoded, err := base64.StdEncoding.DecodeString(req.RootPublicKey)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid root_public_key")
			return
		}
		pub = ed25519.PublicKey(decoded)
	}

	c := chain.NewIdentityChain(req.OwnerID, pub)
	if err := s.store.Create(context.Background(), c); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	resp := createChainResponse{Chain: c}
	if len(priv) > 0 {
		resp.RootPrivateKey = base64.StdEncoding.EncodeToString(priv)
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) listChains(w http.ResponseWriter, _ *http.Request) {
	items, err := s.store.List(context.Background())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chains": items})
}

func (s *Server) getChain(w http.ResponseWriter, _ *http.Request, ownerID string) {
	c, err := s.store.Get(context.Background(), ownerID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "chain not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) appendEvent(w http.ResponseWriter, r *http.Request, ownerID string) {
	var req appendEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := auth.ValidateRootAction(req.Auth); err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	privBytes, err := base64.StdEncoding.DecodeString(req.SignerPrivateKey)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid signer_private_key")
		return
	}
	priv := ed25519.PrivateKey(privBytes)

	c, err := s.store.Get(context.Background(), ownerID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "chain not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.Event.Timestamp.IsZero() {
		req.Event.Timestamp = time.Now().UTC()
	}
	if req.Event.ActorID == "" {
		req.Event.ActorID = ownerID
	}

	if err := chain.AppendEvent(c, req.Event, priv); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.store.Update(context.Background(), c); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

func (s *Server) verifyChain(w http.ResponseWriter, _ *http.Request, ownerID string) {
	c, err := s.store.Get(context.Background(), ownerID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "chain not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := chain.Verify(c); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"valid": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"valid": true})
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
