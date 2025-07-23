package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	ICINGA_OK       int = 0
	ICINGA_WARNING  int = 1
	ICINGA_CRITICAL int = 2
	ICINGA_UNKNOWN  int = 3
)

func main() {
	certPath := flag.String("p", "", "Path to certificate file")
	daysWarning := flag.Int("w", 21, "Warning threshold days validity")
	daysCritical := flag.Int("c", 14, "Critical threshold days validity")
	flag.Parse()

	certFile, err := os.ReadFile(*certPath)
	if err != nil {
		fmt.Println("[ERROR] Reading PEM file:", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	pemBlock, _ := pem.Decode(certFile)
	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		fmt.Printf("[ERROR] Failed to decode PEM block containing certificate")
		os.Exit(ICINGA_UNKNOWN)
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse certificate: %v", err)
		os.Exit(ICINGA_UNKNOWN)
	}

	now := time.Now()
	daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)

	if daysRemaining <= *daysCritical {
		fmt.Printf("[CRITICAL] Certificate only valid for %d days", daysRemaining)
		os.Exit(ICINGA_CRITICAL)

	} else if *daysCritical <= daysRemaining && daysRemaining <= *daysWarning {
		fmt.Printf("[Warning] Certificate only valid for %d days", daysRemaining)
		os.Exit(ICINGA_WARNING)
	} else {
		fmt.Printf("[OK] Certificate will expire on %v", cert.NotAfter)
		os.Exit(ICINGA_OK)
	}
}
