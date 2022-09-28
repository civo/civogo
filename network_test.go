package civogo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestNetwork_List(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allNetworks := []Network{
			{
				ID:      "0facdbcf-4bd9-4fe4-ae5a-92214d1073af",
				Name:    "test-network",
				Default: true,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network",
				Status:  "active",
			},
			{
				ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
				Name:    "test-network-2",
				Default: false,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network 2",
				Status:  "active",
			},
		}
		value := toJSON(t, allNetworks)
		fmt.Fprint(w, value)
	})

	net, meta, err := client.Network.List(ctx)
	if err != nil {
		t.Errorf("Network.List returned error: %v", err)
	}

	expectedNetwork := []Network{{
		ID:      "0facdbcf-4bd9-4fe4-ae5a-92214d1073af",
		Name:    "test-network",
		Default: true,
		CIDR:    "192.168.1.1/24",
		Label:   "Test Network",
		Status:  "active",
	},
		{
			ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
			Name:    "test-network-2",
			Default: false,
			CIDR:    "192.168.1.1/24",
			Label:   "Test Network 2",
			Status:  "active",
		}}
	if !reflect.DeepEqual(net, expectedNetwork) {
		t.Errorf("Network.List returned keys %+v, expected %+v", net, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.List returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_GetDefault(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allNetworks := []Network{
			{
				ID:      "0facdbcf-4bd9-4fe4-ae5a-92214d1073af",
				Name:    "test-network",
				Default: true,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network",
				Status:  "active",
			},
			{
				ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
				Name:    "test-network-2",
				Default: false,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network 2",
				Status:  "active",
			},
		}
		value := toJSON(t, allNetworks)
		fmt.Fprint(w, value)
	})

	net, meta, err := client.Network.GetDefault(ctx)
	if err != nil {
		t.Errorf("Network.GetDefault returned error: %v", err)
	}

	expectedNetwork := &Network{
		ID:      "0facdbcf-4bd9-4fe4-ae5a-92214d1073af",
		Name:    "test-network",
		Default: true,
		CIDR:    "192.168.1.1/24",
		Label:   "Test Network",
		Status:  "active",
	}
	if !reflect.DeepEqual(net, expectedNetwork) {
		t.Errorf("Network.GetDefault returned keys %+v, expected %+v", net, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.GetDefault returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_GetByID(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks/21c4f72a-3dec-4940-81fd-8ca4af0b3c0e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allNetworks := Network{
			ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
			Name:    "test-network",
			Default: true,
			CIDR:    "192.168.1.1/24",
			Label:   "Test Network",
			Status:  "active",
		}
		value := toJSON(t, allNetworks)
		fmt.Fprint(w, value)
	})

	net, meta, err := client.Network.GetByID(ctx, "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e")
	if err != nil {
		t.Errorf("Network.GetByID returned error: %v", err)
	}

	expectedNetwork := &Network{
		ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
		Name:    "test-network",
		Default: true,
		CIDR:    "192.168.1.1/24",
		Label:   "Test Network",
		Status:  "active",
	}
	if !reflect.DeepEqual(net, expectedNetwork) {
		t.Errorf("Network.GetByID returned keys %+v, expected %+v", net, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.GetByID returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_Find(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		allNetworks := []Network{
			{
				ID:      "0facdbcf-4bd9-4fe4-ae5a-92214d1073af",
				Name:    "test-network",
				Default: true,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network",
				Status:  "active",
			},
			{
				ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
				Name:    "test-network-2",
				Default: false,
				CIDR:    "192.168.1.1/24",
				Label:   "Test Network 2",
				Status:  "active",
			},
		}
		value := toJSON(t, allNetworks)
		fmt.Fprint(w, value)
	})

	net, meta, err := client.Network.Find(ctx, "test-network-2")
	if err != nil {
		t.Errorf("Network.Find returned error: %v", err)
	}

	expectedNetwork := &Network{
		ID:      "21c4f72a-3dec-4940-81fd-8ca4af0b3c0e",
		Name:    "test-network-2",
		Default: false,
		CIDR:    "192.168.1.1/24",
		Label:   "Test Network 2",
		Status:  "active",
	}
	if !reflect.DeepEqual(net, expectedNetwork) {
		t.Errorf("Network.Find returned keys %+v, expected %+v", net, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.Find returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_Create(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Name:   "test-network",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	newNetwork := &NetworkCreateRequest{
		Label:  "test-network",
		Region: "TEST",
	}

	result, meta, err := client.Network.Create(ctx, newNetwork)
	if err != nil {
		t.Errorf("SSHKey.Create returned error: %v", err)
	}

	expectedKeys := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Name: "test-network", Result: "success"}
	if !reflect.DeepEqual(result, expectedKeys) {
		t.Errorf("SSHKey.Create returned keys %+v, expected %+v", result, expectedKeys)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("SSHKey.Create returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_Update(t *testing.T) {
	initServer()
	defer downServer()

	updateNetwork := &NetworkUpdateRequest{
		Label: "test-update-network",
		Region: "TEST",
	}

	mux.HandleFunc("/v2/networks/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		expectedUpdateNetwork := map[string]interface{}{
			"label": "test-update-network",
			"region": "TEST",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expectedUpdateNetwork) {
			t.Errorf("Request body = %#v, expected %#v", v, expectedUpdateNetwork)
		}

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Name:   "test-update-network",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.Network.Update(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", updateNetwork)
	if err != nil {
		t.Errorf("Network.Update returned error: %v", err)
	}

	expectedNetwork := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Name: "test-update-network", Result: "success"}
	if !reflect.DeepEqual(result, expectedNetwork) {
		t.Errorf("Network.Update returned keys %+v, expected %+v", result, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.Update returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}

func TestNetwork_Delete(t *testing.T) {
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/networks/78f64e5c-abd3-4f4d-85c8-ac63b50caa55", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		respondSimple := SimpleResponse{
			ID:     "78f64e5c-abd3-4f4d-85c8-ac63b50caa55",
			Result: "success",
		}
		value := toJSON(t, respondSimple)
		fmt.Fprint(w, value)
	})

	result, meta, err := client.Network.Delete(ctx, "78f64e5c-abd3-4f4d-85c8-ac63b50caa55")
	if err != nil {
		t.Errorf("Network.Delete returned error: %v", err)
	}

	expectedNetwork := &SimpleResponse{ID: "78f64e5c-abd3-4f4d-85c8-ac63b50caa55", Result: "success"}
	if !reflect.DeepEqual(result, expectedNetwork) {
		t.Errorf("Network.Delete returned keys %+v, expected %+v", result, expectedNetwork)
	}

	// compare status code
	if meta.StatusCode != http.StatusOK {
		t.Errorf("Network.Delete returned status code %d, expected %d", meta.StatusCode, http.StatusOK)
	}
}
