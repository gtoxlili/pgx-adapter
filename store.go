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
	s.tableName = lo.SnakeCase(tableName)
}

func (s *store) initTable(ctx context.Context) error {
	return s.createTable(ctx)
}

func (s *store) createTable(ctx context.Context) error {
	vColumns := lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i)
	})
	sqlSeq := strings.SplitSeq(fmt.Sprintf(createTable, s.tableName, strings.Join(lo.Map(vColumns, func(v string, _ int) string {
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

func (s *store) dropTable(ctx context.Context) error {
	sql := fmt.Sprintf(deleteAll, s.tableName)
	_, err := s.db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to drop table: %v", err)
	}
	return nil
}

func (s *store) insertRow(ctx context.Context, args ...string) error {
	if len(args) != s.fieldCount+1 {
		return fmt.Errorf("args length %d is not equal to field count %d", len(args)-1, s.fieldCount)
	}

	sql := fmt.Sprintf(insertRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i)
	}), ","), strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "$" + strconv.Itoa(i+1)
	}), ","))

	_, err := s.db.Exec(ctx, sql, lo.ToAnySlice(args)...)
	if err != nil {
		return fmt.Errorf("failed to insert row: %v", err)
	}
	return nil
}

func (s *store) selectAll(ctx context.Context) ([][]string, error) {
	return s.selectWhere(ctx, "", 0)
}

func (s *store) selectWhere(ctx context.Context, ptype string, startIdx int, args ...string) ([][]string, error) {
	sql := fmt.Sprintf(selectSQL, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i)
	}), ","), ptype)

	var conditions []string
	if lo.IsNotEmpty(ptype) {
		conditions = append(conditions, "ptype = $1")
		args = append([]string{ptype}, args...)
	}
	conditions = append(conditions, lo.Map(lo.Filter(lo.Map(args, func(arg string, i int) string {
		return "v" + strconv.Itoa(i+startIdx)
	}), func(_ string, i int) bool {
		return lo.IsEmpty(args[i])
	}), func(arg string, i int) string {
		return arg + " = $" + strconv.Itoa(i+1+len(conditions))
	})...)

	rows, err := s.db.Query(ctx, sql, lo.ToAnySlice(lo.Compact(args))...)
	if err != nil {
		return nil, fmt.Errorf("failed to select where: %v", err)
	}
	defer rows.Close()

	var result [][]string
	for rows.Next() {
		row := lo.ToSlicePtr(make([]string, s.fieldCount+1))
		err := rows.Scan(lo.ToAnySlice(row)...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result = append(result, lo.FromSlicePtr(row))
	}
	return result, nil
}
