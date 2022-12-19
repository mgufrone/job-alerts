package wrapper

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io/ioutil"
	"math/rand"
	worker2 "mgufrone.dev/job-alerts/pkg/worker"
	"net/http"
)

var (
	agents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
	}
)

func RandomUserAgent() string {
	randID := rand.Int31n(int32(len(agents)))
	return agents[randID]
}

type HTTPClient struct {
	agent  *http.Client
	tracer opentracing.Tracer
}
type TracerRoundtripper struct {
	tracer   opentracing.Tracer
	original http.RoundTripper
}

func (t *TracerRoundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	agent := req.UserAgent()
	if agent == "" {
		req.Header.Add("User-Agent", agent)
	}
	ctx := req.Context()
	spRef := opentracing.SpanFromContext(ctx).Context()
	sp := t.tracer.StartSpan(fmt.Sprintf("httpClient/%s", req.Method), opentracing.ChildOf(spRef))
	sp.SetTag("http.req.url", req.URL.String())
	sp.SetTag("http.req.method", req.Method)
	sp.LogKV("http.req.query", req.URL.Query().Encode())
	sp.LogKV("http.agent", agent)
	if req.GetBody != nil {
		if bd, err := req.GetBody(); err == nil && bd != nil {
			str, _ := ioutil.ReadAll(bd)
			sp.LogKV("http.req.body", string(str))
		}
	}

	rt := t.original
	if rt == nil {
		rt = http.DefaultTransport
	}
	res, err := rt.RoundTrip(req)
	sp.LogKV("http.res.code", res.StatusCode)
	if err != nil {
		sp.LogKV("http.res.error", err)
	}
	sp.Finish()
	return res, err
}

type ClientOption func(agent *HTTPClient)

func WithTracer(tracer opentracing.Tracer) ClientOption {
	return func(agent *HTTPClient) {
		agent.tracer = tracer
		original := agent.agent.Transport
		agent.agent.Transport = &TracerRoundtripper{tracer: tracer, original: original}
	}
}

func (h *HTTPClient) ToHTTPClient() *http.Client {
	return h.agent
}

func NewHTTPClient(agent *http.Client, opts ...ClientOption) worker2.IHTTPClient {
	if agent == nil {
		return nil
	}
	cli := &HTTPClient{agent: agent}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func (h *HTTPClient) Do(req *http.Request) (res *http.Response, err error) {
	if req.Header.Get("user-agent") == "" {
		ua := RandomUserAgent()
		req.Header.Set(
			"user-agent",
			ua,
		)
	}
	res, err = h.agent.Do(req)
	return
}
