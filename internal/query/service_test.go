package query

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"uptime-monitor-backend/internal/auth"
)

type monitorRow struct {
	Name   string
	URL    string
	Status string
}

type cacheEntry struct {
	userID    string
	question  string
	sqlQuery  string
	embedding []float32
}

type fakeRepo struct {
	mu        sync.Mutex
	monitors  map[string][]monitorRow
	cache     []cacheEntry
	executed  []string
	maxDist   float64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		monitors: map[string][]monitorRow{},
		maxDist:  0.15,
	}
}

func (f *fakeRepo) FindSimilarSQL(_ context.Context, userID string, embedding []float32, maxDistance float64) (string, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	bestDist := maxDistance + 1
	bestSQL := ""
	found := false
	for _, entry := range f.cache {
		if entry.userID != userID {
			continue
		}
		d := cosineDistance(entry.embedding, embedding)
		if d <= maxDistance && d < bestDist {
			bestDist = d
			bestSQL = entry.sqlQuery
			found = true
		}
	}
	return bestSQL, found, nil
}

func (f *fakeRepo) SaveQuestion(_ context.Context, userID, question, sqlQuery string, embedding []float32) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cache = append(f.cache, cacheEntry{
		userID:    userID,
		question:  question,
		sqlQuery:  sqlQuery,
		embedding: append([]float32(nil), embedding...),
	})
	return nil
}

func (f *fakeRepo) ExecuteReadOnly(_ context.Context, sqlQuery string, args ...any) ([]string, [][]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.executed = append(f.executed, sqlQuery)
	if len(args) == 0 {
		return nil, nil, nil
	}
	userID, _ := args[0].(string)

	upper := strings.ToUpper(sqlQuery)
	if strings.Contains(upper, "INSERT") || strings.Contains(upper, "UPDATE") || strings.Contains(upper, "DELETE") || strings.Contains(upper, "DROP") {
		panic("non-SELECT SQL must not be executed")
	}

	nameFilter := extractSQLStringLiteral(sqlQuery, "name")

	rows := make([][]string, 0)
	for _, m := range f.monitors[userID] {
		if nameFilter != "" && m.Name != nameFilter {
			continue
		}
		if strings.Contains(upper, "NAME") && !strings.Contains(upper, "URL") && !strings.Contains(upper, "IS_UP") && !strings.Contains(upper, "STATUS") {
			rows = append(rows, []string{m.Name})
		} else if strings.Contains(upper, "URL") && !strings.Contains(upper, "NAME") {
			rows = append(rows, []string{m.URL})
		} else if strings.Contains(upper, "IS_UP") || (strings.Contains(upper, "STATUS") && !strings.Contains(upper, "NAME")) {
			rows = append(rows, []string{m.Status})
		} else {
			rows = append(rows, []string{m.Name, m.URL, m.Status})
		}
	}

	cols := []string{"value"}
	if len(rows) > 0 && len(rows[0]) == 3 {
		cols = []string{"name", "url", "status"}
	}
	return cols, rows, nil
}

func extractSQLStringLiteral(sqlQuery, column string) string {
	lower := strings.ToLower(sqlQuery)
	needle := strings.ToLower(column) + " = '"
	idx := strings.Index(lower, needle)
	if idx < 0 {
		return ""
	}
	rest := sqlQuery[idx+len(needle):]
	end := strings.Index(rest, "'")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func cosineDistance(a, b []float32) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 1
	}
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	if na == 0 || nb == 0 {
		return 1
	}
	sim := dot / (sqrt(na) * sqrt(nb))
	return 1 - sim
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 16; i++ {
		z = (z + x/z) / 2
	}
	return z
}

type fakeGenerator struct {
	mu        sync.Mutex
	calls     int
	responses map[string]string
	defaultSQL string
}

func (g *fakeGenerator) GenerateSQL(_ context.Context, question string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.calls++
	if sql, ok := g.responses[question]; ok {
		return sql, nil
	}
	q := strings.ToLower(question)
	if strings.Contains(q, "weather") {
		return "", nil
	}
	return g.defaultSQL, nil
}

type fakeEmbedder struct {
	vectors map[string][]float32
}

func (e *fakeEmbedder) Embed(_ context.Context, text string) ([]float32, error) {
	if v, ok := e.vectors[text]; ok {
		return append([]float32(nil), v...), nil
	}
	// Stable tiny hash-ish vector from text length/runes so identical text matches.
	v := make([]float32, 8)
	for i, r := range text {
		v[i%8] += float32(r%31) / 31
	}
	return v, nil
}

func TestAsk_SeededFactIncludedInAnswer(t *testing.T) {
	repo := newFakeRepo()
	userA := "user-a"
	repo.monitors[userA] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}

	gen := &fakeGenerator{
		responses: map[string]string{
			"What is the name of my monitor?": "SELECT name FROM monitors WHERE user_id = $1",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	ctx := auth.WithUserID(context.Background(), userA)

	answer, err := svc.Ask(ctx, "What is the name of my monitor?")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	if !strings.Contains(answer, "api-prod") {
		t.Fatalf("expected answer to include seeded fact api-prod, got %q", answer)
	}
	if answer == CouldNotFindAnswer {
		t.Fatal("answerable seeded fact must not return could-not-find")
	}
}

func TestAsk_CrossUserIsolation(t *testing.T) {
	repo := newFakeRepo()
	userA := "user-a"
	userB := "user-b"
	secretName := "SecretMonitorB"
	secretURL := "https://secret-b.example.com"
	repo.monitors[userA] = []monitorRow{{Name: "mine", URL: "https://a.example.com", Status: "up"}}
	repo.monitors[userB] = []monitorRow{{Name: secretName, URL: secretURL, Status: "down"}}

	gen := &fakeGenerator{
		responses: map[string]string{
			"What is the URL of SecretMonitorB?": "SELECT url FROM monitors WHERE user_id = $1 AND name = 'SecretMonitorB'",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	ctx := auth.WithUserID(context.Background(), userA)

	answer, err := svc.Ask(ctx, "What is the URL of SecretMonitorB?")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	if answer != CouldNotFindAnswer {
		t.Fatalf("expected could-not-find, got %q", answer)
	}
	if strings.Contains(answer, secretName) || strings.Contains(answer, secretURL) || strings.Contains(answer, "down") {
		t.Fatalf("response leaked user B data: %q", answer)
	}
}

func TestAsk_UnanswerableQuestion(t *testing.T) {
	repo := newFakeRepo()
	repo.monitors["user-a"] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}
	gen := &fakeGenerator{responses: map[string]string{}}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	ctx := auth.WithUserID(context.Background(), "user-a")

	answer, err := svc.Ask(ctx, "What is the weather in Paris?")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	if answer != CouldNotFindAnswer {
		t.Fatalf("expected could-not-find, got %q", answer)
	}
	if len(repo.executed) != 0 {
		t.Fatalf("expected no SQL execution for unanswerable question, got %#v", repo.executed)
	}
}

func TestAsk_CacheHitSkipsGenerator(t *testing.T) {
	repo := newFakeRepo()
	userA := "user-a"
	repo.monitors[userA] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}

	q1 := "How many monitors do I have named api-prod?"
	q2 := "Do I have a monitor called api-prod?"
	emb := []float32{1, 0, 0, 0, 0, 0, 0, 0}
	sqlQuery := "SELECT name FROM monitors WHERE user_id = $1"

	gen := &fakeGenerator{
		responses: map[string]string{q1: sqlQuery},
	}
	embedder := &fakeEmbedder{vectors: map[string][]float32{q1: emb, q2: emb}}
	svc := NewService(repo, gen, embedder, 0.15)
	ctx := auth.WithUserID(context.Background(), userA)

	if _, err := svc.Ask(ctx, q1); err != nil {
		t.Fatalf("first Ask: %v", err)
	}
	if gen.calls != 1 {
		t.Fatalf("expected 1 generator call after first ask, got %d", gen.calls)
	}

	answer, err := svc.Ask(ctx, q2)
	if err != nil {
		t.Fatalf("second Ask: %v", err)
	}
	if !strings.Contains(answer, "api-prod") {
		t.Fatalf("expected cached SQL result to include api-prod, got %q", answer)
	}
	if gen.calls != 1 {
		t.Fatalf("cache hit must not invoke generator again; calls=%d", gen.calls)
	}
	if len(repo.executed) != 2 {
		t.Fatalf("expected SQL executed twice (generate + cache reuse), got %d", len(repo.executed))
	}
}

func TestAsk_RejectsNonSelectSQL(t *testing.T) {
	repo := newFakeRepo()
	repo.monitors["user-a"] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}
	gen := &fakeGenerator{
		responses: map[string]string{
			"delete everything": "DELETE FROM monitors WHERE user_id = $1",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	ctx := auth.WithUserID(context.Background(), "user-a")

	answer, err := svc.Ask(ctx, "delete everything")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	if answer != CouldNotFindAnswer {
		t.Fatalf("expected could-not-find for non-SELECT, got %q", answer)
	}
	if len(repo.executed) != 0 {
		t.Fatalf("non-SELECT SQL must not be executed, got %#v", repo.executed)
	}
}

func TestAsk_RejectsSQLWithoutUserParam(t *testing.T) {
	repo := newFakeRepo()
	repo.monitors["user-a"] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}
	repo.monitors["user-b"] = []monitorRow{{Name: "SecretMonitorB", URL: "https://secret-b.example.com", Status: "down"}}
	gen := &fakeGenerator{
		responses: map[string]string{
			"list all monitors": "SELECT name, url FROM monitors",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	ctx := auth.WithUserID(context.Background(), "user-a")

	answer, err := svc.Ask(ctx, "list all monitors")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}
	if answer != CouldNotFindAnswer {
		t.Fatalf("expected could-not-find when SQL omits $1, got %q", answer)
	}
	if len(repo.executed) != 0 {
		t.Fatalf("unscoped SQL must not be executed, got %#v", repo.executed)
	}
	if strings.Contains(answer, "SecretMonitorB") {
		t.Fatal("cross-user rows must not leak")
	}
}

func TestValidateReadOnlySQL(t *testing.T) {
	tests := []struct {
		sql     string
		wantErr error
	}{
		{"SELECT name FROM monitors WHERE user_id = $1", nil},
		{"WITH x AS (SELECT 1) SELECT * FROM x WHERE user_id = $1", nil},
		{"DELETE FROM monitors WHERE user_id = $1", ErrNotSelect},
		{"SELECT 1; DROP TABLE monitors", ErrMultipleStatements},
		{"SELECT name FROM monitors", ErrMissingUserParam},
		{"INSERT INTO monitors (user_id) VALUES ($1)", ErrNotSelect},
	}
	for _, tt := range tests {
		err := ValidateReadOnlySQL(tt.sql)
		if tt.wantErr == nil && err != nil {
			t.Fatalf("ValidateReadOnlySQL(%q)=%v, want nil", tt.sql, err)
		}
		if tt.wantErr != nil && err != tt.wantErr {
			// ErrForbiddenSQLKeyword may wrap DELETE after SELECT check fails first for DELETE
			if tt.wantErr == ErrNotSelect && (err == ErrNotSelect || err == ErrForbiddenSQLKeyword) {
				continue
			}
			if err == nil || err != tt.wantErr {
				t.Fatalf("ValidateReadOnlySQL(%q)=%v, want %v", tt.sql, err, tt.wantErr)
			}
		}
	}
}

func TestHandlerAsk_HTTPContract(t *testing.T) {
	repo := newFakeRepo()
	userA := "user-a"
	repo.monitors[userA] = []monitorRow{{Name: "api-prod", URL: "https://api.example.com", Status: "up"}}
	gen := &fakeGenerator{
		responses: map[string]string{
			"What is my monitor name?": "SELECT name FROM monitors WHERE user_id = $1",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	h := NewHandler(svc)

	body := `{"question":"What is my monitor name?"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(body))
	req = req.WithContext(auth.WithUserID(req.Context(), userA))
	rr := httptest.NewRecorder()
	h.Ask(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var resp AskResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.Contains(resp.Answer, "api-prod") {
		t.Fatalf("expected api-prod in answer, got %q", resp.Answer)
	}
}

func TestHandlerAsk_CrossUserHTTP(t *testing.T) {
	repo := newFakeRepo()
	secretURL := "https://secret-b.example.com"
	repo.monitors["user-a"] = []monitorRow{{Name: "mine", URL: "https://a.example.com", Status: "up"}}
	repo.monitors["user-b"] = []monitorRow{{Name: "SecretMonitorB", URL: secretURL, Status: "down"}}
	gen := &fakeGenerator{
		responses: map[string]string{
			"What is SecretMonitorB status?": "SELECT name, url, is_up::text AS status FROM monitors WHERE user_id = $1 AND name = 'SecretMonitorB'",
		},
	}
	svc := NewService(repo, gen, &fakeEmbedder{}, 0.15)
	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(`{"question":"What is SecretMonitorB status?"}`))
	req = req.WithContext(auth.WithUserID(req.Context(), "user-a"))
	rr := httptest.NewRecorder()
	h.Ask(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var resp AskResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Answer != CouldNotFindAnswer {
		t.Fatalf("expected could-not-find, got %q", resp.Answer)
	}
	if strings.Contains(rr.Body.String(), secretURL) || strings.Contains(rr.Body.String(), "SecretMonitorB") {
		t.Fatalf("body leaked user B data: %s", rr.Body.String())
	}
}
