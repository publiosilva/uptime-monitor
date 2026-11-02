package query

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrNotSelect           = errors.New("sql must be a single SELECT statement")
	ErrMultipleStatements  = errors.New("sql must not contain multiple statements")
	ErrForbiddenSQLKeyword = errors.New("sql contains a forbidden keyword")
	ErrMissingUserParam    = errors.New("sql must bind the authenticated user via $1")
)

var forbiddenKeyword = regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|GRANT|REVOKE|COPY|EXECUTE|CALL|MERGE|VACUUM|REINDEX|CLUSTER|COMMENT|SECURITY|OWNER|INTO)\b`)

// ValidateReadOnlySQL ensures sql is a single SELECT/CTE that binds $1 as the user scope parameter.
func ValidateReadOnlySQL(sql string) error {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return ErrNotSelect
	}

	trimmed = strings.TrimSuffix(trimmed, ";")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return ErrNotSelect
	}

	if strings.Contains(trimmed, ";") {
		return ErrMultipleStatements
	}

	upper := strings.ToUpper(trimmed)
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return ErrNotSelect
	}

	if forbiddenKeyword.MatchString(trimmed) {
		return ErrForbiddenSQLKeyword
	}

	if !strings.Contains(trimmed, "$1") {
		return ErrMissingUserParam
	}

	return nil
}
