package searches

import "github.com/rendyfutsuybase-case-courses/helpers/request"

type CourseSearchHelper struct{ request.SearchPredefineBase }

func (CourseSearchHelper) GetSearchColumns() []string {
	return []string{
		"c.title",
		"c.short_description",
	}
}

func (CourseSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = CourseSearchHelper{}

func NewCourseSearchHelper() CourseSearchHelper {
	t := 0.75
	return CourseSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
