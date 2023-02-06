package common

const (
	Schema      = "schema"
	Collection  = "collection"
	Version     = "version"
	Uid         = "_id"
	MergeFields = "merge_fields"
	QueryOR     = "queryOr"
	Sort        = "querySort"
)

type StoreType string

const StoreTypeMongo StoreType = "mongo"
