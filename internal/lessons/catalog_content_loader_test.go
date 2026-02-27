package lessons

import "testing"

func TestLoadDefaultCatalogContent(t *testing.T) {
	content, err := loadDefaultCatalogContent()
	if err != nil {
		t.Fatalf("load default catalog content: %v", err)
	}

	if got := len(content.CPU); got != 8 {
		t.Fatalf("cpu lesson specs=%d want=8", got)
	}
	if got := len(content.Memory.Lessons); got != 7 {
		t.Fatalf("memory lesson specs=%d want=7", got)
	}
	if got := len(content.Concurrency.Lessons); got != 6 {
		t.Fatalf("concurrency lesson specs=%d want=6", got)
	}
	if got := len(content.Persistence.Lessons); got != 7 {
		t.Fatalf("persistence lesson specs=%d want=7", got)
	}
}
