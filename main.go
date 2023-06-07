package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

func main() {
	fmt.Println("Heynow!")

	httpin.UseGochiURLParam("path", chi.URLParam)

	r := chi.NewRouter()
	r.Get("/{name}", WithErrorHandler(handleIndex))

	http.ListenAndServe(":8080", r)
}

// more-or-less the same custom error handler logic
type handlerWithError func(w http.ResponseWriter, r *http.Request) error
func WithErrorHandler(h handlerWithError) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		var invalidFieldError *httpin.InvalidFieldError
		if errors.As(err, &invalidFieldError) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(invalidFieldError)
			return
		}

		if err != nil {
			fmt.Printf("error: %+v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ruh-roh"))
			return
		}
	}
}

type PaginationParams struct {
	Page int `in:"query=page;default=0"`
	PageSize int `in:"query=page_size;default=20"`
}

type IndexInput struct {
	PaginationParams

	// other fields in query params
	ShouldError bool `in:"query=should_error;required"`

	// URL params
	Name string `in:"path=name"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	var input IndexInput
	if err := httpin.Decode(r, &input); err != nil {
		return err
	}

	fmt.Printf("input struct: %+v\n", input)

	if input.ShouldError {
		return errors.New("erroring")
	}

	fmt.Fprint(w, "all good")

	return nil
}
