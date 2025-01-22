package location

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStoreDistances(t *testing.T) {

	// Royal Armories Museum, Leeds, S10 1LT as start point
	// Results expected to be
	// []location.StoreWithDistance{
	//		location.StoreWithDistance{StoreID:145, StoreName:"Walthamstow", RegionName:"London and the South-East of England", Latitude:51.583371, Longitude:-0.023809, DistanceMiles:139.1034687246046},
	//		location.StoreWithDistance{StoreID:193, StoreName:"Woolwich", RegionName:"London and the South-East of England", Latitude:51.491979, Longitude:0.064665, DistanceMiles:146.44476176156368},
	//		location.StoreWithDistance{StoreID:3058, StoreName:"Havant", RegionName:"London and the South-East of England", Latitude:50.852325, Longitude:-0.982041, DistanceMiles:176.34720204432986}}
	//	}
	nsd := NewStoreDistances()
	sd, err := nsd.Distances("S10 1LT", []string{"Walthamstow", "Woolwich", "Havant"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(sd), 3; got != want {
		fmt.Printf("got %d want %d results", got, want)
	}
	if got, want := sd[len(sd)-1].StoreName, "Havant"; got != want {
		fmt.Printf("got %s want %s as furtherest store", got, want)
	}
	if got, want := int(math.Round(sd[len(sd)-1].DistanceMiles)), 176; got != want {
		fmt.Printf("got %d want %d for Havent distance", got, want)
	}
	// fmt.Printf("%#v\n", sd)

}

func TestPrintStoredDistances(t *testing.T) {

	tests := []struct {
		StoreID       int
		StoreName     string
		DistanceMiles float64
		expected      string
	}{
		{
			StoreID:       0,
			StoreName:     "store 0",
			DistanceMiles: 9999.9999,
			expected:      "store 0",
		},
		{
			StoreID:       1,
			StoreName:     "store 1",
			DistanceMiles: 1.22567,
			expected:      "store 1 (1.2mi)",
		},
		{
			StoreID:       2,
			StoreName:     "store two",
			DistanceMiles: 3.5111,
			expected:      "store two (3.5mi)",
		},
		{
			StoreID:       3,
			StoreName:     "store Three",
			DistanceMiles: 10.9999,
			expected:      "store Three (11mi)",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			swd := StoreWithDistance{
				StoreID:       tt.StoreID,
				StoreName:     tt.StoreName,
				DistanceMiles: tt.DistanceMiles,
			}
			output := fmt.Sprint(swd)
			if got, want := output, tt.expected; got != want {
				t.Errorf("got %s want %s", got, want)
			}
		})
	}
}

func TestStoreSorting(t *testing.T) {

	tests := []struct {
		swd               []StoreWithDistance
		expectedNameSlice []string
	}{
		{
			swd: []StoreWithDistance{
				StoreWithDistance{StoreName: "z"},
				StoreWithDistance{StoreName: "a"},
				StoreWithDistance{StoreName: "A"},
			},
			expectedNameSlice: []string{"A", "a", "z"},
		},
		{
			swd: []StoreWithDistance{
				StoreWithDistance{StoreName: "zzz"},
				StoreWithDistance{StoreName: "1"},
				StoreWithDistance{StoreName: ""},
			},
			expectedNameSlice: []string{"", "1", "zzz"},
		},
		{
			swd: []StoreWithDistance{
				StoreWithDistance{StoreName: "z", DistanceMiles: 1},
				StoreWithDistance{StoreName: "a", DistanceMiles: 2},
				StoreWithDistance{StoreName: "A", DistanceMiles: 3},
			},
			expectedNameSlice: []string{"z", "a", "A"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			storeSorter(tt.swd)
			got := []string{}
			for _, g := range tt.swd {
				got = append(got, g.StoreName)
			}
			want := tt.expectedNameSlice
			if diff := cmp.Diff(
				got,
				want,
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
