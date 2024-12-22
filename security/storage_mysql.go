package security

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	mysqlDbNameConfigPattern = "security.database.group"
	tableName                = "c_security"
	createTableDDL           = `create table if not exists test.c_security
(
    id      int auto_increment primary key,
    type    varchar(20) not null,
    name    varchar(20) null,
    content mediumtext  null
);`
	storageObjectNamePrivateKey = "private_key"
	storageObjectNamePublicKey  = "public_key"
	storageObjectTypeRsa        = "rsa"
)

type (
	mysqlStorage struct {
		db string
	}
	mysqlStorageObject struct {
		Id      int    `orm:"id" json:"id"`
		Type    string `orm:"type" json:"type"`
		Name    string `orm:"name" json:"name"`
		Content string `orm:"content" json:"content"`
	}
)

func NewMysqlStorage(ctx context.Context) Storage {
	s := &mysqlStorage{}
	v, err := g.Cfg().Get(ctx, mysqlDbNameConfigPattern)
	if err == nil {
		if name := v.String(); name != "" {
			s.db = name
		}
	}
	if err = s.createTableIfNotExists(ctx); err != nil {
		panic(err)
	}
	return s
}

func (m *mysqlStorage) createTableIfNotExists(ctx context.Context) (err error) {
	tables, err := g.DB(m.db).Ctx(ctx).Tables(ctx)
	if err != nil {
		return
	}
	for _, table := range tables {
		if table == tableName {
			return
		}
	}
	_, err = g.DB(m.db).Ctx(ctx).Exec(ctx, createTableDDL)
	return
}

func (m *mysqlStorage) StorePublicKey(ctx context.Context, data []byte) (err error) {
	return m.upsert(ctx, data, storageObjectNamePublicKey, storageObjectTypeRsa)

}

func (m *mysqlStorage) StorePrivateKey(ctx context.Context, data []byte) (err error) {
	return m.upsert(ctx, data, storageObjectNamePrivateKey, storageObjectTypeRsa)
}

func (m *mysqlStorage) upsert(ctx context.Context, data []byte, name, typ string) (err error) {
	r, err := g.DB(m.db).Ctx(ctx).Model(tableName).Where("type = ? and name = ?", typ, name).One(ctx)
	switch {
	case r == nil || errors.Is(err, sql.ErrNoRows):
		// insert
		_, err = g.DB(m.db).Ctx(ctx).Model(tableName).Insert(mysqlStorageObject{
			Type:    typ,
			Name:    name,
			Content: string(gbase64.Encode(data)),
		})
	case err == nil:
		// update
		_, err = g.DB(m.db).Ctx(ctx).Model(tableName).Where("type = ? and name = ?", typ, name).Update(g.Map{
			"content": string(gbase64.Encode(data)),
		})
	default:
	}
	return
}

func (m *mysqlStorage) LoadPublicKey(ctx context.Context) (err error) {
	o, err := m.load(ctx, storageObjectNamePublicKey, storageObjectTypeRsa)
	if err == nil && len(o.Content) > 0 {
		var content []byte
		if content, err = gbase64.DecodeString(o.Content); err != nil {
			return
		}
		publicKey, err = decodePublicKeyPem(content)
	}
	return
}

func (m *mysqlStorage) LoadPrivateKey(ctx context.Context) (err error) {
	o, err := m.load(ctx, storageObjectNamePrivateKey, storageObjectTypeRsa)
	if err == nil && len(o.Content) > 0 {
		var content []byte
		if content, err = gbase64.DecodeString(o.Content); err != nil {
			return
		}
		privateKey, err = decodePrivateKeyPem(content)
	}
	return
}

func (m *mysqlStorage) load(ctx context.Context, name, typ string) (o *mysqlStorageObject, err error) {
	o = new(mysqlStorageObject)
	r, err := g.DB(m.db).Ctx(ctx).Model(tableName).Where("type = ? and name = ?", typ, name).One(ctx)
	if err != nil && !errors.Is(gerror.Cause(err), sql.ErrNoRows) {
		return
	}
	err = r.Struct(o)
	if err != nil && errors.Is(gerror.Cause(err), sql.ErrNoRows) {
		err = nil
	}
	return
}
