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
	//usage += fmt.Sprintf("The subnet %s is not valid\n", subnet)
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
	var subnetSet = false
	argLen := len(os.Args)
	// Check if the subnet was given as the first parameter
	// This is due to a limitation with the flag module that only accepts positional arguments after the flags
	// Will pull out the second index from args if it is not a flag, merge either side of the os.Args slice and then pass it to flags.Parse()
	if argLen >= 2 && os.Args[1][0] != '-' {
		subnetSet = true
		subnet = os.Args[1]
		if argLen >= 3 {
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		}
	}
	flag.StringVar(&count, "c", "1", "Number of packets to send, equivalent to the same flag for the ping utility\n")
	flag.StringVar(&wait, "w", "1", "How many seconds to wait for a response, equivalent to the -W flag for the ping utility\n")
	flag.Parse()
	if !subnetSet {
		args := flag.Args()
		if len(args) == 1 {
			subnet = args[0]
		} else {
			if len(args) < 2 {
				printUsage(subnet, MissingArgument)
				os.Exit(1)
			}
			printUsage(subnet, TooManyArguments)
			os.Exit(1)
		}
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
