package health

import (
	"net/http"
	"sync"

	"github.com/qdm12/golibs/logging"
)

func newHandler(logger logging.Logger) *handler {
	return &handler{
		logger: logger,
	}
}

type handler struct {
	logger      logging.Logger
	healthErr   error
	healthErrMu sync.RWMutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet || (r.RequestURI != "" && r.RequestURI != "/") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err := h.getHealthErr()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) setHealthErr(err error) {
	h.healthErrMu.Lock()
	defer h.healthErrMu.Unlock()
	h.healthErr = err
}

func (h *handler) getHealthErr() (err error) {
	h.healthErrMu.RLock()
	defer h.healthErrMu.RUnlock()
	return h.healthErr
}
