package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// fakeDoer records requests and replays canned responses in order.
type fakeDoer struct {
	responses []fakeResponse
	requests  []recordedRequest
}

type fakeResponse struct {
	status int
	body   string
}

type recordedRequest struct {
	method string
	url    string
	auth   string
	body   string
}

func (f *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{
		method: req.Method,
		url:    req.URL.String(),
		auth:   req.Header.Get("Authorization"),
	}
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		rec.body = string(b)
	}
	f.requests = append(f.requests, rec)

	if len(f.responses) == 0 {
		return nil, io.EOF
	}
	next := f.responses[0]
	f.responses = f.responses[1:]
	return &http.Response{
		StatusCode: next.status,
		Body:       io.NopCloser(strings.NewReader(next.body)),
	}, nil
}

func newTestClient(responses ...fakeResponse) (*Client, *fakeDoer) {
	doer := &fakeDoer{responses: responses}
	c := NewClient("https://api.test.local", "ASOT-TEST-KEY", "user@example.com")
	c.HTTP = doer
	return c, doer
}

const activateOK = `{"token":"jwt-1","user":{"id":"u1","email":"user@example.com","plan":"pro"}}`

func TestGetActivatesBeforeFirstRequest(t *testing.T) {
	c, doer := newTestClient(
		fakeResponse{200, activateOK},
		fakeResponse{200, `[{"id":"a1","name":"シンプル録音"}]`},
	)

	var apps []App
	if err := c.Get(context.Background(), "/api/apps", &apps); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(doer.requests) != 2 {
		t.Fatalf("リクエスト2回のはず: got %d", len(doer.requests))
	}
	if got := doer.requests[0].url; !strings.HasSuffix(got, "/api/licenses/activate") {
		t.Errorf("1回目は activate のはず: %s", got)
	}
	// activate はライセンスキーを `key` で送る (`license_key` ではない)
	var body map[string]string
	if err := json.Unmarshal([]byte(doer.requests[0].body), &body); err != nil {
		t.Fatalf("activate body: %v", err)
	}
	if body["key"] != "ASOT-TEST-KEY" || body["email"] != "user@example.com" {
		t.Errorf("activate のボディが違う: %v", body)
	}
	if doer.requests[0].auth != "" {
		t.Error("activate に Authorization を付けてはいけない (401リトライが再帰する)")
	}
	if doer.requests[1].auth != "Bearer jwt-1" {
		t.Errorf("2回目に取得済みトークンを使っていない: %q", doer.requests[1].auth)
	}
	if len(apps) != 1 || apps[0].Name != "シンプル録音" {
		t.Errorf("デコード結果が違う: %+v", apps)
	}
}

func TestGetReactivatesOnceOn401(t *testing.T) {
	c, doer := newTestClient(
		fakeResponse{200, activateOK},                    // 初回 activate
		fakeResponse{401, `{"error":"unauthorized"}`},    // トークン失効
		fakeResponse{200, `{"token":"jwt-2","user":{}}`}, // 再アクティベート
		fakeResponse{200, `[]`},                          // 新トークンでリトライ成功
	)

	var apps []App
	if err := c.Get(context.Background(), "/api/apps", &apps); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(doer.requests) != 4 {
		t.Fatalf("activate→401→activate→retry の4回のはず: got %d", len(doer.requests))
	}
	if doer.requests[1].auth != "Bearer jwt-1" {
		t.Errorf("失効前は jwt-1 のはず: %q", doer.requests[1].auth)
	}
	if !strings.HasSuffix(doer.requests[2].url, "/api/licenses/activate") {
		t.Errorf("401後に activate していない: %s", doer.requests[2].url)
	}
	if doer.requests[3].auth != "Bearer jwt-2" {
		t.Errorf("リトライで新トークンを使っていない: %q", doer.requests[3].auth)
	}
}

func TestGetDoesNotRetryTwiceOn401(t *testing.T) {
	c, doer := newTestClient(
		fakeResponse{200, activateOK},
		fakeResponse{401, `{"error":"unauthorized"}`},
		fakeResponse{200, `{"token":"jwt-2","user":{}}`},
		fakeResponse{401, `{"error":"unauthorized"}`}, // リトライも401
	)

	err := c.Get(context.Background(), "/api/apps", &[]App{})
	if err == nil {
		t.Fatal("2回目も401ならエラーを返すべき")
	}
	if len(doer.requests) != 4 {
		t.Errorf("無限リトライしてはいけない: got %d requests", len(doer.requests))
	}
}

func TestGetSurfacesStatusCodeForNon401(t *testing.T) {
	c, _ := newTestClient(
		fakeResponse{200, activateOK},
		fakeResponse{403, `{"error":"keyword gap analysis requires Pro plan"}`},
	)

	err := c.Get(context.Background(), "/api/apps/a1/competitors/keyword-gap", &[]any{})
	if err == nil {
		t.Fatal("403 はエラーになるべき")
	}
	apiErr, ok := err.(*apiError)
	if !ok {
		t.Fatalf("*apiError のはず: %T", err)
	}
	if apiErr.Status != 403 {
		t.Errorf("ステータスが違う: %d", apiErr.Status)
	}
	if !strings.Contains(apiErr.Body, "Pro plan") {
		t.Errorf("プラン制限の理由が失われている: %q", apiErr.Body)
	}
}

func TestActivateRejectsResponseWithoutToken(t *testing.T) {
	c, _ := newTestClient(fakeResponse{200, `{"user":{"id":"u1"}}`})

	err := c.Get(context.Background(), "/api/apps", &[]App{})
	if err == nil || !strings.Contains(err.Error(), "no token") {
		t.Fatalf("トークン無しレスポンスを弾くべき: %v", err)
	}
}

func TestNewClientDefaultsBaseURLAndTrimsSlash(t *testing.T) {
	if got := NewClient("", "k", "e").BaseURL; got != DefaultBaseURL {
		t.Errorf("既定のベースURLが違う: %s", got)
	}
	if got := NewClient("https://example.com/", "k", "e").BaseURL; got != "https://example.com" {
		t.Errorf("末尾スラッシュが除去されていない: %s", got)
	}
}
