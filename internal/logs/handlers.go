package logs

import (
	"net/http"
)

func HandleMessage(w http.ResponseWriter, r *http.Request, code int, msg string) {
	logEntry := Log.
		WithField("status", http.StatusText(code)).
		WithField("code", code).
		WithField("method", r.Method).
		WithField("path", r.URL.Path)

	if code != http.StatusOK {
		logEntry.Error(msg)
	} else {
		logEntry.Info(msg)
	}

	w.WriteHeader(code)

	_, err := w.Write([]byte(msg))
	if err != nil {
		logEntry.WithError(err).Error("writing response failed")
	}
}
