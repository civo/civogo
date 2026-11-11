# gencmd

`gencmd` generates the resource layer of civogo v3 from a Civo v3 OpenAPI
document. It is not a general OpenAPI generator: it relies on the shapes
grpc-gateway emits from our protos, which keeps it small (one file,
`main.go`).

```
compute.openapi3.json
  │  1. parse into a minimal model (paths, operations, schemas)
  │  2. select operations by path prefix   → /v3alpha1/compute/instances*
  │  3. walk $refs to find reachable types → only what those calls use
  ├─ compute/client.go     package doc and Client{Instances, …}
  ├─ compute/instances.go  InstancesService: one method per operation, plus List/All
  ├─ compute/types.go      one struct or enum per reachable schema
  └─ compute/reasons.go    ErrorReason enum → Reason* constants
  then gofmt each file; abort if it does not compile
```

## Running it

Run it manually, against the spec version that is live in production, not
the latest unreleased one:

```sh
go run ./internal/gencmd \
  -spec path/to/compute.openapi3.json \
  -proto-commit <commit the spec was generated from>
```

The output is deterministic: running it twice on the same spec gives a zero
diff. Never edit generated files by hand.

## How it works

### Packages and services

One spec (one proto package, e.g. `civo.compute.v3alpha1`) becomes one Go
package. Proto message names are unique within a proto package, so the Go
types never collide, and types shared between services (`PlacementRule`,
error reasons) live in one place.

The `packages` map in `main.go` declares each package, selected with
`-pkg` (default `compute`), and its services:

| Field        | Example                        | Purpose                                  |
|--------------|--------------------------------|------------------------------------------|
| `name`       | `Instances`                    | Field on the package `Client`; `InstancesService` type |
| `entity`     | `Instance`                     | Item type yielded by `All`               |
| `pathPrefix` | `/v3alpha1/compute/instances`  | Operations whose path starts with this   |
| `strip`      | `Instance`, `ToInstance`, …    | Tokens removed from method names         |

Adding a Compute service (SSH keys, snapshots, …) means adding an entry;
adding another API means adding a package and running with its spec.

### Method names

grpc-gateway names operations `ComputeService_StartInstance`. The generator
drops everything up to `_`, then removes the `strip` tokens:
`StartInstance` → `Start`, `AssignIPToInstance` → `AssignIP`.

### Types (`types.go`)

Starting from the request body and `200` response of every operation in
the package's services, the
generator follows every `$ref` through properties, arrays and maps, and
emits only those schemas. `v1PageInfo` maps to `api.PageInfo`;
`rpcStatus` and `protobufAny` are never emitted.

- **Names**: `v3alpha1Instance` → `Instance`; grpc-gateway body types
  `ComputeServiceUpdateInstanceBody` → `UpdateInstanceRequest`.
- **Enums**: a named string type plus constants (`InstanceStatusActive`).
- **Fields**: Go-cased names (`public_ip_required` → `PublicIPRequired`,
  using the `initialisms` map), JSON tag with the proto name and
  `omitempty`, doc comment from the proto description.
- **Scalars**: string → `string` (date-time → `time.Time`),
  integer → `int`/`int64`, number → `float64`, boolean → `bool`,
  arrays → slices, maps → `map[string]T`. Proto `int64` is a string in
  the spec (protojson encodes it so) and stays a Go `string`.
- **Pointers**, where leaving a field unset must omit it:
  - proto3 `optional` scalars (marked `nullable` by
    `proto3_optional_nullable`) → `*bool`, `*string`, …; `nil` omits the
    field, `civogo.Ptr(false)` sends `false`;
  - nested messages → `*PlacementRule`, … because `omitempty` never omits
    a struct, and the API treats an empty object differently from an
    absent one.

### Methods (`<service>.go`)

Each operation becomes one of three shapes:

- **Simple call**: no body, or a body with only `region`/`account_id`
  (Get, Delete, Start, Stop, …). Signature
  `(ctx, id, opts ...api.CallOption)`; `api.WithRegion` and
  `api.WithAccountID` scope the call, with the client's default region as
  fallback. Sent as query parameters, or as the JSON body for POSTs that
  declare one.
- **Request struct** (Create, Update, Resize, …):
  `(ctx, id, req XRequest)`. An empty `req.Region` is filled with the
  default region; `req` is sent as JSON.
- **List** (operation name starts with `List`): emits `ListInstancesRequest` from
  the query parameters plus `Page`, `List` returning one page with the
  server's pagination, and `All`, an `iter.Seq2` that pages through
  `List` 100 items at a time until `total_pages`.

Path parameters (`{instance_id}`) become string arguments, escaped into
the path. Responses are unwrapped: a wrapper with a single field
(`{"instance": {...}}`) returns `*Instance`; a list field returns
`[]Instance`; an empty response returns only `error`; a wrapper with more
fields is returned whole.

All methods call `api.Transport.Do`. Auth, encoding and error decoding
are hand-written there, never generated.

### Error reasons (`reasons.go`)

When the spec contains `v3alpha1ErrorReason`, each value except
`UNSPECIFIED` becomes an untyped string constant
(`ReasonInstanceNotFound = "INSTANCE_NOT_FOUND"`) for `civogo.HasReason`.
Per-value docs are parsed from the ` - VALUE: text` lines grpc-gateway
appends to the enum description. Specs without the enum produce no file.

### Output

Every file is stamped with the spec name and the `-proto-commit` value,
then formatted with `go/format`. If formatting fails, the source is
written to `<file>.broken` and the generator exits non-zero.

## Assumptions

These hold for Compute today and are the first things to revisit when
adding a service:

- The List template hardcodes the `regions`, `account_id` and
  `pagination.*` query parameters, and expects one array field plus
  `page_info` in the response.
- An operation is a list only if its name starts with `List`.
- A body is "simple" only if all its properties are `region` or
  `account_id`.
- The reasons schema name (`v3alpha1ErrorReason`) is hardcoded; an API on
  another proto version needs it derived from the spec.
- A `$ref` to an undefined schema aborts generation.
