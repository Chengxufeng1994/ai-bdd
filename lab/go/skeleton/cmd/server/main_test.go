package main

import (
	"context"
	"errors"
	"testing"

	apperrors "skeleton/internal/application/errors"
	"skeleton/internal/application/usecase/query"
	"skeleton/internal/interfaces/http/errmap"
	"skeleton/pkg/i18n"
)

// unreadableStamp stands in for out.VersionProvider, failing the way the port
// declares it may. The port has one method, so a hand-written stub is enough —
// see query_test's stubProvider for the same call.
type unreadableStamp struct{}

// Version satisfies out.VersionProvider by always failing, which is the only
// way to make the use case raise its classified error.
func (unreadableStamp) Version(context.Context) (string, error) {
	return "", errors.New("build stamp unreadable")
}

// pkg/i18n renders a key it cannot translate as the key itself, deliberately:
// a gap that reads as "version.unavailable" on screen is a gap somebody
// notices, where a default-locale fallback would ship looking finished. That
// makes a missing translation loud but a *wrong key* silent — rename the key
// at the call site that raises the failure, or the entry in this binary's
// table, and every affected client starts receiving a bare identifier as the
// human-readable title of its error, with nothing failing anywhere.
//
// So this walks the real path end to end: the use case raises the error, and
// the table this binary actually ships renders it. It lives in package main
// because that is where the table lives, and the table lives there because
// which locales a deployment carries is a property of the deployment.
//
// The expected sentence is read from the table rather than repeated here. What
// is being checked is that the two sides agree, and a literal copy would make
// rewording one message a two-file edit while adding no coverage — an
// inconsistent rename fails on the lookup, a consistent one is not a defect.
func TestTheShippedTableTranslatesTheFailureThisBinaryCanRaise(t *testing.T) {
	_, err := query.NewGetVersion(unreadableStamp{}).Handle(context.Background(), query.GetVersion{})
	if err == nil {
		t.Fatal("want the use case to fail, so its classified error can be rendered")
	}

	var classified apperrors.Error
	if !errors.As(err, &classified) {
		t.Fatalf("want a classified apperrors.Error, got %T: %v", err, err)
	}

	want := messages["en"][classified.MessageKey]
	if want == "" {
		t.Fatalf("the shipped table carries no English message for %q, the only key this binary can raise", classified.MessageKey)
	}

	_, problem := errmap.ToProblem(err, "en", i18n.NewBundle(messages))

	if problem.Title == classified.MessageKey {
		t.Errorf("Title rendered as the bare key %q: the table and the call site no longer agree", classified.MessageKey)
	}
	if problem.Title != want {
		t.Errorf("Title: want the shipped English sentence %q, got %q", want, problem.Title)
	}
}
