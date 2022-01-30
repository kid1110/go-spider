package exception

import (
	"encoding/json"
	"log"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func ErrWrapper(handler appHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic:%v", r)
				http.Error(
					writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		err := handler(writer, request)

		if err != nil {
			log.Printf("Error occurred handling request: %s", err.Error())
			if userErr, ok := err.(NewErrorException); ok {
				json.NewEncoder(writer).Encode(&userErr.res)
				return
			}

		}

	}
}
