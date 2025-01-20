package location

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetStores(t *testing.T) {

	testdata, err := os.ReadFile("testdata/stores.json")
	if err != nil {
		t.Fatal(err)
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(testdata))
	}))

	// repoint url
	storeURL = svr.URL

	err = getStoreLocations()
	if err != nil {
		t.Fatal(err)
	}

	for k := range Stores {
		fmt.Println(k)
	}

	if got, want := len(Stores), 4; got != want {
		t.Errorf("got %d want %d stores", got, want)
	}

	w1, ok := Stores["London - W1 Rathbone Place"]
	if !ok {
		t.Error("expected value for London - W1 Rathbone Place")
		return
	}
	if got, want := w1.StoreID, 2; got != want {
		t.Errorf("got %d want %d storeid", got, want)
	}

	_, ok = Stores["London W1 Rathbone"]
	if !ok {
		t.Error("expected alias value for London - W1 Rathbone Place")
		return
	}

}
