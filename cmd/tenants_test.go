package main

import (
	"reflect"
	"testing"
)

func TestTenantHosts(t *testing.T) {
	cases := []struct {
		host, hosts string
		want        []string
		invalid     bool
	}{
		{host: "sso.st29.ru", want: []string{"sso.st29.ru"}},
		{hosts: "SSO.st29.ru, timetrack.st29.ru", want: []string{"sso.st29.ru", "timetrack.st29.ru"}},
		{host: "a", hosts: "b", invalid: true},
		{hosts: "a,A", invalid: true},
		{hosts: "a,", invalid: true},
	}
	for _, tc := range cases {
		got, err := tenantHosts(tc.host, tc.hosts)
		if (err != nil) != tc.invalid || !tc.invalid && !reflect.DeepEqual(got, tc.want) {
			t.Errorf("tenantHosts(%q, %q) = %v, %v", tc.host, tc.hosts, got, err)
		}
	}
}
