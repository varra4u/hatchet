/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * stats_handler.go
 */

package hatchet

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
)

// StatsHandler responds to API calls
func StatsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	/** APIs
	 * /hatchets/{hatchet}/stats/audit
	 * /hatchets/{hatchet}/stats/slowops
	 */
	hatchetName := params.ByName("hatchet")
	attr := params.ByName("attr")
	dbase, err := GetDatabase(hatchetName)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
		return
	}
	defer dbase.Close()
	if dbase.GetVerbose() {
		log.Println("StatsHandler", r.URL.Path, hatchetName, attr)
	}
	info := dbase.GetHatchetInfo()
	summary := GetHatchetSummary(info)
	download := r.URL.Query().Get("download")

	if attr == "audit" {
		data, err := dbase.GetAuditData()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		templ, err := GetAuditTablesTemplate()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		doc := map[string]interface{}{"Hatchet": hatchetName, "Info": info, "Summary": summary, "Data": data}
		if err = templ.Execute(w, doc); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		return
	} else if attr == "slowops" {
		collscan := false
		var ns string
		var op string
		if r.URL.Query().Get(COLLSCAN) == "true" {
			collscan = true
		}
		if len(r.URL.Query().Get("ns")) > 0 {
			ns = r.URL.Query().Get("ns")
		}
		if len(r.URL.Query().Get("op")) > 0 {
			op = r.URL.Query().Get("op")
		}
		var order, orderBy string
		orderBy = r.URL.Query().Get("orderBy")
		if orderBy == "" {
			orderBy = "avg_ms"
		} else if orderBy == "index" || orderBy == "_index" {
			orderBy = "_index"
		}
		order = r.URL.Query().Get("order")
		if order == "" {
			if orderBy == "op" || orderBy == "ns" {
				order = "ASC"
			} else {
				order = "DESC"
			}
		}
		//ops, err := dbase.GetSlowOps(orderBy, order, collscan)
		ops, err := dbase.GetSlowOpsV2(orderBy, order, collscan, ns, handleOpParamValues(op))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		templ, err := GetStatsTableTemplate(collscan, orderBy, ns, op, download)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		doc := map[string]interface{}{"Hatchet": hatchetName, "Ops": ops, "Summary": summary}
		if err = templ.Execute(w, doc); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": 0, "error": err.Error()})
			return
		}
		return
	}
}

func handleOpParamValues(op string) []string {
	var result []string
	if len(strings.Trim(op, "")) > 0 {
		ops := strings.Split(strings.ReplaceAll(op, " ", ""), ",")
		for i := 0; i < len(ops); i++ {
			opParam := strings.Trim(ops[i], " ")
			if len(opParam) > 0 && !strings.HasPrefix(opParam, "-") {
				result = append(result, strings.ToLower(opParam))
			}
		}
	}
	return result
}
