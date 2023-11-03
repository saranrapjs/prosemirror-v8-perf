package main

import (
	"encoding/json"
	"fmt"
	"os"

	"rogchap.com/v8go"
	"net/http"
)

type (
	Steps = json.RawMessage
	Doc   = json.RawMessage
)

type Payload struct {
	Steps Steps `json:"steps"`
	Doc   Doc   `json:"doc"`
}

type Response struct {
	Doc Doc `json:"doc"`
}

func applySteps(ctx *v8go.Context, doc Doc, steps Steps) (Doc, error) {
	obj := ctx.Global()
	obj.Set("doc", string(doc))
	obj.Set("steps", string(steps))
	result, err := ctx.RunScript("applySteps(doc, steps)", "execution.js")
	if err != nil {
		return doc, err
	}
	return Doc(result.String()), nil
}

func main() {
	file, _ := os.ReadFile("prosemirror-go-out.js")
	iso := v8go.NewIsolate()
	defer iso.Dispose()
	ctx := v8go.NewContext(iso)
	defer ctx.Close()
	_, err := ctx.RunScript(string(file), "script.js")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON payload into a Payload struct
		var payload Payload
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
			return
		}

		result, err := applySteps(ctx, payload.Doc, payload.Steps)
		if err != nil {
			http.Error(w, "error with execution", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(&Response{ Doc: result })
	})
	fmt.Println("go server listening...")
	http.ListenAndServe(":8080", nil)
}
