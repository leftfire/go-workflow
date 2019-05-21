package service

import (
	"errors"
	"sync"
	"time"

	"github.com/mumushuiding/go-workflow/workflow-engine/flow"

	"github.com/mumushuiding/util"

	"github.com/mumushuiding/go-workflow/workflow-engine/model"
)

var saveLock sync.Mutex

// Procdef 流程定义表
type Procdef struct {
	Name string `json:"name"`
	// 流程定义json字符串
	Resource *flow.Node `json:"resource"`
	// 用户id
	Userid string `jsong:"userid"`
	// 用户所在公司
	Company   string `json:"company"`
	PageSize  int    `json:"pageSize"`
	PageIndex int    `json:"pageIndex"`
}

// GetProcdefByID 根据流程定义id获取流程定义
func GetProcdefByID(id int) (*model.Procdef, error) {
	return model.GetProcdefByID(id)
}

// GetProcdefLatestByNameAndCompany GetProcdefLatestByNameAndCompany
// 根据流程定义名字和公司查询流程定义
func GetProcdefLatestByNameAndCompany(name, company string) (*model.Procdef, error) {
	return model.GetProcdefLatestByNameAndCompany(name, company)
}

// GetResourceByNameAndCompany GetResourceByNameAndCompany
// 获取流程定义配置信息
func GetResourceByNameAndCompany(name, company string) (*flow.Node, int, string, error) {
	prodef, err := GetProcdefLatestByNameAndCompany(name, company)
	if err != nil {
		return nil, 0, "", err
	}
	node := &flow.Node{}
	err = util.Str2Struct(prodef.Resource, node)
	return node, prodef.ID, prodef.Name, err
}

// GetResourceByID GetResourceByID
// 根据id查询流程定义
func GetResourceByID(id int) (*flow.Node, int, error) {
	prodef, err := GetProcdefByID(id)
	if err != nil {
		return nil, 0, err
	}
	node := &flow.Node{}
	err = util.Str2Struct(prodef.Resource, node)
	return node, prodef.ID, err
}

// ExistsProcdefByNameAndCompany if exists
// 查询流程定义是否存在
func ExistsProcdefByNameAndCompany(name, company string) (yes bool, version int, err error) {
	p, err := model.GetProcdefLatestByNameAndCompany(name, company)
	if p == nil {
		return false, 1, err
	}
	version = p.Version + 1
	return true, version, err
}

// SaveProcdef 保存
func (p *Procdef) SaveProcdef() (id int, err error) {
	resource, err := util.ToJSONStr(p.Resource)
	if err != nil {
		return 0, err
	}
	// fmt.Println(resource)
	var procdef = model.Procdef{
		Name:     p.Name,
		Userid:   p.Userid,
		Company:  p.Company,
		Resource: resource,
	}
	return SaveProcdef(&procdef)
}

// SaveProcdef 保存
func SaveProcdef(p *model.Procdef) (id int, err error) {
	// 参数是否为空判定
	if len(p.Userid) == 0 || len(p.Name) == 0 || len(p.Company) == 0 || len(p.Resource) == 0 {
		return 0, errors.New("userid，name,company,resource can not be null  不能为空")
	}
	saveLock.Lock()
	defer saveLock.Unlock()
	// check if exists
	// 判断是否存在
	_, version, err := ExistsProcdefByNameAndCompany(p.Name, p.Company)
	if err != nil {
		return 0, err
	}

	p.Version = version
	p.DeployTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	return p.Save()
}

// FindAllPageAsJSON find by page and  transform result to string
// 分页查询并将结果转换成 json 字符串
func (p *Procdef) FindAllPageAsJSON() (string, error) {
	datas, count, err := p.FindAll()
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, p.PageIndex, p.PageSize)
}

// FindAll FindAll
func (p *Procdef) FindAll() ([]*model.Procdef, int, error) {
	var page = util.Page{}
	page.PageRequest(p.PageIndex, p.PageSize)
	maps := p.getMaps()
	return model.FindProcdefsWithCountAndPaged(page.PageIndex, page.PageSize, maps)
}
func (p *Procdef) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if len(p.Name) > 0 {
		maps["Name"] = p.Name
	}
	if len(p.Company) > 0 {
		maps["company"] = p.Company
	}
	return maps
}

// DelProcdefByID del by id
func DelProcdefByID(id int) error {
	return model.DelProcdefByID(id)
}
