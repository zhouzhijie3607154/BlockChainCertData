package controllers

import (
	"BlockChainCertDataPorject/models"
	"BlockChainCertDataPorject/utils_BCCDP"
	"github.com/astaxie/beego"
	"time"
)

type SendSmsController struct {
	beego.Controller
}

//发送短信验证码的功能
func (s *SendSmsController) Post(){
	var smsLogin models.SmsLogin
	err := s.ParseForm(&smsLogin)
	if err != nil{
		s.Ctx.WriteString("发送验证码数据解析失败，请重试")
		return
	}
	phone := smsLogin.Phone
	code := utils_BCCDP.GenRandCode(6)//返回一个六位的随机数
	result, err := utils_BCCDP.SendSms(phone,code,utils_BCCDP.SMS_TLP_REGISTER)
	if err != nil{
		s.Ctx.WriteString("发送验证码失败，请重试")
		return
	}
	if len(result.BizId) == 0{
		s.Ctx.WriteString(result.Message)
		return
	}
	//验证码发送成功
	smsRecord := models.SmsRecord{
		BizId:     result.BizId,
		Phone:     phone,
		Code:      code,
		Status:    result.Code,
		Message:   result.Message,
		TimeStamp: time.Now().Unix(),
	}
	_, err = smsRecord.SaveSmsRecord()
	if err != nil{
		s.Ctx.WriteString("抱歉获取验证码失败，请重试")
		return
	}
	//保存成功
	s.Data["Phone"] = smsLogin.Phone
	s.Data["BizId"] = smsRecord.BizId
	//验证码登录
	s.TplName = "login_sms_second.html"
}
