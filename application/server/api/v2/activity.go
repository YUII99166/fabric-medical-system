package v2

import (
	"application/pkg/app"
	"application/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllianceActivities 获取联盟活动列表
func GetAllianceActivities(c *gin.Context) {
	appG := app.Gin{C: c}

	type QueryParams struct {
		Page         int    `json:"page" form:"page"`
		PageSize     int    `json:"pageSize" form:"pageSize"`
		Organization string `json:"organization" form:"organization"`
		ActivityType string `json:"activityType" form:"activityType"`
		StartDate    string `json:"startDate" form:"startDate"`
		EndDate      string `json:"endDate" form:"endDate"`
	}

	params := QueryParams{
		Page:     1,
		PageSize: 20,
	}

	if err := c.ShouldBind(&params); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数解析错误: %s", err.Error()))
		return
	}

	// 调用service层获取活动数据
	activities, total, err := service.GetAllianceActivities(
		params.Page,
		params.PageSize,
		params.Organization,
		params.ActivityType,
		params.StartDate,
		params.EndDate,
	)

	if err != nil {
		fmt.Printf("获取联盟活动失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("获取活动失败: %v", err))
		return
	}

	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"activities": activities,
		"total":      total,
	})
}

// GetActivityStatistics 获取活动统计数据
func GetActivityStatistics(c *gin.Context) {
	appG := app.Gin{C: c}

	stats, err := service.GetActivityStatistics()
	if err != nil {
		fmt.Printf("获取活动统计失败: %v\n", err)
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("获取统计失败: %v", err))
		return
	}

	appG.Response(http.StatusOK, "成功", stats)
}
