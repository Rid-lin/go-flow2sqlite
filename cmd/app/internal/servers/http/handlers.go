package serverHTTP

// import (
// 	"encoding/json"
// 	"net/http"
// 	"time"
// )

// // renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
// func renderJSON(w http.ResponseWriter, v interface{}) {
// 	js, err := json.Marshal(v)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	_, _ = w.Write(js)
// }

// func (s *Server) handleUpdateDevices() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		s.timerUpdatedevice.Stop()
// 		s.timerUpdatedevice.Reset(1 * time.Second)
// 	}
// }

// func (s *Server) handleGetDevices() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		arr, err := s.GetDevices()
// 		if err != nil {
// 			return
// 		}
// 		renderJSON(w, arr)
// 	}
// }
