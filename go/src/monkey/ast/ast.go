package ast

import (
    "bytes"
    "strings"

    "monkey/token"
)

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) > 0 {
        return p.Statements[0].TokenLiteral()
    } else {
        return "0"
    }
}

func (p *Program) String() string {
    var out bytes.Buffer

    for _, s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

// Statements
type LetStatement struct {
    Token token.Token
    Name  *Identifier
    Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")

    return out.String()
}

func (ls *LetStatement) statementNode() { }

type ReturnStatement struct {
    Token       token.Token
    ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    out.WriteString(";")

    return out.String()
}

func (rs *ReturnStatement) statementNode() { }

type ExpressionStatement struct {
    Token       token.Token
    Expression  Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
    var out bytes.Buffer

    if es.Expression != nil {
        out.WriteString(es.Expression.String())
    }

    return out.String()
}

func (es *ExpressionStatement) statementNode() { }

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}

func (i *Identifier) String() string {
    return i.Value
}

func (i *Identifier) expressionNode() { }

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() { }

type StringLiteral struct {
    Token token.Token
    Value string
}

func (sl *StringLiteral) TokenLiteral() string {
    return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
    return sl.Token.Literal
}

func (sl *StringLiteral) expressionNode() {}

type Boolean struct {
    Token token.Token
    Value bool
}

func (bl *Boolean) TokenLiteral() string {
    return bl.Token.Literal
}

func (bl *Boolean) String() string {
    return bl.Token.Literal
}

func (bl *Boolean) expressionNode() { }

type PrefixExpression struct {
    Token    token.Token
    Operator string
    Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

func (pe *PrefixExpression) expressionNode() { }

type InfixExpression struct {
    Token    token.Token
    Left     Expression
    Operator string
    Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(ie.Left.String())
    out.WriteString(" " + ie.Operator + " ")
    out.WriteString(ie.Right.String())
    out.WriteString(")")

    return out.String()
}

func (ie *InfixExpression) expressionNode() { }

type IfExpression struct {
    Token       token.Token // the if token
    Condition   Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}

func (ie *IfExpression) String() string {
    var out bytes.Buffer

    out.WriteString(ie.TokenLiteral())
    out.WriteString(ie.Condition.String())
    out.WriteString(" ")
    out.WriteString(ie.Consequence.String())
    
    if ie.Alternative != nil {
        out.WriteString("else")
        out.WriteString(ie.Alternative.String())
    }

    return out.String()
}

func (ie *IfExpression) expressionNode() { }

type FunctionLiteral struct {
    Token      token.Token
    Parameters []*Identifier
    Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
    return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range fl.Parameters {
        params = append(params, p.String())
    }

    out.WriteString(fl.TokenLiteral())
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") ")
    out.WriteString(fl.Body.String())
    

    return out.String()
}

func (fl *FunctionLiteral) expressionNode() { }

type CallExpression struct {
    Token     token.Token // the ( token
    Function  Expression
    Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}

func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }
    
    out.WriteString(ce.Function.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}

func (ce *CallExpression) expressionNode() { }

type BlockStatement struct {
    Token token.Token // the { token
    Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
    var out bytes.Buffer

    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

func (bs *BlockStatement) statementNode() {}

type ArrayLiteral struct {
    Token token.Token // the '[' token
    Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
    return al.Token.Literal
}

func (al *ArrayLiteral) String() string {
    var out bytes.Buffer

    elements := []string{}
    for _, e := range al.Elements {
        elements = append(elements, e.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

func (al *ArrayLiteral) expressionNode() {}