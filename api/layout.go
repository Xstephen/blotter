package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/register"
	"go.mongodb.org/mongo-driver/bson"
)

type LayoutResponse struct {
	Menus    []Menu `json:"menus"`
	View     int    `json:"view"`
	Beian    string `json:"beian"`
	BlogName string `json:"blog_name"`
}

func Layout(context *register.HandleContext) (err error) {
	res := LayoutResponse{}

	res.Menus, err = getMenus()
	if err != nil {
		return
	}

	m := make(VariablesResponse)
	if m, err = getVariables("beian", "view", "blog_name"); err != nil {
		return
	}
	res.View = int(m["view"].(float64))
	res.Beian = m["beian"].(string)
	res.BlogName = m["blog_name"].(string)

	go func() {
		mongo.Update(
			"blotter", "variables", bson.M{"key": "view"},
			bson.M{"$inc": bson.M{"value": 1}}, nil,
		)
	}()

	context.ReturnJSON(res)
	return
}
