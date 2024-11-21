package transport

import (
	"encoding/json"
	"log"
	"my-first-api/internal/todo"
	"net/http"
)

type TodoItem struct {
	Item string `json:"item"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(todoSrv *todo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		todoItems, err := todoSrv.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(todoItems)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	})

	mux.HandleFunc("POST /todo", func(writer http.ResponseWriter, r *http.Request) {
		var t TodoItem
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = todoSrv.Add(t.Item)
		if err != nil {
			log.Println(err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusCreated)
		return
	})

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		results, err := todoSrv.Search(query)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
