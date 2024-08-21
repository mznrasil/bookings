package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	xe := e[field]
	if len(xe) > 0 {
		return xe[0]
	}
	return ""
}
