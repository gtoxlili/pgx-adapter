package pgxadapter

const (
	createTable = `
		create table if not exists %[1]s (
			id int generated always as identity primary key,
			ptype text not null,
			%[2]s
		);
		create unique index if not exists uk_%[1]s on %[1]s (ptype, %[3]s)
	`

	insertRow = "insert into %[1]s (ptype, %[2]s) values (%[3]s) on conflict (ptype, %[2]s) do nothing"

	updateRow = "update %[1]s set ptype = $1, %[2]s where ptype = $3 and %[3]s"

	deleteAll = "truncate table %[1]s restart identity"

	deleteRow = "delete from %[1]s where ptype = $1 and %[2]s"

	deleteByPType = "delete from %[1]s where ptype = $1"

	selectSQL = "select ptype, %[2]s from %[1]s"
)
