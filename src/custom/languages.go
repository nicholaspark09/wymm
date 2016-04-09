package custom

type Language struct{
	Name string
	Locale string
	all map[string]map[string]string
}

func (l *Language) Init(locale string){
	l.Locale = locale
	l.all = make(map[string]map[string]string)
	l.all["en_US"] = make(map[string]string)
	l.all["es_ES"] = make(map[string]string)

}