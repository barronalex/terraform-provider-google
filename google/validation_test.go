package google

import (
	"fmt"
	"testing"
)

func TestValidateGCPName(t *testing.T) {
	x := []GCPNameTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "f"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a capital", Value: "Foobar", ExpectError: true},
		{TestName: "starts with a number", Value: "1foobar", ExpectError: true},
		{TestName: "has an underscore", Value: "foo_bar", ExpectError: true},
		{TestName: "too long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoob", ExpectError: true},
	}

	es := testGCPNames(x)
	if len(es) > 0 {
		t.Errorf("Failed to validate GCP names: %v", es)
	}
}

func TestValidateRFC1918Network(t *testing.T) {
	x := []RFC1918NetworkTestCase{
		// No errors
		{TestName: "valid 10.x", CIDR: "10.0.0.0/8", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid 172.x", CIDR: "172.16.0.0/16", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid 192.x", CIDR: "192.168.0.0/32", MinPrefix: 0, MaxPrefix: 32},
		{TestName: "valid, bounded 10.x CIDR", CIDR: "10.0.0.0/8", MinPrefix: 8, MaxPrefix: 32},
		{TestName: "valid, bounded 172.x CIDR", CIDR: "172.16.0.0/16", MinPrefix: 12, MaxPrefix: 32},
		{TestName: "valid, bounded 192.x CIDR", CIDR: "192.168.0.0/32", MinPrefix: 16, MaxPrefix: 32},

		// With errors
		{TestName: "empty CIDR", CIDR: "", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "missing mask", CIDR: "10.0.0.0", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "invalid CIDR", CIDR: "10.1.0.0/8", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
		{TestName: "valid 10.x CIDR with lower bound violation", CIDR: "10.0.0.0/8", MinPrefix: 16, MaxPrefix: 32, ExpectError: true},
		{TestName: "valid 10.x CIDR with upper bound violation", CIDR: "10.0.0.0/24", MinPrefix: 8, MaxPrefix: 16, ExpectError: true},
		{TestName: "valid public CIDR", CIDR: "8.8.8.8/32", MinPrefix: 0, MaxPrefix: 32, ExpectError: true},
	}

	es := testRFC1918Networks(x)
	if len(es) > 0 {
		t.Errorf("Failed to validate RFC1918 Networks: %v", es)
	}
}

type GCPNameTestCase struct {
	TestName    string
	Value       string
	ExpectError bool
}

type RFC1918NetworkTestCase struct {
	TestName    string
	CIDR        string
	MinPrefix   int
	MaxPrefix   int
	ExpectError bool
}

func testGCPNames(cases []GCPNameTestCase) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testGCPName(c)...)
	}

	return es
}

func testGCPName(testCase GCPNameTestCase) []error {
	_, es := validateGCPName(testCase.Value, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		} else {
			return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with string \"%s\"", testCase.TestName, testCase.Value)}
		}
	}

	return es
}

func testRFC1918Networks(cases []RFC1918NetworkTestCase) []error {
	es := make([]error, 0)
	for _, c := range cases {
		es = append(es, testRFC1918Network(c)...)
	}

	return es
}

func testRFC1918Network(testCase RFC1918NetworkTestCase) []error {
	f := validateRFC1918Network(testCase.MinPrefix, testCase.MaxPrefix)
	_, es := f(testCase.CIDR, testCase.TestName)
	if testCase.ExpectError {
		if len(es) > 0 {
			return nil
		}
		return []error{fmt.Errorf("Didn't see expected error in case \"%s\" with CIDR=\"%s\" MinPrefix=%v MaxPrefix=%v",
			testCase.TestName, testCase.CIDR, testCase.MinPrefix, testCase.MaxPrefix)}
	}

	return es
}
