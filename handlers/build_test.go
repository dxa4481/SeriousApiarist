package handlers

//
// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )
//
// func TestBuild(t *testing.T) {
// 	req, err := http.NewRequest("POST", "/build/sevoma/SeriousApiarist", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if val, ok := r.Context().Value("app.req.id").(string); !ok {
// 			t.Errorf("app.req.id not in request context: got %q", val)
// 		}
// 	})
//
// 	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Build)
//
// 	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
// 	// directly and pass in our Request and ResponseRecorder.
// 	handler.ServeHTTP(rr, req)
//
// 	// Check the status code is what we expect.
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
//
// 	// Check the response body is what we expect.
// 	expected := `{"alive": true}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
