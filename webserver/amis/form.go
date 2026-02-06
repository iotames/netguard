package amis

type Form struct {
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	SubmitText string      `json:"submitText"`
	Mode       string      `json:"mode"`
	Api        string      `json:"api"`
	Body       []*FormItem `json:"body"`
}

func NewForm(apiurl string) *Form {
	ele := &Form{
		Type:       "form",
		SubmitText: "提交",
		Mode:       "horizontal",
		Api:        apiurl,
		Body:       make([]*FormItem, 0),
	}
	return ele
}

func (f *Form) SetTitle(title string) *Form {
	f.Title = title
	return f
}

func (f *Form) SetSubmitText(text string) *Form {
	f.SubmitText = text
	return f
}

func (f *Form) AddItem(item *FormItem) *Form {
	f.Body = append(f.Body, item)
	return f
}
