package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"wbl0/internal/db"
	"wbl0/internal/model"
	"wbl0/pkg/cash"
)

func HandleOrder(db db.OrderRepository, cash *cash.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("internal/static/page.html"))
		uid := r.FormValue("uid")
		var order model.Order
		res, ok := cash.Get(uid)
		if ok {
			order = res.(model.Order)
			fmt.Println("Order from cash: ", order)
			err := tmpl.Execute(w, order)
			if err != nil {
				fmt.Errorf("Error executing template: %s", err)
			}
		} else {
			orderFromDb, err := db.GetOrderByUid(uid)
			if err != nil {
				fmt.Errorf("Error getting order by uid: %s", err)
			}
			cash.Set(uid, orderFromDb, 0)
			res, ok = cash.Get(uid)
			if !ok {
				fmt.Errorf("Error getting order from cash: %s", err)
			}
			order = res.(model.Order)
			fmt.Println("Order from db: ", order)
			err = tmpl.Execute(w, order)
			if err != nil {
				fmt.Errorf("Error executing template: %s", err)
			}

		}

	}
}
