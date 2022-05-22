package net

import "github.com/gin-gonic/gin"

type PageInfo struct {
	Size     int    `form:"size" json:"size"`
	Current  int    `form:"current" json:"current"`
	Criteria string `form:"criteria" json:"criteria"`
}

type PageResult struct {
	PageInfo
	Total   int         `form:"total" json:"total"`
	Records interface{} `form:"records" json:"records"`
}

func NewPageInfo(size, curpage int, criteria string) *PageInfo {
	if size <= 0 {
		size = 10
	}
	if curpage <= 0 {
		curpage = 1
	}
	return &PageInfo{size, curpage, criteria}
}
func NewPageResult(size, total, curpage int, criteria string, data interface{}) *PageResult {
	return &PageResult{*NewPageInfo(size, curpage, criteria), total, data}
}

func GetPageInfo(c *gin.Context) (PageInfo, error) {
	var page PageInfo
	err := c.ShouldBindQuery(&page)
	if err != nil {
		return page, err
	} else {
		if page.Size <= 0 {
			page.Size = 10
		}
		if page.Current <= 0 {
			page.Current = 1
		}
		return page, nil
	}
}
