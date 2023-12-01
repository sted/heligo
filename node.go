package heligo

type node struct {
	text       string
	children   []*node
	childColon *node
	childStar  *node
	handler    Handler
	param      string
}

func (n *node) nextNode(s string) *node {
	slen := len(s)
	if slen == 1 {
		if s[0] == COLON {
			if n.childColon == nil {
				n.childColon = &node{text: string(COLON)}
			}
			return n.childColon
		} else if s[0] == STAR {
			if n.childStar == nil {
				n.childStar = &node{text: string(STAR)}
			}
			return n.childStar
		}
	}
	for i := 0; i < len(n.children); i++ {
		minlen := slen
		child := n.children[i]
		clen := len(child.text)
		if clen < slen {
			minlen = clen
		}
		idx := minlen
		for j := 0; j < minlen; j++ {
			if s[j] != child.text[j] {
				idx = j
				break
			}
		}
		switch idx {
		case 0:
			continue
		case minlen:
			if clen < slen {
				return child.nextNode(s[idx:])
			} else if clen == slen {
				return child
			} else {
				n.children[i] = &node{text: s, children: []*node{child}}
				child.text = child.text[idx:]
				return n.children[i]
			}
		default:
			n.children[i] = &node{text: s[:idx], children: []*node{child, {text: s[idx:]}}}
			child.text = child.text[idx:]
			return n.children[i].children[1]
		}
	}
	newNode := &node{text: s}
	n.children = append(n.children, newNode)
	return newNode
}

func (n *node) findNode(s string, offset int, p *params) *node {
	var child *node
	slen := len(s)
	for i := 0; i < len(n.children); i++ {
		child = n.children[i]
		clen := len(child.text)
		if clen > slen || child.text != s[0:clen] {
			continue
		}
		if clen < slen {
			child = child.findNode(s[clen:], offset+clen, p)
			if child == nil {
				// backtrack
				break
			}
		}
		return child
	}
	if n.childColon != nil || n.childStar != nil {
		c := p.count
		if n.childColon != nil {
			child = n.childColon
			p.names[c] = &child.param
			k := 0
			for {
				if k == slen {
					p.valueBeg[c] = uint16(offset)
					p.valueEnd[c] = 0
					p.count++
					return child
				}
				if s[k] == SLASH {
					p.valueBeg[c] = uint16(offset)
					p.valueEnd[c] = uint16(k)
					p.count++
					child = child.findNode(s[k:], offset+k, p)
					if child != nil {
						return child
					} else {
						// backtrack
						p.count--
						break
					}
				}
				k++
			}
		}
		if n.childStar != nil {
			child = n.childStar
			p.names[c] = &child.param
			p.valueBeg[c] = uint16(offset)
			p.valueEnd[c] = 0
			p.count++
			return child
		}
	}
	return nil
}
