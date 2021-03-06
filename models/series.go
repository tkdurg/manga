package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Series struct {
	Id            int              `orm:"auto"`
	Name          string           `orm:"column(name)"`
	Description   string           `orm:"column(description)"`
	CoverImage    string           `orm:"column(cover_image)"`
	TypeName      string           `orm:"column(type_name)"`
	TypeDemonym   string           `orm:"column(type_demonym)"`
	Status        *Statuses        `orm:"column(status);rel(fk)"`
	SeriesTags    []*SeriesTags    `orm:"reverse(many)"`
	SeriesAliases []*SeriesAliases `orm:"reverse(many)"`
}

func (t *Series) TableName() string {
	return "series"
}

func init() {
	orm.RegisterModel(new(Series))
}

// GetAllSeries gets all languages and returns an array on success
func GetAllSeriesArray() (series []*Series, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("series").All(&series)
	return
}

// AddSeries insert a new Series into database and returns
// last inserted Id on success.
func AddSeries(m *Series) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

//GetRecentSeriesByChapter Gets Series in order of chapter updates
func GetRecentSeriesByChapter(offset int, limit int) (b []SeriesChapters, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("series_chapters").OrderBy("time_uploaded").Limit(limit, offset).RelatedSel().All(&b)
	return
}

//  checks to see if username exists
// returns bool
func SeriesNameExists(name string) (exists bool) {
	o := orm.NewOrm()
	exists = o.QueryTable("series").Filter("Name", name).Exist()
	return
}

// GetSeriesByName retrieves Series by Name. Returns error if
// Id doesn't exist
func GetSeriesByName(name string) (v Series, err error) {
	o := orm.NewOrm()
	err = o.QueryTable("series").Filter("name__icontains", name).RelatedSel().One(&v)
	return
}

// GetSeriesById retrieves Series by Id. Returns error if
// Id doesn't exist
func GetSeriesById(id int) (v *Series, err error) {
	var temp Series

	o := orm.NewOrm()
	err = o.QueryTable("series").Filter("id", id).One(&temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

// GetAllSeries retrieves all Series matches certain condition. Returns empty list if
// no records exist
func GetAllSeries(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Series))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Series
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).RelatedSel().All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateSeries updates Series by Id and returns error if
// the record to be updated doesn't exist
func UpdateSeriesById(m *Series) (err error) {
	o := orm.NewOrm()
	v := Series{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSeries deletes Series by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSeries(id int) (err error) {
	o := orm.NewOrm()
	v := Series{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Series{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
