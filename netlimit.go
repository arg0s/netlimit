package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type debugging bool

type Provider struct {
  name string
	URL  string
	reg1 *regexp.Regexp
	reg2 *regexp.Regexp
  max string
  used string
}

func first(x *regexp.Regexp, y error) *regexp.Regexp {
	return x
}


var (
	debug   debugging
	verbose debugging
	info    debugging

	config = []Provider{
		{
      "ACT Fibrenet", "http://portal.acttv.in/index.php/mypackage",
      first(regexp.Compile(`(\d+.\d+)\sGB`)),
      first(regexp.Compile(`(\d+.\d+)&nbsp;GB`)),
      "",
      "",
      },
		{
      "Airtel",
      "http://122.160.230.125:8080/gbod/gb_on_demand.do",
      first(regexp.Compile(`Balance.+(\d+.\d+)&nbsp;GB`)),
      first(regexp.Compile(`High.+(\d+.\d+)&nbsp;GB`)),
      "",
      "",
      },
	}
)

func neterr(err error) {
  debug.Printf("ERROR - %s", err)
  if err != nil {
    info.Printf("%s", err)
    panic("Oops, we are unable to fetch the current network stats")
  }
}

func (d debugging) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}

func netstats(c chan Provider) {
  // Grab the provider, URL, hit it and fetch the body

  p:= <- c
  debug.Printf("%s", p)
	resp, err := http.Get(p.URL)
	neterr(err)

	defer resp.Body.Close()

  // Load up the body and panic if we error out
	body, err := ioutil.ReadAll(resp.Body)
	neterr(err)

  debug.Printf("...parsed %s message body of length %d", p.name, len(body))

  // Run the regexps to pull submatches against the patterns
	current := p.reg1.FindSubmatch(body)
	max := p.reg2.FindSubmatch(body)

  debug.Printf("%s %s",current, max)

  // If we've hit gold, send word back via the channel
	if current != nil && max != nil{
    p.max = string(max[0])
    p.used = string(current[0])
    c <- p
	} else {
    debug.Printf("%s is not currently your provider", p.name)
		return
	}
}

func main() {
	info = true
  debug = false
  info.Printf("Pulling details from various providers, this may take a minute.")
  c := make(chan Provider)
  defer func() {
      if r := recover(); r != nil {
          debug.Printf("Recovered from panic >> ", r)
      }
  }()
	for _, element := range config {
    // Fire off a goroutine for each provider
		go netstats(c)
    c <- element
	}
  // We expect only one to call home via the channel with the results
  p := <- c
	info.Printf("%s: %s of %s", p.name, p.used, p.max)
}
