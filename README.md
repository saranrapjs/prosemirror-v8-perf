# prosemirror-v8-perf

## What's this about?

ProseMirror's support for collaborative editing requires an "authority server". From the [ProseMirror documentation](https://prosemirror.net/docs/guide/#collab):

> The role of the central authority is actually rather simple. It must...
> - Track a current document version
> - Accept changes from editors, and when these can be applied, add them to its list of changes
> - Provide a way for editors to receive changes since a given version

The "track a current document version" and "Provide a way for editors to receive changes" aspects of what an Authority Server does are relatively runtime agnostic — a server (or even an API), written in any runtime, can perform these without any specific knowledge of ProseMirror.

The "Accept changes from editors" part _does_ require some component to be ProseMirror-aware. ProseMirror's libraries are written in Javascript, targeting the browser and node.js. The [reference implementation](https://github.com/ProseMirror/website/tree/master/src/collab/server) of an authority server stores the current document and steps in-memory, and uses functions exported by ProseMirror libraries to apply steps to a document. This is probably how many authority server implementations work.

This repo tries to evaluate alternatives to running an authority server in node. There may be many reasons to do this — a service already written in an existing runtime, performance considerations surrounding node vs. other language runtimes — that are beyond the scope of this example.

## How to run this?

This repo has a node.js and Go implementation of an HTTP endpoint that receives a ProseMirror document and array of steps, hydrates them using a specific ProseMirror schema, applies the steps to the document, and returns the results. The benchmark scripts start each server, and send requests to them, collecting latency information.

Requires: go, node, and k6 (for load testing).

### `server-node.js`

The node.js implementation imports ProseMirror libraries directly, and is quite simple.

Test using:

```shell
./benchmark-node.sh
```

### `server-go.go`

The Go implementation configures and runs the ProseMirror libraries indirectly, using the [`v8go`](https://github.com/rogchap/v8go) package, which allows Go to create and run v8 isolates (which are tiny Javascript VM's).

Test using:

```shell
./benchmark-go.sh
```

## Results

Node:

```
  scenarios: (100.00%) 1 scenario, 40 max VUs, 1m0s max duration (incl. graceful stop):
           * default: 40 looping VUs for 30s (gracefulStop: 30s)


     data_received..................: 352 kB 12 kB/s
     data_sent......................: 311 kB 10 kB/s
     http_req_blocked...............: avg=73.38µs min=1µs   med=4µs    max=2.35ms  p(90)=9µs    p(95)=14µs   
     http_req_connecting............: avg=38.23µs min=0s    med=0s     max=1.38ms  p(90)=0s     p(95)=0s     
     http_req_duration..............: avg=3.9ms   min=325µs med=3.17ms max=19.39ms p(90)=6.84ms p(95)=11.83ms
       { expected_response:true }...: avg=3.9ms   min=325µs med=3.17ms max=19.39ms p(90)=6.84ms p(95)=11.83ms
     http_req_failed................: 0.00%  ✓ 0         ✗ 1200
     http_req_receiving.............: avg=39.45µs min=7µs   med=32µs   max=455µs   p(90)=66µs   p(95)=83.04µs
     http_req_sending...............: avg=25.12µs min=5µs   med=17µs   max=1.2ms   p(90)=35µs   p(95)=54µs   
     http_req_tls_handshaking.......: avg=0s      min=0s    med=0s     max=0s      p(90)=0s     p(95)=0s     
     http_req_waiting...............: avg=3.84ms  min=298µs med=3.1ms  max=19.31ms p(90)=6.79ms p(95)=11.8ms 
     http_reqs......................: 1200   39.815906/s
     iteration_duration.............: avg=1s      min=1s    med=1s     max=1.02s   p(90)=1s     p(95)=1.01s  
     iterations.....................: 1200   39.815906/s
     vus............................: 40     min=40      max=40
     vus_max........................: 40     min=40      max=40
```

Go:

```
  scenarios: (100.00%) 1 scenario, 40 max VUs, 1m0s max duration (incl. graceful stop):
           * default: 40 looping VUs for 30s (gracefulStop: 30s)

^C
     data_received..................: 23 kB 19 kB/s
     data_sent......................: 21 kB 17 kB/s
     http_req_blocked...............: avg=1.23ms   min=2µs    med=721.5µs max=2.69ms  p(90)=2.6ms   p(95)=2.63ms  
     http_req_connecting............: avg=669.55µs min=0s     med=336.5µs max=1.81ms  p(90)=1.43ms  p(95)=1.48ms  
     http_req_duration..............: avg=5.61ms   min=1.36ms med=3.96ms  max=13.26ms p(90)=12.24ms p(95)=12.85ms 
       { expected_response:true }...: avg=5.61ms   min=1.36ms med=3.96ms  max=13.26ms p(90)=12.24ms p(95)=12.85ms 
     http_req_failed................: 0.00% ✓ 0         ✗ 80  
     http_req_receiving.............: avg=21.78µs  min=7µs    med=16µs    max=154µs   p(90)=40.2µs  p(95)=50µs    
     http_req_sending...............: avg=56.73µs  min=4µs    med=13µs    max=775µs   p(90)=93.3µs  p(95)=411.74µs
     http_req_tls_handshaking.......: avg=0s       min=0s     med=0s      max=0s      p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=5.54ms   min=1.34ms med=3.95ms  max=13.21ms p(90)=11.84ms p(95)=12.82ms 
     http_reqs......................: 80    64.521176/s
     iteration_duration.............: avg=1s       min=1s     med=1s      max=1s      p(90)=1s      p(95)=1s      
     iterations.....................: 40    32.260588/s
     vus............................: 40    min=40      max=40
     vus_max........................: 40    min=40      max=40
```