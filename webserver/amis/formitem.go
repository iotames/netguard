package amis

type FormItem map[string]any

func NewFormItem() *FormItem {
	item := make(FormItem)
	return &item
}

func (f *FormItem) Set(k string, v any) *FormItem {
	(*f)[k] = v // 需要解引用
	return f
}

// item1 := amis.NewFormItem().Set("label", "客户简称").Set("type", "select").Set("name", "inputtpl").Set("value", PoCustomers[0].Code).Set("source", "/api/customer/list")
// item2 := amis.NewFormItem().Set("type", "input-file").Set("name", "inputfile").Set("accept", ".xlsx").Set("label", "上传.xlsx文件").Set("maxSize", 10048576).Set("receiver", "/api/uploadfile")
