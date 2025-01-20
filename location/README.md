# Location

Add locations for stores in Cexfind.

A user on reddit suggested adding location searching, to allow the list
of stores holding item to show how far they are away from a known
location.

## Location of Cex stores

Cex holds a list of stores with their geolocation data at

https://wss2.cex.uk.webuy.io/v3/stores

This provides a list of stores as follows:

```
{
  "response": {
    "ack": "Success",
    "data": {
      "stores": [
        {
          "storeId": 1,
          "storeName": "London - W1 Tottenham Crt Rd",
          "regionName": "London and the South-East of England",
          "latitude": 51.520383,
          "longitude": -0.134501,
          "regionImageUrls": null,
          "phoneNumber": null,
          "closingTime": "18:00"
        },
        {
          "storeId": 2,
          "storeName": "London - W1 Rathbone Place",
          "regionName": "London and the South-East of England",
          "latitude": 51.51764,
          "longitude": -0.134483,
          "regionImageUrls": null,
          "phoneNumber": null,
          "closingTime": "18:00"
        },
      ]
    },
    "error": {
      "code": "",
      "internal_message": "",
      "moreInfo": []
    }
  }
}
```

The fields of interest are storeId (for joining to results for
computers), 

## Where are you?

There is a wonderful opensource (and free) service for postcodes
called https://postcodes.io/. 

api.postcodes.io/postcodes?q=<postcode>

The results are as follows:

```
{
    "status": 200,
    "result": [
        {
            "postcode": "NW1 6LG",
            "quality": 1,
            "eastings": 527310,
            "northings": 182155,
            "country": "England",
            "nhs_ha": "London",
            "longitude": -0.166312,
            "latitude": 51.523969,
            "european_electoral_region": "London",
            "primary_care_trust": "Westminster",
            "region": "London",
            "lsoa": "Westminster 008A",
            "msoa": "Westminster 008",
            "incode": "6LG",
            "outcode": "NW1",
            "parliamentary_constituency": "Cities of London and Westminster",
            "parliamentary_constituency_2024": "Cities of London and Westminster",
            "admin_district": "Westminster",
            "parish": "Westminster, unparished area",
            "admin_county": null,
            "date_of_introduction": "198001",
            "admin_ward": "Regent's Park",
            "ced": null,
            "ccg": "NHS North West London",
            "nuts": "Westminster",
            "pfa": "Metropolitan Police",
            "codes": {
                "admin_district": "E09000033",
                "admin_county": "E99999999",
                "admin_ward": "E05013805",
                "parish": "E43000236",
                "parliamentary_constituency": "E14001172",
                "parliamentary_constituency_2024": "E14001172",
                "ccg": "E38000256",
                "ccg_id": "W2U3Z",
                "ced": "E99999999",
                "nuts": "TLI32",
                "lsoa": "E01004659",
                "msoa": "E02000967",
                "lau2": "E09000033",
                "pfa": "E23000001"
            }
        }
    ]
}
```

The fields of interest are "longitude", "latitude", "admin_district".

Another potential source of geolocation data is nominatim. Eg

https://nominatim.openstreetmap.org/search?q=NW1%206LG&format=geojson.

## Calculation of distance

An invaluable resource is https://www.movable-type.co.uk/scripts/latlong.html which provides formulae for the haversine function and simpler spherical law of cosines function for calculating distance.

For a haversine formula in js, Chris Veness suggests:

```js
const R = 6371e3; // metres
const φ1 = lat1 * Math.PI/180; // φ, λ in radians
const φ2 = lat2 * Math.PI/180;
const Δφ = (lat2-lat1) * Math.PI/180;
const Δλ = (lon2-lon1) * Math.PI/180;

const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
          Math.cos(φ1) * Math.cos(φ2) *
          Math.sin(Δλ/2) * Math.sin(Δλ/2);
const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));

const d = R * c; // in metres
```

An example implementation of the haversine function in go is at https://github.com/umahmood/haversine/blob/master/haversine.go
