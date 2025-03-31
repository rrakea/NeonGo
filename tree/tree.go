package tree


type tree interface {
	Child (string) tree
	Children () []*tree
	Name() string
	Val() any
}

func Finc_direct(t tree, name string) *tree {
	for _, c := range  t.Children() {
		if (*c).Name() == name {
			return c
		}	
	}
	return nil
}

func Finc_all(t tree, name string) []*tree {
	nodes := []*tree{}
	for _, c := range  t.Children() {
			if (*c).Name() == name {
				nodes = append(nodes, c)
			}
			nodes = append(nodes, Finc_all(*c, name)...)
	}
	return nodes
}