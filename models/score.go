package models

import (
  "github.com/astaxie/beego/orm"
)

type Score struct {
  ID        string `orm:"pk"`
	Score			float64 `orm:"index"`
}

func SaveScore(score Score) {
    if(dbUse == "true") {
      o := orm.NewOrm()
      o.Using(dbAlias)
      updateScore := Score{ID: score.ID}
      if o.Read(&updateScore) == nil {
        o.Update(&score)
      } else {
        o.Insert(&score)
      }
    }
}