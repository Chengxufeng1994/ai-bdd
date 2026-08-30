// Package application orchestrates use cases on top of domain.
//
// It depends on domain only. Whatever it needs from the outside world —
// storage, clocks, external services — it declares as an interface (a port) and
// receives by injection. It never imports infrastructure.
//
// The full request-to-response chain is in ../../docs/DATAFLOW.md.
//
// # Structure
//
//	application/
//	├── errors/                        a failure classified by Kind, with no
//	│   └── errors.go                  transport code anywhere in it
//	├── port/
//	│   ├── in/                        driving: what the outside may ask of us
//	│   │   └── version_service.go     an interface; its methods are capabilities
//	│   └── out/                       driven: what we need from the outside
//	│       ├── version_provider.go    the only driven port that exists
//	│       ├── workout_repository.go  write side, traffics in domain types
//	│       ├── volume_reader.go       read side, traffics in read models
//	│       └── mocks/                 generated; see ../../.mockery.yml
//	├── usecase/
//	│   ├── usecase.go                 QueryHandler / CommandHandler, generic
//	│   ├── query/                     reads: query, result and handler together
//	│   │   ├── get_version.go         the only use case that exists
//	│   │   └── workout_volume.go
//	│   └── command/                   writes: command, result and handler
//	│       └── record_workout.go
//	└── service/                       implements port/in by delegating
//	    └── version_service.go
//
// errors/, port/in, port/out, usecase/, usecase/query and service/ exist.
// Everything else above shows where the next thing goes.
//
// # Three packages that are deliberately absent
//
// dto/, assembler/ and readmodel/ are named in most descriptions of this layout
// and in none of this tree. Where their contents go instead:
//
//   - There is no dto/. A result is declared beside the query or command it
//     answers — query.GetVersionResult sits in usecase/query/get_version.go,
//     and a command's result will sit beside its command in usecase/command/.
//     The next section is why that does not close a cycle.
//
//   - assembler/ holds domain -> result projections, and arrives with the first
//     aggregate there is something to project from. domain/ is empty by rule,
//     so there is no such function to write and no package to hold it. When one
//     appears its signature is most of its specification: it takes an
//     aggregate, returns a result, returns no error, and calls no port.
//
//   - readmodel/ holds a denormalised projection shaped for display, filled
//     from its own table or cache and read through its own reader in port/out.
//     The line between it and a result is checkable: a read model exists
//     independently of any single query, and several queries may serve slices
//     of one; a result is one use case's output shape and nothing else's.
//     query.GetVersionResult is the latter.
//
// usecase/command/ is absent for a duller reason: git does not track empty
// directories, and no command exists yet.
//
// # A result lives beside its query, and no cycle follows
//
// An earlier arrangement gave results a package of their own and defended it as
// load-bearing: dto imported nothing, which was what let port/out and query both
// depend on it, whereas declaring the view inside query would have forced
// port/out to import query in order to declare a reader returning it — while
// query already imported port/out for the interface. A cycle.
//
// That cycle needs both edges. Only one of them exists now, and the other is
// ruled out by a rule worth stating on its own:
//
//	MUST NOT: declare a driven port that returns a use case's result type.
//
// usecase/query imports port/out, as it must — a use case calls the ports it
// needs. port/out does not import usecase/query, because what a driven port
// hands back is a domain type or a read model shaped by its store, never the
// shape one use case chose for its caller. out.VersionProvider returns a
// string, not a query.GetVersionResult.
//
// What the separate package cost outweighed what it bought: a use case's input,
// its output and its implementation sat in three packages, so reading one meant
// opening three files. They are one file now.
//
// # Ports have two sides
//
//	adapters ──▶ port/in ──▶ [ application ] ──▶ port/out ──▶ adapters
//	(driving)    declared      implemented        declared     implemented
//	             by us         by us              by us        by infrastructure
//
// A driving port is the application's published API, and it is an interface
// whose methods are capabilities. The usual Go advice — let the consumer declare
// the interface it needs — is right when there is one consumer. This one has
// four or five: HTTP, gRPC, GraphQL, CLI, Lambda. Each declaring its own copy
// duplicates one contract five times and lets the copies drift.
//
// An interface rather than a struct of handler fields, because only an interface
// can be substituted. An adapter handed a struct can fill in its fields; it
// cannot replace the whole thing with a decorated implementation. "The adapter
// depends on the port" would have been true only in the way a struct literal is
// true.
//
// A driven port is a requirement: the application states what it needs and
// infrastructure complies. An interface declared in infrastructure and imported
// here has the dependency backwards, and will drag database concepts into use
// cases.
//
// Keep ports small. One with fifteen methods is a database driver wearing a
// costume; one with two can be faked in a test without a mocking library at all.
// Mocks are generated into a mocks package beside each interface — see
// ../../.mockery.yml.
//
// # Commands and queries
//
// The point of the split is not the two directories. It is that the read side
// does not go through the aggregate: a query has its own result declared beside
// it and its own reader in port/out, served from a denormalised table or a
// cache, shaped for what the caller displays rather than for what the business
// rules need.
//
// A query handler that loads an aggregate and reads its fields has gained
// nothing except two extra directories. If every query looks like that, the
// separation is decoration and should be removed.
//
// Commands change state and return as little as possible — an identifier. A
// command that returns a populated view is doing a query's job, which makes it
// impossible to route writes and reads differently later.
//
// # Service implements the driving port
//
// Every use case has the same shape, declared once generically for every use
// case there will ever be:
//
//	// usecase/usecase.go
//	type CommandHandler[C, R any] interface {
//		Handle(ctx context.Context, cmd C) (R, error)
//	}
//	type QueryHandler[Q, R any] interface {
//		Handle(ctx context.Context, q Q) (R, error)
//	}
//
// A service holds one such handler per capability and satisfies the driving port
// with one method per capability:
//
//	// service/version_service.go
//	type VersionService struct {
//		getVersion usecase.QueryHandler[query.GetVersion, query.GetVersionResult]
//	}
//
//	func (s VersionService) GetVersion(ctx context.Context, q query.GetVersion) (query.GetVersionResult, error) {
//		return s.getVersion.Handle(ctx, q)
//	}
//
// The field is typed by the generic interface rather than by a concrete use
// case, which is what lets a decorated use case — one wrapped with logging or
// tracing — be substituted for a plain one without the service noticing.
//
// An earlier draft made the driving port a struct of exported handler fields and
// forbade the service any method at all:
//
//	MUST NOT: give a service a method.
//
// That was the right rule for what it constrained. With the port a struct, a
// method could only forward, so every use case cost two mechanical edits in two
// files — and mechanical duplication is what drifts. The premise is gone: a
// struct is not a port, so the port became an interface, and an implementation
// of an interface has methods.
//
// The intent — that no logic accumulates in the service — survives, and needs
// restating with the same diff-visible quality:
//
//	MUST: a service method is a single delegation to a use case, and nothing
//	else.
//
// One line, one call. No branch, no fallback, no computed field. If a method
// needs a decision, that decision is a use case and belongs in usecase/. "One
// line per method" is as checkable in review as "no methods" was, and unlike
// "keep it thin" it is visible in a diff.
//
// Two honest limits:
//
//   - CommandHandler and QueryHandler are structurally identical, so nothing
//     stops a query being held in a field meant for a command. The distinction
//     is intent for readers, enforced in review, not by the compiler.
//
//   - Commands always return a result type, even when empty. A
//     DeleteWorkoutResult with no fields costs one line and keeps every use case
//     the same shape — which is what lets one decorator wrap all of them:
//     WithLogging[C, R](h usecase.CommandHandler[C, R], …) usecase.CommandHandler[C, R].
//
// # Do not leak domain objects outward
//
// Use cases accept commands and queries, and return results. Never a domain
// aggregate: an adapter holding a *workout.Workout can mutate it past its
// invariants, and a domain refactor then breaks every protocol at once.
//
// A result is deliberately incomplete — each one carries what a single use case
// needs and nothing more. A result that mirrors an aggregate field for field
// re-creates the coupling it was introduced to break, and becomes the model
// everyone actually passes around. Complete, bidirectional mapping belongs to
// the persistence model, private to a repository; see ../infrastructure/doc.go.
//
// # Failures are classified here, not translated
//
// A use case that fails returns an errors.Error: a Kind, and the error it wraps.
// It cannot pick a status, because it does not know whether it is being served
// over HTTP, gRPC or a CLI — so it classifies, and each adapter owns the table
// from Kind to its own vocabulary. HTTP's copy lives in
// ../interfaces/http/errmap and is described in ../interfaces/doc.go.
//
// Two properties of that type are load-bearing and easy to undo by accident:
//
//   - KindUnclassified is the zero value. A failure nobody classified, or one a
//     refactor forgot to classify, reaches an adapter as unrecognised and
//     becomes a 500 — rather than silently becoming, say, a NotFound.
//     Default-deny starts at the zero value, not at a review comment.
//
//   - Err is wrapped, never discarded, and callers add context with %w on the
//     way out. By the time an adapter sees a failure it has been wrapped at
//     least once, so only errors.As can still find the Kind.
//
// Which kinds exist, and why the list stops where it does, is argued where they
// are declared — see errors/errors.go. Two boundaries from that argument are
// worth knowing before reaching for a seventh kind. There is no kind for 400:
// shape validation belongs to the interfaces layer, where the request is
// rejected against api/openapi.yaml before a handler runs, so a use case can
// never produce "I cannot read this". And success is not classified here at all:
// 200, 201 and 204 come from an operation's contract and its presenter.
//
// # What exists
//
// One use case: GetVersion, the read side of the walking skeleton. It encodes no
// business rule — it reads a build stamp through a driven port and returns a
// result. Its purpose is structural: it proves the ports, the generic handlers
// and the service described above actually compile and connect.
//
// Nothing yet builds an errors.Error. The only failure this layer can currently
// produce is a driven port failing, and out.VersionProvider never does. The
// vocabulary is declared and tested ahead of its first caller so that the first
// use case with a real failure has a kind to reach for rather than an invention
// to make.
//
// Which commands and queries the tracker needs is a question CLARIFY answers.
// Everything named in the structure above beyond GetVersion is an illustration,
// not a plan.
package application
