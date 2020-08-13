package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hoisie/mustache"
	"github.com/mitchellh/go-homedir"
)

const (
	appName        = "rest-client"
	defaultEnvFile = "rest-client.env.json"
	httpFileExt    = ".http"
	reqNamePrefix  = "#:name"
	reqDescPrefix  = "#:desc"
)

var (
	// global options
	cfgFile   string
	envsFile  string
	httpFiles []string
	verbose   bool

	// exec command options
	envName  string
	reqNames []string

	// http methods supported
	validMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
)

// Envs is a mapping from environment name to its variables.
type Envs = map[string]Vars

// Vars is a mapping from environment variable name to its value.
type Vars = map[string]string

func parseEnvs(r io.Reader) (Envs, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var envs Envs
	err = json.Unmarshal(buf, &envs)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

// Req is data ready from a request in an http requests file.
type Req struct {
	File    string
	Name    string
	Desc    string
	Method  string
	URL     string
	Headers []string
	Body    []string
}

// Res is data returned from executing a request.
type Res struct {
	Req Req
	Res *http.Response
}

func (r Req) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("File: %s\n", r.File))
	sb.WriteString(fmt.Sprintf("Name: %s\n", r.Name))
	sb.WriteString(fmt.Sprintf("Desc: %s\n\n", r.Desc))
	sb.WriteString(fmt.Sprintf("%s %s\n", r.Method, r.URL))
	for _, header := range r.Headers {
		sb.WriteString(fmt.Sprintf("%s\n", header))
	}
	if len(r.Body) > 0 {
		sb.WriteString("\n")
		for _, line := range r.Body {
			sb.WriteString(fmt.Sprintf("%s\n", line))
		}
	}
	return sb.String()
}

// Execute executes the http request and returns a response.
func (r Req) Execute(env map[string]string) (*Res, error) {
	// Build URL
	url := expandString(r.URL, env)

	// Build body.
	var sb strings.Builder
	for _, line := range r.Body {
		line = expandString(line, env)
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	body := strings.NewReader(sb.String())

	// Create request.
	request, err := http.NewRequest(r.Method, url, body)

	// Add request headers.
	for _, line := range r.Headers {
		line = expandString(line, env)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(line, ":")
		k := strings.Trim(parts[0], " ")
		v := strings.Trim(parts[1], " ")
		request.Header.Set(k, v)
	}

	// Execute request and return response.
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	var response *http.Response
	if response, err = client.Do(request); err != nil {
		return nil, err
	}

	// Success!
	return &Res{Req: r, Res: response}, nil
}

func newReq(name string, file string) *Req {
	return &Req{
		Name: name,
		File: file,
	}
}

func loadReqs(paths []string) ([]*Req, error) {
	if len(paths) == 0 {
		// Default to looking in all *.http files in the working directory.
		cwd, _ := os.Getwd()
		paths, _ = listReqFiles(cwd)
	}

	allReqs := make([]*Req, 0)
	for _, path := range paths {
		path, err := homedir.Expand(path)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reqs, err := parseReqs(f, path)
		if err != nil {
			return nil, err
		}
		allReqs = append(allReqs, reqs...)
	}
	return allReqs, nil
}

func parseReqs(r io.Reader, path string) ([]*Req, error) {
	reqs := make([]*Req, 0)
	var req *Req

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, reqNamePrefix) && req == nil {
			// Start building a new Req object.
			name := strings.Split(line, " ")[1]
			req = newReq(name, path)
		} else if strings.HasPrefix(line, reqDescPrefix) && req != nil {
			// Update Req object with description.
			line = strings.TrimPrefix(line, reqDescPrefix)
			req.Desc = strings.TrimSpace(line)
		} else if isReq(line) && req != nil {
			// Update Req object with URL method and path.
			parts := strings.Split(line, " ")
			req.Method = parts[0]
			req.URL = parts[1]
		} else if len(strings.Split(line, ":")) == 2 && req != nil && req.Method != "" {
			// Update Req object with header.
			req.Headers = append(req.Headers, line)
		} else if line == "" {
			// Skip blank lines.
			continue
		} else if strings.HasPrefix(line, "#") && req != nil {
			// Finish building new Req object and add to Reqs to return.
			reqs = append(reqs, req)
			req = nil
		} else if req != nil {
			// Update Request body with line.
			req.Body = append(req.Body, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if req != nil {
		// Handle case where file does not contain "###" finalizer sigil.
		reqs = append(reqs, req)
	}

	return reqs, nil
}

func isReq(line string) bool {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return false
	}

	method := parts[0]
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

func filterReqs(reqs []*Req, names []string) []*Req {
	filtered := make([]*Req, 0)
	for _, req := range reqs {
		if containsString(names, req.Name) {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

func execReqs(reqs []*Req, env map[string]string) ([]*Res, error) {
	responses := make([]*Res, 0, len(reqs))
	for _, req := range reqs {
		response, err := req.Execute(env)
		if err != nil {
			return responses, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func containsString(ss []string, s string) bool {
	for _, _s := range ss {
		if _s == s {
			return true
		}
	}
	return false
}

func expandString(s string, env map[string]string) string {
	return mustache.Render(s, env)
}

func renderResponses(responses []*Res) string {
	var sb strings.Builder
	for i, response := range responses {
		sb.WriteString(renderResponse(response))
		if i < len(responses) - 1 {
			sb.WriteString("\n###\n\n")
		}
	}
	return sb.String()
}

func renderResponse(response *Res) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n", response.Res.Proto, response.Res.Status))

	var headers []string
	for k, v := range response.Res.Header {
		headers = append(headers, fmt.Sprintf("%s: %s\n", k, v))
	}
	sort.Strings(headers)
	for _, header := range headers {
		sb.WriteString(header)
	}

	defer response.Res.Body.Close()
	buf, err := ioutil.ReadAll(response.Res.Body)
	if err != nil {
		panic(err)
	}

	if len(buf) > 0 {
		sb.WriteString(fmt.Sprintf("\n%s", string(buf)))
	}
	return sb.String()
}

func listReqFiles(dir string) ([]string, error) {
	var files []string
	if err := filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == httpFileExt {
				files = append(files, f.Name())
			}
		}
		return nil
	}); err != nil {
		return files, err
	}
	return files, nil
}