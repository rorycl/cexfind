package location

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetStores(t *testing.T) {

	testdata, err := os.ReadFile("testdata/stores.json")
	if err != nil {
		t.Fatal(err)
	}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(testdata))
	}))

	stores := newStores(nil)
	if stores.isInitialised() {
		t.Fatal("initialisation should have failed")
	}

	stores = newStores(nil)
	stores.URL = svr.URL // repoint url to test server

	err = stores.initialise()
	if err != nil {
		t.Fatal(err)
	}
	if !stores.isInitialised() {
		t.Fatal("initialisation should be ok")
	}

	if got, want := stores.length(), 4; got != want {
		t.Errorf("got %d want %d stores", got, want)
		for i, s := range stores.list() {
			fmt.Printf("%5d : %v\n", i, *s)
		}
	}

	w1, ok := stores.get("London - W1 Rathbone Place")
	if !ok {
		t.Error("expected value for London - W1 Rathbone Place")
		return
	}
	if got, want := w1.StoreID, 2; got != want {
		t.Errorf("got %d want %d storeid", got, want)
	}

	_, ok = stores.get("London W1 Rathbone")
	if !ok {
		t.Error("expected alias value for London - W1 Rathbone Place")
		return
	}

	stores.periodicallyReinitialise()

	stores.Lock()
	stores.initialised = false
	stores.update.Reset(time.Millisecond * 40)
	stores.Unlock()
	time.Sleep(time.Millisecond * 50)
	if !stores.isInitialised() {
		t.Error("re-initialisation failed")
	}

}
