# Usage Notes

## Pagination

Many `list` and `find` commands support:

- `--limit` / `--offset`: fetch a single page.
- `--all`: auto-paginates and returns all pages (defaults to `--limit 1000` unless you set `--limit` explicitly).

Example:

```bash
# Fetch all campaigns (auto-pagination)
aads campaigns list --all

# Fetch all ad groups in a campaign
aads adgroups list --campaign-id 12345 --all
```

## Currency For Money Fields

When the CLI builds request bodies from flags (campaign budgets, ad group bids, keyword bids), it needs a currency code.

Resolution order:

1. `--currency`
2. `default_currency` in `~/.aads/config.yaml` (or env vars `AADS_DEFAULT_CURRENCY` / `AADS_CURRENCY`)
3. inferred from `GET /acls` (matched by `org_id`)

If currency can't be resolved unambiguously, the command errors (instead of silently sending `USD`).

## Retries And Timeouts

- All API requests use an HTTP client timeout.
- Retries:
  - Always retries `429 Too Many Requests` (respects `Retry-After` when present).
  - Retries `5xx` only for requests that are safe to retry (`GET`, `PUT`, `DELETE`, selector `POST .../find`, `POST /reports/...`, and bulk delete `POST .../delete/bulk`).
  - On `401 Unauthorized`, the client forces a token refresh and retries once.

## Documentation

- Full Apple Ads docs index (link-only): `docs/OFFICIAL_APPLE_ADS_DOC_INDEX.md`
- CLI-to-endpoint mapping: `docs/OFFICIAL_APPLE_ADS_DOCS.md`
- Command reference (generated): `docs/commands/aads.md`

