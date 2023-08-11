<h1 align="center">
    getsan
  <br>
</h1>

<h4 align="center">A utility to fetch and display dns names from the SSL/TLS cert data</h4>


<p align="center">
  <a href="#installation">🏗️ Installation</a>  
  <a href="#usage">⛏️ Usage</a> 
  <br>
</p>


![TLSDomains](https://example.com/tlsdomains.png)

# Installation
```sh
git clone https://github.com/yourusername/tlsdomains
cd tlsdomains
go build
```

# Usage

- Fetches and displays dns names from the SSL/TLS cert data
- Uses concurrency for efficient and fast lookups

```sh
⚓ echo "cdn.syndication.twitter.com" | getsan | jq

{
  "domain": "cdn.syndication.twitter.com",
  "common_name": "syndication.twitter.com",
  "org": [
    "Twitter, Inc."
  ],
  "dns_names": [
    "syndication.twitter.com",
    "syndication.twimg.com",
    "cdn.syndication.twitter.com",
    "cdn.syndication.twimg.com",
    "syndication-o.twitter.com",
    "syndication-o.twimg.com"
  ]
}
```

