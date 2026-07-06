package opentable

import "testing"

func TestExtractSha256Hash(t *testing.T) {
	fresh := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "compact persisted-query body",
			body: `{"operationName":"RestaurantsAvailability","variables":{"rid":1389331},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"` + fresh + `"}}}`,
			want: fresh,
		},
		{
			name: "no hash present",
			body: `{"operationName":"RestaurantsAvailability","variables":{}}`,
			want: "",
		},
		{
			name: "malformed short hash",
			body: `{"extensions":{"persistedQuery":{"sha256Hash":"deadbeef"}}}`,
			want: "",
		},
		{
			name: "empty body",
			body: "",
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := extractSha256Hash(tc.body); got != tc.want {
				t.Fatalf("extractSha256Hash(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}
