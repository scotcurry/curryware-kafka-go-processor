package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"
	"net"
	"strings"
	"time"
)

func ValidateDNSResolution(domain string, port string) ([]string, error) {

	logging.LogInfo(fmt.Sprintf("Validating DNS resolution for %s", domain))
	if strings.Contains(domain, ":") {
		parts := strings.Split(domain, ":")
		domain = parts[0]
		port = parts[1]
	}

	addresses, err := net.LookupHost(domain)
	if err != nil {
		logging.LogError(fmt.Sprintf("ValidateDNSResolution - DNS resolution failed for %s: %v", domain, err))
		return nil, fmt.Errorf("DNS resolution failed for %s: %v", domain, err)
	} else {
		logging.LogInfo(fmt.Sprintf("ValidateDNSResolution - DNS resolution successful for %s", domain))
		_, portError := validatePortOpen(domain, port, 3*time.Second)
		if portError != nil {
			logging.LogError(fmt.Sprintf("ValidateDNSResolution - Port %s is not open for %s: %v", port, domain, portError))
			return nil, fmt.Errorf("ValidateDNSResolution - port %s is not open for %s: %v", port, domain, portError)
		} else {
			logging.LogInfo(fmt.Sprintf("ValidateDNSResolution - Port %s is open for %s", port, domain))
		}
	}
	return addresses, nil
}

func validatePortOpen(domain string, port string, timeout time.Duration) (bool, error) {

	address := net.JoinHostPort(domain, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false, nil
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logging.LogError("ValidateDNSResolution - validatePortOpen - Error closing connection")
		}
	}(conn)
	return true, nil
}
