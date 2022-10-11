package civogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSSHKey_List(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/sshkeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allSSHKeys := []SSHKey{
			{ID: "a5bd357a-8afd-4f60-b055-ece013646f55"},
			{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad"},
		}
		value := toJSON(t, allSSHKeys)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.SSHKey().List(ctx)
	if err != nil {
		t.Errorf("SSHKey.List returned error: %v", err)
	}

	expectedKeys := []SSHKey{{ID: "a5bd357a-8afd-4f60-b055-ece013646f55"}, {ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad"}}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SSHKey.List returned keys %+v, expected %+v", keys, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.List returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestSSHKey_GetByID(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/sshkeys/a5bd357a-8afd-4f60-b055-ece013646f55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		oneSSHKey := SSHKey{ID: "a5bd357a-8afd-4f60-b055-ece013646f55"}
		value := toJSON(t, oneSSHKey)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.SSHKey().GetByID(ctx, "a5bd357a-8afd-4f60-b055-ece013646f55")
	if err != nil {
		t.Errorf("SSHKey.GetByID returned error: %v", err)
	}

	expectedKeys := &SSHKey{ID: "a5bd357a-8afd-4f60-b055-ece013646f55"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SSHKey.GetByID returned keys %+v, expected %+v", keys, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.GetByID returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestSSHKey_Find(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/sshkeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allSSHKeys := []SSHKey{
			{ID: "a5bd357a-8afd-4f60-b055-ece013646f55"},
			{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad"},
		}
		value := toJSON(t, allSSHKeys)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.SSHKey().Find(ctx, "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad")
	if err != nil {
		t.Errorf("SSHKey.Find returned error: %v", err)
	}

	expectedKeys := &SSHKey{ID: "2ecf2e0a-629c-4d16-9cb9-aa81059c4bad"}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("SSHKey.Find returned keys %+v, expected %+v", keys, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.Find returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestSSHKey_Create(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/sshkeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	newSSHkey := &SSHKeyCreateRequest{
		Name:      "test",
		PublicKey: "testkey",
	}

	result, meta, err := client.SSHKey().Create(ctx, newSSHkey)
	if err != nil {
		t.Errorf("SSHKey.Create returned error: %v", err)
	}

	expectedKeys := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Result: "success"}
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("SSHKey.Create returned keys %+v, expected %+v", result, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.Create returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestSSHKey_Update(t *testing.T) {
	initServer()
	defer downServer()

	updateSSHKey := &SSHKeyUpdateRequest{
		Name: "test",
	}

	mux.HandleFunc("/v2/sshkeys/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		expectedUpdateSSHKey := map[string]interface{}{
			"name": "test",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expectedUpdateSSHKey) {
			t.Errorf("Request body = %#v, expected %#v", v, expectedUpdateSSHKey)
		}

		respondSimple := SSHKey{
			ID:   "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Name: "test",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.SSHKey().Update(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", updateSSHKey)
	if err != nil {
		t.Errorf("SSHKey.Update returned error: %v", err)
	}

	expectedKeys := &SSHKey{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Name: "test"}
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("SSHKey.Update returned keys %+v, expected %+v", result, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.Update returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestSSHKey_Delete(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/sshkeys/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.SSHKey().Delete(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55")
	if err != nil {
		t.Errorf("SSHKey.Delete returned error: %v", err)
	}

	expectedKeys := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Result: "success"}
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("SSHKey.Delete returned keys %+v, expected %+v", result, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.Delete returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}
