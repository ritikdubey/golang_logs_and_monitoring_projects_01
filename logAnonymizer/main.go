package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
var panRegex = regexp.MustCompile(`(?i)\b[A-Z]{5}\d{4}[A-Z]\b`)
var aadhaarRegex = regexp.MustCompile(`\b\d{4}[- ]?\d{4}[- ]?\d{4}\b|\b\d{4}[- ]?\d{4}[- ]?\d{1,4}\b`)
var creditCardRegex = regexp.MustCompile(`\b(?:\d{4}[- ]?\d{6}[- ]?\d{5}|\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4})\b`)
var ipv6LooseRegex = regexp.MustCompile(`(?i)\b[0-9a-f:]+:[0-9a-f:]+\b`)
var phoneRegex = regexp.MustCompile(`(?:\+?\d{1,3}[- ]?)?\(?\d{3}\)?[- ]?\d{3}[- ]?\d{4}`)

func main() {

	inputFile, err := os.Open("application.log")
	if err != nil {
		fmt.Printf("error while reading file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("RedactedLogs.log")
	if err != nil {
		fmt.Printf("error while creating output file: %v", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)

		line = replaceIPsInLine(line)
		line = replaceEmailInLine(line)
		line = replacePANInLine(line)
		line = replacePhoneNumberInLine(line)
		line = replaceCreditCardInLine(line)
		line = replaceAdhaarInLine(line)
		line = replaceIPV6InLine(line)

		fmt.Fprintf(writer, "%v\n", line)

	}

	if err = scanner.Err(); err != nil {
		fmt.Printf("error while scan: %v", err)
	}

	if err = writer.Flush(); err != nil {
		fmt.Printf("error while flush: %v", err)
	}

}

func replaceIPsInLine(line string) string {
	return ipRegex.ReplaceAllStringFunc(line, func(matched string) string {
		if net.ParseIP(matched) != nil {
			return "[REDACTED IP ADDRESS]"
		}
		return matched
	})
}

func replaceEmailInLine(line string) string {
	return emailRegex.ReplaceAllStringFunc(line, func(matched string) string {
		index := strings.Index(matched, "@")
		return "[REDACTED EMAIL]" + matched[index:]
	})
}

func replacePANInLine(line string) string {
	return panRegex.ReplaceAllStringFunc(line, func(matched string) string {
		return "[READACTED PAN NUMBER]"
	})
}

func replaceAdhaarInLine(line string) string {
	return aadhaarRegex.ReplaceAllStringFunc(line, func(matched string) string {
		return "[REDACTED AADHAAR NUMBER]"
	})
}

func replaceCreditCardInLine(line string) string {
	return creditCardRegex.ReplaceAllStringFunc(line, func(matched string) string {
		return "[REDACTED CREDIT CARD NUMBER]"
	})
}

func replaceIPV6InLine(line string) string {
	return ipv6LooseRegex.ReplaceAllStringFunc(line, func(matched string) string {
		ip := net.ParseIP(matched)

		if ip != nil && ip.To4() == nil {
			return "[REDACTED IPV6 ADDRESS]"
		}
		return matched
	})
}

func replacePhoneNumberInLine(line string) string {
	return phoneRegex.ReplaceAllStringFunc(line, func(matched string) string {
		return "[REDACTED PHONE NUMBER]"
	})
}
