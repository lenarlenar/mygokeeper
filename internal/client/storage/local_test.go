package storage

import (
	"path/filepath"
	"testing"
)

func TestSaveAndLoadRecords(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test_store.enc")
	pass := "secret"

	original := []LocalRecord{
		{Type: "password", Data: "bXk=", Meta: "site1"},
		{Type: "card", Data: "Y2FyZA==", Meta: "bank"},
	}

	err := SaveRecordsTo(file, pass, original)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := LoadRecordsFrom(file, pass)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(loaded) != len(original) {
		t.Fatalf("record count mismatch: got %d, want %d", len(loaded), len(original))
	}
	for i := range loaded {
		if loaded[i] != original[i] {
			t.Errorf("record %d mismatch: got %+v, want %+v", i, loaded[i], original[i])
		}
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test_store.enc")

	data := []LocalRecord{
		{Type: "text", Data: "aGVsbG8=", Meta: "note"},
	}
	err := SaveRecordsTo(file, "correct-pass", data)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	_, err = LoadRecordsFrom(file, "wrong-pass")
	if err == nil {
		t.Fatal("expected error on wrong passphrase")
	}
}
