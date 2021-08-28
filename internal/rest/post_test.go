package rest_test

import (
	"art_space/internal/rest"
	"art_space/internal/rest/resttesting"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp"
)

type test struct {
	expected interface{}
	target   interface{}
}

func TestPost_CreateHandler(t *testing.T) {
	t.Parallel()

	type output struct {
		expectedStatus int
		expected       interface{}
		target         interface{}
	}

	tests := []struct {
		name   string
		setup  func(*resttesting.FakePostService)
		input  []byte
		output output
	}{
		{
			"OK: 201",
			func(s *resttesting.FakePostService) {
				s.CreatePostReturns(3, nil)
			},
			func() []byte {
				b, _ := json.Marshal(&rest.CreateUpdatePostRequest{
					Text:     "невероятный текст",
					AuthorId: 2,
				})

				return b
			}(),
			output{
				http.StatusCreated,
				rest.JsonResponse{
					Error:  nil,
					Result: 3,
				},
				&rest.JsonResponse{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New()
			svc := &resttesting.FakePostService{}
			tt.setup(svc)

			rest.NewPostHandler(svc).RegisterRoutes(app)

			resp, err := doRequest(app, httptest.NewRequest(fiber.MethodPost, "/posts", bytes.NewReader(tt.input)))
			if err != nil {
				t.Fatalf("expected no errro, got %v", err)
			}

			assertResponse(t, resp, test{tt.output.expected, tt.output.target})

			if tt.output.expectedStatus != resp.StatusCode {
				t.Fatalf("expected code %d, actual %d", tt.output.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func doRequest(app *fiber.App, req *http.Request) (*http.Response, error) {
	resp, err := app.Test(req, -1)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func assertResponse(t *testing.T, res *http.Response, test test) {
	t.Helper()

	if err := json.NewDecoder(res.Body).Decode(test.target); err != nil {
		t.Fatalf("couldn't decode %s", err)
	}
	defer res.Body.Close()

	if !cmp.Equal(test.expected, test.target) {
		t.Fatalf("expected results don't match: %s", cmp.Diff(test.expected, test.target))
	}
}
