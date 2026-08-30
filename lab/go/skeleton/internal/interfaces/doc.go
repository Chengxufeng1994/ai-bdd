// Package interfaces adapts the outside world to application use cases.
//
// It holds no business rules — a handler that decides something has put a rule
// in the wrong layer.
//
// # What an adapter may import
//
// The driving ports in application/port/in, the command and query types their
// methods name, the results those return, and the application's error kinds.
// Nothing else.
//
//	✓ in.VersionService          the driving port — an interface of capabilities
//	✓ query.GetVersion           a type its methods name
//	✓ query.GetVersionResult     what that query returns
//	✓ errors.Kind                how the application classified a failure
//	✗ service.VersionService         the concrete implementation of the port
//	✗ the use case behind            unexported; the compiler enforces this row
//	  NewGetVersion
//	✗ workout.Workout                the aggregate
//
// The port is an interface, so holding it grants the capabilities it names and
// nothing more, and a decorated implementation substitutes for a plain one
// without the adapter noticing. Depending on the concrete service instead would
// bypass the contract and lose that.
//
// The use-case row is not merely discouraged; the application layer keeps it
// unnameable. usecase/query holds a query, its result and the use case
// answering it in one package, so an adapter must import that package to name
// query.GetVersion at all — which would ordinarily expose the use case too. Its
// type is therefore unexported: NewGetVersion is exported and returns
// usecase.QueryHandler, so an adapter can still construct one — it just cannot
// name the concrete type back. pkg/log.New returns an interface for the same
// reason: the interface is the contract and the concrete type has nothing
// extra to offer, so there is nothing to lose by hiding it.
//
// Depending on the port rather than the concrete service is what keeps five
// adapters from growing five different ideas of what the application does; and
// an adapter that never holds an aggregate cannot mutate one past its
// invariants. Both are enforced by import, not by review.
//
// bootstrap/ is the exception: as the composition root it constructs the
// concrete services and hands each adapter the port.
//
// # Structure
//
// One package per protocol, each self-contained:
//
//	interfaces/
//	├── http/              REST; gin + oapi-codegen
//	│   ├── apigen/        generated request and response types
//	│   ├── mapper/        HTTP request -> command / query
//	│   ├── presenter/     result -> HTTP response
//	│   ├── errmap/        application error kind -> status and Problem document
//	│   ├── router.go      engine, middleware, route registration
//	│   ├── server.go      holds the driving port; satisfies the strict interface
//	│   └── version.go     one operation, as a method on Server
//	├── grpc/
//	│   ├── pb/            generated from proto/
//	│   └── …              same mapper / presenter / errmap split
//	├── graphql/
//	│   └── gqlgen/        generated from schema.graphql
//	├── cli/
//	│   ├── request/       hand-written: no schema to generate from
//	│   └── response/
//	└── lambda/
//
// There is no handler/ package, and adding one would not buy what it looks like
// it buys. oapi-codegen's strict interface requires every operation to be a
// method on one type satisfying apigen.StrictServerInterface, so a separate
// package could only hang those methods on Server anyway, or make Server forward
// into it — and a forwarding layer is what the application layer rejects for its
// own facade, for the same reason: two mechanical edits per operation, and
// somewhere for a decision to hide. Operations are therefore methods on Server,
// one file per operation, named after the operation: http/version.go.
//
// Request and response types come from the generator wherever a schema exists.
// Hand-writing request/ and response/ beside apigen/ would create a second
// source of truth that drifts the moment the spec changes. Only protocols with
// no schema — the CLI, a Lambda event shape — declare them by hand.
//
// # Mapper in, presenter out
//
// Both are mechanical translation; they are named apart so the direction is
// visible at a glance rather than inferred from a signature.
//
//	mapper     apigen request  ->  command / query      inbound
//	presenter  result          ->  apigen response      outbound
//
// A handler does nothing else: take the parameters, call the mapper, call the
// port, call the presenter, and translate any error through errmap. If a handler
// contains a decision, that decision is a business rule sitting in the wrong
// layer.
//
// Presenters may drop fields — each protocol exposes as little as it likes. If a
// protocol needs a field the result does not carry, either the result should
// carry it or that protocol needs a different query. Neither is a reason to
// reach into the domain.
//
// # What this is not
//
// Not the Clean Architecture presenter, which receives output through an
// inverted port so that the use case never returns. In Go a returned value
// already keeps the dependency pointing inward — application hands back its own
// result type and knows nothing about HTTP — so an output port would buy an
// inversion that already exists, at the cost of an interface and an indirection.
//
// A presenter here stays mechanical until formatting acquires rules of its own:
// converting kilograms to pounds for a locale, rounding, collapsing large
// numbers, or producing output incrementally for streaming and cursor
// pagination.
//
// When that happens, settle which side of the line the rule is on first. "The
// unit this user prefers" is a business concept if the system remembers it and
// compares against it; it is presentation only if it decides how this one
// response is rendered. Getting that wrong hides a business rule inside an
// adapter, where no scenario will ever reach it.
//
// # Every adapter is a translation layer
//
// Read protocol input, build an application command or query, hand back the
// result in protocol terms. Nothing else. The test: deleting a protocol package
// should lose an entry point and not a single business rule.
//
// This is what makes five protocols affordable. Each one is thin, and the
// behaviour they expose is specified once — in features/, at the application
// boundary — rather than five times over.
//
// # Errors
//
// An application error carries a kind, not a transport code, because the
// application does not know which protocol is serving it. Every adapter owns the
// mapping from kind to its own vocabulary. The kinds are declared in
// application/errors; this is HTTP's column, implemented in http/errmap:
//
//	kind          HTTP  gRPC                 GraphQL extension   CLI exit
//	NotFound      404   NOT_FOUND            NOT_FOUND           4
//	Invalid       422   INVALID_ARGUMENT     BAD_USER_INPUT      2
//	Conflict      409   ABORTED              CONFLICT            3
//	Unauthorized  401   UNAUTHENTICATED      UNAUTHENTICATED     5
//	Forbidden     403   PERMISSION_DENIED    FORBIDDEN           6
//	Unavailable   503   UNAVAILABLE          SERVICE_UNAVAILABLE 7
//
// The mapping tables differ; the classification must not. Keeping the kinds in
// the application and the tables in the adapters is what stops the same failure
// from being a 404 over HTTP and a 500 over gRPC.
//
// Only the HTTP column exists as code. The other three are what the table would
// look like, kept here so the first adapter that needs one starts from the same
// classification rather than inventing a second.
//
//	MUST: recognise, then default. Look for a classification; map the kinds
//	this table knows; send everything else to 500 with a generic body.
//
// The order is the whole security property. The mirror image — map the ones we
// know, pass the rest through — reads as equivalent and is not: the first
// unclassified error to arrive renders err.Error() to whoever made the request,
// and stack traces, SQL fragments and internal hostnames arrive with it.
// errmap_test.go guards the default arm by asserting on the rendered JSON rather
// than on one field, so a new field cannot open a second leak quietly.
//
// A second rule keeps that redaction reachable at all:
//
//	MUST: a handler returns a typed response and a nil error, never a
//	returned error.
//
// oapi-codegen's generated strict handler has a default error path that
// writes gin.H{"msg": err.Error()} — the HandlerErrorFunc default in
// apigen/api.gen.go. router.go replaces that default with one that renders
// errmap's generic Problem body instead, so breaking the MUST no longer leaks;
// it only costs a less precise status than the handler's own kind would have
// chosen. That backstop is not a reason to lean on it: rendering the failure
// into a declared response type and returning nil is still what keeps the
// contract and the redaction in the same place, and what lets a handler answer
// with more than a flat 500. Every protocol's generated glue has an equivalent
// default; the rule and the backstop are the same in each.
//
// Three details of the recognising half:
//
//   - errors.As, not a type assertion. A use case wraps its failures with %w, so
//     the error reaching an adapter is not the errors.Error itself but something
//     around it. A type assertion sees the wrapper, finds nothing, and every
//     classified failure silently becomes a 500.
//
//   - KindUnclassified is the zero value, and the default arm catches it. A kind
//     someone forgot to set falls to 500 rather than becoming a NotFound.
//
//   - The list of kinds is closed on purpose, and completing it from
//     api/openapi.yaml's components/responses would be wrong. There is no kind
//     for 400: OapiRequestValidator rejects a malformed request in router.go
//     before any handler runs, so the application cannot produce one. Success
//     statuses are not classified at all — 200, 201 and 204 come from the
//     operation's contract and its presenter.
//
// http/version.go's GetVersion comment carries today's concrete case:
// errmap.ToProblem always calls StatusFor, so the function has a production
// caller now, but /version's contract has only one failure response, so
// GetVersion coerces the body's status back to 500 regardless of what
// StatusFor classified rather than switching on it. It is written once there
// rather than repeated here or in DATAFLOW.md's kind section, so a second
// failure response only invalidates one copy.
//
// # Generated code
//
// Three of these protocols are generated from a schema — OpenAPI, protobuf,
// GraphQL SDL. Each generator writes into its own subpackage (apigen, pb,
// gqlgen) so that hand-written code never shares a namespace with generated
// code, and so a directory can be deleted and rebuilt to prove the pipeline
// still works. All directives live in ../../generate.go.
//
// # Not yet populated
//
// Only http/ exists. It serves the version endpoint and carries the mapper,
// presenter and errmap packages described above — trivially small for this one
// operation, but present so that the shape is the one the next protocol and the
// next endpoint copy. errmap carries the whole kind table, but /version can only
// answer 200 or 500, so the table is exercised by its tests and by nothing in
// production yet.
//
// The rest arrive when something needs them — a protocol added before it has a
// caller is a guess about an interface nobody has asked for.
package interfaces
