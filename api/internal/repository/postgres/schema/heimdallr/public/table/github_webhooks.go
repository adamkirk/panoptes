//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var GithubWebhooks = newGithubWebhooksTable("public", "github_webhooks", "")

type githubWebhooksTable struct {
	postgres.Table

	// Columns
	ID         postgres.ColumnString
	OccurredAt postgres.ColumnTimestampz
	Payload    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type GithubWebhooksTable struct {
	githubWebhooksTable

	EXCLUDED githubWebhooksTable
}

// AS creates new GithubWebhooksTable with assigned alias
func (a GithubWebhooksTable) AS(alias string) *GithubWebhooksTable {
	return newGithubWebhooksTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new GithubWebhooksTable with assigned schema name
func (a GithubWebhooksTable) FromSchema(schemaName string) *GithubWebhooksTable {
	return newGithubWebhooksTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new GithubWebhooksTable with assigned table prefix
func (a GithubWebhooksTable) WithPrefix(prefix string) *GithubWebhooksTable {
	return newGithubWebhooksTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new GithubWebhooksTable with assigned table suffix
func (a GithubWebhooksTable) WithSuffix(suffix string) *GithubWebhooksTable {
	return newGithubWebhooksTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newGithubWebhooksTable(schemaName, tableName, alias string) *GithubWebhooksTable {
	return &GithubWebhooksTable{
		githubWebhooksTable: newGithubWebhooksTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newGithubWebhooksTableImpl("", "excluded", ""),
	}
}

func newGithubWebhooksTableImpl(schemaName, tableName, alias string) githubWebhooksTable {
	var (
		IDColumn         = postgres.StringColumn("id")
		OccurredAtColumn = postgres.TimestampzColumn("occurred_at")
		PayloadColumn    = postgres.StringColumn("payload")
		allColumns       = postgres.ColumnList{IDColumn, OccurredAtColumn, PayloadColumn}
		mutableColumns   = postgres.ColumnList{OccurredAtColumn, PayloadColumn}
	)

	return githubWebhooksTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		OccurredAt: OccurredAtColumn,
		Payload:    PayloadColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}