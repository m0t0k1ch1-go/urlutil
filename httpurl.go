package urlutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
)

var (
	_ fmt.Stringer             = HTTPURL{}
	_ driver.Valuer            = HTTPURL{}
	_ sql.Scanner              = &HTTPURL{}
	_ encoding.TextMarshaler   = HTTPURL{}
	_ json.MarshalerTo         = HTTPURL{}
	_ json.Marshaler           = HTTPURL{}
	_ encoding.TextUnmarshaler = &HTTPURL{}
	_ json.UnmarshalerFrom     = &HTTPURL{}
	_ json.Unmarshaler         = &HTTPURL{}
)

// HTTPURL represents an HTTP(S) URL.
type HTTPURL struct {
	u url.URL
}

// NewHTTPURL returns a new [HTTPURL] from a [url.URL].
func NewHTTPURL(u *url.URL) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setURL(u); err != nil {
		return HTTPURL{}, fmt.Errorf("invalid url: %w", err)
	}

	return hu, nil
}

// MustNewHTTPURL is like [NewHTTPURL] but panics if the input is invalid.
func MustNewHTTPURL(u *url.URL) HTTPURL {
	hu, err := NewHTTPURL(u)
	if err != nil {
		panic(err)
	}

	return hu
}

func (hu *HTTPURL) setURL(u *url.URL) error {
	if u == nil {
		return errors.New("nil")
	}
	if u.Host == "" {
		return errors.New("invalid host: empty")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid scheme: must be http or https")
	}

	hu.u = *u

	return nil
}

// NewHTTPURLFromString returns a new [HTTPURL] from a string.
func NewHTTPURLFromString(s string) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setString(s); err != nil {
		return HTTPURL{}, fmt.Errorf("invalid url string: %w", err)
	}

	return hu, nil
}

// MustNewHTTPURLFromString is like [NewHTTPURLFromString] but panics if the input is invalid.
func MustNewHTTPURLFromString(s string) HTTPURL {
	hu, err := NewHTTPURLFromString(s)
	if err != nil {
		panic(err)
	}

	return hu
}

func (hu *HTTPURL) setString(s string) error {
	if len(s) == 0 {
		return errors.New("empty")
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	return hu.setURL(u)
}

// URL returns a copy of the underlying [url.URL].
func (hu HTTPURL) URL() *url.URL {
	return hu.u.Clone()
}

// String implements [fmt.Stringer].
// It encodes hu as a string.
func (hu HTTPURL) String() string {
	return hu.u.String()
}

// Value implements [driver.Valuer].
// It encodes hu as a string.
func (hu HTTPURL) Value() (driver.Value, error) {
	return hu.String(), nil
}

// Scan implements [sql.Scanner].
// It decodes a string or []byte into hu.
func (hu *HTTPURL) Scan(src any) error {
	if src == nil {
		return errors.New("unsupported source: nil")
	}

	switch src := src.(type) {
	case string:
		if err := hu.setString(src); err != nil {
			return fmt.Errorf("invalid string source: %w", err)
		}

		return nil

	case []byte:
		if err := hu.setString(string(src)); err != nil {
			return fmt.Errorf("invalid bytes source: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", src)
	}
}

// MarshalText implements [encoding.TextMarshaler].
// It encodes hu as a string.
func (hu HTTPURL) MarshalText() ([]byte, error) {
	return []byte(hu.String()), nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes hu as a quoted string and writes it to enc.
func (hu HTTPURL) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, hu.String())
}

// MarshalJSON implements [json.Marshaler].
// It is like [HTTPURL.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (hu HTTPURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(hu)
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It decodes a string into hu.
func (hu *HTTPURL) UnmarshalText(text []byte) error {
	if err := hu.setString(string(text)); err != nil {
		return fmt.Errorf("invalid string: %w", err)
	}

	return nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a quoted string from dec into hu.
func (hu *HTTPURL) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindString:
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		return hu.UnmarshalText([]byte(s))

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [HTTPURL.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (hu *HTTPURL) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, hu)
}
