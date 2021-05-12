package notebooks

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/influxdata/influxdb/v2/kit/feature"
	kithttp "github.com/influxdata/influxdb/v2/kit/transport/http"
	"github.com/influxdata/influxdb/v2/notebooks/service"
	"go.uber.org/zap"
)

const (
	prefixNotebooks = "/api/v2/notebooks"
)

// NotebookTransport is the handler for the notebook service
type NotebookTransport struct {
	chi.Router
}

func (t *NotebookTransport) Prefix() string {
	return prefixNotebooks
}

func NewNotebookTransport(log *zap.Logger) *NotebookTransport {
	t := &NotebookTransport{}

	r := chi.NewRouter()
	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		MwNotebookFlag(kithttp.NewAPI(kithttp.WithLog(log))), // temporary, remove when feature flag for notebooks is removed
		MwSetOrgID,
	)

	r.Mount("/", service.NewHandler(log))

	t.Router = r

	return t
}

// MwNotebookFlag is middleware for returning no content if the notebooks feature
// flag is set to false. remove this middleware when the feature flag is removed.
func MwNotebookFlag(a *kithttp.API) func(http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			flags := feature.FlagsFromContext(r.Context())

			if !flags["notebooks"].(bool) || !flags["notebooksApi"].(bool) {
				a.Respond(w, r, http.StatusNoContent, nil)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return mw
}

// MwSetOrgID sets the orgID as provided by a query parameter on the request context.
// this will only be used when getting a list of all notebooks.
func MwSetOrgID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		i := r.URL.Query().Get("orgID")
		ctx := context.WithValue(r.Context(), service.ContextKeyOrgID, i)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
