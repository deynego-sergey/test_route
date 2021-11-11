package router

import (
	"fmt"
	"strings"
	"testing"
)

var routes = []string{
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/leaf7/#`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/leaf7/+`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/leaf7/leaf8`,
	`tp/binance/trading/leaf2/+leaf1/leaf3//leaf4/leaf5+/+/leaf7/#`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/#/leaf5+/+/leaf7/+`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/+leaf5+/+/leaf7/leaf8`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/++/leaf7/#`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/#/leaf7/`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/+/+`,
}
var topic = []string{
	`tp/binance/trading/leaf2/AAleaf1/leaf3/leaf4/leaf5BB/leaf6/leaf7/CCC/DDD`,
	`tp/binance/trading/leaf2/AAA/leaf3/leaf4/leaf5BB/leaf6/leaf7/CCC/DDD`,
	`tp/binance/trading/leaf2/AAleaf1/leaf3/leaf4/BBB/leaf6/leaf7/CCC/DDD`,
	`tp/binance/trading/leaf2/AAleaf1/leaf3/leaf4/leaf5BBB/leaf6/leaf7`,
}

var subscribe = []string{
	`tp/binance/trading/leaf2/#`,
	`tp/binance/trading/leaf2/XXX/#`,
	`tp/binance/trading/leaf2/+/leaf3/leaf4/+/leaf6/leaf7/aaa`,
	`tp/binance/trading/leaf2/AAleaf1/leaf3/leaf4/+ddd/leaf6/leaf7/aaa`,
	`tp/binance/trading/leaf2/AAleaf1/leaf3/leaf4/ccc+/leaf6/leaf7/aaa`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/+/leaf6/leaf7/aaa`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/ccc+/+/leaf7/aaa`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+aaa/leaf7/aaa`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/bbb+/leaf7/aaa`,
	`tp/binance/trading/#leaf2/+leaf1/leaf3/leaf4/leaf5+/bbb+/leaf7/aaa`,
	`tp/binance/trading/++/+leaf1/leaf3/leaf4/leaf5+/bbb+/leaf7/aaa`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+`,
	`tp/binance/trading/+/+/+/+/+/+/+`,
	`tp/binance/trading/leaf2/+leaf1/leaf3/leaf4/leaf5+/+/leaf7/leaf8/aaa`,
}

//
//
//
func Test_NewRouterPattern(t *testing.T) {

	for _, v := range routes {

		if r, e := NewRoutePattern("", v); e != nil {
			fmt.Println(v, e)
			continue
		} else {
			for _, s := range subscribe {
				if e = r.Match(s); e == nil {
					fmt.Println(r.Pattern(), s, true)
				} else if strings.Contains(e.Error(), errNotMatched.Error()) {
					fmt.Println(r.Pattern(), s, false)
				} else {

					fmt.Println(r.Pattern(), s, e)
				}
			}
		}
	}
}
