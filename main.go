package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bobmaertz/test-lsp/pkg/analysis"
	"github.com/bobmaertz/test-lsp/pkg/lsp"
	"github.com/bobmaertz/test-lsp/pkg/rpc"
)

const name = "educational-lsp" 

func main() {
	l := getLogger("/Users/bob/personal/education-lsp/out.log")
	l.Println("Starting LSP")

	state := analysis.NewState()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitFunc)

	for scanner.Scan() {
		method, contents, err := rpc.DecodeMessage(scanner.Bytes())
		if err != nil {
			l.Printf("error decoding message: %v\n", err)
			continue
		}
		handleMessage(l, state, method, contents)
	}
}

func handleMessage(l *log.Logger, state analysis.State, method string, contents []byte) {
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			l.Printf("unable to unmarshal initialize request: %v\n", err)
			return
		}
		l.Printf("connected to client %s - %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		response := lsp.NewInitializeResponse(1)
		out := rpc.EncodeMessage(response)
		fmt.Print(out)
	case "textDocument/didOpen":
		var notification lsp.DidOpenNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			l.Printf("unable to unmarshal textDocument/didOpen notification: %v\n", err)
			return
		}
		l.Printf("didOpen> %v\n", notification.Params.TextDocument.Uri)
		state.OpenDocument(notification.Params.TextDocument.Uri, notification.Params.TextDocument.Text)
	case "textDocument/didChange":
        var notification lsp.DidChangeNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			l.Printf("unable to unmarshal textDocument/didChange notification: %v\n", err)
			return
		}

		l.Printf("didChange> %v\n", notification.Params.TextDocument.Uri)
        for _, change := range notification.Params.ContentChanges {
	    	state.UpdateDocument(notification.Params.TextDocument.Uri, change.Text)
        }
    case "textDocument/willSave": 
        l.Printf("will Save: %v", string(contents)) 
    case "textDocument/didSave":
        l.Printf("did Save: %v", string(contents)) 
    case "textDocument/formatting":
        l.Printf("formatting: %v", "<>") 
	case "textDocument/completion":
		var request lsp.TextCompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			l.Printf("unable to unmarshal textdocument/completion request: %v\n", err)
			return
		}
		response := lsp.NewTextCompletionResponse(request.Id)
		out := rpc.EncodeMessage(response)
		fmt.Print(out)
	default:
		l.Printf("received method: %s, message: %s\n", method, contents)
        l.Printf("state: %v", state)
	}
}

func getLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic("error opening log file")
	}

	return log.New(file, "["+name+"] ", log.Ldate|log.Ltime|log.Lshortfile)
}
