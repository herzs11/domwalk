<h2 id="domwalk">domwalk</h2>
<p>CLI tool to find and store domain relationships</p>
<h3 id="synopsis">Synopsis</h3>
<p>domwalk is a CLI tool to find and store domain relationships. It is
written in Go and uses GORM and a local SQLite backend. Data is stored
in a local SQLite database and can be pushed to BigQuery.</p>
<p>Currently, the tool can enrich domains with the following relationships:</p>
<pre><code>- Certificate Subject Alternative Names (SANs)
- Web Redirects
- Sitemap Web Domains
- Sitemap Contact Page Domains
</code></pre>
<p>The tool can also enrich domains with DNS data. In a future version, this dns data will be used to form additional domain relationships</p>
<pre><code>domwalk [flags]</code></pre>
<h3 id="examples">Examples</h3>
<pre><code>domwalk -d unum.com,coloniallife.com --workers 20 --cert-sans --web-redirects --sitemaps-web --sitemaps-contact --dns --gorm-db domwalk.db</code></pre>
<h3 id="options">Options</h3>
<pre><code>      --cert-sans              Enrich domains with cert SANs
      --dns                    Enrich domains with dns data
  -d, --domains strings        Domains to process
      --file string            File with domains to process
  -f, --found-domains-only     Only process domains from database
      --gorm-db string         GORM SQLite database name, can also set &#39;GORM_SQLITE_NAME&#39; environment variable
      --header                 File with domains to process has a header row (default true)
  -h, --help                   help for domwalk
  -l, --limit int              Limit of domains to process (default 3000)
      --min-freshness string   Minimum date to refresh relationships, (YYYY-MM-DD) (default &quot;0001-01-01&quot;)
  -s, --offset int             Offset of domains to process
      --sitemaps-contact       Enrich domains with sitemap contact page scraped domains
      --sitemaps-web           Enrich domains with sitemap web domains
      --web-redirects          Enrich domains with web redirects
  -w, --workers int            Number of concurrent workers to use (default 15)</code></pre>
<h1 id="domwalk-completion">domwalk completion</h1>
<p>Generate the autocompletion script for the specified shell</p>
<h3 id="synopsis-1">Synopsis</h3>
<p>Generate the autocompletion script for domwalk for the specified
shell. See each sub-command’s help for details on how to use the
generated script.</p>
<h2 id="domwalk-completion-bash">domwalk completion bash</h2>
<p>Generate the autocompletion script for bash</p>
<h3 id="synopsis-2">Synopsis</h3>
<p>Generate the autocompletion script for the bash shell.</p>
<p>This script depends on the ‘bash-completion’ package. If it is not
installed already, you can install it via your OS’s package manager.</p>
<p>To load completions in your current shell session:</p>
<pre><code>source &lt;(domwalk completion bash)</code></pre>
<p>To load completions for every new session, execute once:</p>
<h4 id="linux">Linux:</h4>
<pre><code>domwalk completion bash &gt; /etc/bash_completion.d/domwalk</code></pre>
<h4 id="macos">macOS:</h4>
<pre><code>domwalk completion bash &gt; $(brew --prefix)/etc/bash_completion.d/domwalk</code></pre>
<p>You will need to start a new shell for this setup to take effect.</p>
<pre><code>domwalk completion bash</code></pre>
<h3 id="options-2">Options</h3>
<pre><code>  -h, --help              help for bash
      --no-descriptions   disable completion descriptions</code></pre>
<h2 id="domwalk-completion-fish">domwalk completion fish</h2>
<p>Generate the autocompletion script for fish</p>
<h3 id="synopsis-3">Synopsis</h3>
<p>Generate the autocompletion script for the fish shell.</p>
<p>To load completions in your current shell session:</p>
<pre><code>domwalk completion fish | source</code></pre>
<p>To load completions for every new session, execute once:</p>
<pre><code>domwalk completion fish &gt; ~/.config/fish/completions/domwalk.fish</code></pre>
<p>You will need to start a new shell for this setup to take effect.</p>
<pre><code>domwalk completion fish [flags]</code></pre>
<h3 id="options-3">Options</h3>
<pre><code>  -h, --help              help for fish
      --no-descriptions   disable completion descriptions</code></pre>
<h2 id="domwalk-completion-powershell">domwalk completion
powershell</h2>
<p>Generate the autocompletion script for powershell</p>
<h3 id="synopsis-4">Synopsis</h3>
<p>Generate the autocompletion script for powershell.</p>
<p>To load completions in your current shell session:</p>
<pre><code>domwalk completion powershell | Out-String | Invoke-Expression</code></pre>
<p>To load completions for every new session, add the output of the
above command to your powershell profile.</p>
<pre><code>domwalk completion powershell [flags]</code></pre>
<h3 id="options-4">Options</h3>
<pre><code>  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions</code></pre>
<h2 id="domwalk-completion-zsh">domwalk completion zsh</h2>
<p>Generate the autocompletion script for zsh</p>
<h3 id="synopsis-5">Synopsis</h3>
<p>Generate the autocompletion script for the zsh shell.</p>
<p>If shell completion is not already enabled in your environment you
will need to enable it. You can execute the following once:</p>
<pre><code>echo &quot;autoload -U compinit; compinit&quot; &gt;&gt; ~/.zshrc</code></pre>
<p>To load completions in your current shell session:</p>
<pre><code>source &lt;(domwalk completion zsh)</code></pre>
<p>To load completions for every new session, execute once:</p>
<h4 id="linux-1">Linux:</h4>
<pre><code>domwalk completion zsh &gt; &quot;${fpath[1]}/_domwalk&quot;</code></pre>
<h4 id="macos-1">macOS:</h4>
<pre><code>domwalk completion zsh &gt; $(brew --prefix)/share/zsh/site-functions/_domwalk</code></pre>
<p>You will need to start a new shell for this setup to take effect.</p>
<pre><code>domwalk completion zsh [flags]</code></pre>
<h3 id="options-5">Options</h3>
<pre><code>  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions</code></pre>
<h2 id="domwalk-db">domwalk db</h2>
<p>Pushing and pulling data to and from BigQuery</p>
<h3 id="synopsis-6">Synopsis</h3>
<p>This command is used to push and pull data to and from BigQuery.</p>
<pre><code>The --push flag snapshots the current BQ dataset into the domwalk_snapshots dataset, then overwrites the current BQ dataset with the data from the local SQLite database</code></pre>
<pre><code>domwalk db [flags]</code></pre>
<h3 id="examples-1">Examples</h3>
<pre><code>domwalk db --push --gorm-db domwalk.db --bq-dataset domwalk</code></pre>
<h3 id="options-6">Options</h3>
<pre><code>      --bq-dataset string   BQ dataset to sync to, can also set &#39;GORM_BQ_DATASET&#39; environment variable
      --cert-sans           Sync cert sans domains
      --dns                 Sync DNS data
      --domains             Sync domains
  -h, --help                help for db
      --pull                Pull data from BigQuery
      --push                Push data to BigQuery
      --sitemaps            Sync sitemaps
      --snapshot            Snapshot domains (default true)
      --web-redirects       Sync web redirect domains
      --gorm-db string      GORM SQLite database name, can also set &#39;GORM_SQLITE_NAME&#39; environment variable
</code></pre>
