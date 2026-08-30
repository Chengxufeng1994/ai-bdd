package in

// VersionService bundles the use cases that report build information.
//
// It lives here, beside the port it bundles, rather than in a separate
// package: it is the application's published API, so an adapter holding it
// depends on contracts and nothing else. It is a struct with fields and no
// methods because every field is already a driving port — there is no
// behaviour left to abstract, and an interface here would only add a
// GetVersion() method and an indirection between the field and the call.
//
// Adding a use case is one field and one line at the composition root. The
// forwarding-method alternative costs two edits in two files every time, all
// of it mechanical — and mechanical duplication is what drifts.
type VersionService struct {
	// GetVersion reports the version the running service was built from.
	GetVersion GetVersionUseCase
}
