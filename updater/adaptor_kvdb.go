package updater

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/kvdb"
)

var (
	kvDatabasePrefix = "updater_record/"
)

type KVDatabaseAdaptor struct {
	db kvdb.Database // updater_record/{name}_{type}: Record
}

func NewKVDatabaseAdaptor(db kvdb.Database) (a RecordAdaptor) {
	return &KVDatabaseAdaptor{db: db}
}

func (k KVDatabaseAdaptor) Store(ctx context.Context, record *Record) (err error) {
	return k.db.Set(ctx, k.getKey(record), record, 0)
}

func (k KVDatabaseAdaptor) Load(ctx context.Context, params *RecordQueryParams) (res *RecordQueryResult, err error) {
	res = &RecordQueryResult{}
	kvs, err := k.db.GetPrefix(ctx, kvDatabasePrefix)
	if err != nil {
		return
	}

	for _, v := range kvs {
		record := new(Record)
		if err := v.Value.Struct(&record); err != nil {
			g.Log().Warningf(ctx, "failed to parse record: key=%s err=%v", v.Key, err)
			continue
		}
		if params.Name != nil {
			if &record.Name != params.Name {
				continue
			}
		}
		if params.Type != nil {
			if &record.Type != params.Type {
				continue
			}
		}
		res.Records = append(res.Records, record)
	}
	res.Total = len(res.Records)
	return
}

func (k KVDatabaseAdaptor) getKey(record *Record) string {
	return fmt.Sprintf("%s%s_%v", kvDatabasePrefix, record.Name, record.Type)
}
