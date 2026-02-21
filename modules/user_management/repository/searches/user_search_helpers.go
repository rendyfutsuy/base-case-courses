package searches

import "github.com/rendyfutsuybase-case-courses/helpers/request"

// initialize, value for search and map the function & variable need for it
type UserSearchHelper struct{ request.SearchPredefineBase }

func (UserSearchHelper) GetSearchColumns() []string {
	return []string{
		"usr.full_name",
		"usr.gender",
		"usr.email",
		"rl.name",
		"usr.username",
	}
}
func (UserSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = UserSearchHelper{}

func NewUserSearchHelper() UserSearchHelper {
	return UserSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: nil}}
}
