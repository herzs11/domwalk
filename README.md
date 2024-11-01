# domwalk

Domwalk 🖇 is a tool to find and store domain relationships

It exists currently as a Cloud function that can be interacted with by a CLI tool.

## Cloud Function


### Data Flow
:::mermaid
stateDiagram-v2
    Client --> CloudFunction: Request enrichment of given domains
    CloudFunction --> DomainEnrichment
    state DomainEnrichment {
        [*] --> DNSData: Get MX records, SOA, NS, A, and AAAA records
        [*] --> CertData: Get certificate SANs
        [*] --> WebRedirect: Get web redirects
        [*] --> Sitemap: Get sitemap web domains and contact emails scraped from contact pages
    }
    DomainEnrichment --> BQTable: Upsert domains into domwalk.domains
    DomainEnrichment --> Client: Return enriched domains as JSON
:::

### Installation

```
git clone https://dev.azure.com/Unum/Mkt_Analytics/_git/domwalk
cd domwalk
make
```

## CLI Tool

CLI tool to find and store domain relationships

### Synopsis

domwalk is a CLI tool to find and store domain relationships.
It is written in Go and acts as a client for a domain enrichment Cloud Function. The cloud function url is defined in the ENRICH_DOMAIN_CF_URL environment variable.

Currently, the tool can enrich domains with the following relationships:
- Certificate Subject Alternative Names (SANs)
- Web Redirects
- SitemapLoc Web Domains
- SitemapLoc Contact Page Domains

The tool can also enrich domains with DNS data. In a future version, this dns data will be used to form additional domain relationships


### Examples

```
domwalk domains -d unum.com,coloniallife.com --workers 20 --cert-sans --web-redirects --sitemaps --dns
```

### Options

```
      --cert-sans              Enrich domains with cert SANs
      --dns                    Enrich domains with dns data
  -h, --help                   help for domwalk
      --min-freshness string   Minimum date to refresh relationships, (YYYY-MM-DD) (default "0001-01-01")
  -q, --no-return              Do not return results
  -m, --only-matched           Only return matched domains
  -o, --output string          Output JSON file for results, cannot be used with --no-return
      --sitemaps               Enrich domains with sitemap web domains
      --web-redirects          Enrich domains with web redirects
  -w, --workers int            Number of concurrent workers to use (default 15)
```

### SEE ALSO

* [domwalk completion](docs/domwalk_completion.md)	 - Generate the autocompletion script for the specified shell
* [domwalk domains](docs/domwalk_domains.md)	 - Enrich domains from a list of domain names
* [domwalk file](docs/domwalk_file.md)	 - Enrich domains from file



