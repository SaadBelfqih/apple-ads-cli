package cmd

import (
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

const defaultPageSize = 1000

func effectivePageSize(pageSize int) int {
	if pageSize > 0 {
		return pageSize
	}
	return defaultPageSize
}

func collectAllOffsetPaginated[T any](pageSize, startOffset int, fetch func(limit, offset int) ([]T, *types.PageDetail, error)) ([]T, error) {
	limit := effectivePageSize(pageSize)
	offset := startOffset

	var out []T
	for {
		page, pag, err := fetch(limit, offset)
		if err != nil {
			return nil, err
		}
		out = append(out, page...)

		if len(page) == 0 {
			break
		}

		// If the API tells us total results, use it as the primary termination condition.
		if pag != nil && pag.TotalResults > 0 {
			if offset+len(page) >= pag.TotalResults {
				break
			}
			offset += len(page)
			continue
		}

		// Otherwise, stop when we get a partial page.
		if len(page) < limit {
			break
		}

		offset += len(page)
	}
	return out, nil
}

func collectAllSelectorPaginated[T any](sel *types.Selector, pageSize int, fetch func(*types.Selector) ([]T, *types.PageDetail, error)) ([]T, error) {
	if sel == nil {
		sel = &types.Selector{}
	}

	limit := effectivePageSize(pageSize)
	offset := 0
	if sel.Pagination != nil {
		offset = sel.Pagination.Offset
		if sel.Pagination.Limit > 0 {
			limit = sel.Pagination.Limit
		}
	}
	if limit <= 0 {
		return nil, fmt.Errorf("invalid page size")
	}

	var out []T
	for {
		reqSel := *sel // shallow copy so we don't mutate caller's selector
		reqSel.Pagination = &types.Pagination{Limit: limit, Offset: offset}

		page, pag, err := fetch(&reqSel)
		if err != nil {
			return nil, err
		}
		out = append(out, page...)

		if len(page) == 0 {
			break
		}

		if pag != nil && pag.TotalResults > 0 {
			if offset+len(page) >= pag.TotalResults {
				break
			}
			offset += len(page)
			continue
		}

		if len(page) < limit {
			break
		}

		offset += len(page)
	}

	return out, nil
}
