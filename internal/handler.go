package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/yezihack/github-webhook/logger"
	"github.com/yezihack/github-webhook/util"
)

type WebHookHandler func(eventName string, payload *GitHubRepo, req *http.Request) error

var log logger.Logger

func Handler(secret string, l logger.Logger, fn WebHookHandler) http.HandlerFunc {
	log = l
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			util.Response(w, 405, "Method Not Allowed")
			return
		}
		fmt.Println(r.Header)
		event := r.Header.Get("x-gitlab-event")
		token := r.Header.Get("x-gitlab-token")
		contentType := r.Header.Get("Content-Type")
		log.Printf("event:%s, token:%s, contentType:%s \n", event, token, contentType)
		// Utility funcs
		_fail := func(err error) {
			fail(w, event, err)
		}
		_succeed := func() {
			succeed(w, event)
		}

		// Ensure headers are all there
		if event == "" || token == "" {
			_fail(fmt.Errorf("missing x-gitlab-event and x-gitlab-token headers"))
			return
		}

		// Read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Verbose("Read request body:%s\n", string(body))
			_fail(err)
			return
		}
		fmt.Println("body", string(body))
		log.Verbose("RequestBody:%s\n", string(body))
		// Validate payload (only when secret is provided)
		if secret != token {
			_fail(fmt.Errorf("token can not match"))
			return
		}
		repo := GitHubRepo{}
		result := gjson.ParseBytes(body)
		repo.Name = result.Get("repository.name").Str
		repo.FullName = result.Get("repository.full_name").Str
		repo.CloneURL = result.Get("repository.clone_url").Str
		repo.CommitID = result.Get("head_commit.id").Str
		repo.CommitName = result.Get("head_commit.committer.name").Str
		repo.CommitEmail = result.Get("head_commit.committer.email").Str
		repo.CommitAt = result.Get("head_commit.timestamp").Time().Format("2006-01-02 15:04:05")
		repo.BranchName = result.Get("ref").Str

		// Do something with payload
		if err := fn(event, &repo, r); err == nil {
			_succeed()
		} else {
			_fail(err)
		}
	}
}

func validePayloadSignature(secret, signatureHeader string, body []byte) error {
	// Check header is valid
	signature_parts := strings.SplitN(signatureHeader, "=", 2)
	if len(signature_parts) != 2 {
		return fmt.Errorf("invalid signature header: '%s' does not contain two parts (hash type and hash)", signatureHeader)
	}

	// Ensure secret is a sha1 hash
	signature_type := signature_parts[0]
	signature_hash := signature_parts[1]
	if signature_type != "sha1" {
		return fmt.Errorf("signature should be a 'sha1' hash not '%s'", signature_type)
	}

	// Check that payload came from github
	// skip check if empty secret provided
	if !IsValidPayload(secret, signature_hash, body) {
		return fmt.Errorf("payload did not come from GitHub")
	}
	return nil
}

func succeed(w http.ResponseWriter, event string) {
	log.Verbose("http-code:%d, event:%s\n", 200, event)
	w.WriteHeader(200)
	render(w, PayloadPong{
		Ok:    true,
		Event: event,
	})
}

func fail(w http.ResponseWriter, event string, err error) {
	log.Verbose("http-code:%d, event:%s, err:%s\n", 500, event, err)
	w.WriteHeader(500)
	render(w, PayloadPong{
		Ok:    false,
		Event: event,
		Error: err.Error(),
	})
}
func render(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(data)
}
