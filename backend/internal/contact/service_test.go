package contact

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"

	contactv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/contact/v1"
)

func newTestService(t *testing.T, handler http.HandlerFunc) *Service {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	svc := NewService(Config{
		APIToken:    "mock_token",
		FromAddress: "admin@example.com",
		BCCAddress:  "bcc@example.com",
	})
	svc.sendURL = server.URL

	return svc
}

func submit(t *testing.T, svc *Service, req *contactv1.SubmitContactRequest) (*contactv1.SubmitContactResponse, error) {
	t.Helper()
	resp, err := svc.SubmitContact(context.Background(), connect.NewRequest(req))
	if resp == nil {
		return nil, err
	}
	return resp.Msg, err
}

func TestSubmitContactSuccess(t *testing.T) {
	var capturedBody map[string]any
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if got := r.Header.Get("Api-Token"); got != "mock_token" {
			t.Fatalf("Api-Token header = %q, want %q", got, "mock_token")
		}
		w.WriteHeader(http.StatusOK)
	})

	msg, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:  "氏名",
		Email: "test@example.com",
		Text:  "これはテストメッセージです。",
	})
	if err != nil {
		t.Fatalf("SubmitContact returned error: %v", err)
	}
	if len(msg.GetErrors()) != 0 {
		t.Fatalf("errors = %v, want empty", msg.GetErrors())
	}

	from := capturedBody["from"].(map[string]any)
	if from["email"] != "admin@example.com" || from["name"] != "Hiroto Sasagawa" {
		t.Fatalf("from = %v", from)
	}
	to := capturedBody["to"].([]any)[0].(map[string]any)
	if to["email"] != "test@example.com" || to["name"] != "氏名" {
		t.Fatalf("to = %v", to)
	}
	bcc := capturedBody["bcc"].([]any)[0].(map[string]any)
	if bcc["email"] != "bcc@example.com" {
		t.Fatalf("bcc = %v", bcc)
	}
	if capturedBody["subject"] != "お問い合わせを受け付けました" {
		t.Fatalf("subject = %v", capturedBody["subject"])
	}
}

func TestSubmitContactRobotCheckFails(t *testing.T) {
	called := false
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	msg, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:    "氏名",
		Email:   "test@example.com",
		Text:    "これはテストメッセージです。",
		ImRobot: true,
	})
	if err != nil {
		t.Fatalf("SubmitContact returned error: %v", err)
	}
	if got := msg.GetErrors()["imRobot"]; got != "Botによるメッセージ送信はできません" {
		t.Fatalf(`errors["imRobot"] = %q, want "Botによるメッセージ送信はできません"`, got)
	}
	if called {
		t.Fatal("Mailtrap should not be called when robot check fails")
	}
}

func TestSubmitContactMissingRequiredFields(t *testing.T) {
	called := false
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	msg, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:  "",
		Email: "test@example.com",
		Text:  "これはテストメッセージです。",
	})
	if err != nil {
		t.Fatalf("SubmitContact returned error: %v", err)
	}
	if got := msg.GetErrors()["name"]; got != "氏名は必須です" {
		t.Fatalf(`errors["name"] = %q, want "氏名は必須です"`, got)
	}
	if called {
		t.Fatal("Mailtrap should not be called on validation failure")
	}
}

func TestSubmitContactInvalidEmailFormat(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {})

	msg, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:  "氏名",
		Email: "invalid-email",
		Text:  "これはテストメッセージです。",
	})
	if err != nil {
		t.Fatalf("SubmitContact returned error: %v", err)
	}
	if got := msg.GetErrors()["email"]; got != "メールアドレスの形式が不適切です" {
		t.Fatalf(`errors["email"] = %q, want "メールアドレスの形式が不適切です"`, got)
	}
}

func TestSubmitContactMultipleValidationErrors(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {})

	msg, err := submit(t, svc, &contactv1.SubmitContactRequest{})
	if err != nil {
		t.Fatalf("SubmitContact returned error: %v", err)
	}
	errs := msg.GetErrors()
	want := map[string]string{
		"name":  "氏名は必須です",
		"email": "メールアドレスは必須です",
		"text":  "本文は必須です",
	}
	for field, wantMsg := range want {
		if got := errs[field]; got != wantMsg {
			t.Fatalf("errors[%q] = %q, want %q", field, got, wantMsg)
		}
	}
}

func TestSubmitContactMailAPIFailure(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:  "氏名",
		Email: "test@example.com",
		Text:  "これはテストメッセージです。",
	})
	if err == nil {
		t.Fatal("expected an error when Mailtrap returns a non-2xx status, got nil")
	}
	if connect.CodeOf(err) != connect.CodeUnavailable {
		t.Fatalf("error code = %v, want %v", connect.CodeOf(err), connect.CodeUnavailable)
	}
}

func TestSubmitContactNetworkError(t *testing.T) {
	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {})
	// Point at a URL nothing is listening on to force a network-level error.
	svc.sendURL = "http://127.0.0.1:0"

	_, err := submit(t, svc, &contactv1.SubmitContactRequest{
		Name:  "氏名",
		Email: "test@example.com",
		Text:  "これはテストメッセージです。",
	})
	if err == nil {
		t.Fatal("expected an error when the network call fails, got nil")
	}
}
