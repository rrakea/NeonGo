package parser

type ast struct {
	children *[]*ast
	val *ast_node
}

type ast_node struct {

}

func Parse() *ast{
	return &ast{}
}

func (t *ast) Get_child(name string) *ast {
	return &ast{}
}

func (t *ast) Get_children() *[]*ast {
	return &[]*ast{}
}

func (t *ast) Get_val() *ast_node {
	return t.val
}