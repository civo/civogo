# civogo v2 → v3

What changes for someone using the SDK. v3 covers Compute instances only
for now; everything else is still v2.

## At a glance

| | v2 | v3 |
|---|---|---|
| Import | `github.com/civo/civogo` | `github.com/civo/civogo/v3` (+ `/compute`) |
| Client | `civogo.NewClient(key, "LON1")` | `civogo.New(key, civogo.WithDefaultRegion("lon1"))` |
| Calls | `client.GetInstance(id)` | `client.Compute.Instances.Get(ctx, id)` |
| Context / cancellation | none | `ctx` on every call |
| Region | fixed per client | default per client, `civogo.WithRegion` per call |
| Multi-region list | one call per region | one call, many `WithRegion` |
| Pagination | `ListInstances(page, perPage)` or `ListAllInstances()` (slice, capped at 10k) | `List` (one page) or `All` (lazy iterator, no cap) |
| Results | often `*SimpleResponse{Result: "success"}` | the updated resource, or just `error` |
| Errors | ~250 sentinel errors matched by string | one `*civogo.Error` with reason, request ID, violations |
| Optional fields | sent always, sentinel tricks (`notes_delete`) | pointers: `nil` = unchanged |
| Concurrency | client mutates itself per request | immutable, safe to share |
| Dependencies | k8s.io/api, apimachinery, … | standard library only |

## Client and regions

v2 binds the client to one region; another region means another client
(or mutating `client.Region`, which races):

```go
client, _ := civogo.NewClient(key, "LON1")
inst, _ := client.GetInstance(id) // always LON1
```

v3 has a default region you can override on any call:

```go
client, _ := civogo.New(key, civogo.WithDefaultRegion("lon1"))
inst, _ := client.Compute.Instances.Get(ctx, id)                           // lon1
inst, _ = client.Compute.Instances.Get(ctx, id, civogo.WithRegion("fra1")) // fra1
```

Other options: `WithBaseURL`, `WithHTTPClient` (your own timeouts, proxy,
retries; v2 overwrote the transport on every request), `WithUserAgent`,
`WithAccountID` (JWT auth only).

## Context

v2 calls can't be cancelled or given a deadline. In v3 every call takes a
`context.Context`, so timeouts, cancellation and tracing work as usual.

## Listing

v2:

```go
page, _ := client.ListInstances(1, 50)  // one page, one region
all, _ := client.ListAllInstances()     // everything in memory, errors past 10k items
```

v3:

```go
page, _ := client.Compute.Instances.List(ctx, compute.ListInstancesRequest{
    Page: api.PageParams{PerPage: 50},
}, civogo.WithRegion("lon1"), civogo.WithRegion("fra1")) // both regions, one call

for inst, err := range client.Compute.Instances.All(ctx, compute.ListInstancesRequest{}) {
    if err != nil { ... }
    // pages are fetched as you go; break stops fetching
}
```

`page.PageInfo` carries the server's real totals (`TotalItemsCount`,
`TotalPages`).

## Return values

v2 actions often return `*SimpleResponse` (`Result: "success"`), so you
re-fetch the resource to see its new state. v3 returns what the API
returns: `Update`, `Resize`, `Delete`, … give back the `*Instance`;
fire-and-forget actions (`Start`, `Stop`, …) return only `error`.

## Errors

v2 turns API errors into one of ~250 sentinel values by matching codes and
strings, and you check with `errors.Is(err, civogo.DatabaseInstanceNotFoundError)`
and friends. Unknown errors lose detail.

v3 returns a single `*civogo.Error` decoded from the API's structured error:

```go
_, err := client.Compute.Instances.Get(ctx, id)
switch {
case civogo.HasReason(err, compute.ReasonInstanceNotFound):
    // gone
case civogo.IsRateLimited(err):
    // back off
case err != nil:
    var apiErr *civogo.Error
    if errors.As(err, &apiErr) {
        log.Print(apiErr.RequestID)          // quote this to support
        for _, v := range apiErr.FieldViolations {
            log.Printf("%s: %s", v.Field, v.Description)
        }
    }
}
```

- Reasons are generated constants (`compute.Reason*`), so the compiler
  catches typos and they stay in sync with the API.
- Generic helpers: `IsNotFound`, `IsAlreadyExists`, `IsAuthFailure`,
  `IsRateLimited`, `IsInvalidArgument`.
- Every error carries the request ID and, for validation failures, which
  field was wrong.

## Optional fields and updates

v2 sends every field, so "leave unchanged" vs "clear" needs special cases
(`notes_delete`, `PublicIPRequired: "true"` as a string, …).

v3 uses pointers for optional fields: `nil` leaves the field alone, a value
sets it, including `false` or `""`:

```go
client.Compute.Instances.Update(ctx, id, compute.UpdateInstanceRequest{
    PublicIPRequired: civogo.Ptr(false), // turn it off
    Notes:            civogo.Ptr(""),    // clear notes
    // everything else untouched
})
```

## Concurrency

v2's `Client` writes to itself on every request (`LastJSONResponse`, the
HTTP transport), so sharing it across goroutines races. v3's client is
immutable after `New`; share one everywhere.

## Not in v3 yet

- Resources other than Compute instances (they'll arrive as more packages
  under `v3/`, e.g. `client.Kubernetes`).
- Helpers like `FindInstance` (partial name/ID match) and
  `NewInstanceConfig` (pre-filled defaults).
- A fake client for tests (v2's `FakeClient` / `Clienter`); for now use
  `civogo.WithBaseURL` with an `httptest.Server`.
