package urlutil

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
)

var (
	_ fmt.Stringer         = HTTPURL{}
	_ driver.Valuer        = HTTPURL{}
	_ sql.Scanner          = (*HTTPURL)(nil)
	_ json.MarshalerTo     = HTTPURL{}
	_ json.Marshaler       = HTTPURL{}
	_ json.UnmarshalerFrom = (*HTTPURL)(nil)
	_ json.Unmarshaler     = (*HTTPURL)(nil)
)

// HTTPURL represents an HTTP(S) URL.
type HTTPURL struct {
	u url.URL
}

// NewHTTPURL returns a new [HTTPURL] from a [url.URL].
func NewHTTPURL(u *url.URL) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setURL(u); err != nil {
		return HTTPURL{}, err
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
		return errors.New("invalid url: nil")
	}
	if u.Host == "" {
		return errors.New("invalid url: invalid host: empty")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid url: invalid scheme: must be http or https")
	}

	hu.u = *u

	return nil
}

// NewHTTPURLFromString returns a new [HTTPURL] from a string.
func NewHTTPURLFromString(s string) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setString(s); err != nil {
		return HTTPURL{}, err
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
		return errors.New("invalid url string: empty")
	}

	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid url string: %w", err)
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
		return errors.New("invalid source: nil")
	}

	var s string
	{
		switch src := src.(type) {
		case string:
			s = src
		case []byte:
			s = string(src)
		default:
			return fmt.Errorf("unsupported source type: %T", src)
		}
	}

	return hu.setString(s)
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes hu as a quoted string and writes it to enc.
func (hu HTTPURL) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, hu.String())
}

// MarshalJSON implements [json.Marshaler].
// It is like [HTTPURL.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (hu HTTPURL) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := hu.MarshalJSONTo(jsontext.NewEncoder(&buf)); err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a JSON string from dec into hu.
func (hu *HTTPURL) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == jsontext.KindNull {
		return errors.New("invalid json string: null")
	}

	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return fmt.Errorf("invalid json string: %w", err)
	}

	return hu.setString(s)
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [HTTPURL.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (hu *HTTPURL) UnmarshalJSON(b []byte) error {
	return hu.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(b)))
}
