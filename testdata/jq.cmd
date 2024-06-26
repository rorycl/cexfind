# investigate json with jq

cat terminator.json | jq .results[0].hits[0] | jq keys | egrep "box"
cat terminator.json | jq .results[0].hits[0] | jq '{"boxId","boxName","sellPrice","cashPriceCalculated","exchangePriceCalculated"}'
