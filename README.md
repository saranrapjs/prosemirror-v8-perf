# prosemirror-v8-perf

## What's this about?

ProseMirror's support for collaborative editing requires an "authority server". From the [ProseMirror documentation](https://prosemirror.net/docs/guide/#collab):

> The role of the central authority is actually rather simple. It must...
> - Track a current document version
> - Accept changes from editors, and when these can be applied, add them to its list of changes
> - Provide a way for editors to receive changes since a given version

The "track a current document version" and "Provide a way for editors to receive changes" aspects of what an Authority Server does are relatively runtime agnostic — a server (or even someone else's API), written in any runtime, can perform these without any specific knowledge of ProseMirror.

The "Accept changes from editors" part _does_ require some component to be ProseMirror-aware. ProseMirror's libraries are written in Javascript, targeting the browser and node.js. The [reference implementation](https://github.com/ProseMirror/website/tree/master/src/collab/server) of an authority server stores the current document and steps in-memory in a node.js server, and uses functions exported by ProseMirror libraries to apply steps to a document. This is probably how many authority server implementations work.

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
     http_req_blocked...............: avg=70.88µs min=0s    med=4µs    max=2.31ms  p(90)=7.1µs  p(95)=11.04µs
     http_req_connecting............: avg=38.79µs min=0s    med=0s     max=1.33ms  p(90)=0s     p(95)=0s     
     http_req_duration..............: avg=3.45ms  min=344µs med=3.08ms max=16.55ms p(90)=5.25ms p(95)=6.22ms 
       { expected_response:true }...: avg=3.45ms  min=344µs med=3.08ms max=16.55ms p(90)=5.25ms p(95)=6.22ms 
     http_req_failed................: 0.00%  ✓ 0         ✗ 1200
     http_req_receiving.............: avg=28.4µs  min=7µs   med=25µs   max=287µs   p(90)=45µs   p(95)=56µs   
     http_req_sending...............: avg=17.19µs min=4µs   med=14µs   max=174µs   p(90)=31µs   p(95)=40.04µs
     http_req_tls_handshaking.......: avg=0s      min=0s    med=0s     max=0s      p(90)=0s     p(95)=0s     
     http_req_waiting...............: avg=3.4ms   min=325µs med=3.03ms max=16.5ms  p(90)=5.21ms p(95)=6.2ms  
     http_reqs......................: 1200   39.840717/s
     iteration_duration.............: avg=1s      min=1s    med=1s     max=1.01s   p(90)=1s     p(95)=1s     
     iterations.....................: 1200   39.840717/s
     vus............................: 40     min=40      max=40
     vus_max........................: 40     min=40      max=40
```

Go:

```
  scenarios: (100.00%) 1 scenario, 40 max VUs, 1m0s max duration (incl. graceful stop):
           * default: 40 looping VUs for 30s (gracefulStop: 30s)


     data_received..................: 220 kB 7.3 kB/s
     data_sent......................: 311 kB 10 kB/s
     http_req_blocked...............: avg=71µs    min=1µs   med=4µs    max=2.34ms  p(90)=11µs   p(95)=16µs    
     http_req_connecting............: avg=37.26µs min=0s    med=0s     max=1.46ms  p(90)=0s     p(95)=0s      
     http_req_duration..............: avg=2.64ms  min=199µs med=2.38ms max=13.39ms p(90)=5.04ms p(95)=6.03ms  
       { expected_response:true }...: avg=2.64ms  min=199µs med=2.38ms max=13.39ms p(90)=5.04ms p(95)=6.03ms  
     http_req_failed................: 0.00%  ✓ 0         ✗ 1200
     http_req_receiving.............: avg=40.38µs min=4µs   med=26µs   max=4.12ms  p(90)=49µs   p(95)=82.04µs 
     http_req_sending...............: avg=54.41µs min=6µs   med=16µs   max=6.05ms  p(90)=85µs   p(95)=160.14µs
     http_req_tls_handshaking.......: avg=0s      min=0s    med=0s     max=0s      p(90)=0s     p(95)=0s      
     http_req_waiting...............: avg=2.55ms  min=174µs med=2.29ms max=11.86ms p(90)=4.88ms p(95)=5.84ms  
     http_reqs......................: 1200   39.871942/s
     iteration_duration.............: avg=1s      min=1s    med=1s     max=1.01s   p(90)=1s     p(95)=1s      
     iterations.....................: 1200   39.871942/s
     vus............................: 40     min=40      max=40
     vus_max........................: 40     min=40      max=40
```
