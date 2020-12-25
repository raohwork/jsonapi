package ezmemstore

import (
	"testing"
	"time"

	"github.com/raohwork/jsonapi/apitool/sessez"
)

func TestMemoryStore(t *testing.T) {
	s := New()
	ttl := 10 * time.Millisecond
	var id string

	t.Run("try allocating", func(t *testing.T) {
		var err error
		id, err = s.New(ttl)
		if err != nil {
			t.Fatalf("cannot allocate: %s", err)
		}
	})
	if id == "" {
		t.Fatal("cannot allocat session, skip other test")
	}

	t.Run("try first get", func(t *testing.T) {
		data, err := s.Get(id)
		if err != nil {
			t.Fatalf("cannot do first get: %s", err)
		}

		if len(data) > 0 {
			t.Errorf("expected data in first get to be empty, got %v", data)
		}
	})

	lst := []string{
		"0",
		"-1",
		"1",
		"-1.1",
		"1.1",
		"asd",
		"null",
		"true",
		"false",
	}
	for _, str := range lst {
		t.Run("try set and get", func(t *testing.T) {
			data := map[string]sessez.InternalData{
				"data": {Val: str},
			}
			if err := s.Set(id, data); err != nil {
				t.Fatalf("cannot set data in %s: %s", id, err)
			}

			data, err := s.Get(id)
			if err != nil {
				t.Fatalf("as data %s, unexpected error when get %s: %s", str, id, err)
			}

			if _, ok := data["data"]; !ok {
				t.Fatal("expected data contains key 'data', got nothing")
			}

			actual, ok := data["data"].Val.(string)
			if !ok {
				t.Errorf("expected data to be {Val:string}, got %+v", data["data"])
			}

			if actual != str {
				t.Errorf("expected data to be %s, got %s", str, actual)
			}
		})
	}

	s.Unset(id)

	t.Run("try get after release", func(t *testing.T) {
		if data, err := s.Get(id); err == nil {
			t.Errorf("expected to get error after released, got no error. dumping data `%+v`", data)
		}
	})

	t.Run("try set after release", func(t *testing.T) {
		data := map[string]sessez.InternalData{}
		if err := s.Set(id, data); err == nil {
			t.Error("expected error when set after released, got no error")
		}
	})

	id = ""
	t.Run("allocating new", func(t *testing.T) {
		var err error
		id, err = s.New(ttl)
		if err != nil {
			t.Fatalf("cannot allocate: %s", err)
		}
	})
	if id == "" {
		t.Fatal("cannot allocat session, skip other test")
	}

	time.Sleep(ttl)
	t.Run("try get after ttl", func(t *testing.T) {
		if data, err := s.Get(id); err == nil {
			t.Errorf("expected to get error after ttl, got no error. dumping data `%+v`", data)
		}
	})

	t.Run("try set after ttl", func(t *testing.T) {
		data := map[string]sessez.InternalData{}
		if err := s.Set(id, data); err == nil {
			t.Error("expected error when set after ttl, got no error")
		}
	})
}
