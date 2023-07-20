package webhook

import (
	"github.com/go-resty/resty/v2"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/cetus/xnet"

	"github.com/gotomicro/cetus/xcheck/model/dto"
)

func Webhook(fw dto.Webhook, attach dto.AttachInfo) {
	client := resty.New()
	r := client.R().SetHeader("Content-Type", "application/json")
	for k, v := range fw.Headers {
		r.SetHeader(k, v)
	}
	if _, ok := fw.Body["mode"]; ok {
		ip, err := xnet.GetOutBoundIP()
		if err != nil {
			elog.Error("forward", elog.FieldEvent("GetOutBoundIP"), l.E(err), l.S("msg", "get outbound ip error"))
			return
		}
		fw.Body["addr"] = ip + ":9003"
		elog.Info("xcheck", l.A("size", attach.CurrentAbs), l.A("attach", attach), l.S("ip", ip), l.A("body", fw.Body))
	}
	r.SetBody(fw.Body)
	resp, err := r.Post(fw.Url)
	if err != nil {
		elog.Error("forward", elog.FieldEvent("post"), l.E(err), l.A("fw", fw), l.S("resp", resp.String()))
		return
	}
}
