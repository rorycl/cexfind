package location

import (
	"fmt"
	"math"
	"testing"
)

func TestStoreDistances(t *testing.T) {

	// Royal Armories Museum, Leeds, S10 1LT as start point
	// Results expected to be
	// []location.StoreWithDistance{
	//		location.StoreWithDistance{StoreID:145, StoreName:"Walthamstow", RegionName:"London and the South-East of England", Latitude:51.583371, Longitude:-0.023809, DistanceMiles:139.1034687246046},
	//		location.StoreWithDistance{StoreID:193, StoreName:"Woolwich", RegionName:"London and the South-East of England", Latitude:51.491979, Longitude:0.064665, DistanceMiles:146.44476176156368},
	//		location.StoreWithDistance{StoreID:3058, StoreName:"Havant", RegionName:"London and the South-East of England", Latitude:50.852325, Longitude:-0.982041, DistanceMiles:176.34720204432986}}
	//	}
	sd, err := StoreDistances("S10 1LT", []string{"Walthamstow", "Woolwich", "Havant"})
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

	// reset stores
	Stores = stores{}
}
