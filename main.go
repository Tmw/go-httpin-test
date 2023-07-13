package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi"
)

type CustomBodyDecoderFunc func(src io.Reader, dst any) error

func (fn CustomBodyDecoderFunc) Decode(src io.Reader, dst any) error {
	return fn(src, dst)
}

func main() {
	httpin.UseGochiURLParam("path", chi.URLParam)

	// Replace body json decoder with one that returns an error
	// when it encounters unknown fields
	httpin.ReplaceBodyDecoder("json", CustomBodyDecoderFunc(func(src io.Reader, dst any) error {
		decoder := json.NewDecoder(src)
		decoder.DisallowUnknownFields()
		return decoder.Decode(dst)
	}))

	r := chi.NewRouter()
	r.Get("/{name}", WithErrorHandler(handleIndex))
	r.Post("/{name}", WithErrorHandler(handleIndex))

	panic(http.ListenAndServe("localhost:8080", r))
}

// more-or-less the same custom error handler logic
type handlerWithError func(w http.ResponseWriter, r *http.Request) error

func WithErrorHandler(h handlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	Page     int `in:"query=page;default=0"`
	PageSize int `in:"query=page_size;default=20"`
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type IndexInput struct {
	PaginationParams

	// URL params
	Name string `in:"path=name"`

	// Body
	Payload *User `in:"body=json"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	var input IndexInput
	if err := httpin.Decode(r, &input); err != nil {
		return err
	}

	fmt.Printf("input struct: %+v\n", input)
	fmt.Printf("payload: %+v\n", *input.Payload)
	fmt.Fprint(w, "OK")

	return nil
}
