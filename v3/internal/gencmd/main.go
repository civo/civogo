// Command gencmd generates the typed service packages of civogo/v3 from
// a Civo v3 OpenAPI document (as emitted by grpc-gateway from the private
// proto repos). It is intentionally NOT a general OpenAPI generator: it
// relies on the uniform shapes grpc-gateway produces.
//
// Usage:
//
//	go run ./internal/gencmd -spec path/to/compute.openapi3.json -proto-commit <sha>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"sort"
	"strings"
)

// serviceDef declares one generated service package. Adding a service to
// a spec means adding a row here.
type serviceDef struct {
	pkg        string   // Go package name and output directory
	entity     string   // primary entity, e.g. "Instance"
	pathPrefix string   // spec paths handled by this service
	strip      []string // tokens removed from operation names
}

var services = []serviceDef{
	{
		pkg:        "instances",
		entity:     "Instance",
		pathPrefix: "/v3alpha1/compute/instances",
		strip:      []string{"ToInstance", "FromInstance", "Instances", "Instance"},
	},
}

// reasonsSchema is the spec enum listing the ErrorInfo reasons the API
// reports. Its values are generated as string constants into reasonsPkg,
// for use with civogo.HasReason. Specs without it generate no reasons file.
const (
	reasonsSchema = "v3alpha1ErrorReason"
	reasonsPkg    = "compute"
)

type spec struct {
	Paths      map[string]map[string]*operation `json:"paths"`
	Components struct {
		Schemas map[string]*schema `json:"schemas"`
	} `json:"components"`
}

type operation struct {
	OperationID string       `json:"operationId"`
	Parameters  []*parameter `json:"parameters"`
	RequestBody *struct {
		Content map[string]struct {
			Schema *schema `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
	Responses map[string]struct {
		Content map[string]struct {
			Schema *schema `json:"schema"`
		} `json:"content"`
	} `json:"responses"`
}

type parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Required    bool    `json:"required"`
	Description string  `json:"description"`
	Schema      *schema `json:"schema"`
}

type schema struct {
	Ref                  string             `json:"$ref"`
	Type                 string             `json:"type"`
	Format               string             `json:"format"`
	Description          string             `json:"description"`
	Enum                 []string           `json:"enum"`
	Nullable             bool               `json:"nullable"` // set on proto3 optional fields
	Items                *schema            `json:"items"`
	Properties           map[string]*schema `json:"properties"`
	AdditionalProperties *schema            `json:"additionalProperties"`
}

func (s *schema) refName() string {
	if s == nil || s.Ref == "" {
		return ""
	}
	return s.Ref[strings.LastIndex(s.Ref, "/")+1:]
}

var initialisms = map[string]string{
	"id": "ID", "ip": "IP", "ipv4": "IPv4", "ipv6": "IPv6", "url": "URL",
	"ssh": "SSH", "dns": "DNS", "gpu": "GPU", "cpu": "CPU", "ram": "RAM",
	"mb": "MB", "gb": "GB", "api": "API", "vnc": "VNC", "uri": "URI",
}

// goName converts a snake_case or SCREAMING_SNAKE spec name to Go CamelCase.
func goName(s string) string {
	parts := strings.Split(strings.ToLower(s), "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if u, ok := initialisms[p]; ok {
			b.WriteString(u)
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}

// typeName maps a spec schema name to its generated Go type name.
func typeName(name string) string {
	n := strings.TrimPrefix(name, "v3alpha1")
	// grpc-gateway body types: ComputeServiceStartInstanceBody -> StartInstanceRequest
	if i := strings.Index(n, "Service"); i > 0 && strings.HasSuffix(n, "Body") {
		n = strings.TrimSuffix(n[i+len("Service"):], "Body") + "Request"
	}
	return n
}

// enumConstName maps INSTANCE_STATUS_ACTIVE (type InstanceStatus) to
// InstanceStatusActive.
func enumConstName(typ, val string) string {
	return typ + strings.TrimPrefix(goName(val), typ)
}

type generator struct {
	spec        *spec
	protoCommit string
	specName    string
}

func main() {
	specPath := flag.String("spec", "", "path to the OpenAPI 3 JSON document")
	protoCommit := flag.String("proto-commit", "unknown", "proto repo commit the spec was generated from")
	outDir := flag.String("out", ".", "root of the v3 module")
	flag.Parse()
	if *specPath == "" {
		fmt.Fprintln(os.Stderr, "gencmd: -spec is required")
		os.Exit(2)
	}

	raw, err := os.ReadFile(*specPath)
	if err != nil {
		fatal(err)
	}
	var sp spec
	if err := json.Unmarshal(raw, &sp); err != nil {
		fatal(fmt.Errorf("parsing %s: %w", *specPath, err))
	}

	g := &generator{spec: &sp, protoCommit: *protoCommit, specName: base(*specPath)}
	for _, svc := range services {
		types, methods := g.render(svc)
		writeFile(*outDir, svc.pkg, "types.go", types)
		writeFile(*outDir, svc.pkg, "service.go", methods)
	}
	if src := g.renderReasons(); src != nil {
		writeFile(*outDir, reasonsPkg, "reasons.go", src)
	}
}

func base(p string) string { return p[strings.LastIndex(p, "/")+1:] }

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "gencmd:", err)
	os.Exit(1)
}

func writeFile(outDir, pkg, name string, src []byte) {
	formatted, err := format.Source(src)
	if err != nil {
		// Write the unformatted source anyway so the error is debuggable.
		_ = os.WriteFile(outDir+"/"+pkg+"/"+name+".broken", src, 0o644)
		fatal(fmt.Errorf("%s/%s does not compile: %w", pkg, name, err))
	}
	if err := os.MkdirAll(outDir+"/"+pkg, 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(outDir+"/"+pkg+"/"+name, formatted, 0o644); err != nil {
		fatal(err)
	}
}

func (g *generator) header(pkg string) string {
	return fmt.Sprintf(`// Code generated from %s (proto commit %s) by gencmd. DO NOT EDIT.

package %s
`, g.specName, g.protoCommit, pkg)
}

// opInfo is one operation resolved against the spec.
type opInfo struct {
	name       string
	httpMethod string
	path       string // spec path, e.g. /v3alpha1/compute/instances/{instance_id}
	pathParams []string
	query      []*parameter
	bodyRef    string // schema name of the request body, "" if none
	respRef    string // schema name of the 200 response
}

func (g *generator) collectOps(svc serviceDef) []opInfo {
	var ops []opInfo
	for path, byMethod := range g.spec.Paths {
		if !strings.HasPrefix(path, svc.pathPrefix) {
			continue
		}
		for httpMethod, op := range byMethod {
			name := op.OperationID[strings.Index(op.OperationID, "_")+1:]
			for _, tok := range svc.strip {
				name = strings.ReplaceAll(name, tok, "")
			}
			info := opInfo{name: name, httpMethod: strings.ToUpper(httpMethod), path: path}
			for _, p := range op.Parameters {
				switch p.In {
				case "path":
					info.pathParams = append(info.pathParams, p.Name)
				case "query":
					info.query = append(info.query, p)
				}
			}
			if op.RequestBody != nil {
				info.bodyRef = op.RequestBody.Content["application/json"].Schema.refName()
			}
			if ok, exists := op.Responses["200"]; exists {
				info.respRef = ok.Content["application/json"].Schema.refName()
			}
			ops = append(ops, info)
		}
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].name < ops[j].name })
	return ops
}

// reachableTypes walks $refs from the operations' request/response schemas.
func (g *generator) reachableTypes(ops []opInfo) []string {
	seen := map[string]bool{}
	var visit func(name string)
	visit = func(name string) {
		switch name {
		case "", "rpcStatus", "protobufAny", "v1PageInfo", "v1PageOptions":
			return // mapped or error-path types, never emitted
		}
		if seen[name] {
			return
		}
		seen[name] = true
		sch := g.spec.Components.Schemas[name]
		if sch == nil {
			fatal(fmt.Errorf("schema %s referenced but not defined", name))
		}
		for _, prop := range sch.Properties {
			visit(deepRef(prop))
		}
	}
	for _, op := range ops {
		if !g.trivialBody(op.bodyRef) { // trivial bodies are api.CallParams
			visit(op.bodyRef)
		}
		visit(op.respRef)
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func deepRef(s *schema) string {
	for s != nil {
		if s.Ref != "" {
			return s.refName()
		}
		switch {
		case s.Items != nil:
			s = s.Items
		case s.AdditionalProperties != nil:
			s = s.AdditionalProperties
		default:
			return ""
		}
	}
	return ""
}

// goType maps a schema to a Go type expression.
func (g *generator) goType(s *schema) string {
	if s == nil {
		return "json.RawMessage"
	}
	if ref := s.refName(); ref != "" {
		switch ref {
		case "v1PageInfo":
			return "api.PageInfo"
		case "protobufAny":
			return "json.RawMessage"
		}
		return typeName(ref)
	}
	switch s.Type {
	case "string":
		if s.Format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if s.Format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		return "[]" + g.goType(s.Items)
	case "object":
		if s.AdditionalProperties != nil {
			return "map[string]" + g.goType(s.AdditionalProperties)
		}
		return "json.RawMessage"
	}
	return "json.RawMessage"
}

// fieldType is goType for a struct field, made a pointer where the field
// has presence, so that leaving it unset omits it from the JSON:
//   - a proto3 optional scalar (nullable in the spec), so its zero value
//     can be sent: nil omits it, a pointer to false or "" sends it;
//   - a nested message, since omitempty never omits a struct and an empty
//     object is not the same request as an absent one.
func (g *generator) fieldType(s *schema) string {
	t := g.goType(s)
	if s.Nullable && s.Ref == "" && s.Type != "array" && s.Type != "object" {
		return "*" + t
	}
	if g.isMessageRef(s) {
		return "*" + t
	}
	return t
}

// isMessageRef reports whether s references a generated message type
// (not an enum, and not a type mapped to the hand-written core).
func (g *generator) isMessageRef(s *schema) bool {
	switch ref := s.refName(); ref {
	case "", "v1PageInfo", "protobufAny":
		return false
	default:
		target := g.spec.Components.Schemas[ref]
		return target != nil && len(target.Enum) == 0
	}
}

func (g *generator) render(svc serviceDef) (types, methods []byte) {
	ops := g.collectOps(svc)
	return g.renderTypes(svc, ops), g.renderService(svc, ops)
}

func (g *generator) renderTypes(svc serviceDef, ops []opInfo) []byte {
	var b strings.Builder
	b.WriteString(g.header(svc.pkg))

	body := &strings.Builder{}
	for _, name := range g.reachableTypes(ops) {
		sch := g.spec.Components.Schemas[name]
		tn := typeName(name)
		if len(sch.Enum) > 0 {
			writeDoc(body, tn, sch.Description)
			fmt.Fprintf(body, "type %s string\n\nconst (\n", tn)
			for _, v := range sch.Enum {
				fmt.Fprintf(body, "%s %s = %q\n", enumConstName(tn, v), tn, v)
			}
			body.WriteString(")\n\n")
			continue
		}
		writeDoc(body, tn, sch.Description)
		fmt.Fprintf(body, "type %s struct {\n", tn)
		for _, prop := range sortedKeys(sch.Properties) {
			ps := sch.Properties[prop]
			if ps.Description != "" {
				fmt.Fprintf(body, "// %s\n", strings.ReplaceAll(strings.TrimSpace(ps.Description), "\n", "\n// "))
			}
			fmt.Fprintf(body, "%s %s `json:\"%s,omitempty\"`\n", goName(prop), g.fieldType(ps), prop)
		}
		body.WriteString("}\n\n")
	}

	src := body.String()
	writeImports(&b, src)
	b.WriteString(src)
	return []byte(b.String())
}

func writeDoc(b *strings.Builder, name, desc string) {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		fmt.Fprintf(b, "// %s is defined in the Civo v3 API.\n", name)
		return
	}
	if !strings.HasPrefix(desc, name) {
		desc = name + " — " + desc
	}
	fmt.Fprintf(b, "// %s\n", strings.ReplaceAll(desc, "\n", "\n// "))
}

func writeImports(b *strings.Builder, src string) {
	var imports []string
	if strings.Contains(src, "time.Time") {
		imports = append(imports, `"time"`)
	}
	if strings.Contains(src, "json.RawMessage") {
		imports = append(imports, `"encoding/json"`)
	}
	if strings.Contains(src, "api.") {
		imports = append(imports, `"github.com/civo/civogo/v3/api"`)
	}
	if strings.Contains(src, "iter.Seq2") {
		imports = append(imports, `"iter"`)
	}
	if strings.Contains(src, "url.Values") || strings.Contains(src, "url.PathEscape") {
		imports = append(imports, `"net/url"`)
	}
	if strings.Contains(src, "strings.Replace") {
		imports = append(imports, `"strings"`)
	}
	if strings.Contains(src, "context.Context") {
		imports = append(imports, `"context"`)
	}
	if strings.Contains(src, "fmt.") {
		imports = append(imports, `"fmt"`)
	}
	if len(imports) == 0 {
		return
	}
	b.WriteString("import (\n" + strings.Join(imports, "\n") + "\n)\n\n")
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// trivialBody reports whether a request body schema carries nothing but
// call-scoping fields (region/account_id); such operations get a
// variadic-option signature instead of a request struct.
func (g *generator) trivialBody(ref string) bool {
	if ref == "" {
		return true
	}
	sch := g.spec.Components.Schemas[ref]
	for prop := range sch.Properties {
		if prop != "region" && prop != "account_id" {
			return false
		}
	}
	return true
}

// unwrap resolves a response wrapper to (goType, fieldJSONName).
// Empty responses return ("", ""); multi-field responses return the
// wrapper struct type itself with field "".
func (g *generator) unwrap(ref string) (string, string) {
	if ref == "" {
		return "", ""
	}
	sch := g.spec.Components.Schemas[ref]
	if len(sch.Properties) == 0 {
		return "", ""
	}
	if len(sch.Properties) == 1 {
		for name, prop := range sch.Properties {
			if deepRef(prop) != "" {
				return g.goType(prop), name
			}
		}
	}
	return typeName(ref), ""
}

func (g *generator) renderService(svc serviceDef, ops []opInfo) []byte {
	var b strings.Builder
	b.WriteString(g.header(svc.pkg))

	body := &strings.Builder{}
	fmt.Fprintf(body, `// Service provides access to the %s API.
type Service struct {
	t *api.Transport
}

// NewService is called by the root civogo package; use civogo.New.
func NewService(t *api.Transport) *Service { return &Service{t: t} }

`, svc.pkg)

	for _, op := range ops {
		g.renderMethod(body, svc, op)
	}

	src := body.String()
	writeImports(&b, src)
	b.WriteString(src)
	return []byte(b.String())
}

func (g *generator) renderMethod(b *strings.Builder, svc serviceDef, op opInfo) {
	retType, retField := g.unwrap(op.respRef)
	isList := strings.HasPrefix(op.name, "List")

	sig := []string{"ctx context.Context"}
	var pathExpr = fmt.Sprintf("%q", op.path)
	for _, pp := range op.pathParams {
		arg := lowerFirst(goName(pp))
		sig = append(sig, arg+" string")
		pathExpr = fmt.Sprintf("strings.Replace(%s, \"{%s}\", url.PathEscape(%s), 1)", pathExpr, pp, arg)
	}

	reqStruct := ""
	trivial := g.trivialBody(op.bodyRef)
	switch {
	case isList:
		reqStruct = "ListRequest"
		sig = append(sig, "req ListRequest")
	case !trivial:
		reqStruct = typeName(op.bodyRef)
		sig = append(sig, "req "+reqStruct)
	default:
		sig = append(sig, "opts ...api.CallOption")
	}
	params := strings.Join(sig, ", ")

	rets := "error"
	if retType != "" {
		star := ""
		if !strings.HasPrefix(retType, "[]") {
			star = "*"
		}
		rets = fmt.Sprintf("(%s%s, error)", star, retType)
	}

	if isList {
		g.renderListMethod(b, svc, op)
		return
	}

	writeMethodDoc(b, op, retType)
	fmt.Fprintf(b, "func (s *Service) %s(%s) %s {\n", op.name, params, rets)

	// Query and body assembly. Trivial ops carry region/account_id in
	// the query (GET/DELETE, or bodyless POSTs) or as the JSON body.
	switch {
	case trivial && op.bodyRef == "":
		b.WriteString("q := api.ResolveCallParams(s.t, opts...).Query()\n")
	case trivial:
		b.WriteString("body := api.ResolveCallParams(s.t, opts...)\n")
	case reqStruct != "":
		b.WriteString("if req.Region == \"\" {\nreq.Region = s.t.DefaultRegion()\n}\n")
	}

	errRet := "return err"
	okRet := "return nil"
	out := "nil"
	if retType != "" {
		fmt.Fprintf(b, "var out struct {\nV %s `json:\"%s\"`\n}\n", retType, retField)
		out = "&out"
		if !strings.HasPrefix(retType, "[]") {
			okRet = "return &out.V, nil"
		} else {
			okRet = "return out.V, nil"
		}
		errRet = "return nil, err"
	}

	query, bodyArg := "nil", "nil"
	if trivial && op.bodyRef == "" {
		query = "q"
	} else if trivial {
		bodyArg = "body"
	} else if reqStruct != "" {
		bodyArg = "req"
	}

	fmt.Fprintf(b, "if err := s.t.Do(ctx, %q, %s, %s, %s, %s); err != nil {\n%s\n}\n%s\n}\n\n",
		op.httpMethod, pathExpr, query, bodyArg, out, errRet, okRet)
}

// renderListMethod emits ListRequest, List and All.
func (g *generator) renderListMethod(b *strings.Builder, svc serviceDef, op opInfo) {
	respType := typeName(op.respRef)

	b.WriteString("// ListRequest selects what List returns.\ntype ListRequest struct {\n")
	for _, p := range op.query {
		switch p.Name {
		case "pagination.page", "pagination.per_page":
			continue // covered by Page
		}
		if d := strings.TrimSpace(p.Description); d != "" {
			fmt.Fprintf(b, "// %s\n", d)
		}
		fmt.Fprintf(b, "%s %s\n", goName(strings.ReplaceAll(p.Name, ".", "_")), g.goType(p.Schema))
	}
	b.WriteString("// Page selects one page; the zero value uses server defaults.\nPage api.PageParams\n}\n\n")

	fmt.Fprintf(b, "// List calls %s %s, returning one page.\n", op.httpMethod, op.path)
	fmt.Fprintf(b, `func (s *Service) List(ctx context.Context, req ListRequest) (*%s, error) {
	q := url.Values{}
	regions := req.Regions
	if len(regions) == 0 && s.t.DefaultRegion() != "" {
		regions = []string{s.t.DefaultRegion()}
	}
	q["regions"] = regions
	if req.AccountID != "" {
		q.Set("account_id", req.AccountID)
	}
	if req.Page.Page > 0 {
		q.Set("pagination.page", fmt.Sprint(req.Page.Page))
	}
	if req.Page.PerPage > 0 {
		q.Set("pagination.per_page", fmt.Sprint(req.Page.PerPage))
	}
	var out %s
	if err := s.t.Do(ctx, "GET", %q, q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

`, respType, respType, op.path)

	itemsField := ""
	sch := g.spec.Components.Schemas[op.respRef]
	for _, name := range sortedKeys(sch.Properties) {
		if sch.Properties[name].Type == "array" {
			itemsField = goName(name)
		}
	}
	fmt.Fprintf(b, `// All lazily iterates every %s matching req across all pages.
func (s *Service) All(ctx context.Context, req ListRequest) iter.Seq2[%s, error] {
	return api.All(func(p api.PageParams) ([]%s, *api.PageInfo, error) {
		req.Page = p
		resp, err := s.List(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		return resp.%s, &resp.PageInfo, nil
	})
}

`, svc.entity, svc.entity, svc.entity, itemsField)
}

func writeMethodDoc(b *strings.Builder, op opInfo, retType string) {
	verb := op.name
	fmt.Fprintf(b, "// %s calls %s %s.\n", verb, op.httpMethod, op.path)
	if retType == "" && op.respRef != "" {
		b.WriteString("// The API returns no payload for this call.\n")
	}
}

func lowerFirst(s string) string {
	if s == "ID" {
		return "id"
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// enumValueDoc matches the " - VALUE: description" lines grpc-gateway
// appends to an enum's description, one per documented value.
var enumValueDoc = regexp.MustCompile(`(?m)^ - ([A-Z0-9_]+): `)

// splitEnumDoc separates an enum description into its overview and the
// per-value descriptions grpc-gateway appends to it.
func splitEnumDoc(desc string) (overview string, values map[string]string) {
	values = map[string]string{}
	locs := enumValueDoc.FindAllStringSubmatchIndex(desc, -1)
	if len(locs) == 0 {
		return strings.TrimSpace(desc), values
	}
	overview = strings.TrimSpace(desc[:locs[0][0]])
	for i, loc := range locs {
		end := len(desc)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		values[desc[loc[2]:loc[3]]] = strings.TrimSpace(desc[loc[1]:end])
	}
	return overview, values
}

func (g *generator) renderReasons() []byte {
	sch := g.spec.Components.Schemas[reasonsSchema]
	if sch == nil {
		return nil
	}
	overview, docs := splitEnumDoc(sch.Description)

	b := &strings.Builder{}
	b.WriteString(strings.Replace(g.header(reasonsPkg), "\npackage ",
		"\n// Package "+reasonsPkg+" holds the error reasons the Compute v3 API reports.\n//\n// "+
			strings.ReplaceAll(overview, "\n", "\n// ")+"\npackage ", 1))
	b.WriteString("\n// Reasons reported in Error.Reason; match them with civogo.HasReason.\n")
	b.WriteString("// The API may add reasons, so never treat this set as exhaustive.\nconst (\n")
	for _, v := range sch.Enum {
		if strings.HasSuffix(v, "_UNSPECIFIED") {
			continue
		}
		name := "Reason" + goName(v)
		if d := docs[v]; d != "" {
			fmt.Fprintf(b, "// %s — %s\n", name, strings.ReplaceAll(d, "\n", "\n// "))
		}
		fmt.Fprintf(b, "%s = %q\n", name, v)
	}
	b.WriteString(")\n")
	return []byte(b.String())
}
