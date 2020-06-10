package utils_test

import (
	"kube-review/search"
	"kube-review/utils"
	"testing"
)

func TestStuff(t *testing.T) {
	schema, _ := utils.Load("../search/queryschema.json")
	actual := map[string]search.QueryData{}
	err := utils.LoadJSON("../querylist.json", &actual, schema)
	t.Errorf("%v", actual)
	if err != nil {
		t.Errorf(err.Error())
	}
}
