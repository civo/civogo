# civogo v3 (preview)

Go SDK for the Civo v3 API. Design notes live in the vNEXT epic
(civo/team/core&3); this module is the Compute PoC (&5).

## Usage

```go
client, err := civogo.New(token, civogo.WithDefaultRegion("lon1"))
if err != nil { ... }

// Multi-region listing with real server-side pagination.
page, err := client.Instances.List(ctx, instances.ListRequest{
    Regions: []string{"lon1", "fra1"},
    Page:    api.PageParams{PerPage: 50},
})

// Lazy iteration over every page.
for inst, err := range client.Instances.All(ctx, instances.ListRequest{}) {
    if err != nil { ... }
    fmt.Println(inst.Hostname, inst.Region)
}

// Per-call region override; the client default applies otherwise.
inst, err := client.Instances.Get(ctx, id, api.WithRegion("fra1"))

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

- `civogo.go`, `api/` — hand-written core: client construction and
  options, the single request path, the `Error` type and predicates,
  pagination primitives.
- `instances/` — **generated** from the Compute v3 OpenAPI document.
  Do not edit; regenerate instead.
- `compute/` — **generated** error reason constants, from the spec's
  `ErrorReason` enum (only emitted once the pinned spec includes it).
- `internal/gencmd/` — the generator. Deliberately specialised to the
  OpenAPI shapes grpc-gateway emits from the private proto repos.

## Regenerating

Regeneration is manual and must be run against the spec version that is
**live in production** (not the proto repo's HEAD — see the epic for the
drift rationale):

```sh
go run ./internal/gencmd \
  -spec ~/dev/civo/compute/proto/gen/openapi/civo/compute/v3alpha1/compute.openapi3.json \
  -proto-commit "$(git -C ~/dev/civo/compute/proto rev-parse HEAD)"
```

The proto commit is stamped into every generated file header. Running
the generator twice against the same spec produces a zero diff; MRs that
touch `instances/` by hand will not survive the next regeneration.
