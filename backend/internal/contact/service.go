package contact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"connectrpc.com/connect"

	contactv1 "github.com/nagutabby/sveltekit-blog/backend/gen/blog/contact/v1"
)

var emailRegexp = regexp.MustCompile(
	`^[a-zA-Z0-9_+-]+(\.[a-zA-Z0-9_+-]+)*@([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.)+[a-zA-Z]{2,}$`,
)

const mailtrapSendURL = "https://send.api.mailtrap.io/api/send"

// Config holds the secrets and addressing needed to send contact-form
// notifications through the Mailtrap API.
type Config struct {
	APIToken   string
	FromAddress string
	BCCAddress  string
}

// Service implements the blog.contact.v1.ContactService Connect RPC
// service. It ports the validation and Mailtrap call previously done in
// web/src/routes/contact/+page.server.ts.
type Service struct {
	config     Config
	httpClient *http.Client
	sendURL    string
}

func NewService(config Config) *Service {
	return &Service{
		config:     config,
		httpClient: http.DefaultClient,
		sendURL:    mailtrapSendURL,
	}
}

type emailPayload struct {
	From    emailAddress   `json:"from"`
	To      []emailAddress `json:"to"`
	BCC     []emailAddress `json:"bcc"`
	Subject string         `json:"subject"`
	HTML    string         `json:"html"`
}

type emailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *Service) SubmitContact(
	ctx context.Context,
	req *connect.Request[contactv1.SubmitContactRequest],
) (*connect.Response[contactv1.SubmitContactResponse], error) {
	msg := req.Msg

	if errs := validate(msg); len(errs) > 0 {
		return connect.NewResponse(&contactv1.SubmitContactResponse{Errors: errs}), nil
	}

	payload := emailPayload{
		From: emailAddress{Email: s.config.FromAddress, Name: "Hiroto Sasagawa"},
		To: []emailAddress{
			{Email: msg.GetEmail(), Name: msg.GetName()},
		},
		BCC: []emailAddress{
			{Email: s.config.BCCAddress, Name: "Hiroto Sasagawa"},
		},
		Subject: "お問い合わせを受け付けました",
		HTML: fmt.Sprintf(
			`<!DOCTYPE HTML><html><p>お問い合わせ内容は以下の通りです。</p><ul><li>氏名: %s</li><li>メールアドレス: %s</li><li>本文: %s</li></ul><p>返信まで数日かかる場合がございます。予めご了承ください。</p></html>`,
			msg.GetName(), msg.GetEmail(), msg.GetText(),
		),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.sendURL, bytes.NewReader(body))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Api-Token", s.config.APIToken)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("mailtrap API error: status %d", resp.StatusCode))
	}

	return connect.NewResponse(&contactv1.SubmitContactResponse{}), nil
}

func validate(msg *contactv1.SubmitContactRequest) map[string]string {
	errs := map[string]string{}

	if msg.GetImRobot() {
		errs["imRobot"] = "Botによるメッセージ送信はできません"
	}

	if msg.GetName() == "" {
		errs["name"] = "氏名は必須です"
	}

	switch {
	case msg.GetEmail() == "":
		errs["email"] = "メールアドレスは必須です"
	case !emailRegexp.MatchString(msg.GetEmail()):
		errs["email"] = "メールアドレスの形式が不適切です"
	}

	if msg.GetText() == "" {
		errs["text"] = "本文は必須です"
	}

	return errs
}
