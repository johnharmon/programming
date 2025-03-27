package main

import (
	"flag"
	"fmt"
	ni "net/netip"
	"os"
	"os/exec"
	"sync"
)

type IPResult struct {
	ip        ni.Addr
	responded bool
}

type ImproperUsage int

const (
	InvalidSubnet ImproperUsage = iota
	TooManyArguments
	MissingArgument
)

type ProcessingError int

const (
	GenericError ProcessingError = iota + 1
)

func (iu ImproperUsage) String(subnet string) string {
	switch iu {
	case InvalidSubnet:
		return fmt.Sprintf("The subnet %s is not valid\n", subnet)
	case TooManyArguments:
		return "Too many arguments/flags provided\n"
	case MissingArgument:
		return "Missing required argument <subnet>\n"
	default:
		return "Error processing command line arguments\n"
	}
}

func printUsage(subnet string, reason ImproperUsage) {
	usage := "Usage:\n"
	usage += "    ping <subnet>\n"
	usage += "    <subnet> is a valid IPv4 CIDR notation subnet\n"
	usage += "    No other positional arguments or flags are supported\n"
	usage += "    " + reason.String(subnet)
	fmt.Print(usage)
}

func isValidSubnet(ip string) (prefix ni.Prefix, valid bool) {
	prefix, err := ni.ParsePrefix(ip)
	if err != nil {
		return prefix, false
	}
	return prefix, true
}

func getSubnetIPs(prefix ni.Prefix) (ips []*IPResult) {
	for ip := prefix.Addr(); prefix.Contains(ip); ip = ip.Next() {
		ips = append(ips, &IPResult{ip: ip, responded: false})
	}
	return ips
}

func pingIP(ipResult *IPResult, count string, wait string, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ping", "-c", count, "-W", wait, ipResult.ip.String())
	_, err := cmd.Output()
	//output, err := cmd.Output()
	//if strings.Contains(string(output), fmt.Sprintf("%s packets transmitted, %s packets received", count, count)) && err == nil {
	if err == nil {
		ipResult.responded = true
	}
}

func scanSubnet(prefix ni.Prefix, count string, wait string) (results []*IPResult) {
	ips := getSubnetIPs(prefix)
	wg := sync.WaitGroup{}
	for _, ip := range ips {
		wg.Add(1)
		go pingIP(ip, count, wait, &wg)
	}
	wg.Wait()
	return ips
}

func main() {
	var subnet string
	var count string
	var wait string
	fs := flag.NewFlagSet("default", flag.ExitOnError)
	fs.StringVar(&count, "c", "1", "Number of packets to send, equivalent to the same flag for the ping utility\n")
	fs.StringVar(&wait, "w", "1", "How many seconds to wait for a response, equivalent to the -W flag for the ping utility\n")
	fs.Parse(os.Args[1:])
	if fs.NArg() > 0 {
		subnet = fs.Arg(0)
	}
	prefix, valid := isValidSubnet(subnet)
	if !valid {
		printUsage(subnet, InvalidSubnet)
		os.Exit(2)
	}
	results := scanSubnet(prefix, count, wait)
	neutralString := "\033[00m"
	for _, result := range results {
		colorString := "\033[01;31m"
		if result.responded {
			colorString = "\033[01;32m"
		}
		fmt.Printf("%s%s%s\n", colorString, result.ip.String(), neutralString)
	}
}
