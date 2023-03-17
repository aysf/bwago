package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aysf/bwago/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var testCase = []struct {
	name               string
	method             string
	path               string
	expectedStatusCode int
}{
	{"home", "GET", "/", http.StatusOK},
	{"about", "GET", "/about", http.StatusOK},
	{"contact", "GET", "/contact", http.StatusOK},
	{"generals-quarters", "GET", "/generals-quarters", http.StatusOK},
	{"majors-suite", "GET", "/majors-suite", http.StatusOK},
	{"search-availability", "GET", "/search-availability", http.StatusOK},

	// {"make-reservation", "GET", "/make-reservation", []postData{}, http.StatusOK},
	// {"reservation-summary", "GET", "/reservation-summary", []postData{}, http.StatusOK},
	// {"search-availability", "POST", "/search-availability", []postData{
	// 	{key: "start", value: "01-02-2021"},
	// 	{key: "end", value: "01-03-2021"},
	// }, http.StatusOK},
	// {"search-availability-json", "POST", "/search-availability-json", []postData{
	// 	{key: "start", value: "01-02-2021"},
	// 	{key: "end", value: "01-03-2021"},
	// }, http.StatusOK},
	// {"make-reservation", "POST", "/make-reservation", []postData{
	// 	{key: "first_name", value: "Tra"},
	// 	{key: "last_name", value: "Guy"},
	// 	{key: "phone", value: "555-555"},
	// 	{key: "email", value: "ansuf@gmail.com"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoute()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range testCase {
		var resp *http.Response
		var err error

		resp, err = ts.Client().Get(ts.URL + e.path)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected status code %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// request_recorder is basically simulating what we get from the request
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 99
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	// test 1
	reqBody := "start_data=2050-01-01"
	reqBody += "&" + "end_date=2050-01-02"
	reqBody += "&" + "first_name=John"
	reqBody += "&" + "last_name=Smith"
	reqBody += "&" + "email=John@smith.com"
	reqBody += "&" + "phone=123123123"
	reqBody += "&" + "room_id=1"

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test 2

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Remove(ctx, "reservation")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test 3: parse form failure

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test 4: invalid form
	reqBody = "start_data=2050-01-01"
	reqBody += "&" + "end_date=2050-01-02"
	reqBody += "&" + "first_name=John"
	reqBody += "&" + "last_name=Smith"
	reqBody += "&" + "email=Joh#^n@smith.com"
	reqBody += "&" + "phone=123123123"
	reqBody += "&" + "room_id=1"

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test 5: insert db failure
	reqBody = "start_data=2050-01-01"
	reqBody += "&" + "end_date=2050-01-02"
	reqBody += "&" + "first_name=John"
	reqBody += "&" + "last_name=Five"
	reqBody += "&" + "email=John@smith.com"
	reqBody += "&" + "phone=123123123"
	reqBody += "&" + "room_id=2"

	reservation = models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Major Suite",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test 5: insert db restriction failure
	reqBody = "start_data=2050-01-01"
	reqBody += "&" + "end_date=2050-01-02"
	reqBody += "&" + "first_name=John"
	reqBody += "&" + "last_name=Five"
	reqBody += "&" + "email=John@smith.com"
	reqBody += "&" + "phone=123123123"
	reqBody += "&" + "room_id=999"

	reservation = models.Reservation{
		RoomID: 999,
		Room: models.Room{
			ID:       999,
			RoomName: "Major Suite",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	testCases := []map[string]interface{}{
		{
			"name":    "The 1st case of AvailabilityJSON",
			"reqBody": "start=2040-01-01&end=2040-01-01&room_id=1",
		},
		{
			"name":    "The 2nd case of AvailabilityJSON",
			"reqBody": nil,
		},
		{
			"name":    "The 3rd case of AvailabilityJSON",
			"reqBody": "start=2040-01-01&end=2040-01-01&room_id=999",
		},
	}

	for _, tc := range testCases {
		log.Println("test:", tc["name"].(string))

		var req *http.Request
		var ctx context.Context
		var handler http.HandlerFunc
		var j jsonResponse

		rr := httptest.NewRecorder()

		if tc["reqBody"] != nil {
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(tc["reqBody"].(string)))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)

		}

		ctx = getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler = http.HandlerFunc(Repo.AvailabilityJson)

		handler.ServeHTTP(rr, req)
		err := json.Unmarshal(rr.Body.Bytes(), &j)
		if err != nil {
			t.Errorf("failed to parse json")
		}

	}

}

func TestRepository_ChooseRoom(t *testing.T) {

	testCases := []map[string]interface{}{
		{
			"name":       "The 1st case of ChooseRoom Handler",
			"urlRequest": "/choose-room/2",
			"expect":     http.StatusSeeOther,
		},
		{
			"name":       "The 2nd case of ChooseRoom Handler",
			"urlRequest": "/choose-room/a",
			"expect":     http.StatusBadRequest,
		},
		{
			"name":       "The 3rd case of ChooseRoom Handler",
			"urlRequest": "/choose-room/1",
			"expect":     http.StatusInternalServerError,
			"session":    true,
		},
	}

	for _, tc := range testCases {
		log.Println(tc["name"].(string))
		var req *http.Request
		var ctx context.Context
		var handler http.HandlerFunc
		rr := httptest.NewRecorder()

		req = httptest.NewRequest("GET", tc["urlRequest"].(string), nil)

		ctx = getCtx(req)
		req = req.WithContext(ctx)

		reservation := models.Reservation{
			RoomID: 2,
			Room: models.Room{
				ID:       2,
				RoomName: "Major Suite",
			},
		}
		session.Put(ctx, "reservation", reservation)

		if tc["session"] != nil && tc["session"].(bool) {
			session.Remove(ctx, "reservation")
		}

		handler = http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)
		if rr.Code != tc["expect"].(int) {
			t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, tc["expect"].(int))
		}

	}

}

func TestRepository_PostAvailability(t *testing.T) {

	testCases := []map[string]interface{}{
		{
			"name":    "PostAvailability Case 1: Error parse form",
			"bodyReq": nil,
			"expect":  http.StatusBadRequest,
		},
		{
			"name":    "PostAvailability Case 2: Status Ok",
			"bodyReq": "start=2040-01-01&end=2040-01-01",
			"expect":  http.StatusOK,
		},
		{
			"name":    "PostAvailability Case 3: Invalid Start Date",
			"bodyReq": "start=<invalid>&end=2040-01-01",
			"expect":  http.StatusBadRequest,
		},
		{
			"name":    "PostAvailability Case 4: Invalid End Date",
			"bodyReq": "start=2040-01-01&end=<invalid>",
			"expect":  http.StatusBadRequest,
		},
		{
			"name":    "PostAvailability Case 5: error searching data",
			"bodyReq": "start=2099-12-31&end=2040-01-01",
			"expect":  http.StatusInternalServerError,
		},
		{
			"name":    "PostAvailability Case 6: data not found",
			"bodyReq": "start=2000-12-31&end=2040-01-01",
			"expect":  http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		log.Println(tc["name"].(string))

		var req *http.Request
		rr := httptest.NewRecorder()

		if tc["bodyReq"] != nil {
			req = httptest.NewRequest("POST", "/make-reservation", strings.NewReader(tc["bodyReq"].(string)))
		} else {
			req, _ = http.NewRequest("POST", "/make-reservation", nil)
		}

		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)
		if rr.Code != tc["expect"].(int) {
			t.Errorf("PostAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, tc["expect"].(int))
		}

	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	testCases := []map[string]interface{}{
		{
			"name":    "PostAvailability Case 1: session is available",
			"session": true,
			"expect":  http.StatusOK,
		},
		{
			"name":    "PostAvailability Case 2: Status Ok",
			"session": false,
			"expect":  http.StatusTemporaryRedirect,
		},
	}

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Major Suite",
		},
	}

	for _, tc := range testCases {
		log.Println(tc["name"].(string))

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/reservation-summary", nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if tc["session"].(bool) {
			session.Put(ctx, "reservation", reservation)
		} else {
			session.Remove(ctx, "reservation")
		}

		handler := http.HandlerFunc(Repo.ReservationSummary)
		handler.ServeHTTP(rr, req)

		if rr.Code != tc["expect"].(int) {
			t.Errorf("ReservationSummary result code got %d, expect %d", rr.Code, tc["expect"].(int))
		}
	}

}

// func TestRepository_BookRoom(t *testing.T) {
// 	testCases := []map[string]interface{}{
// 		{
// 			"name":   "BookRoom Case 1: status Ok",
// 			"expect": http.StatusOK,
// 		},
// 		{
// 			"name":   "BookRoom Case 2: error parse room id",
// 			"expect": http.StatusBadRequest,
// 		},
// 		{
// 			"name":   "BookRoom Case 3: error parse start date",
// 			"expect": http.StatusBadRequest,
// 		},
// 		{
// 			"name":   "BookRoom Case 4: error parse end date",
// 			"expect": http.StatusBadRequest,
// 		},
// 		{
// 			"name":   "BookRoom Case 5: error get room by id",
// 			"expect": http.StatusBadRequest,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		log.Println(tc["name"].(string))

// 		rr :=

// 	}

// }

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("error get context:", err)
	}
	return ctx
}
