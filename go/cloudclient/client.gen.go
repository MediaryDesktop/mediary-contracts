package cloudclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	Admin     UserBodyRole = "admin"
	Moderator UserBodyRole = "moderator"
	User      UserBodyRole = "user"
)

func (e UserBodyRole) Valid() bool {
	switch e {
	case Admin:
		return true
	case Moderator:
		return true
	case User:
		return true
	default:
		return false
	}
}

type ClientRequirements struct {
	MinSupported *string `json:"min_supported,omitempty"`

	Recommended *string `json:"recommended,omitempty"`

	UpdateUrl *string `json:"update_url,omitempty"`
}

type ErrorDetail struct {
	Location *string `json:"location,omitempty"`

	Message *string `json:"message,omitempty"`

	Value interface{} `json:"value,omitempty"`
}

type ErrorModel struct {
	Schema *string `json:"$schema,omitempty"`

	Detail *string `json:"detail,omitempty"`

	Errors *[]ErrorDetail `json:"errors,omitempty"`

	Instance *string `json:"instance,omitempty"`

	Status *int64 `json:"status,omitempty"`

	Title *string `json:"title,omitempty"`

	Type *string `json:"type,omitempty"`
}

type HealthOutputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Built string `json:"built"`

	Client ClientRequirements `json:"client"`

	Commit string `json:"commit"`

	Status string `json:"status"`

	Version string `json:"version"`
}

type LoginInputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Email openapi_types.Email `json:"email"`

	Password string `json:"password"`
}

type OAuthCallbackInputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Code string `json:"code"`

	State string `json:"state"`
}

type OAuthProvidersOutputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Providers *[]string `json:"providers"`
}

type OAuthStartOutputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Url string `json:"url"`
}

type PingOutputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Module string `json:"module"`

	Ready bool `json:"ready"`
}

type PingOutputBody1 struct {
	Schema *string `json:"$schema,omitempty"`

	Module string `json:"module"`

	Ready bool `json:"ready"`
}

type PingOutputBody2 struct {
	Schema *string `json:"$schema,omitempty"`

	Module string `json:"module"`

	Ready bool `json:"ready"`
}

type PingOutputBody3 struct {
	Schema *string `json:"$schema,omitempty"`

	Module string `json:"module"`

	Ready bool `json:"ready"`
}

type RegisterInputBody struct {
	Schema *string `json:"$schema,omitempty"`

	DisplayName string `json:"display_name"`

	Email openapi_types.Email `json:"email"`

	Password string `json:"password"`
}

type RevokedOutputBody struct {
	Schema *string `json:"$schema,omitempty"`

	Revoked int64 `json:"revoked"`
}

type SessionBody struct {
	Schema *string `json:"$schema,omitempty"`

	ExpiresAt time.Time `json:"expires_at"`

	Token string `json:"token"`

	User UserBody `json:"user"`
}

type UserBody struct {
	Schema *string `json:"$schema,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	DisplayName string `json:"display_name"`

	Email openapi_types.Email `json:"email"`

	EmailVerified bool `json:"email_verified"`

	Id openapi_types.UUID `json:"id"`

	Role UserBodyRole `json:"role"`
}

type UserBodyRole string

type CompleteOAuthParams struct {
	UserAgent *string `json:"User-Agent,omitempty"`
}

type RegisterAccountParams struct {
	UserAgent *string `json:"User-Agent,omitempty"`
}

type CreateSessionParams struct {
	UserAgent *string `json:"User-Agent,omitempty"`
}

type CompleteOAuthJSONRequestBody = OAuthCallbackInputBody

type RegisterAccountJSONRequestBody = RegisterInputBody

type CreateSessionJSONRequestBody = LoginInputBody

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Server string

	Client HttpRequestDoer

	RequestEditors []RequestEditorFn
}

type ClientOption func(*Client) error

func NewClient(server string, opts ...ClientOption) (*Client, error) {

	client := Client{
		Server: server,
	}

	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

type ClientInterface interface {
	GetCurrentAccount(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	ListOAuthProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	CompleteOAuthWithBody(ctx context.Context, provider string, params *CompleteOAuthParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CompleteOAuth(ctx context.Context, provider string, params *CompleteOAuthParams, body CompleteOAuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	StartOAuth(ctx context.Context, provider string, reqEditors ...RequestEditorFn) (*http.Response, error)

	RegisterAccountWithBody(ctx context.Context, params *RegisterAccountParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	RegisterAccount(ctx context.Context, params *RegisterAccountParams, body RegisterAccountJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteAllSessions(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSessionWithBody(ctx context.Context, params *CreateSessionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateSession(ctx context.Context, params *CreateSessionParams, body CreateSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteCurrentSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	PingCatalog(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	PingMetadata(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	PingModeration(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	PingSocial(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetCurrentAccount(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCurrentAccountRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListOAuthProviders(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListOAuthProvidersRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CompleteOAuthWithBody(ctx context.Context, provider string, params *CompleteOAuthParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCompleteOAuthRequestWithBody(c.Server, provider, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CompleteOAuth(ctx context.Context, provider string, params *CompleteOAuthParams, body CompleteOAuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCompleteOAuthRequest(c.Server, provider, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) StartOAuth(ctx context.Context, provider string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewStartOAuthRequest(c.Server, provider)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RegisterAccountWithBody(ctx context.Context, params *RegisterAccountParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegisterAccountRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RegisterAccount(ctx context.Context, params *RegisterAccountParams, body RegisterAccountJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegisterAccountRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteAllSessions(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteAllSessionsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSessionWithBody(ctx context.Context, params *CreateSessionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSessionRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateSession(ctx context.Context, params *CreateSessionParams, body CreateSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateSessionRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteCurrentSession(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteCurrentSessionRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PingCatalog(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingCatalogRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PingMetadata(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingMetadataRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PingModeration(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingModerationRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PingSocial(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingSocialRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewGetCurrentAccountRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/me")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewListOAuthProvidersRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/oauth/providers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewCompleteOAuthRequest(server string, provider string, params *CompleteOAuthParams, body CompleteOAuthJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCompleteOAuthRequestWithBody(server, provider, params, "application/json", bodyReader)
}

func NewCompleteOAuthRequestWithBody(server string, provider string, params *CompleteOAuthParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "provider", provider, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/oauth/%s/callback", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.UserAgent != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithOptions("simple", false, "User-Agent", *params.UserAgent, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return nil, err
			}

			req.Header.Set("User-Agent", headerParam0)
		}

	}

	return req, nil
}

func NewStartOAuthRequest(server string, provider string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "provider", provider, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: ""})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/oauth/%s/start", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewRegisterAccountRequest(server string, params *RegisterAccountParams, body RegisterAccountJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRegisterAccountRequestWithBody(server, params, "application/json", bodyReader)
}

func NewRegisterAccountRequestWithBody(server string, params *RegisterAccountParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/register")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.UserAgent != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithOptions("simple", false, "User-Agent", *params.UserAgent, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return nil, err
			}

			req.Header.Set("User-Agent", headerParam0)
		}

	}

	return req, nil
}

func NewDeleteAllSessionsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/sessions")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewCreateSessionRequest(server string, params *CreateSessionParams, body CreateSessionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateSessionRequestWithBody(server, params, "application/json", bodyReader)
}

func NewCreateSessionRequestWithBody(server string, params *CreateSessionParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/sessions")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	if params != nil {

		if params.UserAgent != nil {
			var headerParam0 string

			headerParam0, err = runtime.StyleParamWithOptions("simple", false, "User-Agent", *params.UserAgent, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationHeader, Type: "string", Format: ""})
			if err != nil {
				return nil, err
			}

			req.Header.Set("User-Agent", headerParam0)
		}

	}

	return req, nil
}

func NewDeleteCurrentSessionRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/account/sessions/current")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewPingCatalogRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/catalog/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewPingMetadataRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/metadata/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewPingModerationRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/moderation/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewPingSocialRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/social/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

type ClientWithResponses struct {
	ClientInterface
}

func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

type ClientWithResponsesInterface interface {
	GetCurrentAccountWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCurrentAccountResponse, error)

	ListOAuthProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListOAuthProvidersResponse, error)

	CompleteOAuthWithBodyWithResponse(ctx context.Context, provider string, params *CompleteOAuthParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CompleteOAuthResponse, error)

	CompleteOAuthWithResponse(ctx context.Context, provider string, params *CompleteOAuthParams, body CompleteOAuthJSONRequestBody, reqEditors ...RequestEditorFn) (*CompleteOAuthResponse, error)

	StartOAuthWithResponse(ctx context.Context, provider string, reqEditors ...RequestEditorFn) (*StartOAuthResponse, error)

	RegisterAccountWithBodyWithResponse(ctx context.Context, params *RegisterAccountParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RegisterAccountResponse, error)

	RegisterAccountWithResponse(ctx context.Context, params *RegisterAccountParams, body RegisterAccountJSONRequestBody, reqEditors ...RequestEditorFn) (*RegisterAccountResponse, error)

	DeleteAllSessionsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteAllSessionsResponse, error)

	CreateSessionWithBodyWithResponse(ctx context.Context, params *CreateSessionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSessionResponse, error)

	CreateSessionWithResponse(ctx context.Context, params *CreateSessionParams, body CreateSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSessionResponse, error)

	DeleteCurrentSessionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteCurrentSessionResponse, error)

	PingCatalogWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingCatalogResponse, error)

	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)

	PingMetadataWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingMetadataResponse, error)

	PingModerationWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingModerationResponse, error)

	PingSocialWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingSocialResponse, error)
}

type GetCurrentAccountResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *UserBody

	ApplicationproblemJSON401 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r GetCurrentAccountResponse) GetJSON200() *UserBody {
	return r.JSON200
}

func (r GetCurrentAccountResponse) GetApplicationproblemJSON401() *ErrorModel {
	return r.ApplicationproblemJSON401
}

func (r GetCurrentAccountResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r GetCurrentAccountResponse) GetBody() []byte {
	return r.Body
}

func (r GetCurrentAccountResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetCurrentAccountResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetCurrentAccountResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type ListOAuthProvidersResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *OAuthProvidersOutputBody

	ApplicationproblemJSONDefault *ErrorModel
}

func (r ListOAuthProvidersResponse) GetJSON200() *OAuthProvidersOutputBody {
	return r.JSON200
}

func (r ListOAuthProvidersResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r ListOAuthProvidersResponse) GetBody() []byte {
	return r.Body
}

func (r ListOAuthProvidersResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r ListOAuthProvidersResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r ListOAuthProvidersResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type CompleteOAuthResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *SessionBody

	ApplicationproblemJSON401 *ErrorModel

	ApplicationproblemJSON404 *ErrorModel

	ApplicationproblemJSON409 *ErrorModel

	ApplicationproblemJSON422 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel

	ApplicationproblemJSON502 *ErrorModel
}

func (r CompleteOAuthResponse) GetJSON201() *SessionBody {
	return r.JSON201
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON401() *ErrorModel {
	return r.ApplicationproblemJSON401
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON404() *ErrorModel {
	return r.ApplicationproblemJSON404
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON409() *ErrorModel {
	return r.ApplicationproblemJSON409
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON422() *ErrorModel {
	return r.ApplicationproblemJSON422
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r CompleteOAuthResponse) GetApplicationproblemJSON502() *ErrorModel {
	return r.ApplicationproblemJSON502
}

func (r CompleteOAuthResponse) GetBody() []byte {
	return r.Body
}

func (r CompleteOAuthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r CompleteOAuthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r CompleteOAuthResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type StartOAuthResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *OAuthStartOutputBody

	ApplicationproblemJSON404 *ErrorModel

	ApplicationproblemJSON422 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r StartOAuthResponse) GetJSON200() *OAuthStartOutputBody {
	return r.JSON200
}

func (r StartOAuthResponse) GetApplicationproblemJSON404() *ErrorModel {
	return r.ApplicationproblemJSON404
}

func (r StartOAuthResponse) GetApplicationproblemJSON422() *ErrorModel {
	return r.ApplicationproblemJSON422
}

func (r StartOAuthResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r StartOAuthResponse) GetBody() []byte {
	return r.Body
}

func (r StartOAuthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r StartOAuthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r StartOAuthResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type RegisterAccountResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *SessionBody

	ApplicationproblemJSON422 *ErrorModel

	ApplicationproblemJSON429 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r RegisterAccountResponse) GetJSON201() *SessionBody {
	return r.JSON201
}

func (r RegisterAccountResponse) GetApplicationproblemJSON422() *ErrorModel {
	return r.ApplicationproblemJSON422
}

func (r RegisterAccountResponse) GetApplicationproblemJSON429() *ErrorModel {
	return r.ApplicationproblemJSON429
}

func (r RegisterAccountResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r RegisterAccountResponse) GetBody() []byte {
	return r.Body
}

func (r RegisterAccountResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r RegisterAccountResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r RegisterAccountResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type DeleteAllSessionsResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *RevokedOutputBody

	ApplicationproblemJSON401 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r DeleteAllSessionsResponse) GetJSON200() *RevokedOutputBody {
	return r.JSON200
}

func (r DeleteAllSessionsResponse) GetApplicationproblemJSON401() *ErrorModel {
	return r.ApplicationproblemJSON401
}

func (r DeleteAllSessionsResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r DeleteAllSessionsResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteAllSessionsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r DeleteAllSessionsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r DeleteAllSessionsResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type CreateSessionResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *SessionBody

	ApplicationproblemJSON401 *ErrorModel

	ApplicationproblemJSON403 *ErrorModel

	ApplicationproblemJSON422 *ErrorModel

	ApplicationproblemJSON429 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r CreateSessionResponse) GetJSON201() *SessionBody {
	return r.JSON201
}

func (r CreateSessionResponse) GetApplicationproblemJSON401() *ErrorModel {
	return r.ApplicationproblemJSON401
}

func (r CreateSessionResponse) GetApplicationproblemJSON403() *ErrorModel {
	return r.ApplicationproblemJSON403
}

func (r CreateSessionResponse) GetApplicationproblemJSON422() *ErrorModel {
	return r.ApplicationproblemJSON422
}

func (r CreateSessionResponse) GetApplicationproblemJSON429() *ErrorModel {
	return r.ApplicationproblemJSON429
}

func (r CreateSessionResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r CreateSessionResponse) GetBody() []byte {
	return r.Body
}

func (r CreateSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r CreateSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r CreateSessionResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type DeleteCurrentSessionResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	ApplicationproblemJSON401 *ErrorModel

	ApplicationproblemJSON500 *ErrorModel
}

func (r DeleteCurrentSessionResponse) GetApplicationproblemJSON401() *ErrorModel {
	return r.ApplicationproblemJSON401
}

func (r DeleteCurrentSessionResponse) GetApplicationproblemJSON500() *ErrorModel {
	return r.ApplicationproblemJSON500
}

func (r DeleteCurrentSessionResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteCurrentSessionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r DeleteCurrentSessionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r DeleteCurrentSessionResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type PingCatalogResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *PingOutputBody

	ApplicationproblemJSONDefault *ErrorModel
}

func (r PingCatalogResponse) GetJSON200() *PingOutputBody {
	return r.JSON200
}

func (r PingCatalogResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r PingCatalogResponse) GetBody() []byte {
	return r.Body
}

func (r PingCatalogResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r PingCatalogResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r PingCatalogResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *HealthOutputBody

	ApplicationproblemJSONDefault *ErrorModel
}

func (r GetHealthResponse) GetJSON200() *HealthOutputBody {
	return r.JSON200
}

func (r GetHealthResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r GetHealthResponse) GetBody() []byte {
	return r.Body
}

func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetHealthResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type PingMetadataResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *PingOutputBody1

	ApplicationproblemJSONDefault *ErrorModel
}

func (r PingMetadataResponse) GetJSON200() *PingOutputBody1 {
	return r.JSON200
}

func (r PingMetadataResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r PingMetadataResponse) GetBody() []byte {
	return r.Body
}

func (r PingMetadataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r PingMetadataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r PingMetadataResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type PingModerationResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *PingOutputBody3

	ApplicationproblemJSONDefault *ErrorModel
}

func (r PingModerationResponse) GetJSON200() *PingOutputBody3 {
	return r.JSON200
}

func (r PingModerationResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r PingModerationResponse) GetBody() []byte {
	return r.Body
}

func (r PingModerationResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r PingModerationResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r PingModerationResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type PingSocialResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *PingOutputBody2

	ApplicationproblemJSONDefault *ErrorModel
}

func (r PingSocialResponse) GetJSON200() *PingOutputBody2 {
	return r.JSON200
}

func (r PingSocialResponse) GetApplicationproblemJSONDefault() *ErrorModel {
	return r.ApplicationproblemJSONDefault
}

func (r PingSocialResponse) GetBody() []byte {
	return r.Body
}

func (r PingSocialResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r PingSocialResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r PingSocialResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

func (c *ClientWithResponses) GetCurrentAccountWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCurrentAccountResponse, error) {
	rsp, err := c.GetCurrentAccount(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCurrentAccountResponse(rsp)
}

func (c *ClientWithResponses) ListOAuthProvidersWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListOAuthProvidersResponse, error) {
	rsp, err := c.ListOAuthProviders(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListOAuthProvidersResponse(rsp)
}

func (c *ClientWithResponses) CompleteOAuthWithBodyWithResponse(ctx context.Context, provider string, params *CompleteOAuthParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CompleteOAuthResponse, error) {
	rsp, err := c.CompleteOAuthWithBody(ctx, provider, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCompleteOAuthResponse(rsp)
}

func (c *ClientWithResponses) CompleteOAuthWithResponse(ctx context.Context, provider string, params *CompleteOAuthParams, body CompleteOAuthJSONRequestBody, reqEditors ...RequestEditorFn) (*CompleteOAuthResponse, error) {
	rsp, err := c.CompleteOAuth(ctx, provider, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCompleteOAuthResponse(rsp)
}

func (c *ClientWithResponses) StartOAuthWithResponse(ctx context.Context, provider string, reqEditors ...RequestEditorFn) (*StartOAuthResponse, error) {
	rsp, err := c.StartOAuth(ctx, provider, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseStartOAuthResponse(rsp)
}

func (c *ClientWithResponses) RegisterAccountWithBodyWithResponse(ctx context.Context, params *RegisterAccountParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RegisterAccountResponse, error) {
	rsp, err := c.RegisterAccountWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegisterAccountResponse(rsp)
}

func (c *ClientWithResponses) RegisterAccountWithResponse(ctx context.Context, params *RegisterAccountParams, body RegisterAccountJSONRequestBody, reqEditors ...RequestEditorFn) (*RegisterAccountResponse, error) {
	rsp, err := c.RegisterAccount(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegisterAccountResponse(rsp)
}

func (c *ClientWithResponses) DeleteAllSessionsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteAllSessionsResponse, error) {
	rsp, err := c.DeleteAllSessions(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteAllSessionsResponse(rsp)
}

func (c *ClientWithResponses) CreateSessionWithBodyWithResponse(ctx context.Context, params *CreateSessionParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateSessionResponse, error) {
	rsp, err := c.CreateSessionWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSessionResponse(rsp)
}

func (c *ClientWithResponses) CreateSessionWithResponse(ctx context.Context, params *CreateSessionParams, body CreateSessionJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateSessionResponse, error) {
	rsp, err := c.CreateSession(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateSessionResponse(rsp)
}

func (c *ClientWithResponses) DeleteCurrentSessionWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*DeleteCurrentSessionResponse, error) {
	rsp, err := c.DeleteCurrentSession(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteCurrentSessionResponse(rsp)
}

func (c *ClientWithResponses) PingCatalogWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingCatalogResponse, error) {
	rsp, err := c.PingCatalog(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingCatalogResponse(rsp)
}

func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

func (c *ClientWithResponses) PingMetadataWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingMetadataResponse, error) {
	rsp, err := c.PingMetadata(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingMetadataResponse(rsp)
}

func (c *ClientWithResponses) PingModerationWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingModerationResponse, error) {
	rsp, err := c.PingModeration(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingModerationResponse(rsp)
}

func (c *ClientWithResponses) PingSocialWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingSocialResponse, error) {
	rsp, err := c.PingSocial(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingSocialResponse(rsp)
}

func ParseGetCurrentAccountResponse(rsp *http.Response) (*GetCurrentAccountResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCurrentAccountResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest UserBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParseListOAuthProvidersResponse(rsp *http.Response) (*ListOAuthProvidersResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListOAuthProvidersResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest OAuthProvidersOutputBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

func ParseCompleteOAuthResponse(rsp *http.Response) (*CompleteOAuthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CompleteOAuthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest SessionBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 502:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON502 = &dest

	}

	return response, nil
}

func ParseStartOAuthResponse(rsp *http.Response) (*StartOAuthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &StartOAuthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest OAuthStartOutputBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParseRegisterAccountResponse(rsp *http.Response) (*RegisterAccountResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RegisterAccountResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest SessionBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParseDeleteAllSessionsResponse(rsp *http.Response) (*DeleteAllSessionsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteAllSessionsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RevokedOutputBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParseCreateSessionResponse(rsp *http.Response) (*CreateSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest SessionBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 422:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON422 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParseDeleteCurrentSessionResponse(rsp *http.Response) (*DeleteCurrentSessionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteCurrentSessionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case rsp.StatusCode == 204:
		break

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSON500 = &dest

	}

	return response, nil
}

func ParsePingCatalogResponse(rsp *http.Response) (*PingCatalogResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingCatalogResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PingOutputBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest HealthOutputBody
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

func ParsePingMetadataResponse(rsp *http.Response) (*PingMetadataResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingMetadataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PingOutputBody1
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

func ParsePingModerationResponse(rsp *http.Response) (*PingModerationResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingModerationResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PingOutputBody3
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}

func ParsePingSocialResponse(rsp *http.Response) (*PingSocialResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingSocialResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PingOutputBody2
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest ErrorModel
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.ApplicationproblemJSONDefault = &dest

	}

	return response, nil
}
