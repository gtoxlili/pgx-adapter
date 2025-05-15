package pgxadapter

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

const (
	// Table Name. The default table name should be casbin_rule.
	defaultTableName = "casbin_rule"
	// Ptype Column. Name of this column should be ptype instead of p_type or Ptype.
	defaultPtypeColumn = "ptype"
	// Database Name. The default database name should be casbin.
	defaultDatabaseName = "casbin"
	// Data Structure. Adapter should support reading at least six columns.
	defaultFieldCount = 6
)

type store struct {
	db         *pgxpool.Pool
	tableName  string
	fieldCount int
}

func newStore(db *pgxpool.Pool) *store {
	return &store{db: db, fieldCount: defaultFieldCount, tableName: defaultTableName}
}

func (s *store) setFieldCount(fieldCount int) {
	s.fieldCount = fieldCount
}

func (s *store) setTableName(tableName string) {
	s.tableName = tableName
}

func (s *store) initTable(ctx context.Context) error {
	vColumns := lo.Times(20, func(i int) string {
		return "v" + strconv.Itoa(i)
	})
	sqlSeq := strings.SplitSeq(fmt.Sprintf(createTableSQL, lo.SnakeCase(s.tableName), strings.Join(lo.Map(vColumns, func(v string, _ int) string {
		return v + " text not null"
	}), ","), strings.Join(vColumns, ",")), ";")

	batch := &pgx.Batch{}
	for sql := range sqlSeq {
		batch.Queue(sql)
	}

	br := s.db.SendBatch(ctx, batch)
	defer br.Close()

	// batch.Len()
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch: %v", err)
		}
	}
	return nil
}
