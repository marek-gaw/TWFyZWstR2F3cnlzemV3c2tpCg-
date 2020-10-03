package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// RESTy routes for "fetch" resource
	r.Route("/api/fetcher", func(r chi.Router) {
		r.Post("/", CreateFetch) // POST /api/fetcher
		r.Get("/", GetAllFetch)  // GET /api/fetcher

		r.Route("/{fetchID}", func(r chi.Router) {
			r.Use(FetchCtx)             // Load the *Fetch on the request context
			r.Get("/history", GetFetch) // GET /api/fetcher/123/history
			r.Delete("/", DeleteFetch)  // DELETE /api/fetcher/123
		})

	})

	http.ListenAndServe(":8080", r)
}

func GetAllFetch(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewFetchListResponse(fetchs)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// FetchCtx middleware is used to load a Fetch object from
// the URL parameters passed through as the request. In case
// the Fetch could not be found, we stop here and return a 404.
func FetchCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fetch *Fetch
		var err error

		if fetchID := chi.URLParam(r, "fetchID"); fetchID != "" {
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "fetch", fetch)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateFetch persists the posted Fetch and returns it
// back to the client as an acknowledgement.
func CreateFetch(w http.ResponseWriter, r *http.Request) {
	data := &FetchRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	fetch := data.Fetch
	dbNewFetch(fetch)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewFetchResponse(fetch))
}

// GetFetch returns the specific Fetch. You'll notice it just
// fetches the Fetch right off the context, as its understood that
// if we made it this far, the Fetch must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetFetch(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the fetch
	// context because this handler is a child of the FetchCtx
	// middleware. The worst case, the recoverer middleware will save us.
	fetch := r.Context().Value("fetch").(*Fetch)

	if err := render.Render(w, r, NewFetchResponse(fetch)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// UpdateFetch updates an existing Fetch in our persistent store.
func UpdateFetch(w http.ResponseWriter, r *http.Request) {
	fetch := r.Context().Value("fetch").(*Fetch)

	data := &FetchRequest{Fetch: fetch}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	fetch = data.Fetch
	dbUpdateFetch(fetch.ID, fetch)

	render.Render(w, r, NewFetchResponse(fetch))
}

// DeleteFetch removes an existing Fetch from our persistent store.
func DeleteFetch(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the fetch
	// context because this handler is a child of the FetchCtx
	// middleware. The worst case, the recoverer middleware will save us.
	fetch := r.Context().Value("fetch").(*Fetch)

	fetch, err = dbRemoveFetch(fetch.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewFetchResponse(fetch))
}

// This is entirely optional, but I wanted to demonstrate how you could easily
// add your own logic to the render.Respond method.
func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

//--
// Request and Response payloads for the REST api.
//
// The payloads embed the data model objects an
//
// In a real-world project, it would make sense to put these payloads
// in another file, or another sub-package.
//--

// FetchRequest is the request payload for Fetch data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type FetchRequest struct {
	*Fetch
}

func (a *FetchRequest) Bind(r *http.Request) error {
	// a.Fetch is nil if no Fetch fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Fetch == nil {
		return errors.New("missing required Fetch fields.")
	}

	// just a post-process after a decode..
	a.Fetch.URL = strings.ToLower(a.Fetch.URL) // as an example, we down-case
	return nil
}

// FetchResponse is the response payload for the Fetch data model.
// See NOTE above in FetchRequest as well.
//
// In the FetchResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type FetchResponse struct {
	*Fetch
}

func NewFetchResponse(fetch *Fetch) *FetchResponse {
	resp := &FetchResponse{Fetch: fetch}

	return resp
}

func (rd *FetchResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewFetchListResponse(fetchs []*Fetch) []render.Renderer {
	list := []render.Renderer{}
	for _, fetch := range fetchs {
		list = append(list, NewFetchResponse(fetch))
	}
	return list
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

//--
// Data model objects and persistence mocks:
//--

// Fetch data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Fetch struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Interval  int    `json:"interval"`
	Response  string `json:"response"`
	Duration  int    `json:"duration"`
	CreatedAt int    `json:"created_at"`
}

// Fetch fixture data
var fetchs = []*Fetch{}

func dbNewFetch(fetch *Fetch) (string, error) {
	fetch.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	fetchs = append(fetchs, fetch)
	return fetch.ID, nil
}

func dbGetFetch(id string) (*Fetch, error) {
	for _, f := range fetchs {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("fetch not found.")
}

func dbUpdateFetch(id string, fetch *Fetch) (*Fetch, error) {
	for i, f := range fetchs {
		if f.ID == id {
			fetchs[i] = fetch
			return fetch, nil
		}
	}
	return nil, errors.New("fetch not found.")
}

func dbRemoveFetch(id string) (*Fetch, error) {
	for i, a := range fetchs {
		if a.ID == id {
			fetchs = append((fetchs)[:i], (fetchs)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("fetch not found.")
}
