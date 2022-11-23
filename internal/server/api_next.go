package server

import "net/http"

func (s *Server) next(w http.ResponseWriter, r *http.Request) {

	buf, err := s.cfg.Picker.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			s.log.Error(writeErr.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(buf)
	if err != nil {
		s.log.Error(err.Error())
	}
}
