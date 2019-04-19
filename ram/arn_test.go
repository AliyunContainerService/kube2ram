package ram

import (
	"testing"
)

func TestIsValidBaseARN(t *testing.T) {
	arns := []string{
		"acs:ram::123456789012:role/path",
		"acs:ram::123456789012:role/path/",
		"acs:ram::123456789012:role/path/sub-path",
		"acs:ram::123456789012:role/path/sub_path",
		"acs:ram::123456789012:role/subdomain.domain",
		"acs:ram::123456789012:role",
		"acs:ram::123456789012:role/",
		"acs:ram::123456789012:role-part",
		"acs:ram::123456789012:role_part",
		"acs:ram::123456789012:role_123",
	}
	for _, arn := range arns {
		if !IsValidBaseARN(arn) {
			t.Errorf("%s is a valid base arn", arn)
		}
	}
}

func TestIsValidBaseARNWithInvalid(t *testing.T) {
	arns := []string{
		"acs:ram::123456789012::role/path",
		"acs:ram:us-east-1:123456789012:role/path",
		"acs:oss::123456789012:role/path",
		"acs:ram::123456789012:role/$",
		"acs:ram::12345-6789012:role/",
		"acs:ram::abcdef:role/",
		"acs:ram:::role",
	}
	for _, arn := range arns {
		if IsValidBaseARN(arn) {
			t.Errorf("%s is not a valid base arn", arn)
		}
	}
}
