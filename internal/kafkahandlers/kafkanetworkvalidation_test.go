package kafkahandlers

import (
	"testing"
)

func TestValidateDNSResolution(t *testing.T) {

	domain := "kafka.curryware.org:9092"
	_, err := ValidateDNSResolution(domain, "")
	if err != nil {
		t.Errorf("Error validating DNS resolution: %s", err)
	}

	domain = "postgres.curryware.org"
	postgresPort := "5432"
	_, err = ValidateDNSResolution(domain, postgresPort)
	if err != nil {
		t.Errorf("Error validating DNS resolution: %s", err)
	}
}
