package mysql

import (
	"fmt"
	"strings"

	"github.com/Zhan9Yunhua/blog-svr/utils"
)

func (m *db) Insert(stut interface{}) error {
	ks := utils.HandleStructTag(stut, "db")
	mm, err := utils.Struct2Json(stut)
	if err != nil {
		return err
	}
	var syb []string
	vs := make([]interface{}, 0)
	for _, v := range mm {
		syb = append(syb, "?")
		vs = append(vs, v)
	}

	sql := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES (%s)", m.table, strings.Join(ks, ","), strings.Join(syb, ","))
	stmt, err := m.DB.Prepare(sql)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(vs...); err != nil {
		return err
	}
	return nil
}
