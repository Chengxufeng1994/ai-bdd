// Package i18n renders message keys into human-readable text.
//
// It sits in pkg/ for the same reason pkg/log does: translation is a concern
// with no layer of its own. Adapters need it because they are the only place
// that knows a request's locale; the application layer must never see one.
//
// # Why a map and not x/text
//
// The Translator interface is what the architecture depends on; this Bundle is
// the simplest thing that satisfies it. golang.org/x/text/message is
// catalog-driven and would become a direct dependency with its own conventions
// — a choice this project has not made. Swapping Bundle for it later touches
// this file and nothing else.
//
// # A missing translation renders as its key
//
// Not as a default-locale fallback. A fallback looks like a finished message,
// so a gap ships unnoticed; "workout.not_found" on screen never does. The same
// rule covers a template whose params are missing: emitting the template with
// an unresolved placeholder would read as a product bug rather than a gap at
// the call site.
package i18n

import (
	"fmt"
	"strings"
)

// Translator renders a message key in a locale.
//
// It never fails. Every failure mode — an unknown locale, an unknown key, a
// missing parameter — resolves to the key, so a caller rendering an error can
// never itself fail and re-enter error rendering.
type Translator interface {
	Translate(locale, key string, params map[string]any) string
}

// Bundle is a Translator backed by an in-memory table of locale to messages.
type Bundle struct {
	messages map[string]map[string]string
}

// Compile-time check that Bundle still satisfies the contract.
var _ Translator = Bundle{}

// NewBundle returns a Bundle over the given locale-to-messages table.
func NewBundle(messages map[string]map[string]string) Bundle {
	return Bundle{messages: messages}
}

// Translate renders key in locale, interpolating params into {placeholders}.
//
// It returns key unchanged when the locale is unknown, the key is unknown, or
// the template needs a parameter that params does not supply.
func (b Bundle) Translate(locale, key string, params map[string]any) string {
	template, ok := b.messages[locale][key]
	if !ok {
		return key
	}

	// The scan walks the template rather than the rendered result. Checking the
	// result instead would read a brace that arrived inside a parameter's value
	// as an unresolved placeholder, discard a correctly rendered message, and
	// return the key for input that was in fact complete. Walking the template
	// once also removes any dependence on map iteration order.
	var out strings.Builder
	for {
		open := strings.Index(template, "{")
		if open < 0 {
			out.WriteString(template)
			return out.String()
		}

		closeAt := strings.Index(template[open:], "}")
		if closeAt < 0 {
			// An unterminated brace is a malformed template, which is a gap in
			// the message table rather than at the call site — but it is still
			// a gap, and this package's answer to every gap is the key.
			return key
		}
		closeAt += open

		value, ok := params[template[open+1:closeAt]]
		if !ok {
			return key
		}

		out.WriteString(template[:open])
		out.WriteString(toString(value))
		template = template[closeAt+1:]
	}
}

// toString renders a param value without pulling in fmt's reflection for the
// common case of a string.
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
