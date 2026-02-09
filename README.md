# aads - Apple Ads CLI

A command-line interface for the [Apple Ads Campaign Management API 5](https://developer.apple.com/documentation/apple_ads). Manage campaigns, ad groups, keywords, creatives, reports, and more from your terminal.

**~70 API endpoints** covered across 17 command groups. Pure Go (minimal deps), safe retries, auto-pagination (`--all`), and multiple output formats.

Documentation:
- Overview + usage notes: [`docs/README.md`](docs/README.md)
- Command reference (generated): [`docs/commands/aads.md`](docs/commands/aads.md)

Official docs:
- Full index (link-only): [`docs/OFFICIAL_APPLE_ADS_DOC_INDEX.md`](docs/OFFICIAL_APPLE_ADS_DOC_INDEX.md)
- CLI mapping: [`docs/OFFICIAL_APPLE_ADS_DOCS.md`](docs/OFFICIAL_APPLE_ADS_DOCS.md)

## Installation

### From source

```bash
git clone https://github.com/SaadBelfqih/apple-ads-cli.git
cd apple-ads-cli
make install
```

This installs `aads` to `$(go env GOBIN)` (or `$(go env GOPATH)/bin` if `GOBIN` is unset).

Or build locally:

```bash
make build
./aads --help
```

Requires Go 1.25+.

## Configuration

### Interactive setup

```bash
aads configure
```

This prompts for your Apple Ads API credentials and saves them to `~/.aads/config.yaml`:

```yaml
client_id: "SEARCHADS.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
team_id: "SEARCHADS.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
key_id: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
org_id: "1234567"
private_key_path: "/path/to/private_key.pem"
default_currency: "USD" # optional (if omitted, inferred from `aads acls list`)
```

### Environment variables

All config values can be overridden with environment variables:

| Variable | Description |
|---|---|
| `AADS_CLIENT_ID` | OAuth2 client ID |
| `AADS_TEAM_ID` | Team ID |
| `AADS_KEY_ID` | Private key ID |
| `AADS_ORG_ID` | Organization ID |
| `AADS_PRIVATE_KEY_PATH` | Path to EC P-256 private key PEM |
| `AADS_DEFAULT_CURRENCY` | Default currency for Money fields built from flags (e.g., USD) |
| `AADS_CURRENCY` | Alias for `AADS_DEFAULT_CURRENCY` |

### Getting credentials

1. Go to [Apple Search Ads](https://searchads.apple.com) > Settings > API
2. Create an API user with appropriate role
3. Generate an EC P-256 private key and upload the public key
4. Note your Client ID, Team ID, Key ID, and Org ID

## Global Flags

```
-o, --output string   Output format: json, table, yaml (default "json")
-v, --verbose         Verbose HTTP request/response logging
    --org-id string   Override org ID from config
    --fields string   Comma-separated fields for partial fetch
    --currency string Override currency for money fields (e.g., USD)
```

Many `list` and `find` commands also support `--all` to auto-paginate through all results.

## Commands

### Campaigns

```bash
# List all campaigns
aads campaigns list

# List with partial fields
aads campaigns list --fields id,name,status

# Get a specific campaign
aads campaigns get --id 12345

# Create a campaign
aads campaigns create \
  --name "My Campaign" \
  --adam-id 123456789 \
  --countries "US,GB,CA" \
  --budget "1000" \
  --daily-budget "50" \
  --status ENABLED

# Create from JSON
aads campaigns create --from-json @campaign.json

# Find campaigns with simple filter
aads campaigns find --field name --op STARTSWITH --values "My"

# Find with full selector JSON
aads campaigns find --selector-json '{"conditions":[{"field":"status","operator":"EQUALS","values":["ENABLED"]}]}'

# Update a campaign
aads campaigns update --id 12345 --daily-budget "75" --status PAUSED

# Delete a campaign
aads campaigns delete --id 12345
```

### Ad Groups

```bash
# List ad groups in a campaign
aads adgroups list --campaign-id 12345

# Get a specific ad group
aads adgroups get --campaign-id 12345 --id 67890

# Create an ad group
aads adgroups create \
  --campaign-id 12345 \
  --name "Exact Match Keywords" \
  --default-bid "1.50" \
  --search-match=false

# Create with targeting dimensions
aads adgroups create --campaign-id 12345 --from-json @adgroup.json

# Find ad groups across all campaigns (org-level)
aads adgroups find-all --selector-json '{"conditions":[{"field":"status","operator":"EQUALS","values":["ENABLED"]}]}'

# Update an ad group
aads adgroups update --campaign-id 12345 --id 67890 --default-bid "2.00"

# Delete an ad group
aads adgroups delete --campaign-id 12345 --id 67890
```

### Keywords (Targeting)

```bash
# List keywords in an ad group
aads keywords list --campaign-id 12345 --adgroup-id 67890

# Create a keyword
aads keywords create \
  --campaign-id 12345 \
  --adgroup-id 67890 \
  --text "grammar checker" \
  --match-type EXACT \
  --bid "1.25"

# Bulk create from JSON
aads keywords create --campaign-id 12345 --adgroup-id 67890 \
  --from-json '[{"text":"grammar","matchType":"BROAD"},{"text":"spell check","matchType":"EXACT","bidAmount":{"amount":"2.00","currency":"USD"}}]'

# From file
aads keywords create --campaign-id 12345 --adgroup-id 67890 --from-json @keywords.json

# From stdin
cat keywords.json | aads keywords create --campaign-id 12345 --adgroup-id 67890 --from-json @-

# Find keywords in an ad group
aads keywords find --campaign-id 12345 --adgroup-id 67890 --selector-json '{"conditions":[{"field":"matchType","operator":"EQUALS","values":["EXACT"]}]}'

# Find keywords across all ad groups in a campaign
aads keywords find-campaign --campaign-id 12345

# Update keywords (bulk)
aads keywords update --campaign-id 12345 --adgroup-id 67890 \
  --from-json '[{"id":111,"bidAmount":{"amount":"3.00","currency":"USD"}}]'

# Delete keywords (bulk)
aads keywords delete --campaign-id 12345 --adgroup-id 67890 --ids 111,222,333

# Delete a single keyword
aads keywords delete-one --campaign-id 12345 --adgroup-id 67890 --id 111
```

### Negative Keywords

Campaign-level and ad group-level negative keywords share the `negatives` command with prefixed subcommands:

```bash
# Campaign-level negatives
aads negatives campaign-create --campaign-id 12345 --text "free" --match-type BROAD
aads negatives campaign-list --campaign-id 12345
aads negatives campaign-find --campaign-id 12345 --selector-json '...'
aads negatives campaign-update --campaign-id 12345 --from-json '[...]'
aads negatives campaign-delete --campaign-id 12345 --ids 111,222

# Ad group-level negatives
aads negatives adgroup-create --campaign-id 12345 --adgroup-id 67890 --text "free" --match-type EXACT
aads negatives adgroup-list --campaign-id 12345 --adgroup-id 67890
aads negatives adgroup-find --campaign-id 12345 --adgroup-id 67890 --selector-json '...'
aads negatives adgroup-update --campaign-id 12345 --adgroup-id 67890 --from-json '[...]'
aads negatives adgroup-delete --campaign-id 12345 --adgroup-id 67890 --ids 111,222

# Bulk create from JSON
aads negatives campaign-create --campaign-id 12345 \
  --from-json '[{"text":"cheap","matchType":"BROAD"},{"text":"free download","matchType":"EXACT"}]'
```

### Ads

```bash
# List ads in an ad group
aads ads list --campaign-id 12345 --adgroup-id 67890

# Create an ad
aads ads create --campaign-id 12345 --adgroup-id 67890 \
  --creative-id 99999 --name "Default Ad" --status ENABLED

# Find ads across all campaigns (org-level)
aads ads find-all --selector-json '...'

# Update an ad
aads ads update --campaign-id 12345 --adgroup-id 67890 --id 11111 --status PAUSED

# Delete an ad
aads ads delete --campaign-id 12345 --adgroup-id 67890 --id 11111
```

### Creatives

Creatives are org-level resources (not nested under campaigns):

```bash
# List all creatives
aads creatives list

# Get a creative
aads creatives get --id 99999

# Create a creative
aads creatives create --adam-id 123456789 --product-page-id "pp-12345" --name "Holiday CPP"

# Find creatives
aads creatives find --selector-json '...'
```

### Reports

```bash
# Campaign-level report
aads reports campaigns --start-time 2025-01-01 --end-time 2025-01-31

# With weekly granularity
aads reports campaigns --start-time 2025-01-01 --end-time 2025-01-31 --granularity WEEKLY

# Grouped by country
aads reports campaigns --start-time 2025-01-01 --end-time 2025-01-31 --group-by countryOrRegion

# Ad group-level report
aads reports adgroups --campaign-id 12345 --start-time 2025-01-01 --end-time 2025-01-31

# Keyword-level report (optionally scoped to ad group)
aads reports keywords --campaign-id 12345 --start-time 2025-01-01 --end-time 2025-01-31
aads reports keywords --campaign-id 12345 --adgroup-id 67890 --start-time 2025-01-01 --end-time 2025-01-31

# Search term report
aads reports searchterms --campaign-id 12345 --start-time 2025-01-01 --end-time 2025-01-31

# Ad-level report
aads reports ads --campaign-id 12345 --start-time 2025-01-01 --end-time 2025-01-31

# With selector for filtering
aads reports campaigns --start-time 2025-01-01 --end-time 2025-01-31 \
  --selector-json '{"conditions":[{"field":"countryOrRegion","operator":"EQUALS","values":["US"]}],"orderBy":[{"field":"localSpend","sortOrder":"DESCENDING"}]}'
```

Granularity options: `HOURLY`, `DAILY` (default), `WEEKLY`, `MONTHLY`

### Impression Share Reports

```bash
# Create a custom impression share report
aads impression-share create --from-json @report_request.json

# Get a report by ID
aads impression-share get --id 12345

# List all reports
aads impression-share list
```

### Apps

```bash
# Search for apps
aads apps search --query "grammar checker"

# Search returning only owned apps
aads apps search --query "grammar" --return-owned-apps

# Check app eligibility for Apple Ads
aads apps eligibility --adam-id 123456789

# Get app details
aads apps details --adam-id 123456789

# Get localized app details
aads apps localized --adam-id 123456789
```

### Product Pages

```bash
# List product pages for an app
aads product-pages list --adam-id 123456789

# Get a specific product page
aads product-pages get --id "45812c9b-c296-43d3-c6a0-c5a02f74bf6e" --adam-id 123456789

# Get locale details
aads product-pages locales --id "45812c9b-c296-43d3-c6a0-c5a02f74bf6e" --adam-id 123456789

# List supported countries/regions
aads product-pages countries

# Get device size mapping
aads product-pages device-sizes
```

### Ad Rejections

```bash
# Find ad creative rejection reasons
aads ad-rejections find --selector-json '...'

# Get a rejection reason by ID
aads ad-rejections get --id 12345

# Find app assets
aads ad-rejections find-assets --adam-id 123456789
```

### Geolocations

```bash
# Search for geolocations
aads geo search --query "New York"
aads geo search --query "California" --entity AdminArea
aads geo search --query "US" --entity Country --country-code US

# Get geo location by ID
aads geo get --geo-id 123456
```

Entity types: `Country`, `AdminArea`, `Locality`

### Budget Orders

```bash
# List budget orders
aads budgetorders list

# Get a budget order
aads budgetorders get --id 12345

# Create a budget order
aads budgetorders create --from-json @budget_order.json

# Update a budget order
aads budgetorders update --id 12345 --from-json @budget_update.json
```

### ACLs

```bash
# List user ACLs (org access)
aads acls list

# Get caller details
aads acls me
```

## Output Formats

### JSON (default)

```bash
aads campaigns list
```

Pretty-printed JSON output, suitable for piping to `jq`.

### Table

```bash
aads campaigns list -o table
```

ASCII table format for interactive use.

### YAML

```bash
aads campaigns list -o yaml
```

## Selector JSON

Find commands accept `--selector-json` for complex queries. You can pass inline JSON or reference a file with `@`:

```bash
# Inline
aads campaigns find --selector-json '{"conditions":[{"field":"status","operator":"EQUALS","values":["ENABLED"]}]}'

# From file
aads campaigns find --selector-json @selector.json
```

### Selector structure

```json
{
  "conditions": [
    {
      "field": "status",
      "operator": "EQUALS",
      "values": ["ENABLED"]
    },
    {
      "field": "name",
      "operator": "STARTSWITH",
      "values": ["US"]
    }
  ],
  "orderBy": [
    {
      "field": "name",
      "sortOrder": "ASCENDING"
    }
  ],
  "pagination": {
    "offset": 0,
    "limit": 50
  }
}
```

Operators: `EQUALS`, `NOT_EQUALS`, `GREATER_THAN`, `LESS_THAN`, `IN`, `LIKE`, `STARTSWITH`, `CONTAINS`, `ENDSWITH`, `IS`, `BETWEEN`

For simple single-condition queries, use the shorthand flags:

```bash
aads campaigns find --field name --op STARTSWITH --values "US"
```

## JSON Input

Commands that create or update resources in bulk accept `--from-json`:

```bash
# Inline JSON
aads keywords create --campaign-id 12345 --adgroup-id 67890 \
  --from-json '[{"text":"grammar","matchType":"BROAD"}]'

# From file
aads keywords create --campaign-id 12345 --adgroup-id 67890 --from-json @keywords.json

# From stdin
cat keywords.json | aads keywords create --campaign-id 12345 --adgroup-id 67890 --from-json @-
```

## Partial Fetch

GET endpoints support the `--fields` flag to limit returned fields:

```bash
aads campaigns list --fields id,name,status,dailyBudgetAmount
aads campaigns get --id 12345 --fields id,name,budgetAmount
```

## Authentication

The CLI uses Apple's OAuth2 flow with ES256 JWT client credentials:

1. Builds a JWT signed with your EC P-256 private key
2. Exchanges the JWT for an access token at `https://appleid.apple.com/auth/oauth2/token`
3. Caches the token in memory (auto-refreshes about 60 seconds before expiry)
4. Sends `Authorization: Bearer <token>` and `X-AP-Context: orgId=<orgId>` on all API calls

Retry policy:
- Retries `429 Too Many Requests` with exponential backoff (2s, 4s, 8s, 16s max), respecting `Retry-After` when present.
- Retries `5xx` only for requests that are safe to retry (`GET`, `PUT`, `DELETE`, selector `POST .../find`, `POST /reports/...`, and bulk delete `POST .../delete/bulk`).
- On `401 Unauthorized`, forces a token refresh and retries once.

## API Coverage

| Resource | Endpoints | Commands |
|---|---|---|
| Campaigns | 6 | `campaigns {create,get,list,find,update,delete}` |
| Ad Groups | 7 | `adgroups {create,get,list,find,find-all,update,delete}` |
| Keywords | 8 | `keywords {create,get,list,find,find-campaign,update,delete,delete-one}` |
| Campaign Negatives | 6 | `negatives campaign-{create,get,list,find,update,delete}` |
| Ad Group Negatives | 6 | `negatives adgroup-{create,get,list,find,update,delete}` |
| Ads | 7 | `ads {create,get,list,find,find-all,update,delete}` |
| Creatives | 4 | `creatives {create,get,list,find}` |
| Product Pages | 5 | `product-pages {list,get,locales,countries,device-sizes}` |
| Ad Rejections | 3 | `ad-rejections {find,get,find-assets}` |
| Reports | 5 | `reports {campaigns,adgroups,keywords,searchterms,ads}` |
| Impression Share | 3 | `impression-share {create,get,list}` |
| Apps | 4 | `apps {search,eligibility,details,localized}` |
| Geolocations | 2 | `geo {search,get}` |
| Budget Orders | 4 | `budgetorders {create,get,list,update}` |
| ACLs | 2 | `acls {list,me}` |
| **Total** | **~72** | |

## Project Structure

```
apple-ads-cli/
├── main.go
├── go.mod
├── Makefile
├── cmd/                    # Cobra commands (flags + service call + print)
│   ├── root.go             # Root cmd, global flags, output helpers
│   ├── helpers.go          # Selector/JSON parsing utilities
│   ├── configure.go        # Interactive credential setup
│   ├── version.go
│   ├── campaigns.go
│   ├── adgroups.go
│   ├── keywords.go
│   ├── negatives.go
│   ├── ads.go
│   ├── creatives.go
│   ├── reports.go
│   ├── impression_share.go
│   ├── apps.go
│   ├── product_pages.go
│   ├── ad_rejections.go
│   ├── geo.go
│   ├── budgetorders.go
│   └── acls.go
├── internal/
│   ├── api/                # HTTP client + service layers
│   │   ├── client.go       # HTTP client with retry + auth
│   │   ├── auth.go         # ES256 JWT + OAuth2 token exchange
│   │   ├── errors.go       # API error types
│   │   └── *.go            # One service file per resource
│   ├── types/              # Pure data structs
│   │   ├── common.go       # Money, Selector, Pagination
│   │   ├── responses.go    # Generic APIResponse[T]
│   │   └── *.go            # One type file per resource
│   ├── config/
│   │   └── config.go       # ~/.aads/config.yaml + env vars
│   └── output/
│       ├── output.go       # Format dispatcher
│       ├── json.go
│       ├── table.go
│       └── yaml.go
```

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/golang-jwt/jwt/v5` | ES256 JWT signing |
| `gopkg.in/yaml.v3` | Config + YAML output |

No external HTTP client, table, or API client libraries.

## License

MIT
