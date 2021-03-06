package mssqlx

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// check bad connection
func isErrBadConn(err error) bool {
	if err != nil {
		if err == driver.ErrBadConn || err == sql.ErrConnDone || err == mysql.ErrInvalidConn {
			return true
		}

		// Postgres/Mysql Driver returns default driver.ErrBadConn
		// Mysql Driver ("github.com/go-sql-driver/mysql") is not
		s := strings.ToLower(err.Error())
		return s == "invalid connection" || s == "bad connection"
	}
	return false
}

// IsDeadlock ERROR 1213: Deadlock found when trying to get lock
func IsDeadlock(err error) bool {
	return isErrCode(err, 1213)
}

// IsWsrepNotReady ERROR 1047: WSREP has not yet prepared node for application use
func IsWsrepNotReady(err error) bool {
	return isErrCode(err, 1047)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}

	switch mErr := err.(type) {

	case *mysql.MySQLError:
		return mErr.Number == uint16(code)

	default:
		se := strings.ToLower(err.Error())
		return strings.HasPrefix(se, fmt.Sprintf("error %d:", code))
	}
}

func parseError(w *wrapper, err error) error {
	if err == nil {
		return nil
	}

	if w != nil && ping(w) != nil {
		return ErrNetwork
	}

	return err
}

func reportError(v string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("mssqlx;;%s;;%s;;%s;;%s\n", time.Now().Format("2006-01-02 15:04:05"), hostName, v, err.Error()))
	}
}
