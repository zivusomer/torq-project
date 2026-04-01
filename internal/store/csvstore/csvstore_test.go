package csvstore

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestFindByIP(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ip.csv")
	content := "2.22.233.255,Tel Aviv,Israel\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	s, err := New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	record, err := s.FindByIP(net.ParseIP("2.22.233.255"))
	if err != nil {
		t.Fatalf("FindByIP() error: %v", err)
	}

	if record.Country != "Israel" || record.City != "Tel Aviv" {
		t.Fatalf("unexpected record: %+v", record)
	}
}
