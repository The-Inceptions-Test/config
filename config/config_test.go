// Copyright © by Jeff Foley 2017-2023. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"testing"
)

func TestCheckSettings(t *testing.T) {
	c := NewConfig()

	err := c.CheckSettings()

	if err != nil {
		t.Errorf("Error checking settings.\n%v", err)
	}

}
func TestDomainRegex(t *testing.T) {
	c := NewConfig()
	got := c.DomainRegex("owasp.org")

	if got != nil {
		t.Errorf("Error with DomainRegex.\n%v", got)
	}
}

func TestAddDomains(t *testing.T) {
	c := NewConfig()
	example := "owasp.org/test"
	list := []string{"owasp.org", "google.com", "yahoo.com"}
	c.AddDomains(list...)
	got := c.Domains()
	sort.Strings(list)
	sort.Strings(got)
	c.AddDomains(list...)

	if !reflect.DeepEqual(list, got) {
		t.Errorf("Domains do not match.\nWanted:%v\nGot:%v\n", list, got)
	}
	t.Run("Testing AddDomain...", func(t *testing.T) {

		c.AddDomain(example)
		want := true
		got := false
		for _, l := range c.Scope.ProvidedNames {
			if example == l {
				got = true
			}
		}
		if got != want {
			t.Errorf("Expected:%v\nGot:%v", want, got)
		}
		t.Run("Testing Domains...", func(t *testing.T) {
			if c.Domains() == nil {
				t.Errorf("No domains in current configuration.")
			}

			if len(c.Domains()) <= 0 {
				t.Errorf("Failed to populate c.domains.\nLength:%v", len(c.Domains()))
			}

		})

		t.Run("Testing IsDomainInScope...", func(t *testing.T) {

			if !c.IsDomainInScope(example) {
				t.Errorf("Domain is considered out of scope.\nExample:%v\nGot:%v,\nWant:%v", example, got, want)
			}
		})

		t.Run("Testing WhichDomain...", func(t *testing.T) {

			if example != c.WhichDomain(example) {
				t.Errorf("Failed to find example.\nExample:%v\nGot:%v", example, got)
			}
		})
	})
}

func TestIsAddressInScope(t *testing.T) {
	c := NewConfig()
	example := "192.0.2.1"
	c.Scope.Addresses = append(c.Scope.Addresses, net.ParseIP(example).String())
	c.Scope.toIPs(c.Scope.Addresses)
	if !c.IsAddressInScope(example) {
		t.Errorf("Failed to find address %v in scope.\nAddress List:%v", example, c.Scope.ProvidedNames)
	}
}

func TestBlacklist(t *testing.T) {
	c := NewConfig()
	example := "owasp.org"
	c.Scope.Blacklist = append(c.Scope.Blacklist, example)
	got := c.Blacklisted(example)
	want := true

	if got != want {
		t.Errorf("Failed to find %v in blacklist.", example)
	}
}

func TestLoadSettings(t *testing.T) {
	c := NewConfig()
	path := "/home/adem/go/src/amass/examples/config.yaml"
	err := c.LoadSettings(path)
	if err != nil {
		t.Errorf("Config file failed to load.")
		fmt.Println(err)
	}
}

func TestConfigCheckSettings(t *testing.T) {
	type fields struct {
		c *Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "brute-force & passive set",
			fields: fields{
				&Config{BruteForcing: true, Passive: true},
			},
			wantErr: true,
		},
		{
			name: "brute-force & empty wordlist - load default wordlist",
			fields: fields{
				&Config{BruteForcing: true, Passive: false, Wordlist: []string{}},
			},
			wantErr: false,
		},
		{
			name: "active & passive enumeration set",
			fields: fields{
				&Config{Passive: true, Active: true},
			},
			wantErr: true,
		},
		{
			name: "alterations set with empty alt-wordlist - load default alt-wordlist",
			fields: fields{
				&Config{Alterations: true, AltWordlist: []string{}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.c.CheckSettings(); (err != nil) != tt.wantErr {
				t.Errorf("Config.CheckSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigGetListFromFile(t *testing.T) {
	var list = "/home/adem/scripts_and_tools/SecLists/Discovery/DNS/shubs-subdomains.txt"
	if _, err := GetListFromFile(list); err != nil {
		t.Errorf("GetListFromFile() error = %v", err)
	}
}
