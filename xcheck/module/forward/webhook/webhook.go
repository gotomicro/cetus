package webhook

import (
	"github.com/go-resty/resty/v2"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/cetus/xcheck/model/utils"

	"github.com/gotomicro/cetus/xcheck/model/dto"
)

func Webhook(fw dto.Webhook, attach dto.AttachInfo) {
	client := resty.New()
	r := client.R().SetHeader("Content-Type", "application/json")
	for k, v := range fw.Headers {
		r.SetHeader(k, v)
	}
	if _, ok := fw.Body["mode"]; ok {
		ip, err := utils.GetOutBoundIP()
		if err != nil {
			elog.Error("forward", elog.FieldEvent("GetOutBoundIP"), elog.FieldErr(err), elog.String("msg", "get outbound ip error"))
			return
		}
		fw.Body["addr"] = ip + ":9003"
		elog.Info("xcheck", elog.Any("size", attach.CurrentAbs), elog.Any("attach", attach), elog.String("ip", ip), elog.Any("body", fw.Body))
	}
	r.SetBody(fw.Body)
	resp, err := r.Post(fw.Url)
	if err != nil {
		elog.Error("forward", elog.FieldEvent("post"), elog.FieldErr(err), elog.Any("fw", fw), elog.String("resp", resp.String()))
		return
	}
}
