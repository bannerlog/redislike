package main

import (
	"reflect"
	"testing"
)

const succeed = "\u2713"
const failed = "\u2717"

func TestNewStorage(t *testing.T) {
	t.Log("Given the need to test allocation of a new storage")
	{
		t.Log("\tWhen there is no allocated storage")
		{
			s := newStorage()

			if reflect.TypeOf(s) == reflect.TypeOf((*storage)(nil)) {
				t.Logf("\t%s\tShould return pointer to \"storage\"", succeed)
			} else {
				t.Errorf("\t%s\tShould return pointer to \"storage\"", failed)
			}
		}
	}
}

func TestSet(t *testing.T) {

	s := newStorage()

	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "string-key", "string-value"},
		{"list", "list-key", []string{"a", "b", "c"}},
		{"hash", "hash-key", map[string]string{"field-1": "val-1", "field-2": "val-2"}},
	}

	t.Log("Given an empty storage which should be filled with keys")

	for i, tt := range tests {
		tf := func(t *testing.T) {
			t.Logf("\tTest: %d\tWhen adding key %q of type %q", i, tt.key, tt.name)

			tr := false

			if r := s.set(tt.key, tt.value); r {
				e := s.getEntry(tt.key)
				if e != nil && reflect.DeepEqual(e.value, tt.value) {
					tr = true
				}
			}

			if tr {
				t.Logf("\t%s\tShould return TRUE and entry should appear in the storage", succeed)
			} else {
				t.Errorf("\t%s\tShould return TRUE and entry should appear in the storage", failed)
			}
		}

		t.Run(tt.name, tf)
	}

	t.Log("\tTest: 4\tWhen adding key with <nil> value")
	tr := false
	if r := s.set("key", nil); !r {
		if !s.exists("key") {
			tr = true
		}
	}

	if tr {
		t.Logf("\t%s\tShould return FALSE and entry should not appear in the storage", succeed)
	} else {
		t.Errorf("\t%s\tShould return FALSE and entry should not appear in the storage", failed)
	}
}

func TestGet(t *testing.T) {
	badKey := "nonexistent-key"
	goodKey := "existing-key"
	goodValue := "existing-key-value"

	t.Logf("Given storage with one key: %q", goodKey)
	s := newStorage()
	s.set(goodKey, goodValue)

	t.Logf("\tWhen getting existing key %q", goodKey)
	if v := s.get(goodKey); v != nil && reflect.DeepEqual(goodValue, v) {
		t.Logf("\t%s\tShould get value %q", succeed, goodValue)
	} else {
		t.Errorf("\t%s\tShould get value %q", failed, goodValue)
	}

	t.Logf("\tWhen getting nonexistent key %q", badKey)
	if v := s.get(badKey); v == nil {
		t.Logf("\t%s\tShould get nil", succeed)
	} else {
		t.Errorf("\t%s\tShould get nil", failed)
	}
}

func TestDel(t *testing.T) {
	keyToDelete := "key-to-delete"
	keySurvivor := "key-survivor"
	t.Logf("Given storage which contains two keys")
	s := newStorage()
	s.set(keyToDelete, "Delete me if you can.")
	s.set(keySurvivor, "I will live forever!")

	t.Logf("\tWhen calling delete function with key %q", keyToDelete)
	s.del(keyToDelete)
	if !s.exists(keyToDelete) && s.exists(keySurvivor) {
		t.Logf("\t%s\tShould have only %q", succeed, keySurvivor)
	} else {
		t.Errorf("\t%s\tShould have only %q", succeed, keySurvivor)
	}
}
