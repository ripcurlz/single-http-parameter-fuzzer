package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var headers string

// this function generates an array of ints in the given range between (including) min and max. E.g. [10000,10001,10002,.....,99999]
func generateArrayOfIntsInRange(min int, max int) []string {
	intArray := make([]int, max-min+1)
	for i := range intArray {
		intArray[i] = min + i
	}
	//now convert each element from int to string in a newly made string array with the same size...
	stringArray := make([]string, max-min+1)
	for i := range intArray {
		stringArray[i] = strconv.Itoa(intArray[i])
	}
	return stringArray
}

// reads a file and returns a string array
func readFileToStringArray(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func doPostRequest(host string, port int, parametertofuzz string, payload string, headers string) string {

	// set a timeout for the http requests to 5 seconds
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	body := []byte(fmt.Sprintf("%v=%v", parametertofuzz, payload))

	url := "http://" + host + ":" + strconv.Itoa(port)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(body),
	)
	//read in all given headers
	var headerarray []string

	if headers != "" {

		// get tail of post command for all headers
		tail := flag.Args()
		headerarray = append(tail, headers)

		// add tail and everything else to postheaderarray to add all headers to request
		for _, header := range headerarray {

			splitted := strings.Split(header, ":")
			headername := splitted[0]
			headervalue := splitted[1]
			req.Header.Add(headername, headervalue)
		}
	}

	if err != nil {
		log.Fatalf("[!] Unable to generate request: %s\n", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("[!] Unable to process response: %s\n", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[!] Unable to read response body: %s\n", err)
	}

	// convert result response body to string
	bodyString := string(body)

	resp.Body.Close()

	return bodyString

}

func doGetRequest(host string, port int, parametertofuzz string, payload string, headers string) string {

	// set a timeout for the http requests to 5 seconds
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := "http://" + host + ":" + strconv.Itoa(port) + "?" + parametertofuzz + "=" + payload
	req, err := http.NewRequest("GET", url, nil)
	//read in all given headers
	var headerarray []string

	if headers != "" {

		// get tail of post command for all headers
		tail := flag.Args()
		headerarray = append(tail, headers)

		// add tail and everything else to postheaderarray to add all headers to request
		for _, header := range headerarray {

			splitted := strings.Split(header, ":")
			headername := splitted[0]
			headervalue := splitted[1]
			req.Header.Add(headername, headervalue)
		}

	}

	if err != nil {
		log.Fatalf("[!] Unable to generate request: %s\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("[!] Unable to process response: %s\n", err)
	}

	var body []byte

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[!] Unable to read response body: %s\n", err)
	}

	// convert result response body to string
	bodyString := string(body)

	resp.Body.Close()

	return bodyString

}

func printFound(bodyString string, payload string) {

	fmt.Printf("We found the body we are looking for with payload %v\n\n", payload)
	fmt.Printf("The body is:\n")
	fmt.Printf("%v\n\n", bodyString)

}

func checkBodyForStrings(bodyString string, payload string, stringtobeinresponse string, stringnottobeinresponse string) {

	//  if stringtobeinresponse is "", it is assumed it is not used:
	if stringtobeinresponse == "" && stringnottobeinresponse != "" {
		if !strings.Contains(bodyString, stringnottobeinresponse) {
			printFound(bodyString, payload)
		}

	}
	//  if stringnottobeinresponse is "", it is assumed it is not used:
	if stringtobeinresponse != "" && stringnottobeinresponse == "" {
		if strings.Contains(bodyString, stringtobeinresponse) {
			printFound(bodyString, payload)
		}
	}
	// both are not ""
	if stringtobeinresponse != "" && stringnottobeinresponse != "" {
		if strings.Contains(bodyString, stringtobeinresponse) && !strings.Contains(bodyString, stringnottobeinresponse) {
			printFound(bodyString, payload)

		}
	}

}

func main() {

	// example usage:
	// go run main.go -method post -host 10.10.63.138 -port 8085 -stringnottobeinresponse "Oh no" -parametertofuzz number -startnumber 10900 -endnumber 99999 -headers "Content-Type:application/x-www-form-urlencoded" "X-Originating-IP:127.0.0.1" "X-Forwarded-For:127.0.0.1" "X-Remote-IP:127.0.0.1" "X-Remote-Addr:127.0.0.1" "X-Client-IP:127.0.0.1" "X-Host:127.0.0.1" "X-Forwarded-Host:127.0.0.1"
	hostflag := flag.String("host", "", "the IP address or hostname you want to fuzz on")
	portflag := flag.Int("port", 80, "the port you want to use")
	stringtobeinresponseflag := flag.String("stringtobeinresponse", "", "the string you do want to be in the http response body for a match")
	stringnottobeinresponseflag := flag.String("stringnottobeinresponse", "", "the string you do not want to be in the http response body for a match")
	parametertofuzzflag := flag.String("parametertofuzz", "", "the single http parameter you want to fuzz")
	wordlistflag := flag.String("wordlist", "", "path of the wordlist you want to use for fuzzing the parameter, if you do not use startnumber and endnumber")
	startnumberflag := flag.Int("startnumber", 0, "the start number to use for fuzzing the parameter")
	endnumberflag := flag.Int("endnumber", 100, "the end number to use for fuzzing the parameter")
	methodflag := flag.String("method", "", "the http method you want to use, either 'get' or 'post'")
	flag.StringVar(&headers, "headers (optional)", "", "http headers you want to use, in the form of 'header1:value1' 'header2:value2' and so on")

	flag.Parse()
	fmt.Printf("Now starting...\n")

	// if one of the required flags is missing, print usage and exit with code 1
	if *methodflag == "" || *hostflag == "" || *portflag == -1 || (*stringtobeinresponseflag == "" && *stringnottobeinresponseflag == "") || *parametertofuzzflag == "" || (*wordlistflag == "" && (*startnumberflag == -1 || *endnumberflag == 0)) {
		flag.Usage()
		os.Exit(1)
	}

	host := *hostflag
	port := *portflag
	stringtobeinresponse := *stringtobeinresponseflag
	stringnottobeinresponse := *stringnottobeinresponseflag
	parametertofuzz := *parametertofuzzflag
	wordlist := *wordlistflag
	startnumber := *startnumberflag
	endnumber := *endnumberflag
	method := *methodflag

	// now check, if a wordlist shall be used or the number upcounted...

	// case one: wordlist!
	// now read the wordlist from the given file name
	var payloads []string

	if wordlist != "" {

		fmt.Printf("Using wordlist '%v'...\n", wordlist)

		lines, err := readFileToStringArray(wordlist)
		if err != nil {
			log.Fatalf("readFileToStringArray from wordlist: %s", err)
		}

		// instantiate empty string array with the length of the lines of the read file
		stringArray := make([]string, len(lines))

		for i, line := range lines {
			stringArray[i] = line

		}
		payloads = stringArray
	}

	// case two: upcounting a number instead of wordlist!
	if endnumber != 0 && wordlist == "" {
		fmt.Printf("Using startnumber %v and endnumber %v. Will start counting up now...\n", startnumber, endnumber)
		// create payload array from number range
		payloads = generateArrayOfIntsInRange(startnumber, endnumber)
	}

	// counter to be used to print out info at certain points
	infocounter := 1

	// now, do a http request for every number in that range!
	for _, payload := range payloads {

		// log info every 100th loop :)
		if infocounter%100 == 0 {
			fmt.Printf("Currently at payload %v=%v\n\n", parametertofuzz, payload)
		}
		infocounter = infocounter + 1

		var bodyString string

		// if we have a post http method do the following
		if method == "post" {

			bodyString = doPostRequest(host, port, parametertofuzz, payload, headers)

		}

		// if we have a get http method, do the following
		if method == "get" {

			bodyString = doGetRequest(host, port, parametertofuzz, payload, headers)

		}

		checkBodyForStrings(bodyString, payload, stringtobeinresponse, stringnottobeinresponse)

	}

}
