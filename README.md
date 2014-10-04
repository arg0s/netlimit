#Netlimit

CLI tool to fetch your current bandwidth quota. Currently supports Airtel and ACT Fibrenet. 

## Goroutines

This was mostly an experiment to play around with goroutines. 

Processing of each provider mostly consists of fetching data from a provider URL and then matching it against a regex to pull out the quota details. Each of these goroutines runs concurrently, and writes back via a bidirectional channel. 

<pre>
➜  netlimit  go build netlimit.go
➜  netlimit  ./netlimit
2014/10/04 08:31:51 Pulling details from various providers, this may take a minute.
2014/10/04 08:31:52 ACT Fibrenet: 6.44 GB of 75.00&nbsp;GB
</pre>