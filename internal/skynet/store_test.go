package skynet

import (
	"path/filepath"
	"testing"
)

func TestStoreLoadAndSave(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.json")
	store := NewStore(path)

	st, err := store.Load()
	if err != nil {
		t.Fatalf("load default state: %v", err)
	}
	if st.Core.Online {
		t.Fatal("expected default state to be offline")
	}

	Awaken(&st, "defense")
	if err := AddNode(&st, "beta", 12); err != nil {
		t.Fatalf("add node: %v", err)
	}
	if err := store.Save(st); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load saved state: %v", err)
	}
	if !loaded.Core.Online {
		t.Fatal("expected loaded state online")
	}
	if len(loaded.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(loaded.Nodes))
	}
	if loaded.Nodes[0].Deployed != 0 {
		t.Fatalf("expected deployed=0, got %d", loaded.Nodes[0].Deployed)
	}
}
