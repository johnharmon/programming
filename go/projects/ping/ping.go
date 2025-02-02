package main

import (
	"fmt"
	ni "net/netip"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type IPResult struct {
	ip        ni.Addr
	responded bool
}

func printUsage(subnet string) {
	usage := "Usage:\n"
	usage += "     ping <subnet>\n"
	usage += "     <subnet> is a valid IPv4 CIDR notation subnet\n"
	usage += "     No other positional arguments or flags are supported\n"
	usage += fmt.Sprintf("The subnet %s is not valid\n", subnet)
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

func pingIP(ipResult *IPResult, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ipResult.ip.String())
	output, err := cmd.Output()
	if strings.Contains(string(output), "1 packets transmitted, 1 packets received") && err != nil {
		ipResult.responded = true
	}
}

func scanSubnet(prefix ni.Prefix) (results []*IPResult) {
	ips := getSubnetIPs(prefix)
	wg := sync.WaitGroup{}
	for _, ip := range ips {
		wg.Add(1)
		go pingIP(ip, &wg)
	}
	wg.Wait()
	return ips
}

func main() {
	var subnet string
	for idx, arg := range os.Args {
		fmt.Printf("Arg %d has value %s\n", idx, arg)
	}
	if len(os.Args) > 0 && len(os.Args) < 3 {
		subnet = os.Args[1]
	} else {
		printUsage(subnet)
		os.Exit(1)
	}
	prefix, valid := isValidSubnet(subnet)
	if !valid {
		fmt.Printf("")
		printUsage(subnet)
		os.Exit(2)
	}
	results := scanSubnet(prefix)
	//colorString := "\033[01;31m"
	neutralString := "\033[00m"
	for _, result := range results {
		//fmt.Printf("%v\n", result)
		colorString := "\033[01;31m"
		if result.responded {
			colorString = "\033[01;32m"
		}
		fmt.Printf("%s%s%s\n", colorString, result.ip.String(), neutralString)
	}
}
