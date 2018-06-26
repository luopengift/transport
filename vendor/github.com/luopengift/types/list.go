package types

type List []interface{}

func (list *List) Len() int {
	return len(*list)
}
// L.append(object) -- append object to end
func (list *List) Append(v interface{}) {
	*list = append(*list, v)
}

func (list *List) Insert(idx int, v interface{}) {
	rest := (*list)[idx:]
	*list = append((*list)[0:idx], v)
	*list = append(*list, rest...)
}

// 指定下标，删除指定的元素，如果删除一个不存在的元素会报错，默认Pop()删除最后一个元素
func (list *List) Pop(idx ...int) interface{} {
	var index int
	if len(idx) == 0 {
		index = list.Len() - 1
	}else{
		index = idx[0]
	}
	pop := (*list)[index]
	*list = append((*list)[:index], (*list)[index+1:]...)
	return pop
}

// 根据value删除元素
func (list *List) Remove(v interface{}) {
	for idx, value := range *list {
		if value == v {
			*list = append((*list)[:idx], (*list)[idx+1:]...)
		}
	}
}

func (list *List) Index(v interface{}) int {
	for idx, value := range *list {
		if value == v {
			return idx
		}
	}
	return -1
}

// L.count(value) -> integer -- return number of occurrences of value
func (list *List) Count(v interface{}) int {
	cnt := 0
	for _, value := range *list {
		if value == v {
			cnt += 1
		}
	}
	return cnt
}

func (list *List) Contains(v interface{}) bool {
	return list.Count(v) != 0
}

