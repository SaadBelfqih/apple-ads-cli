package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

// parseSelector parses a selector from --selector-json flag or inline flags.
func parseSelector(selectorJSON string, field, op string, values []string, limit, offset int) (*types.Selector, error) {
	if selectorJSON != "" {
		return parseSelectorJSON(selectorJSON)
	}

	sel := &types.Selector{}

	if field != "" && op != "" && len(values) > 0 {
		sel.Conditions = []*types.Condition{
			{Field: field, Operator: op, Values: values},
		}
	}

	if limit > 0 || offset > 0 {
		sel.Pagination = &types.Pagination{Limit: limit, Offset: offset}
	}

	return sel, nil
}

// parseSelectorJSON parses JSON from a string or @file.
func parseSelectorJSON(input string) (*types.Selector, error) {
	var data []byte
	var err error

	if strings.HasPrefix(input, "@") {
		path := input[1:]
		if path == "-" {
			data, err = readStdin()
		} else {
			data, err = os.ReadFile(path)
		}
		if err != nil {
			return nil, fmt.Errorf("read selector: %w", err)
		}
	} else {
		data = []byte(input)
	}

	var sel types.Selector
	if err := json.Unmarshal(data, &sel); err != nil {
		return nil, fmt.Errorf("parse selector JSON: %w", err)
	}
	return &sel, nil
}

// parseJSONInput parses JSON from --from-json flag (inline, @file, or - for stdin).
func parseJSONInput(input string, target any) error {
	var data []byte
	var err error

	if strings.HasPrefix(input, "@") {
		path := input[1:]
		if path == "-" {
			data, err = readStdin()
		} else {
			data, err = os.ReadFile(path)
		}
		if err != nil {
			return fmt.Errorf("read JSON: %w", err)
		}
	} else {
		data = []byte(input)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}
	return nil
}

func readStdin() ([]byte, error) {
	return os.ReadFile("/dev/stdin")
}

var (
	resolvedCurrency   string
	currencyResolved   bool
	currencyResolveErr error
)

func normalizeCurrencyCode(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return v
	}
	return strings.ToUpper(v)
}

// resolveMoneyCurrency resolves the currency to use for Money fields when the CLI builds requests from flags.
// Priority: --currency > config default_currency / env vars > inferred from GET /acls (matched by org id).
func resolveMoneyCurrency() (string, error) {
	if currencyResolved {
		return resolvedCurrency, currencyResolveErr
	}
	currencyResolved = true

	if currencyFlag != "" {
		resolvedCurrency = normalizeCurrencyCode(currencyFlag)
		return resolvedCurrency, nil
	}
	if defaultCurrencyFromConfig != "" {
		resolvedCurrency = normalizeCurrencyCode(defaultCurrencyFromConfig)
		return resolvedCurrency, nil
	}
	if apiClient == nil {
		currencyResolveErr = fmt.Errorf("API client not initialized; pass --currency or set default_currency (or env var AADS_DEFAULT_CURRENCY/AADS_CURRENCY)")
		return "", currencyResolveErr
	}

	acls, err := apiClient.ACLs().List()
	if err != nil {
		currencyResolveErr = fmt.Errorf("infer currency from ACLs: %w", err)
		return "", currencyResolveErr
	}

	if activeOrgID != "" {
		if orgID, err := strconv.ParseInt(activeOrgID, 10, 64); err == nil {
			for _, acl := range acls {
				if acl.OrgID == orgID && acl.Currency != "" {
					resolvedCurrency = normalizeCurrencyCode(acl.Currency)
					return resolvedCurrency, nil
				}
			}
		}
	}

	// If we couldn't match org id, use the currency only if it's unambiguous.
	currencies := map[string]struct{}{}
	for _, acl := range acls {
		if acl.Currency == "" {
			continue
		}
		currencies[normalizeCurrencyCode(acl.Currency)] = struct{}{}
	}
	if len(currencies) == 1 {
		for cur := range currencies {
			resolvedCurrency = cur
			return resolvedCurrency, nil
		}
	}

	currencyResolveErr = fmt.Errorf("unable to determine org currency; pass --currency or set default_currency (or env var AADS_DEFAULT_CURRENCY/AADS_CURRENCY)")
	return "", currencyResolveErr
}

func moneyFromAmount(amount string) (*types.Money, error) {
	cur, err := resolveMoneyCurrency()
	if err != nil {
		return nil, err
	}
	return &types.Money{Amount: amount, Currency: cur}, nil
}
