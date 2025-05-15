package pgxadapter

func genRule(ptype string, rule []string) []string {
	result := make([]string, 1+len(rule))
	result[0] = ptype
	copy(result[1:], rule)
	return result
}
