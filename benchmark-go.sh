npm run build-go-js
go build -o go-server-binary server-go.go
./go-server-binary &
pid=$!
curl --head -X GET --retry 10 --retry-connrefused --retry-delay 5 "http://localhost:8080" > /dev/null 2>&1
echo "server ready"
k6 run --vus 40 --duration 30s k6.js
kill $pid
