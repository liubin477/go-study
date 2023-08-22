package gee

import (
  "fmt"
  "strings"
)

type node struct {
  pattern  string  // 待匹配路由，如 /p/:lang
  part     string  // 路由中的一部分，如 :lang
  children []*node // 子节点，如 [doc, tutorial, intro]
  isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为 true
}

// toString()
func (n *node) String() string {
  return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 查找第一个匹配成功的子节点
func (n *node) matchChild(part string) *node {
  for _, child := range n.children {
    if child.part == part || child.isWild {
      return child
    }
  }
  return nil
}

// 查找所有匹配成功的子节点
func (n *node) matchChildren(part string) []*node {
  nodes := make([]*node, 0)
  for _, child := range n.children {
    if child.part == part || child.isWild {
      nodes = append(nodes, child)
    }
  }
  return nodes
}

// 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
  if len(parts) == height {
    // 只有在叶子节点才赋值pattern属性，用来后续判断是否匹配成功
    n.pattern = pattern
    return
  }

  part := parts[height]
  child := n.matchChild(part)
  if child == nil {
    // 没有匹配到则新建一个
    child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
    n.children = append(n.children, child)
  }
  child.insert(pattern, parts, height+1)
}

// 匹配节点
func (n *node) search(parts []string, height int) *node {
  // 匹配到了`*`或叶子节点就退出
  if len(parts) == height || strings.HasPrefix(n.part, "*") {
    if n.pattern == "" {
      return nil
    }
    return n
  }

  part := parts[height]
  children := n.matchChildren(part)

  for _, child := range children {
    result := child.search(parts, height+1)
    if result != nil {
      return result
    }
  }

  return nil
}

func (n *node) travel(list *([]*node)) {
  if n.pattern != "" {
    *list = append(*list, n)
  }
  for _, child := range n.children {
    child.travel(list)
  }
}
