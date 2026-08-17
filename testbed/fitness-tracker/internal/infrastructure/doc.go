// Package infrastructure implements the ports declared by application:
// persistence, clocks, external clients.
//
// It depends on domain and application. Nothing imports it except cmd, which
// wires the concrete implementations in at startup.
//
// # This is where complete, bidirectional mapping belongs
//
// A repository needs a model that mirrors its aggregate exactly: storage has to
// hold every field and give it all back unchanged. That model lives here, kept
// unexported inside the adapter that uses it:
//
//	// persistence/workout_record.go
//	type workoutRecord struct{ /* every field */ }
//
//	func (r workoutRecord) toDomain() (*workout.Workout, error)
//	func fromDomain(w *workout.Workout) workoutRecord
//
// Because it is private, an aggregate gaining a field touches this one file.
//
// Do not put a complete WorkoutDTO in application. It looks like the same idea
// but does the opposite job:
//
//	                 application DTO                persistence model
//	completeness     partial — one use case's needs  total — every field
//	direction        one way in, or an ID out        both ways
//	count            one per use case                one per aggregate
//	visibility       public, adapters see it         private to the repository
//
// A DTO that mirrors the aggregate changes whenever the aggregate does, which
// re-creates the coupling the DTO was introduced to break — moved one step out
// rather than removed. It also forces every use case to handle fields it does
// not need: recording a workout should not accept an ID the server generates,
// and "must ignore this field" is a rule that gets forgotten.
//
// The slower failure matters more. A complete public DTO becomes the model
// everyone actually passes around, and the aggregate degrades into a formality
// on the way to the database. Anaemic domain models are rarely chosen; they grow
// from exactly this.
//
// # Adapters own their own vocabulary
//
// SQL column names, JSON tags, retry policies and connection settings stay here.
// If a database concept appears in application or domain, the dependency has
// been inverted somewhere.
//
// # Not yet populated
//
// No adapter exists. Ports are declared when a use case needs them, and no use
// case exists until CLARIFY has run. See ../../README.md.
package infrastructure
