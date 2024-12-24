package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/junqirao/gocomponents/auth/dao"
)

type applicationAuthDao = *dao.ApplicationAuthDao

type DataSourceMysql struct {
	applicationAuthDao
}

func NewDataSourceMysql(db gdb.DB, table ...string) *DataSourceMysql {
	tableName := "application_auth"
	if len(table) > 0 && table[0] != "" {
		tableName = table[0]
	}
	return &DataSourceMysql{applicationAuthDao: dao.NewApplicationAuthDao(db, tableName)}
}

func (d DataSourceMysql) List(ctx context.Context, params *ListParams) (res *ListResult, err error) {
	res = new(ListResult)
	query := d.Ctx(ctx)
	if params.AppId != "" {
		if params.Fuzzy {
			query = query.WhereLike(d.Columns().AppId, "%"+params.AppId+"%")
		} else {
			query = query.Where(d.Columns().AppId, params.AppId)
		}
	}
	if params.Name != "" {
		if params.Fuzzy {
			query = query.WhereLike(d.Columns().Name, "%"+params.Name+"%")
		} else {
			query = query.Where(d.Columns().Name, params.Name)
		}
	}
	res.Total, err = query.Count()
	res.List = make([]*AppFullInfo, 0)
	err = query.Scan(&res.List)
	return
}

func (d DataSourceMysql) Store(ctx context.Context, app *AppFullInfo) (err error) {
	v := new(AppInfo)
	err = d.Ctx(ctx).Where(d.Columns().AppId, app.AppId).Scan(&v)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = d.Ctx(ctx).Insert(app)
	case err == nil:
		_, err = d.Ctx(ctx).Where(d.Columns().AppId, app.AppId).Data(g.Map{
			d.Columns().Name:        app.Name,
			d.Columns().Description: app.Description,
			d.Columns().AppKey:      app.AppKey,
			d.Columns().AppSecret:   app.AppSecret,
		}).Update()
	default:
		return err
	}
	return
}

func (d DataSourceMysql) Delete(ctx context.Context, appId string) (err error) {
	_, err = d.Ctx(ctx).Where(d.Columns().AppId, appId).Delete()
	return
}

func (d DataSourceMysql) FindOne(ctx context.Context, appId string) (app *AppInfo, err error) {
	app = new(AppInfo)
	err = d.Ctx(ctx).Where(d.Columns().AppId, appId).Scan(&app)
	return
}
