package testutils

import (
	"context"
	"fmt"

	"github.com/jasonsites/gosk/internal/resolver"
)

// Cleanup deletes all rows on all database tables
func Cleanup(r *resolver.Resolver) error {
	db := r.PostgreSQLClient()

	tables := []string{"resource_entity"}

	for _, t := range tables {
		sql := fmt.Sprintf("DELETE from %s", t)
		_, err := db.Exec(context.TODO(), sql)
		if err != nil {
			return err
		}
	}

	return nil
}
