package urlutil

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/url"
)

// HTTPURL represents a HTTP(S) URL.
type HTTPURL struct {
	u url.URL
}

// NewHTTPURL returns a new HTTPURL.
func NewHTTPURL(u *url.URL) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setURL(u); err != nil {
		return HTTPURL{}, err
	}

	return hu, nil
}

// MustNewHTTPURL panics if the input is invalid.
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

// NewHTTPURLFromString returns a new HTTPURL from a string.
func NewHTTPURLFromString(s string) (HTTPURL, error) {
	var hu HTTPURL
	if err := hu.setString(s); err != nil {
		return HTTPURL{}, err
	}

	return hu, nil
}

// MustNewHTTPURLFromString panics if the input is invalid.
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

// String implements fmt.Stringer.
// It returns the value as a string.
func (hu HTTPURL) String() string {
	return hu.u.String()
}

// Value implements driver.Valuer.
// It returns the value as a string.
func (hu HTTPURL) Value() (driver.Value, error) {
	return hu.String(), nil
}

// Scan implements sql.Scanner.
// It accepts a string or []byte.
func (hu *HTTPURL) Scan(src any) error {
	if src == nil {
		return errors.New("invalid source: nil")
	}

	var s string
	{
		switch v := src.(type) {
		case string:
			s = v
		case []byte:
			s = string(v)
		default:
			return fmt.Errorf("unsupported source type: %T", src)
		}
	}

	if err := hu.setString(s); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	return nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes the value as a JSON string.
func (hu HTTPURL) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, hu.String())
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It accepts a JSON string.
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
