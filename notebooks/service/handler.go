package service

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/influxdata/influxdb/v2/context"
	"github.com/influxdata/influxdb/v2/kit/platform"
	"github.com/influxdata/influxdb/v2/kit/platform/errors"
	kithttp "github.com/influxdata/influxdb/v2/kit/transport/http"
	"go.uber.org/zap"
)

var (
	errBadOrg = &errors.Error{
		Code: errors.EInvalid,
		Msg:  "orgID not valid",
	}

	errBadId = &errors.Error{
		Code: errors.EInvalid,
		Msg:  "notebook id is invalid",
	}
)

// Handler is the handler for the notebook service
type Handler struct {
	chi.Router

	api *kithttp.API
	log *zap.Logger
}

func NewHandler(log *zap.Logger) *Handler {
	h := &Handler{
		log: log,
		api: kithttp.NewAPI(kithttp.WithLog(log)),
	}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.handleGetNotebooks)
		r.Post("/", h.handleCreateNotebook)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.handleGetNotebook)
			r.Delete("/", h.handleDeleteNotebook)
			r.Put("/", h.handleUpdateNotebook)
			r.Patch("/", h.handleUpdateNotebook)
		})
	})

	h.Router = r

	return h
}

// get a list of all notebooks for an org
func (h *Handler) handleGetNotebooks(w http.ResponseWriter, r *http.Request) {
	orgID := r.Context().Value(ContextKeyOrgID).(string)
	id, err := platform.IDFromString(orgID)

	if err != nil {
		h.api.Err(w, r, errBadOrg)
		return
	}

	// Demo data - respond with a list of notebooks for the requesting org
	d := map[string][]Notebook{}
	d["flows"] = demoNotebooks(3, *id)

	h.api.Respond(w, r, http.StatusOK, d)
}

// create a single notebook
func (h *Handler) handleCreateNotebook(w http.ResponseWriter, r *http.Request) {
	// Demo data - just return the body from the request with a generated ID
	b := Notebook{}
	if err := h.api.DecodeJSON(r.Body, &b); err != nil {
		h.api.Err(w, r, err)
		return
	}
	id, _ := platform.IDFromString(strconv.Itoa(1000000000000000 + 1)) // give it an ID from the getNotebooks list so that the UI doesn't break
	b.ID = *id

	h.api.Respond(w, r, http.StatusOK, b)
}

// get a single notebook
func (h *Handler) handleGetNotebook(w http.ResponseWriter, r *http.Request) {
	orgID, err := orgIDFromReq(r)
	if err != nil {
		h.api.Err(w, r, err)
		return
	}

	notebookID, err := platform.IDFromString(chi.URLParam(r, "id"))
	if err != nil {
		h.api.Err(w, r, err)
		return
	}

	// Demo data - return a notebook for the request org and id
	d := demoNotebook(*orgID, *notebookID)

	h.api.Respond(w, r, http.StatusOK, d)
}

// update a single notebook
func (h *Handler) handleUpdateNotebook(w http.ResponseWriter, r *http.Request) {
	id, err := platform.IDFromString(chi.URLParam(r, "id"))
	if err != nil {
		h.api.Err(w, r, errBadId)
		return
	}

	// Demo data - just return the body from the request with the id
	b := Notebook{}
	if err := h.api.DecodeJSON(r.Body, &b); err != nil {
		h.api.Err(w, r, err)
		return
	}
	b.ID = *id

	h.api.Respond(w, r, http.StatusOK, b)
}

// delete a single notebook
func (h *Handler) handleDeleteNotebook(w http.ResponseWriter, r *http.Request) {
	// for now, just respond with 200 unless there is a problem with the notebook ID
	if _, err := platform.IDFromString(chi.URLParam(r, "id")); err != nil {
		h.api.Err(w, r, errBadId)
		return
	}

	h.api.Respond(w, r, http.StatusOK, nil)
}

// this is a placeholder for more complex authorization behavior to come.
// for now, it just returns the first orgID from the user's permission set.
func orgIDFromReq(r *http.Request) (*platform.ID, error) {
	ctx := r.Context()
	a, err := context.GetAuthorizer(ctx)
	if err != nil {
		return nil, err
	}

	p, err := a.PermissionSet()
	if err != nil {
		return nil, err
	}

	for _, pp := range p {
		if pp.Resource.OrgID != nil {
			return pp.Resource.OrgID, nil
		}
	}

	return nil, errBadOrg
}
