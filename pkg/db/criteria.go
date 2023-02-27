package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"strings"
)

type CriteriaBuilder struct {
	conditions []criteria.ICondition
	ands       []criteria.ICriteriaBuilder
	ors        []criteria.ICriteriaBuilder
	nots       []criteria.ICriteriaBuilder
	pagination []int
	sort       [][]string
	prefix     string
}

func NewCriteriaBuilder(prefix string) criteria.ICriteriaBuilder {
	return &CriteriaBuilder{prefix: prefix}
}

func (d CriteriaBuilder) has(slices []criteria.ICriteriaBuilder, other string) bool {
	for _, k := range d.ands {
		if v, ok := k.(CriteriaBuilder); ok && v.prefix == other {
			return true
		}
	}
	return false
}
func (d CriteriaBuilder) Has(other string) bool {
	return d.has(d.ands, other) || d.has(d.ors, other) || d.has(d.nots, other)
}
func (d CriteriaBuilder) Prefix() string {
	return d.prefix
}

func (d CriteriaBuilder) Paginate(page int, perPage int) criteria.ICriteriaBuilder {
	d.pagination = []int{page, perPage}
	return d
}

func (d CriteriaBuilder) Order(field string, direction string) criteria.ICriteriaBuilder {
	d.sort = append(d.sort, []string{field, strings.ToLower(direction)})
	return d
}

func (d CriteriaBuilder) Copy() criteria.ICriteriaBuilder {
	return &CriteriaBuilder{
		conditions: d.conditions,
		ands:       d.ands,
		ors:        d.ors,
		pagination: d.pagination,
		sort:       d.sort,
		prefix:     d.prefix,
	}
}

func (d CriteriaBuilder) Select(fields ...string) criteria.ICriteriaBuilder {
	return CriteriaBuilder{}
}

func (d CriteriaBuilder) Where(condition ...criteria.ICondition) criteria.ICriteriaBuilder {
	for _, r := range condition {
		if r == nil {
			continue
		}
		d.conditions = append(d.conditions, r)
	}
	return d
}

func (d CriteriaBuilder) And(other ...criteria.ICriteriaBuilder) criteria.ICriteriaBuilder {
	d.ands = append(d.ands, other...)
	return d
}

func (d CriteriaBuilder) Or(other ...criteria.ICriteriaBuilder) criteria.ICriteriaBuilder {
	d.ors = append(d.ors, other...)
	return d
}
func (d CriteriaBuilder) Not(other ...criteria.ICriteriaBuilder) criteria.ICriteriaBuilder {
	d.nots = append(d.nots, other...)
	return d
}

func (d CriteriaBuilder) ToString() (res string) {
	concats := make([]string, 0)
	if len(d.pagination) > 0 && (d.pagination[0] > 0 || d.pagination[1] > 0) {
		concats = append(concats, fmt.Sprintf("pagination(%d,%d)", d.pagination[0], d.pagination[1]))
	}
	for _, s := range d.sort {
		concats = append(concats, fmt.Sprintf("sort(%s, %s)", s[0], s[1]))
	}
	for _, r := range d.conditions {
		concats = append(concats, r.ToString())
	}
	if len(d.nots) > 0 {
		var nots []string
		for _, a := range d.nots {
			nots = append(nots, a.ToString())
		}
		concats = append(concats, fmt.Sprintf("nots(%s)", strings.Join(nots, ",")))
	}
	if len(d.ands) > 0 {
		var ands []string
		for _, a := range d.ands {
			ands = append(ands, a.ToString())
		}
		concats = append(concats, fmt.Sprintf("ands(%s)", strings.Join(ands, ",")))
	}
	if len(d.ors) > 0 {
		var ands []string
		for _, a := range d.ors {
			ands = append(ands, a.ToString())
		}
		concats = append(concats, fmt.Sprintf("ors(%s)", strings.Join(ands, ",")))
	}
	res = strings.Join(concats, ";")
	return
}

func operator(condition criteria.ICondition) string {
	switch condition.Operator() {
	case criteria.Eq:
		return "="
	case criteria.Gte:
		return ">="
	case criteria.Gt:
		return ">"
	case criteria.Lt:
		return "<"
	case criteria.Lte:
		return "<="
	case criteria.Not:
		return "!="
	case criteria.Like:
		return "like"
	case criteria.NotLike:
		return "not like"
	case criteria.In:
		return "in"
	case criteria.NotIn:
		return "not in"
	}
	return ""
}
func value(operator criteria.ICondition) interface{} {
	if operator.Operator() == criteria.Like || operator.Operator() == criteria.NotLike {
		return fmt.Sprintf("%%%s%%", operator.Value())
	}
	return operator.Value()
}

func (d CriteriaBuilder) column(name string) string {
	if d.prefix == "" {
		return name
	}
	if strings.Contains(name, "CAST(") {
		return strings.ReplaceAll(name, "CAST(", fmt.Sprintf("CAST(%s.", d.prefix))
	}
	return fmt.Sprintf("%s.%s", d.prefix, name)
}
func (d CriteriaBuilder) Apply(db *gorm.DB) *gorm.DB {
	ori := db
	if len(d.pagination) > 0 {
		db = db.
			Limit(d.pagination[1]).
			Offset((d.pagination[0] - 1) * d.pagination[1])
	}
	if len(d.sort) > 0 {
		for _, srt := range d.sort {
			isDesc := true
			if srt[1] == "asc" {
				isDesc = false
			}
			db = db.Order(clause.OrderByColumn{
				Column:  clause.Column{Name: d.column(srt[0])},
				Desc:    isDesc,
				Reorder: true,
			})
		}
	}
	if len(d.conditions) > 0 {
		ses := db.WithContext(context.TODO())
		for _, r := range d.conditions {
			ses = ses.Where(fmt.Sprintf("%s %s ?", d.column(r.Field()), operator(r)), value(r))
		}
		db = db.Where(ses)
	}
	if len(d.ands) > 0 {
		tx := ori.WithContext(context.TODO())
		for _, a := range d.ands {
			if a == nil {
				continue
			}
			ses := a.(CriteriaBuilder).Apply(ori.WithContext(context.TODO()))
			tx = tx.Where(ses)
		}
		db = db.Where(tx)
	}
	if len(d.nots) > 0 {
		tx := ori.WithContext(context.TODO())
		for _, a := range d.nots {
			if a == nil {
				continue
			}
			ses := a.(CriteriaBuilder).Apply(ori.WithContext(context.TODO()))
			tx = tx.Where(ses)
		}
		db = db.Not(tx)
	}
	if len(d.ors) > 0 {
		tx := ori.WithContext(context.TODO())
		for idx, a := range d.ors {
			if a == nil {
				continue
			}
			ses := a.(CriteriaBuilder).Apply(ori.WithContext(context.TODO()))
			if idx == 0 {
				tx = tx.Where(ses)
			} else {
				tx = tx.Or(ses)
			}
		}
		db = db.Where(tx)
	}
	return db
}
