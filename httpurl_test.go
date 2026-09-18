package urlutil_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/urlutil"
)

func TestHTTPURL(t *testing.T) {
	var hu urlutil.HTTPURL
	require.Implements(t, (*fmt.Stringer)(nil), &hu)
	require.Implements(t, (*driver.Valuer)(nil), &hu)
	require.Implements(t, (*sql.Scanner)(nil), &hu)
	require.Implements(t, (*json.MarshalerTo)(nil), &hu)
	require.Implements(t, (*json.Marshaler)(nil), &hu)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &hu)
	require.Implements(t, (*json.Unmarshaler)(nil), &hu)
}

func TestNewHTTPURL(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *url.URL
			want string
		}{
			{
				"nil",
				nil,
				"invalid url: nil",
			},
			{
				"invalid host: empty",
				&url.URL{},
				"invalid url: invalid host: empty",
			},
			{
				"invalid scheme: ftp",
				&url.URL{
					Scheme: "ftp",
					Host:   "m0t0k1ch1.com",
				},
				"invalid url: invalid scheme: must be http or https",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := urlutil.NewHTTPURL(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *url.URL
			want string
		}{
			{
				"http",
				&url.URL{
					Scheme: "http",
					Host:   "m0t0k1ch1.com",
				},
				"http://m0t0k1ch1.com",
			},
			{
				"https",
				&url.URL{
					Scheme: "https",
					Host:   "m0t0k1ch1.com",
				},
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				hu, err := urlutil.NewHTTPURL(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, hu.String())
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		u := &url.URL{
			Scheme: "http",
			Host:   "m0t0k1ch1.com",
		}
		hu, err := urlutil.NewHTTPURL(u)
		require.NoError(t, err)
		require.Equal(t, "http://m0t0k1ch1.com", hu.String())

		u.Scheme = "https"
		u.Host = "m0t0k1ch2.com"
		require.Equal(t, "https://m0t0k1ch2.com", u.String())

		require.Equal(t, "http://m0t0k1ch1.com", hu.String())
	})
}

func TestMustNewHTTPURL(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *url.URL
			want string
		}{
			{
				"nil",
				nil,
				"invalid url: nil",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.PanicsWithError(t, tc.want, func() {
					urlutil.MustNewHTTPURL(tc.in)
				})
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *url.URL
			want string
		}{
			{
				"http",
				&url.URL{
					Scheme: "http",
					Host:   "m0t0k1ch1.com",
				},
				"http://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				hu := urlutil.MustNewHTTPURL(tc.in)
				require.Equal(t, tc.want, hu.String())
			})
		}
	})
}

func TestNewHTTPURLFromString(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"empty",
				"",
				"invalid url string: empty",
			},
			{
				"missing scheme",
				"://m0t0k1ch1.com",
				"invalid url string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				_, err := urlutil.NewHTTPURLFromString(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"http",
				"http://m0t0k1ch1.com",
				"http://m0t0k1ch1.com",
			},
			{
				"https",
				"https://m0t0k1ch1.com",
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				hu, err := urlutil.NewHTTPURLFromString(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, hu.String())
			})
		}
	})
}

func TestMustNewHTTPURLFromString(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"empty",
				"",
				"invalid url string: empty",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.PanicsWithError(t, tc.want, func() {
					urlutil.MustNewHTTPURLFromString(tc.in)
				})
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"http",
				"http://m0t0k1ch1.com",
				"http://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				hu := urlutil.MustNewHTTPURLFromString(tc.in)
				require.Equal(t, tc.want, hu.String())
			})
		}
	})
}

func TestHTTPURL_URL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   urlutil.HTTPURL
			want string
		}{
			{
				"http",
				urlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"),
				"http://m0t0k1ch1.com",
			},
			{
				"https",
				urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"),
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				u := tc.in.URL()
				require.Equal(t, tc.want, u.String())
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		hu := urlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com")
		u := hu.URL()
		require.Equal(t, "http://m0t0k1ch1.com", u.String())

		u.Scheme = "https"
		u.Host = "m0t0k1ch2.com"
		require.Equal(t, "https://m0t0k1ch2.com", u.String())

		require.Equal(t, "http://m0t0k1ch1.com", hu.String())
	})
}

func TestHTTPURL_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   urlutil.HTTPURL
			want string
		}{
			{
				"http",
				urlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"),
				"http://m0t0k1ch1.com",
			},
			{
				"https",
				urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"),
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestHTTPURL_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"nil",
				nil,
				"unsupported source: nil",
			},
			{
				"bool",
				true,
				"unsupported source type: bool",
			},
			{
				"string: empty",
				"",
				"invalid url string: empty",
			},
			{
				"string: missing scheme",
				"://m0t0k1ch1.com",
				"invalid url string",
			},
			{
				"string: invalid url.URL: invalid host: empty",
				"http://",
				"invalid url: invalid host: empty",
			},
			{
				"string: invalid url.URL: invalid scheme: ftp",
				"ftp://m0t0k1ch1.com",
				"invalid url: invalid scheme: must be http or https",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var hu urlutil.HTTPURL
				err := hu.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"string: http",
				"http://m0t0k1ch1.com",
				"http://m0t0k1ch1.com",
			},
			{
				"[]byte: https",
				[]byte("https://m0t0k1ch1.com"),
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var hu urlutil.HTTPURL
				err := hu.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, hu.String())
			})
		}
	})
}

func TestHTTPURL_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(urlutil.HTTPURL) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(hu urlutil.HTTPURL) ([]byte, error) {
				return json.Marshal(hu)
			},
		},
		{
			"MarshalJSON",
			func(hu urlutil.HTTPURL) ([]byte, error) {
				return hu.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   urlutil.HTTPURL
			want []byte
		}{
			{
				"http",
				urlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"),
				[]byte(`"http://m0t0k1ch1.com"`),
			},
			{
				"https",
				urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"),
				[]byte(`"https://m0t0k1ch1.com"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, enc := range encs {
					t.Run(enc.name, func(t *testing.T) {
						b, err := enc.marshal(tc.in)
						require.NoError(t, err)
						require.Equal(t, tc.want, b)
					})
				}
			})
		}
	})
}

func TestHTTPURL_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *urlutil.HTTPURL) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, hu *urlutil.HTTPURL) error {
				return json.Unmarshal(b, hu)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, hu *urlutil.HTTPURL) error {
				return hu.UnmarshalJSON(b)
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"empty",
				[]byte{},
				"",
			},
			{
				"null",
				[]byte(`null`),
				"unsupported json token kind: null",
			},
			{
				"string: empty",
				[]byte(`""`),
				"invalid url string: empty",
			},
			{
				"string: missing scheme",
				[]byte(`"://m0t0k1ch1.com"`),
				"invalid url string",
			},
			{
				"string: invalid host: empty",
				[]byte(`"http://"`),
				"invalid url: invalid host: empty",
			},
			{
				"string: invalid scheme: ftp",
				[]byte(`"ftp://m0t0k1ch1.com"`),
				"invalid url: invalid scheme: must be http or https",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var hu urlutil.HTTPURL
						err := dec.unmarshal(tc.in, &hu)
						require.ErrorContains(t, err, tc.want)
					})
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"string: http",
				[]byte(`"http://m0t0k1ch1.com"`),
				"http://m0t0k1ch1.com",
			},
			{
				"string: https",
				[]byte(`"https://m0t0k1ch1.com"`),
				"https://m0t0k1ch1.com",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var hu urlutil.HTTPURL
						err := dec.unmarshal(tc.in, &hu)
						require.NoError(t, err)
						require.Equal(t, tc.want, hu.String())
					})
				}
			})
		}
	})
}
