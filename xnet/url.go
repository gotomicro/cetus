package xnet

import (
	"fmt"
	"net/url"
)

func SetQueryValues(original string, adds map[string]string, dels []string) (string, error) {
	u, err := url.Parse(original)
	if err != nil {
		return "", err
	}
	vs := url.Values{}
	vs = u.Query()
	if adds != nil {
		for k, v := range adds {
			if vs.Has(k) {
				vs.Set(k, v)
			} else {
				vs.Add(k, v)
			}
		}
	}
	if dels != nil {
		for _, k := range dels {
			vs.Del(k)
		}
	}
	return fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, vs.Encode()), nil
}
