node server-node.js &
pid=$!
curl --head -X GET --retry 5 --retry-connrefused --retry-delay 1 "http://localhost:8080" > /dev/null 2>&1
echo "server ready"
k6 run --vus 40 --duration 30s k6.js
kill $pid
