// Package interfaces adapts the outside world to application use cases.
//
// It holds no business rules — a handler that decides something has put a rule
// in the wrong layer.
//
// # What an adapter may import
//
// The driving ports in application/port/in, the command and query types they
// name, and the domain's error kinds. Nothing else.
//
//	✓ service.WorkoutService     the handler bundle — fields only, no methods
//	✓ command.RecordWorkout      the types its fields name
//	✓ dto.VolumeView             what a query returns
//	✗ handler.RecordWorkout          the concrete handler
//	✗ workout.Workout                the aggregate
//
// The bundle is a struct rather than an interface because it has no behaviour to
// abstract — every field is already an in.CommandHandler or in.QueryHandler, so
// holding it grants access to contracts and nothing else. Depending on the
// concrete handler instead would bypass those contracts and lose the ability to
// decorate.
//
// Depending on the port rather than the concrete service is what keeps five
// adapters from growing five different ideas of what the application does; and
// an adapter that never holds an aggregate cannot mutate one past its
// invariants. Both are enforced by import, not by review.
//
// cmd/ is the exception: as the composition root it constructs the concrete
// services and hands each adapter the port.
//
// # Structure
//
// One package per protocol, each self-contained:
//
//	interfaces/
//	├── http/              REST; gin + oapi-codegen
//	│   ├── apigen/        generated request and response types
//	│   ├── handler/       take params, call mapper, call port/in, call presenter
//	│   ├── mapper/        HTTP request -> command / query
//	│   ├── presenter/     dto -> HTTP response
//	│   ├── errmap/        domain error kind -> status and Problem document
//	│   ├── router.go      engine, middleware, route registration
//	│   └── server.go      holds the application's service bundle
//	├── grpc/
//	│   ├── pb/            generated from proto/
//	│   └── …              same handler / mapper / presenter / errmap split
//	├── graphql/
//	│   └── gqlgen/        generated from schema.graphql
//	├── cli/
//	│   ├── request/       hand-written: no schema to generate from
//	│   └── response/
//	└── lambda/
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
//	presenter  dto             ->  apigen response      outbound
//
// A handler does nothing else: take the parameters, call the mapper, call the
// port, call the presenter, and translate any error through errmap. If a handler
// contains a decision, that decision is a business rule sitting in the wrong
// layer.
//
// Presenters may drop fields — each protocol exposes as little as it likes. If a
// protocol needs a field the dto does not carry, either the dto should carry it
// or that protocol needs a different query. Neither is a reason to reach into
// the domain.
//
// # What this is not
//
// Not the Clean Architecture presenter, which receives output through an
// inverted port so that the use case never returns. In Go a returned value
// already keeps the dependency pointing inward — application hands back its own
// dto and knows nothing about HTTP — so an output port would buy an inversion
// that already exists, at the cost of an interface and an indirection.
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
// Domain errors carry a kind, not a transport code, because the domain does not
// know which protocol is serving it. Every adapter owns the mapping from kind to
// its own vocabulary:
//
//	kind          HTTP  gRPC                 GraphQL extension   CLI exit
//	NotFound      404   NOT_FOUND            NOT_FOUND           4
//	Invalid       422   INVALID_ARGUMENT     BAD_USER_INPUT      2
//	Conflict      409   ABORTED              CONFLICT            3
//	Unauthorized  401   UNAUTHENTICATED      UNAUTHENTICATED     5
//
// The mapping tables differ; the classification must not. Keeping the kinds in
// the domain and the tables in the adapters is what stops the same failure from
// being a 404 over HTTP and a 500 over gRPC.
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
// Only http/ exists, serving the version endpoint that keeps the generation
// pipeline honest. The rest arrive when something needs them — a protocol added
// before it has a caller is a guess about an interface nobody has asked for.
package interfaces
