// Package filter - Content managed by Project Forge, see [projectforge.md] for details.
package filter

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/util"
)

const MaxRowsDefault = 100

var AllowedColumns = map[string][]string{}

type Params struct {
	Key       string    `json:"key"`
	Orderings Orderings `json:"orderings,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

func ParamsWithDefaultOrdering(key string, params *Params, orderings ...*Ordering) *Params {
	if params == nil {
		return ParamsWithDefaultOrdering(key, &Params{Key: key}, orderings...)
	}
	if len(params.Orderings) == 0 {
		params.Orderings = orderings
	}
	return params
}

func (p *Params) Sanitize(key string, defaultOrderings ...*Ordering) *Params {
	if p == nil {
		return &Params{Key: key, Orderings: defaultOrderings}
	}
	if p.Limit == 0 || p.Limit > 100000 {
		p.Limit = MaxRowsDefault
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	if len(p.Orderings) == 0 {
		return p.CloneOrdering(defaultOrderings...)
	}
	return p
}

func (p *Params) WithLimit(n int) *Params {
	p.Limit = n
	return p
}

func (p *Params) CloneOrdering(orderings ...*Ordering) *Params {
	if p == nil {
		return nil
	}
	return &Params{Key: p.Key, Orderings: orderings, Limit: p.Limit, Offset: p.Offset}
}

func (p *Params) CloneLimit(limit int) *Params {
	if p == nil {
		return nil
	}
	return &Params{Key: p.Key, Orderings: p.Orderings, Limit: limit, Offset: p.Offset}
}

func (p *Params) HasNextPage(count int) bool {
	if p == nil || p.Limit == 0 {
		return false
	}
	return count >= (p.Offset + p.Limit)
}

func (p *Params) NextPage() *Params {
	limit := p.Limit
	if limit == 0 {
		limit = MaxRowsDefault
	}
	offset := p.Offset + limit
	if offset < 0 {
		offset = 0
	}
	return &Params{Key: p.Key, Orderings: p.Orderings, Limit: p.Limit, Offset: offset}
}

func (p *Params) HasPreviousPage() bool {
	return p != nil && p.Offset > 0
}

func (p *Params) PreviousPage() *Params {
	limit := p.Limit
	if limit == 0 {
		limit = MaxRowsDefault
	}
	offset := p.Offset - limit
	if offset < 0 {
		offset = 0
	}
	return &Params{Key: p.Key, Orderings: p.Orderings, Limit: p.Limit, Offset: offset}
}

func (p *Params) GetOrdering(col string) *Ordering {
	return lo.FindOrElse(p.Orderings, nil, func(o *Ordering) bool {
		return o.Column == col
	})
}

func (p *Params) OrderByString() string {
	ret := lo.Map(p.Orderings, func(o *Ordering, _ int) string {
		dir := ""
		if !o.Asc {
			dir = " desc"
		}
		return fmt.Sprintf("%q%s", o.Column, dir)
	})
	return strings.Join(ret, ", ")
}

func (p *Params) Filtered(available []string, logger util.Logger) *Params {
	if available == nil {
		available = AllowedColumns[p.Key]
	}
	if len(available) == 0 {
		logger.Warnf("no columns available for [%s]", p.Key)
	}
	if len(available) == 1 && available[0] == "*" {
		return p
	}
	if len(p.Orderings) > 0 {
		allowed := Orderings{}

		lo.ForEach(p.Orderings, func(o *Ordering, _ int) {
			containsCol := lo.ContainsBy(available, func(c string) bool {
				return c == o.Column
			})
			if containsCol {
				allowed = append(allowed, o)
			} else {
				const msg = "no column [%s] for [%s] available in allowed columns [%s]"
				logger.Warnf(msg, o.Column, p.Key, util.StringArrayOxfordComma(available, "and"))
			}
		})
		return &Params{Key: p.Key, Orderings: allowed, Limit: p.Limit, Offset: p.Offset}
	}
	return p
}

func (p *Params) String() string {
	ol := ""
	if p.Offset > 0 {
		ol += fmt.Sprintf("%d/", p.Offset)
	}
	if p.Limit > 0 {
		ol += fmt.Sprint(p.Limit)
	}
	ord := lo.Map(p.Orderings, func(o *Ordering, _ int) string {
		return o.String()
	})
	return fmt.Sprintf("%s(%s): %s", p.Key, ol, strings.Join(ord, " / "))
}

func (p *Params) ToQueryString(u *fasthttp.URI) string {
	if p == nil {
		return ""
	}

	if u == nil {
		return ""
	}

	ret := &fasthttp.Args{}
	u.QueryArgs().CopyTo(ret)

	ret.Del(p.Key + ".o")
	ret.Del(p.Key + ".l")
	ret.Del(p.Key + ".x")

	lo.ForEach(p.Orderings, func(o *Ordering, _ int) {
		s := o.Column
		if !o.Asc {
			s += ".d"
		}
		ret.Add(p.Key+".o", s)
	})

	if p.Limit != 0 && p.Limit != 1000 {
		ret.Add(p.Key+".l", fmt.Sprint(p.Limit))
	}

	if p.Offset > 0 {
		ret.Add(p.Key+".x", fmt.Sprint(p.Offset))
	}

	return string(ret.QueryString())
}
