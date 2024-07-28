package domains

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/skobelina/currency_converter/internal/constants"
)

type Filter struct {
	// Offset is a number of values that was skipped
	// in:query
	Offset int `json:"offset"`
	// Limit is a count of values
	// in:query
	Limit int `json:"limit"`
	// SortBy
	// in:query
	SortBy []string `json:"sortBy"`
}

func DefaultFilter() Filter {
	return Filter{Limit: constants.MaxSearchLimit}
}

type Pagination struct {
	// Order can be asc or dsc. asc by default
	Order string `json:"order"`
	// Offset is a number of values that was skipped
	Offset int `json:"offset"`
	// Limit is a count of values
	Limit int `json:"limit"`
	// TotalItems is a number of items
	TotalItems *int64 `json:"totalItems"`
}

func (f *Filter) OrderString() string {
	return strings.Join(f.SortBy, ",")
}

func (f *Filter) Validate() {
	f.SortBy = f.validatedSortBy()
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.Limit < 1 || f.Limit > constants.MaxSearchLimit {
		f.Limit = constants.DefaultSearchLimit
	}
}

func (f *Filter) validatedSortBy() []string {
	validatedSortByArray := make([]string, 0, len(f.SortBy))
	for _, sortBy := range f.SortBy {
		tokens := strings.Split(sortBy, " ")
		if len(tokens) != 1 && len(tokens) != 2 {
			continue
		}
		validatedSortBy := tokens[0]
		// update json model fields to database fields
		switch strings.ToLower(validatedSortBy) {
		case "id":
			validatedSortBy = "id"
		case "email":
			validatedSortBy = "email"
		}
		if len(tokens) == 2 {
			switch strings.ToLower(tokens[1]) {
			case constants.OrderASC:
				validatedSortBy += " asc"
			case constants.OrderDESC:
				validatedSortBy += " desc"
			default:
				continue
			}
		}
		validatedSortByArray = append(validatedSortByArray, validatedSortBy)
	}
	return validatedSortByArray
}

func GetFilterFromQuery(r *http.Request) (*Filter, error) {
	var (
		offset int
		limit  int
		err    error
	)
	params := r.URL.Query()
	if len(params["offset"]) != 0 {
		offset, err = strconv.Atoi(params["offset"][0])
		if err != nil {
			return nil, errors.New("cannot parse offset query param")
		}
	}
	if len(params["limit"]) != 0 {
		limit, err = strconv.Atoi(params["limit"][0])
		if err != nil {
			return nil, errors.New("cannot parse limit query param")
		}
	}

	filter := &Filter{
		Offset: offset,
		Limit:  limit,
		SortBy: params["sortBy"],
	}
	filter.Validate()

	return filter, nil
}
