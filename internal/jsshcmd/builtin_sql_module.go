package jsshcmd

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"time"
)

var globalSqlConn *sqlx.DB
var globalSqlConfig typeutil.H

func init() {
	globalSqlConfig = typeutil.H{
		"connMaxLifetime": 60_000,
	}
}

func JsFnSqlSet(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sql.set: missing name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("sql.set: first argument expected string type")
		}
		name := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("sql.set: missing value")
		}
		value := args[1]

		if name == "connMaxLifetime" && !value.IsNumber() {
			return ctx.ThrowTypeError("sql.set: [connMaxLifetime] expected string type")
		}

		v, err := jsexecutor.JSValueToAny(value)
		if err != nil {
			return ctx.ThrowError(err)
		}
		globalSqlConfig[name] = v

		return ctx.Bool(true)
	}
}

func JsFnSqlOpen(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("mysql.open: missing driver name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("mysql.open: first argument expected string type")
		}
		driverName := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("mysql.open: missing data source name")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("mysql.open: second argument expected string type")
		}
		dataSourceName := args[1].String()

		if globalSqlConn != nil {
			return ctx.ThrowInternalError("sql.open: please close the previous connection")
		}

		db, err := sqlx.Connect(driverName, dataSourceName)
		if err != nil {
			return ctx.ThrowError(err)
		}
		globalSqlConn = db

		if connMaxLifetime, ok := globalSqlConfig["connMaxLifetime"].(float64); ok {
			db.SetConnMaxLifetime(time.Millisecond * time.Duration(connMaxLifetime))
		}

		return ctx.Bool(true)
	}
}

func JsFnSqlClose(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if globalSqlConn != nil {
			if err := globalSqlConn.Close(); err != nil {
				errLog.Printf("sql.close: close sql connection fail: %s", err)
			}
		}
		globalSqlConn = nil

		return ctx.Bool(true)
	}
}

func JsFnSqlQuery(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sql.query: missing sql")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("sql.query: first argument expected string type")
		}
		format := args[0].String()

		a := make([]interface{}, 0)
		for _, v := range args[1:] {
			v2, err := jsexecutor.JSValueToAny(v)
			if err != nil {
				return ctx.ThrowError(err)
			}
			a = append(a, v2)
		}

		if globalSqlConn == nil {
			return ctx.ThrowInternalError("sql.query: please open a sql connection")
		}

		rows, err := sqlQueryManyToMap(globalSqlConn, format, a...)
		if err != nil {
			return ctx.ThrowInternalError("sql.query: %s", err)
		}
		return jsexecutor.AnyToJSValue(ctx, rows)
	}
}

func JsFnSqlExec(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sql.exec: missing sql")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("sql.exec: first argument expected string type")
		}
		format := args[0].String()

		a := make([]interface{}, 0)
		for _, v := range args[1:] {
			v2, err := jsexecutor.JSValueToAny(v)
			if err != nil {
				return ctx.ThrowError(err)
			}
			a = append(a, v2)
		}

		if globalSqlConn == nil {
			return ctx.ThrowInternalError("sql.exec: please open a sql connection")
		}

		result, err := globalSqlConn.Exec(format, a...)
		if err != nil {
			return ctx.ThrowInternalError("sql.exec: %s", err)
		}
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return ctx.ThrowInternalError("sql.exec: %s", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return ctx.ThrowInternalError("sql.exec: %s", err)
		}
		return jsexecutor.AnyToJSValue(ctx, typeutil.H{
			"lastInsertId": lastInsertId,
			"rowsAffected": rowsAffected,
		})
	}
}

func sqlQueryManyToMap(tx *sqlx.DB, format string, args ...interface{}) (rows []typeutil.H, err error) {
	result, err := tx.Queryx(format, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer func() {
		if err := result.Close(); err != nil {
			errLog.Printf("sql.query: %s", err)
		}
	}()
	rows = make([]typeutil.H, 0)
	for result.Next() {
		row := make(typeutil.H)
		if err := result.MapScan(row); err != nil {
			return nil, err
		}
		for n, v := range row {
			if b, ok := v.([]byte); ok {
				row[n] = string(b)
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}
