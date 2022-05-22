package app

import (
	"SqlLineage/src/models/do"
)

func QueryDatasetsByDatabaseId(databaseId uint) []do.Dataset {
	datasets := do.SelectDatasetsByDatabaseId(databaseId)
	return datasets
}

func QueryVDatasets() []do.Dataset {
	datasets := do.SelectDatasetsByDatabaseId(0)
	return datasets
}
