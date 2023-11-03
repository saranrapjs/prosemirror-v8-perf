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
     data_received..................: 88 kB 2.9 kB/s
     data_sent......................: 78 kB 2.6 kB/s
     http_req_blocked...............: avg=55.89µs min=2µs   med=5µs    max=1.55ms  p(90)=12µs   p(95)=16.04µs
     http_req_connecting............: avg=17.67µs min=0s    med=0s     max=631µs   p(90)=0s     p(95)=0s     
     http_req_duration..............: avg=3.59ms  min=437µs med=3.46ms max=12.1ms  p(90)=4.91ms p(95)=7.8ms  
       { expected_response:true }...: avg=3.59ms  min=437µs med=3.46ms max=12.1ms  p(90)=4.91ms p(95)=7.8ms  
     http_req_failed................: 0.00% ✓ 0        ✗ 300 
     http_req_receiving.............: avg=40.75µs min=7µs   med=36µs   max=123µs   p(90)=64µs   p(95)=73µs   
     http_req_sending...............: avg=27.8µs  min=9µs   med=22.5µs max=166µs   p(90)=48µs   p(95)=52µs   
     http_req_tls_handshaking.......: avg=0s      min=0s    med=0s     max=0s      p(90)=0s     p(95)=0s     
     http_req_waiting...............: avg=3.53ms  min=396µs med=3.41ms max=12.08ms p(90)=4.86ms p(95)=7.72ms 
     http_reqs......................: 300   9.959356/s
     iteration_duration.............: avg=1s      min=1s    med=1s     max=1.01s   p(90)=1s     p(95)=1s     
     iterations.....................: 300   9.959356/s
     vus............................: 10    min=10     max=10
     vus_max........................: 10    min=10     max=10
```

Go:

```
     data_received..................: 88 kB 2.9 kB/s
     data_sent......................: 78 kB 2.6 kB/s
     http_req_blocked...............: avg=46.49µs min=2µs   med=6µs    max=1.27ms p(90)=11µs   p(95)=20.04µs
     http_req_connecting............: avg=15.06µs min=0s    med=0s     max=505µs  p(90)=0s     p(95)=0s     
     http_req_duration..............: avg=3.18ms  min=489µs med=3.18ms max=6.87ms p(90)=4.57ms p(95)=4.92ms 
       { expected_response:true }...: avg=3.18ms  min=489µs med=3.18ms max=6.87ms p(90)=4.57ms p(95)=4.92ms 
     http_req_failed................: 0.00% ✓ 0       ✗ 300 
     http_req_receiving.............: avg=42.21µs min=7µs   med=37µs   max=380µs  p(90)=66µs   p(95)=77µs   
     http_req_sending...............: avg=29.61µs min=4µs   med=24µs   max=167µs  p(90)=57.1µs p(95)=66µs   
     http_req_tls_handshaking.......: avg=0s      min=0s    med=0s     max=0s     p(90)=0s     p(95)=0s     
     http_req_waiting...............: avg=3.11ms  min=447µs med=3.11ms max=6.83ms p(90)=4.5ms  p(95)=4.87ms 
     http_reqs......................: 300   9.96363/s
     iteration_duration.............: avg=1s      min=1s    med=1s     max=1s     p(90)=1s     p(95)=1s     
     iterations.....................: 300   9.96363/s
     vus............................: 10    min=10    max=10
     vus_max........................: 10    min=10    max=10
```