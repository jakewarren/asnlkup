package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/getlantern/goexpr/isp/ip2location"
	"github.com/muesli/go-app-paths"
	"github.com/ogier/pflag"
)

// Record holds enriched info for an IP
type Record struct {
	IP   string
	ASN  string
	Name string
}

func main() {

	userScope := apppaths.NewScope(apppaths.User, "", "asnlkup")

	cacheDir,_ := userScope.CacheDir()
	cacheDir += "/IP2LOCATION-LITE-ASN.CSV"

	filePath := pflag.StringP("output", "o", "", "output file name")
	dbPath := pflag.StringP("db", "d", cacheDir, "db file name")
	csvOutput := pflag.BoolP("csv", "c", false, "output in CSV format")
	jsonOutput := pflag.BoolP("json", "j", false, "output in JSON format")
	displayHelp := pflag.BoolP("help", "h", false, "display help")

	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	//human-friendly CLI output
	log.SetHandler(cli.New(os.Stderr))

	//set the logging level
	log.SetLevel(log.DebugLevel)

	r := openStdinOrFile()

	scanner := bufio.NewScanner(r)

	ipList := make([]net.IP, 0)
	for scanner.Scan() {
		ip := net.ParseIP(scanner.Text())
		if IsRoutable(ip) {
			ipList = append(ipList, ip)
		} else {
			log.Warnf("non-routable IP: %s", ip)
		}

	}


	resp := make([]Record, 0)

	prov, err := ip2location.NewProvider(*dbPath)
	checkError("unable to create IP2Location provider", err)

	for _, ip := range ipList {
		asn, _ := prov.ASN(ip.String())
		isp, _ := prov.ISP(ip.String())

		r := Record{
			IP:   ip.String(),
			ASN:  fmt.Sprintf("%d", asn),
			Name: isp,
		}
		resp = append(resp, r)
	}

	// if an output file is not provided, write to STDOUT
	var f *os.File
	if *filePath == "" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.Create(*filePath)
		checkError("Cannot create file", err)
		defer f.Close()
	}

	if *csvOutput {
		writeCSV(resp, f)
	} else if *jsonOutput {
		writeJSON(resp, f)
	} else {
		writeHuman(resp, f)
	}

}

// writeJSON outputs the responses as JSON
func writeJSON(records []Record, f *os.File) {

	rec, _ := json.MarshalIndent(records, "", "    ")

	fmt.Fprint(f, string(rec))

}

// writeCSV outputs the response as CSV
func writeCSV(records []Record, f *os.File) {
	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"IP", "ASN", "ISP"})

	for _, i := range records {
		w.Write([]string{i.IP, i.ASN, i.Name})
	}
}

// writeHuman outputs the response as a pretty tabular output
func writeHuman(records []Record, f *os.File) {
	var w *tabwriter.Writer
	w = tabwriter.NewWriter(f, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	fmt.Fprintf(w, "IP\tASN\tISP\n")

	for _, i := range records {
		fmt.Fprintf(w, "%s\t%s\t%s\n", i.IP, i.ASN, i.Name)
	}

}

// print custom usage instead of the default provided by pflag
func displayUsage() {

	fmt.Printf("Usage: bulkiplkup [<flags>] [FILE]\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}


// openStdinOrFile reads from stdin or a file based on what input the user provides
func openStdinOrFile() io.Reader {
	var err error
	r := os.Stdin
	if len(pflag.Args()) >= 1 {
		r, err = os.Open(pflag.Arg(0))
		if err != nil {
			panic(err)
		}
	}
	return r
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

//IPRange stores an IP range
type IPRange struct {
	from, to net.IP
}

// BogonRanges is a subset of the more static IPv4 Bogon/Reserved/Private ranges.
// In other words, these ranges are such fucking bogon's that they aren't even
// out in public.
var BogonRanges = []IPRange{
	{from: net.ParseIP("0.0.0.0"), to: net.ParseIP("0.255.255.255")},
	{from: net.ParseIP("10.0.0.0"), to: net.ParseIP("10.255.255.255")},
	{from: net.ParseIP("100.64.0.0"), to: net.ParseIP("10.127.255.255")},
	{from: net.ParseIP("127.0.0.0"), to: net.ParseIP("127.255.255.255")},
	{from: net.ParseIP("169.254.0.0"), to: net.ParseIP("169.254.255.255")},
	{from: net.ParseIP("172.16.0.0"), to: net.ParseIP("172.31.255.255")},
	{from: net.ParseIP("192.0.0.0"), to: net.ParseIP("192.0.0.255")},
	{from: net.ParseIP("192.0.2.0"), to: net.ParseIP("192.0.2.255")},
	{from: net.ParseIP("192.88.99.0"), to: net.ParseIP("192.88.99.255")},
	{from: net.ParseIP("192.168.0.0"), to: net.ParseIP("192.168.255.255")},
	{from: net.ParseIP("198.18.0.0"), to: net.ParseIP("198.19.255.255")},
	{from: net.ParseIP("198.51.100.0"), to: net.ParseIP("198.51.100.255")},
	{from: net.ParseIP("203.0.113.0"), to: net.ParseIP("203.0.113.255")},
	{from: net.ParseIP("224.0.0.0"), to: net.ParseIP("239.255.255.255")},
	{from: net.ParseIP("240.0.0.0"), to: net.ParseIP("255.255.255.255")},
}

// IsRoutable returns true if the IP is a publicly routable address
func IsRoutable(ip net.IP) bool {
	for _, rr := range BogonRanges {
		if rr.Contains(ip) {
			return false
		}
	}
	return true
}

// Contains checks if a given IP is in the IPRange
func (r *IPRange) Contains(ip net.IP) bool {
	return (bytes.Compare(ip, r.from) >= 0 && bytes.Compare(ip, r.to) <= 0)
}
