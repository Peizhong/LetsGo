package proxy

import "strings"

func parseRequestLine(line string) (method, requestURI, roomId, serviceName string) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	method = line[:s1]
	s2 += s1 + 1
	requestURI = line[s1+1 : s2]
	parseParams := func(param string) (target string) {
		if targetStart := strings.Index(requestURI, param); targetStart >= 0 {
			targetStringGreedy := requestURI[targetStart:]
			if targetEnd := strings.Index(targetStringGreedy, "&"); targetEnd > 0 {
				target = targetStringGreedy[:targetEnd]
			} else {
				target = targetStringGreedy
			}
		}
		lt, lp := len(target), len(param)+1
		if lt > lp {
			return target[lp:]
		}
		return
	}
	roomId = parseParams(RoomIdParam)
	// todo: 从环境变量读取
	serviceName = parseParams(ServiceNameParam)
	return
}
