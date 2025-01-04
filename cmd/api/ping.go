package main

import (
	"fmt"
    "net/http"
    "strconv" 

    "github.com/julienschmidt/httprouter"
)

// Declare a handler which writes a plain-text response with information about the 
// application status, operating environment and version.
func (app *application) pingMetricHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    // Otherwise, interpolate the service ID in a placeholder response.
    fmt.Fprintf(w, "show the details of service ID %d\n", id)
}