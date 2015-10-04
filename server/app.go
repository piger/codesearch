package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/piger/codesearch/index"
	"github.com/piger/codesearch/regexp"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SearchResult struct {
	Filename string
	Match    string
	Line     uint64
}

type SearchOptions struct{}

func searchPattern(idx *index.Index, pattern string, options *SearchOptions, w flushWriter) ([]*SearchResult, error) {
	var results []*SearchResult
	var stdout, stderr bytes.Buffer
	bStdout := bufio.NewWriter(&stdout)
	bStderr := bufio.NewWriter(&stderr)

	grep := regexp.Grep{
		Stdout: bStdout,
		Stderr: bStderr,
		N:      true,
	}

	// grep.AddFlags()
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	grep.Regexp = re
	q := index.RegexpQuery(re.Syntax)
	var post []uint32
	post = idx.PostingQuery(q)

	w.Write([]byte("\"results\": [\n"))

	// This is needed to check whether we need to print a "," to respect JSON
	// syntax
	isFirst := true

	for _, fileid := range post {
		name := idx.Name(fileid)
		grep.File(name)
		bStdout.Flush()
		bStderr.Flush()

		for {
			line, err := stdout.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Error from ReadString: %s", err.Error())
				return nil, err
			}
			// log.Printf("LINE: %s", line)

			fields := strings.SplitN(line, ":", 3)
			ln, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				log.Printf("Error converting line number: %s", err.Error())
				return nil, err
			}
			sr := &SearchResult{
				Filename: fields[0],
				Line:     ln,
				Match:    html.EscapeString(fields[2]),
			}
			results = append(results, sr)
			jr, err := json.Marshal(sr)
			if err != nil {
				log.Printf("JSON error: %s", err.Error())
				return nil, err
			}
			if !isFirst {
				w.Write([]byte(",\n"))
			} else {
				isFirst = false
			}
			w.Write(jr)
		}
	}

	w.Write([]byte("]\n"))

	return results, nil
}

type SearchQuery struct {
	Query string
}

type ResponseError struct {
	Errors []string `json:"errors"`
}

// Serves the app
func indexHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	ctx := TemplateContext{}
	a.RenderTemplate("index.html", ctx, w)
	return 200, nil
}

// API for search
func searchHandler(a *appContext, w http.ResponseWriter, r *http.Request) (int, error) {
	r.ParseForm()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 2048))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err := r.Body.Close(); err != nil {
		return http.StatusInternalServerError, err
	}
	var sq SearchQuery
	if err := json.Unmarshal(body, &sq); err != nil {
		// code 422 ?
		return http.StatusInternalServerError, err
	}

	// validators
	lq := len(sq.Query)
	if lq < 2 || lq > 512 {
		// TODO
		errors := []string{"Invalid query length"}
		rerr := &ResponseError{errors}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		jerr, err := json.Marshal(rerr)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		log.Printf("jerr: %v", string(jerr))
		w.Write(jerr)
		return 200, nil
	}

	// Use non buffered Write for this HTTP connection
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintln(&fw, "{")

	// Lock access to Index
	a.Lock.Lock()
	defer a.Lock.Unlock()

	_, err = searchPattern(a.idx, sq.Query, &SearchOptions{}, fw)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return http.StatusInternalServerError, err
	}

	fmt.Fprintln(&fw, "}")
	return 200, nil
}
