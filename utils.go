package musicmax

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetIntEnv(key string, def int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		slog.Error("GetIntValue error", slog.Any("key", key), slog.Any("error", err))
		return def
	}

	return val
}

func GetPageAndLimit(r *http.Request) (int, int, error) {
	page := 1
	limit := 10
	var err error
	if r.URL.Query().Has("page") {
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			err = fmt.Errorf("GetPageAndLimit page parse error: %w", err)
			return 0, 0, err
		}
	}

	if r.URL.Query().Has("limit") {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			err = fmt.Errorf("GetPageAndLimit limit parse error: %w", err)
			return 0, 0, err
		}
	}

	if page < 1 || limit < 1 {
		err = fmt.Errorf("page or limit is negative\npage = %d; limit = %d", page, limit)
		return 0, 0, err
	}
	return page, limit, nil
}

func GetQueryParams(r *http.Request) (map[string]string, error) {
	slog.Debug("GetQueryParams")
	filters := make(map[string]string)
	query := r.URL.Query().Encode()
	if strings.TrimSpace(query) == "" {
		return nil, nil
	}
	params := strings.Split(query, "&")
	slog.Info("params count after split", slog.Any("params", params), slog.Any("params count", len(params)))
	for _, v := range params {
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return nil, fmt.Errorf("param not valid %s", v)
		}
		key, err := url.QueryUnescape(s[0])
		if err != nil {
			err = fmt.Errorf("error during get filters\nparam: %s error: %w", v, err)
			return nil, err
		}
		value, err := url.QueryUnescape(s[1])
		if err != nil {
			err = fmt.Errorf("error during get filters\nparam: %s error: %w", v, err)
			return nil, err
		}
		filters[key] = value
	}
	return filters, nil
}
