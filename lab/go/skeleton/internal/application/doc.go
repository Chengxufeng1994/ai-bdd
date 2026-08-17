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
//	├── port/
//	│   ├── in/                       driving: what the outside may ask of us
//	│   │   └── handler.go            CommandHandler / QueryHandler, generic
//	│   └── out/                      driven: what we need from the outside
//	│       ├── workout_repository.go write side, traffics in domain types
//	│       ├── volume_reader.go      read side, traffics in dto
//	│       └── clock.go
//	├── command/                      semantic input; no json tags
//	│   └── record_workout.go         the command, and its mapping to domain
//	├── query/
//	│   └── workout_volume.go
//	├── dto/                          what we hand back; never a domain entity
//	│   ├── record_workout_result.go
//	│   └── volume_view.go
//	├── handler/                      implements port/in; orchestrates
//	│   ├── record_workout.go
//	│   └── workout_volume.go
//	├── assembler/                    domain -> dto, when it is more than one line
//	│   └── workout_assembler.go
//	└── service/                      facade: a bundle of handlers
//	    └── workout_service.go
//
// # dto is its own package, and that is load-bearing
//
// dto imports nothing, which is what lets port/out and query both depend on it.
// Putting a view type inside query instead would force port/out to import query
// in order to declare a reader returning it, while query imports port/out to get
// the interface — a cycle. With dto separate, read-side and write-side ports can
// both live in port/out and the classification stays uniform.
//
// # Ports have two sides
//
//	adapters ──▶ port/in ──▶ [ application ] ──▶ port/out ──▶ adapters
//	(driving)    declared      implemented        declared     implemented
//	             by us         by us              by us        by infrastructure
//
// A driving port is the application's published API. The usual Go advice — let
// the consumer declare the interface it needs — is right when there is one
// consumer. This one has four or five: HTTP, gRPC, GraphQL, CLI, Lambda. Each
// declaring its own copy duplicates one contract five times and lets the copies
// drift.
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
// does not go through the aggregate: a query has its own view in dto and its own
// reader in port/out, served from a denormalised table or a cache, shaped for
// what the caller displays rather than for what the business rules need.
//
// A query handler that loads an aggregate and reads its fields has gained
// nothing except two extra directories. If every query looks like that, the
// separation is decoration and should be removed.
//
// Commands change state and return as little as possible — an identifier. A
// command that returns a populated view is doing a query's job, which makes it
// impossible to route writes and reads differently later.
//
// # Service is a bundle of handlers
//
// The driving port is two generic interfaces, declared once for every use case
// there will ever be:
//
//	// port/in/handler.go
//	type CommandHandler[C, R any] interface {
//		Handle(ctx context.Context, cmd C) (R, error)
//	}
//	type QueryHandler[Q, R any] interface {
//		Handle(ctx context.Context, q Q) (R, error)
//	}
//
// A service exposes handlers as fields rather than wrapping each in a forwarding
// method:
//
//	// service/workout_service.go
//	type WorkoutService struct {
//		RecordWorkout in.CommandHandler[command.RecordWorkout, dto.RecordWorkoutResult]
//		DeleteWorkout in.CommandHandler[command.DeleteWorkout, dto.DeleteWorkoutResult]
//		Volume        in.QueryHandler[query.WorkoutVolume, dto.VolumeView]
//	}
//
//	// adapter side
//	res, err := s.RecordWorkout.Handle(ctx, cmd)
//
// Adding a use case is one field and one constructor line, in one file. The
// forwarding-method version costs two edits in two files every time, all of it
// mechanical — and mechanical duplication is what drifts. The field declaration
// is also the contract: it names the input and output types, so it cannot
// disagree with the handler the way a hand-written method can.
//
//	MUST NOT: give a service a method.
//
// That is the whole rule keeping the facade a facade. With no methods there is
// nowhere for logic to accumulate; anything needing a decision is a use case and
// belongs in handler/. Unlike "keep it thin", it is visible in a diff.
//
// Two honest limits:
//
//   - CommandHandler and QueryHandler are structurally identical, so nothing
//     stops a query being assigned to a command field. The distinction is intent
//     for readers, enforced in review, not by the compiler.
//
//   - Commands always return a result type, even when empty. A
//     DeleteWorkoutResult with no fields costs one line and keeps every handler
//     the same shape — which is what lets one decorator wrap all of them:
//     WithLogging[C, R](h in.CommandHandler[C, R], …) in.CommandHandler[C, R].
//
// # Do not leak domain objects outward
//
// Handlers accept commands and queries, and return dto. Never a domain
// aggregate: an adapter holding a *workout.Workout can mutate it past its
// invariants, and a domain refactor then breaks every protocol at once.
//
// dto is deliberately incomplete — each type carries what one use case needs and
// nothing more. A dto that mirrors an aggregate field for field re-creates the
// coupling it was introduced to break, and becomes the model everyone actually
// passes around. Complete, bidirectional mapping belongs to the persistence
// model, private to a repository; see ../infrastructure/doc.go.
//
// # Not yet populated
//
// No use case exists here. Which commands and queries the tracker needs is a
// question CLARIFY answers; the types named above are illustrations, not plans.
// See ../../README.md.
package application
