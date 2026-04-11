package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"go-admin-api/internal/testutil"
)

type stubMessage struct {
	id         string
	stream     string
	values     map[string]interface{}
	prefix     string
	errorCount int
}

func (m *stubMessage) SetID(id string)                            { m.id = id }
func (m *stubMessage) SetStream(stream string)                    { m.stream = stream }
func (m *stubMessage) SetValues(values map[string]interface{})    { m.values = values }
func (m *stubMessage) GetID() string                              { return m.id }
func (m *stubMessage) GetStream() string                          { return m.stream }
func (m *stubMessage) GetValues() map[string]interface{}          { return m.values }
func (m *stubMessage) GetPrefix() string                          { return m.prefix }
func (m *stubMessage) SetPrefix(prefix string)                    { m.prefix = prefix }
func (m *stubMessage) SetErrorCount(count int)                    { m.errorCount = count }
func (m *stubMessage) GetErrorCount() int                         { return m.errorCount }

func TestSaveSysApiRewritesIDPathAndIsIdempotent(t *testing.T) {
	db := testutil.NewTestDB(t, &SysApi{})
	testutil.UseRuntimeDB(t, db)

	tempDir := t.TempDir()
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, "docs"), 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}
	swagger := `{"paths":{"/api/v1/sys-api/{id}":{"get":{"summary":"Get API detail"}}}}`
	if err := os.WriteFile(filepath.Join(tempDir, "docs", "swagger.json"), []byte(swagger), 0o644); err != nil {
		t.Fatalf("write swagger.json: %v", err)
	}

	message := &stubMessage{
		values: map[string]interface{}{
			"List": []runtime.Router{
				{RelativePath: "/api/v1/sys-api/:id", HttpMethod: "GET", Handler: "handler.Get"},
			},
		},
	}

	if err := SaveSysApi(message); err != nil {
		t.Fatalf("SaveSysApi() error = %v", err)
	}
	if err := SaveSysApi(message); err != nil {
		t.Fatalf("SaveSysApi() second call error = %v", err)
	}

	var api SysApi
	if err := db.First(&api).Error; err != nil {
		t.Fatalf("reload api: %v", err)
	}
	if api.Path != "/api/v1/sys-api/:id" {
		t.Fatalf("stored path = %q, want original route path", api.Path)
	}
	if api.Title != "Get API detail" {
		t.Fatalf("stored title = %q, want swagger summary", api.Title)
	}

	var count int64
	if err := db.Model(&SysApi{}).Count(&count).Error; err != nil {
		t.Fatalf("count apis: %v", err)
	}
	if count != 1 {
		t.Fatalf("api count = %d, want 1", count)
	}
}

func TestSaveSysApiCurrentRouteFilterBehavior(t *testing.T) {
	db := testutil.NewTestDB(t, &SysApi{})
	testutil.UseRuntimeDB(t, db)

	tempDir := t.TempDir()
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tempDir, "docs"), 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "docs", "swagger.json"), []byte(`{"paths":{}}`), 0o644); err != nil {
		t.Fatalf("write swagger.json: %v", err)
	}

	message := &stubMessage{
		values: map[string]interface{}{
			"List": []runtime.Router{
				{RelativePath: "/head", HttpMethod: "HEAD", Handler: "handler.Head"},
				{RelativePath: "/swagger/admin/index", HttpMethod: "HEAD", Handler: "handler.Swagger"},
			},
		},
	}

	if err := SaveSysApi(message); err != nil {
		t.Fatalf("SaveSysApi() error = %v", err)
	}

	var apis []SysApi
	if err := db.Order("path").Find(&apis).Error; err != nil {
		t.Fatalf("load apis: %v", err)
	}
	if len(apis) != 1 || apis[0].Path != "/swagger/admin/index" {
		t.Fatalf("current route filter behavior = %+v, want only swagger HEAD route stored", apis)
	}
}
