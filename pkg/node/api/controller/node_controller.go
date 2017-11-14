// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/goodrain/rainbond/pkg/node/api/model"

	httputil "github.com/goodrain/rainbond/pkg/util/http"
)

//NewNode 创建一个节点
func NewNode(w http.ResponseWriter, r *http.Request) {
	var node model.APIHostNode
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &node, nil); !ok {
		return
	}
	if err := nodeService.AddNode(&node); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//NewMultipleNode 多节点添加操作
func NewMultipleNode(w http.ResponseWriter, r *http.Request) {
	var nodes []model.APIHostNode
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &nodes, nil); !ok {
		return
	}
	var successnodes []model.APIHostNode
	for _, node := range nodes {
		if err := nodeService.AddNode(&node); err != nil {
			continue
		}
		successnodes = append(successnodes, node)
	}
	httputil.ReturnSuccess(r, w, successnodes)
}

//GetNodes 获取全部节点
func GetNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := nodeService.GetAllNode()
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nodes)
}

//GetRuleNodes 获取分角色节点
func GetRuleNodes(w http.ResponseWriter, r *http.Request) {
	rule := chi.URLParam(r, "rule")
	if rule != "compute" && rule != "manage" && rule != "storage" {
		httputil.ReturnError(r, w, 400, rule+" rule is not define")
		return
	}
	nodes, err := nodeService.GetAllNode()
	if err != nil {
		err.Handle(r, w)
		return
	}
	var masternodes []*model.HostNode
	for _, node := range nodes {
		if node.Role.HasRule(rule) {
			masternodes = append(masternodes, node)
		}
	}
	httputil.ReturnSuccess(r, w, masternodes)
}

//DeleteRainbondNode 节点删除
func DeleteRainbondNode(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "node_id")
	err := nodeService.DeleteNode(nodeID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//临时存在
func outJSON(w http.ResponseWriter, data interface{}) {
	outJSONWithCode(w, http.StatusOK, data)
}
func outRespSuccess(w http.ResponseWriter, bean interface{}, data []interface{}) {
	outRespDetails(w, 200, "success", "成功", bean, data)
	//m:=model.ResponseBody{}
	//m.Code=200
	//m.Msg="success"
	//m.MsgCN="成功"
	//m.Body.List=data
}
func outRespDetails(w http.ResponseWriter, code int, msg, msgcn string, bean interface{}, data []interface{}) {
	w.Header().Set("Content-Type", "application/json")
	m := model.ResponseBody{}
	m.Code = code
	m.Msg = msg
	m.MsgCN = msgcn
	m.Body.List = data
	m.Body.Bean = bean

	s := ""
	b, err := json.Marshal(m)

	if err != nil {
		s = `{"error":"json.Marshal error"}`
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		s = string(b)
		w.WriteHeader(code)
	}
	fmt.Fprint(w, s)
}
func outJSONWithCode(w http.ResponseWriter, httpCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	s := ""
	b, err := json.Marshal(data)
	fmt.Println(string(b))
	if err != nil {
		s = `{"error":"json.Marshal error"}`
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		s = string(b)
		w.WriteHeader(httpCode)
	}
	fmt.Fprint(w, s)
}