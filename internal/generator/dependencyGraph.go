package generator


type ContentOp int

const (
	OpCreate ContentOp = iota
	OpUpdate
	OpDelete
	OpRename
)

type ContentChange struct {
	Op      ContentOp
	Path  	string
	OldPath string
}

type Content struct {
	Path    string
	Template  string
	IncludeIn []string
}

type DependencyGraph struct {
	Templates map[string]map[string]bool
	Content   map[string]*Content
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Templates: make(map[string]map[string]bool),
		Content:   make(map[string]*Content),
	}
}

func (g *DependencyGraph) AddContent(c *Content) {
	g.Content[c.Path] = c

	if g.Templates[c.Template] == nil {
		g.Templates[c.Template] = make(map[string]bool)
	}
	g.Templates[c.Template][c.Path] = true
}

func (g *DependencyGraph) RemoveContent(source string) {
	content := g.Content[source]
	if content == nil {
		return
	}

	if g.Templates[content.Template] != nil {
		delete(g.Templates[content.Template], source)
	}

	delete(g.Content, source)
}

func (g *DependencyGraph) OnTemplateChange(template string) []string {
	affectedSources := []string{}
	for source := range g.Templates[template] {
		affectedSources = append(affectedSources, source)
	}
	return affectedSources
}
