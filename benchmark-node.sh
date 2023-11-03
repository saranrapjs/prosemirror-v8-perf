node server-node.js &
pid=$!
curl --head -X GET --retry 5 --retry-connrefused --retry-delay 1 "http://localhost:8080" 2>&1 /dev/null
echo "server ready"
k6 run --vus 10 --duration 30s k6.js
