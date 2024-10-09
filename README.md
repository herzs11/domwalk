## domwalk

CLI tool to find and store domain relationships

### Synopsis

domwalk is a CLI tool to find and store domain relationships.
It is written in Go and uses GORM and a local SQLite backend.
Data is stored in a local SQLite database and can be pushed to BigQuery.

	Currently, the tool can enrich domains with the following relationships:
	- Certificate Subject Alternative Names (SANs)
	- Web Redirects
	- Sitemap Web Domains
	- Sitemap Contact Page Domains

	The tool can also enrich domains with DNS data. In a future version, this dns data will be used to form additional domain relationships


```
domwalk [flags]
```

### Installation

```
git clone https://dev.azure.com/Unum/Mkt_Analytics/_git/domwalk
cd domwalk
make
```

#### Pulling latest data from BigQuery
For now, the data is stored on bigquery, but locally it is executed on a SQLite database to optimize the performance. To pull the latest data from BigQuery, you can run the following command:
(You need to have the `gcloud` CLI installed and the GOOGLE_APPLICATION_CREDENTIALS environment variable set)
```
domwalk db --pull
```

### Examples

```
domwalk -d unum.com,coloniallife.com --workers 20 --cert-sans --web-redirects --sitemaps-web --sitemaps-contact --dns --gorm-db domwalk.db
```

### Options

```
      --cert-sans              Enrich domains with cert SANs
      --dns                    Enrich domains with dns data
  -d, --domains strings        Domains to process
      --file string            File with domains to process
  -f, --found-domains-only     Only process domains from database
      --gorm-db string         GORM SQLite database name, can also set 'DOMWALK_SQLITE_NAME' environment variable
      --header                 File with domains to process has a header row (default true)
  -h, --help                   help for domwalk
      --json                   Output results in JSON format
  -l, --limit int              Limit of domains to process (default 3000)
      --min-freshness string   Minimum date to refresh relationships, (YYYY-MM-DD) (default "0001-01-01")
  -s, --offset int             Offset of domains to process
  -o, --output string          Output JSON file for results, only to be used with --json
      --sitemaps-contact       Enrich domains with sitemap contact page scraped domains
      --sitemaps-web           Enrich domains with sitemap web domains
      --web-redirects          Enrich domains with web redirects
  -w, --workers int            Number of concurrent workers to use (default 15)
```

### SEE ALSO

* [domwalk completion](docs/domwalk_completion.md)	 - Generate the autocompletion script for the specified shell
* [domwalk db](docs/domwalk_db.md)	 - Pushing and pulling data to and from BigQuery

