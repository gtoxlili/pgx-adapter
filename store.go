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

func (s *store) insertRow(ctx context.Context, ptype string, args ...string) error {
	if len(args) != s.fieldCount {
		return fmt.Errorf("args length %d is not equal to field count %d", len(args), s.fieldCount)
	}

	sql := fmt.Sprintf(insertRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i)
	}), ","), strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "$" + strconv.Itoa(i+2)
	}), ","))

	_, err := s.db.Exec(ctx, sql, lo.ToAnySlice(genRule(ptype, args))...)
	if err != nil {
		return fmt.Errorf("failed to insert row: %v", err)
	}
	return nil
}

func (s *store) selectAll(ctx context.Context) ([][]string, error) {
	return s.selectWhere(ctx, "", 0)
}

func (s *store) selectWhere(ctx context.Context, ptype string, startIdx int, args ...string) ([][]string, error) {
	// args 需 小于等于 fieldCount - startIdx
	if len(args) > s.fieldCount-startIdx {
		return nil, fmt.Errorf("args length %d is greater than field count %d", len(args), s.fieldCount-startIdx)
	}

	sql := fmt.Sprintf(selectSQL, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i)
	}), ","), ptype)

	tmpArgs := args
	var conditions []string
	if lo.IsNotEmpty(ptype) {
		conditions = append(conditions, "ptype = $1")
		tmpArgs = genRule(ptype, args)
	}
	conditions = append(conditions, lo.Map(lo.Filter(lo.Map(args, func(arg string, i int) string {
		return "v" + strconv.Itoa(i+startIdx)
	}), func(_ string, i int) bool {
		return lo.IsNotEmpty(args[i])
	}), func(arg string, i int) string {
		return arg + " = $" + strconv.Itoa(i+1+len(conditions))
	})...)
	if len(conditions) > 0 {
		sql += " where " + strings.Join(conditions, " and ")
	}

	rows, err := s.db.Query(ctx, sql, lo.ToAnySlice(lo.Compact(tmpArgs))...)
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

// old, updated
func (s *store) updateRow(ctx context.Context, ptype string, old, updated []string) error {
	if len(old) != s.fieldCount || len(updated) != s.fieldCount {
		return fmt.Errorf("args length (old: %d, updated: %d) is not equal to field count %d", len(old), len(updated), s.fieldCount)
	}

	sql := fmt.Sprintf(updateRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+2)
	}), ", "), strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+s.fieldCount+2)
	}), " and "))

	merged := append(old, updated...)
	_, err := s.db.Exec(ctx, sql, lo.ToAnySlice(genRule(ptype, merged))...)
	if err != nil {
		return fmt.Errorf("failed to update row: %v", err)
	}
	return nil
}

func (s *store) deleteRow(ctx context.Context, ptype string, args ...string) error {
	if len(args) != s.fieldCount {
		return fmt.Errorf("args length %d is not equal to field count %d", len(args), s.fieldCount)
	}

	sql := fmt.Sprintf(deleteRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
		return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+2)
	}), " and "))

	_, err := s.db.Exec(ctx, sql, lo.ToAnySlice(genRule(ptype, args))...)
	if err != nil {
		return fmt.Errorf("failed to delete row: %v", err)
	}
	return nil
}

func (s *store) deleteByPType(ctx context.Context, ptype string) error {
	sql := fmt.Sprintf(deleteByPType, s.tableName)
	_, err := s.db.Exec(ctx, sql, ptype)
	if err != nil {
		return fmt.Errorf("failed to delete by ptype: %v", err)
	}
	return nil
}

func (s *store) deleteWhere(ctx context.Context, ptype string, startIdx int, args ...string) error {
	if ptype == "" {
		return fmt.Errorf("ptype is empty")
	}

	// args 需 小于等于 fieldCount - startIdx
	if len(args) > s.fieldCount-startIdx {
		return fmt.Errorf("args length %d is greater than field count %d", len(args), s.fieldCount-startIdx)
	}

	sql := fmt.Sprintf(deleteByPType, s.tableName)

	conditions := strings.Join(lo.Map(lo.Filter(lo.Map(args, func(_ string, i int) string {
		return "v" + strconv.Itoa(i+startIdx)
	}), func(_ string, i int) bool {
		return lo.IsNotEmpty(args[i])
	}), func(arg string, i int) string {
		return arg + " = $" + strconv.Itoa(i+2)
	}), " and ")

	if len(conditions) > 0 {
		sql += " and " + conditions
	}

	_, err := s.db.Exec(ctx, sql, lo.ToAnySlice(lo.Compact(genRule(ptype, args)))...)
	if err != nil {
		return fmt.Errorf("failed to delete where: %v", err)
	}
	return nil
}

func (s *store) deleteAndInsertAll(ctx context.Context, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// 1. 删除现有的所有数据
	_, err = tx.Exec(ctx, fmt.Sprintf(deleteAll, s.tableName))
	if err != nil {
		return fmt.Errorf("failed to drop table: %v", err)
	}

	// 2. 通过 Batch 插入数据
	batch := &pgx.Batch{}
	for _, rule := range rules {
		sql := fmt.Sprintf(insertRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "v" + strconv.Itoa(i)
		}), ","), strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "$" + strconv.Itoa(i+2)
		}), ","))
		batch.Queue(sql, lo.ToAnySlice(rule)...)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err = br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch: %v", err)
		}
	}

	return tx.Commit(ctx)
}

// 批量插入/删除数据
func (s *store) batchInsert(ctx context.Context, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, rule := range rules {
		sql := fmt.Sprintf(insertRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "v" + strconv.Itoa(i)
		}), ","), strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "$" + strconv.Itoa(i+2)
		}), ","))
		batch.Queue(sql, lo.ToAnySlice(genRule(ptype, rule))...)
	}

	br := s.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch: %v", err)
		}
	}
	return nil
}

func (s *store) batchDelete(ctx context.Context, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, rule := range rules {
		sql := fmt.Sprintf(deleteRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+2)
		}), " and "))
		batch.Queue(sql, lo.ToAnySlice(genRule(ptype, rule))...)
	}

	br := s.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch: %v", err)
		}
	}
	return nil
}

func (s *store) batchUpdate(ctx context.Context, ptype string, oldRules, newRules [][]string) error {
	if len(oldRules) == 0 || len(newRules) == 0 {
		return nil
	}
	if len(oldRules) != len(newRules) {
		return fmt.Errorf("oldRules and newRules length mismatch: %d vs %d", len(oldRules), len(newRules))
	}

	batch := &pgx.Batch{}
	for i := 0; i < len(oldRules); i++ {
		if len(oldRules[i]) != s.fieldCount || len(newRules[i]) != s.fieldCount {
			return fmt.Errorf("args[%d] length (old: %d, updated: %d) is not equal to field count %d", i, len(oldRules[i]), len(newRules[i]), s.fieldCount)
		}
		sql := fmt.Sprintf(updateRow, s.tableName, strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+2)
		}), ", "), strings.Join(lo.Times(s.fieldCount, func(i int) string {
			return "v" + strconv.Itoa(i) + " = $" + strconv.Itoa(i+s.fieldCount+2)
		}), " and "))

		merged := append(oldRules[i], newRules[i]...)
		batch.Queue(sql, lo.ToAnySlice(genRule(ptype, merged))...)
	}

	br := s.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch: %v", err)
		}
	}
	return nil
}
