package pgxadapter

const (
	createTableSQL = `
		create table if not exists %[1]s (
			id int generated always as identity primary key,
			ptype text not null,
			%[2]s
		);
		create unique index if not exists uk_%[1]s on %[1]s (ptype, %[3]s);
	`
)
