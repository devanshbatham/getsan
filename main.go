package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

/*
   CertData represents the extracted certificate data for a domain.
*/
type CertData struct {
	Domain     string   `json:"domain"`
	CommonName string   `json:"common_name"`
	Org        []string `json:"org"`
	DNSNames   []string `json:"dns_names"`
}

func main() {
	/*
	   Parse command line flags to determine the number of concurrent tasks.
	*/
	concurrency := flag.Int("c", 50, "Number of concurrent tasks")
	flag.Parse()

	var wg sync.WaitGroup

	// Create channels to communicate between goroutines.
	domains := make(chan string, *concurrency)       // Channel to send domain names.
	results := make(chan *CertData, *concurrency)    // Channel to receive certificate data.

	/*
	   Start worker goroutines.
	*/
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			/*
			   Process domain names from the 'domains' channel.
			*/
			for domain := range domains {
				func() {
					defer recover() // Catch any panics.

					/*
					   Fetch certificate data for the current domain.
					*/
					data, err := getCertificate(domain)
					if err == nil && data != nil {
						results <- data // Send the certificate data to the 'results' channel.
					}
				}()
			}
		}()
	}

	/*
	   Read domain names from standard input and send them to the 'domains' channel.
	*/
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			domain := strings.TrimSpace(scanner.Text())
			domain = strings.TrimPrefix(domain, "http://")
			domain = strings.TrimPrefix(domain, "https://")
			domains <- domain // Send the domain to the 'domains' channel.
		}
		close(domains) // Close the 'domains' channel when all domains are sent.
	}()

	/*
	   Wait for all worker goroutines to finish.
	*/
	go func() {
		wg.Wait()
		close(results) // Close the 'results' channel when all workers are done.
	}()

	first := true

	/*
	   Process certificate data received from the 'results' channel.
	*/
	for result := range results {
		if !first {
			fmt.Println() // Print a newline to separate JSON objects.
		} else {
			first = false
		}
		output, err := json.Marshal(result)
		if err == nil {
			fmt.Print(string(output))
		}
	}
}

/*
   getCertificate fetches SSL/TLS certificate data for the given domain.
*/
func getCertificate(domain string) (*CertData, error) {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	cert := conn.ConnectionState().PeerCertificates[0]
	certData := &CertData{
		Domain:     domain,
		CommonName: cert.Subject.CommonName,
		Org:        cert.Subject.Organization,
		DNSNames:   cert.DNSNames,
	}
	return certData, nil
}
