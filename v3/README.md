# civogo v3 (preview)

Go SDK for the Civo v3 API.

## Usage

```go
client, err := civogo.New(token, civogo.WithDefaultRegion("lon1"))
if err != nil { ... }

// Region is an option on every call; the client default applies otherwise.
inst, err := client.Compute.Instances.Get(ctx, id, civogo.WithRegion("fra1"))
inst, err = client.Compute.Instances.Update(ctx, id, compute.UpdateInstanceRequest{
    Notes: civogo.Ptr("web"),
}, civogo.WithRegion("fra1"))

// Lists take WithRegion several times to query multiple regions, with real
// server-side pagination. Any other call given two regions returns an error.
page, err := client.Compute.Instances.List(ctx, compute.ListInstancesRequest{
    Page: api.PageParams{PerPage: 50},
}, civogo.WithRegion("lon1"), civogo.WithRegion("fra1"))

// Lazy iteration over every page.
for inst, err := range client.Compute.Instances.All(ctx, compute.ListInstancesRequest{}) {
    if err != nil { ... }
    fmt.Println(inst.Hostname, inst.Region)
}

// WithAccountID is only needed with JWT auth; an API key identifies the
// account by itself.

// Errors are one structured type. Reason is the stable condition;
// FieldViolations explain validation failures; RequestID is for support.
if civogo.HasReason(err, compute.ReasonInstanceNotFound) { ... }
if civogo.IsNotFound(err) { ... }
var apiErr *civogo.Error // errors.As for status/reason/violations/request ID
```

The client is immutable and safe for concurrent use. Every method takes
a `context.Context`. The module has no dependencies outside the
standard library.

## Layout

- `civogo.go`, `api/` — hand-written public surface: client construction
  and options, per-call options, the `Error` type and predicates, page
  types.
- `internal/transport/` — hand-written request path shared by the
  generated packages: auth, encoding, error decoding, pagination. Internal
  so it can change without breaking users.
- `compute/` — **generated** from the Compute v3 OpenAPI document: one
  package per API, holding all its services (`Client.Compute.Instances`),
  types and error reasons. Do not edit; regenerate instead. Other APIs
  (Kubernetes, networking, …) will get their own package next to it.
- `internal/gencmd/` — the generator, specialised to the OpenAPI shapes
  grpc-gateway emits. See [how it works](internal/gencmd/README.md).

## Regenerating

Regeneration is manual and must be run against the spec version that is
**live in production**, not the latest unreleased one:

```sh
go run ./internal/gencmd \
  -spec path/to/compute.openapi3.json \
  -proto-commit <commit the spec was generated from>
```

The proto commit is stamped into every generated file header. Running
the generator twice against the same spec produces a zero diff; MRs that
touch `compute/` by hand will not survive the next regeneration.
