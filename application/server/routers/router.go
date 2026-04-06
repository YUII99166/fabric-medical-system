package routers

import (
	v2 "application/api/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRouter 初始化路由信息
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	apiV2 := r.Group("/api/v2")
	{
		apiV2.GET("/hello", v2.Hello)
		apiV2.POST("/createAccountV2", v2.CreateAccountV2)
		apiV2.POST("/queryAccountV2List", v2.QueryAccountV2List)
		apiV2.POST("/queryAccountListFromDB", v2.QueryAccountListFromDB)
		apiV2.POST("/register", v2.Register)
		apiV2.POST("/loginWithPassword", v2.LoginWithPassword)
		apiV2.POST("/getUserInfo", v2.GetUserInfo)
		apiV2.POST("/getUserDetail", v2.GetUserDetail)
		apiV2.POST("/updateUser", v2.UpdateUser)
		apiV2.POST("/deleteUser", v2.DeleteUser)
		apiV2.POST("/restoreUser", v2.RestoreUser)
		apiV2.POST("/batchDeleteUsers", v2.BatchDeleteUsers)
		apiV2.POST("/syncAccountFromBlockchain", v2.SyncAccountFromBlockchain)
		apiV2.POST("/createPrescription", v2.CreatePrescription)
		apiV2.POST("/queryPrescription", v2.QueryPrescriptionList)
		apiV2.POST("/createInsuranceCover", v2.CreateInsuranceCover)
		apiV2.POST("/queryInsuranceCoverList", v2.QueryInsuranceCoverList)
		apiV2.POST("/updateInsuranceCover", v2.UpdateInsuranceCover)
		apiV2.POST("/deleteInsuranceCover", v2.DeleteInsuranceCover)
		apiV2.POST("/createDrugOrder", v2.CreateDrugOrder)
		apiV2.POST("/queryDrugOrderList", v2.QueryDrugOrderList)
		apiV2.POST("/requestAccess", v2.RequestAccess)
		apiV2.POST("/approveAccess", v2.ApproveAccess)
		apiV2.POST("/queryAccessRequests", v2.QueryAccessRequests)
		apiV2.POST("/queryPrescriptionsByPatient", v2.QueryPrescriptionsByPatient)
		apiV2.POST("/addSupplementRecord", v2.AddSupplementRecord)
		apiV2.POST("/querySupplementRecords", v2.QuerySupplementRecords)
		apiV2.POST("/queryFullMedicalHistory", v2.QueryFullMedicalHistory)
		apiV2.GET("/patient/health-profile", v2.GetHealthProfile)
		apiV2.GET("/statistics", v2.GetStatistics)
		apiV2.GET("/recentActivities", v2.GetRecentActivities)
		apiV2.POST("/getAllianceActivities", v2.GetAllianceActivities)
		apiV2.POST("/getActivityStatistics", v2.GetActivityStatistics)
		apiV2.POST("/recordPrescriptionAccess", v2.RecordPrescriptionAccess)
		apiV2.POST("/getMyAccessLogs", v2.GetMyAccessLogs)
		apiV2.GET("/getAccessStatistics", v2.GetAccessStatistics)
	}
	return r
}
