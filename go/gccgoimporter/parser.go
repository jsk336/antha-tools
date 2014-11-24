// antha-tools/go/gccgoimporter/parser.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


package gccgoimporter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antha-lang/antha/token"
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/antha-lang/antha-tools/antha/exact"
	"github.com/antha-lang/antha-tools/antha/types"
)

type parser struct {
	scanner scanner.Scanner
	tok     rune                      // current token
	lit     string                    // literal string; only valid for Ident, Int, String tokens
	pkgpath string                    // package path of imported package
	pkgname string                    // name of imported package
	pkg     *types.Package            // reference to imported package
	imports map[string]*types.Package // package path -> package object
	typeMap map[int]types.Type        // type number -> type
}

func (p *parser) init(filename string, src io.Reader, imports map[string]*types.Package) {
	p.scanner.Init(src)
	p.scanner.Error = func(_ *scanner.Scanner, msg string) { p.error(msg) }
	p.scanner.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanComments | scanner.SkipComments
	p.scanner.Whitespace = 1<<'\t' | 1<<'\n' | 1<<' '
	p.scanner.Filename = filename // for good error messages
	p.next()
	p.imports = imports
	p.typeMap = make(map[int]types.Type)
}

type importError struct {
	pos scanner.Position
	err error
}

func (e importError) Error() string {
	return fmt.Sprintf("import error %s (byte offset = %d): %s", e.pos, e.pos.Offset, e.err)
}

func (p *parser) error(err interface{}) {
	if s, ok := err.(string); ok {
		err = errors.New(s)
	}
	// panic with a runtime.Error if err is not an error
	panic(importError{p.scanner.Pos(), err.(error)})
}

func (p *parser) errorf(format string, args ...interface{}) {
	p.error(fmt.Errorf(format, args...))
}

func (p *parser) expect(tok rune) string {
	lit := p.lit
	if p.tok != tok {
		p.errorf("expected %s, got %s (%s)", scanner.TokenString(tok), scanner.TokenString(p.tok), lit)
	}
	p.next()
	return lit
}

func (p *parser) expectKeyword(keyword string) {
	lit := p.expect(scanner.Ident)
	if lit != keyword {
		p.errorf("expected keyword %s, got %q", keyword, lit)
	}
}

func (p *parser) parseString() string {
	str, err := strconv.Unquote(p.expect(scanner.String))
	if err != nil {
		p.error(err)
	}
	return str
}

// unquotedString     = { unquotedStringChar } .
// unquotedStringChar = <neither a whitespace nor a ';' char> .
func (p *parser) parseUnquotedString() string {
	if p.tok == scanner.EOF {
		p.error("unexpected EOF")
	}
	var buf bytes.Buffer
	buf.WriteString(p.scanner.TokenText())
	// This loop needs to examine each character before deciding whether to consume it. If we see a semicolon,
	// we need to let it be consumed by p.next().
	for ch := p.scanner.Peek(); ch != ';' && ch != scanner.EOF && p.scanner.Whitespace&(1<<uint(ch)) == 0; ch = p.scanner.Peek() {
		buf.WriteRune(ch)
		p.scanner.Next()
	}
	p.next()
	return buf.String()
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
	switch p.tok {
	case scanner.Ident, scanner.Int, scanner.Float, scanner.String, '·':
		p.lit = p.scanner.TokenText()
	default:
		p.lit = ""
	}
}

func (p *parser) parseQualifiedName() (path, name string) {
	return p.parseQualifiedNameStr(p.parseString())
}

func (p *parser) parseUnquotedQualifiedName() (path, name string) {
	return p.parseQualifiedNameStr(p.parseUnquotedString())
}

// qualifiedName = [ ["."] unquotedString "." ] unquotedString .
//
// The above production uses greedy matching.
func (p *parser) parseQualifiedNameStr(unquotedName string) (pkgpath, name string) {
	parts := strings.Split(unquotedName, ".")
	if parts[0] == "" {
		parts = parts[1:]
	}

	switch len(parts) {
	case 0:
		p.errorf("malformed qualified name: %q", unquotedName)
	case 1:
		// unqualified name
		pkgpath = p.pkgpath
		name = parts[0]
	default:
		// qualified name, which may contain periods
		pkgpath = strings.Join(parts[0:len(parts)-1], ".")
		name = parts[len(parts)-1]
	}

	return
}

// getPkg returns the package for a given path. If the package is
// not found but we have a package name, create the package and
// add it to the p.imports map.
//
func (p *parser) getPkg(pkgpath, name string) *types.Package {
	// package unsafe is not in the imports map - handle explicitly
	if pkgpath == "unsafe" {
		return types.Unsafe
	}
	pkg := p.imports[pkgpath]
	if pkg == nil && name != "" {
		pkg = types.NewPackage(pkgpath, name)
		p.imports[pkgpath] = pkg
	}
	return pkg
}

// parseExportedName is like parseQualifiedName, but
// the package path is resolved to an imported *types.Package.
//
// ExportedName = string [string] .
func (p *parser) parseExportedName() (pkg *types.Package, name string) {
	path, name := p.parseQualifiedName()
	var pkgname string
	if p.tok == scanner.String {
		pkgname = p.parseString()
	}
	pkg = p.getPkg(path, pkgname)
	if pkg == nil {
		p.errorf("package %s (path = %q) not found", name, path)
	}
	return
}

// Name = QualifiedName | "?" .
func (p *parser) parseName() string {
	if p.tok == '?' {
		// Anonymous.
		p.next()
		return ""
	}
	// The package path is redundant for us. Don't try to parse it.
	_, name := p.parseUnquotedQualifiedName()
	return name
}

func deref(typ types.Type) types.Type {
	if p, _ := typ.(*types.Pointer); p != nil {
		typ = p.Elem()
	}
	return typ
}

// Field = Name Type [string] .
func (p *parser) parseField(pkg *types.Package) (field *types.Var, tag string) {
	name := p.parseName()
	typ := p.parseType(pkg)
	anon := false
	if name == "" {
		anon = true
		switch typ := deref(typ).(type) {
		case *types.Basic:
			name = typ.Name()
		case *types.Named:
			name = typ.Obj().Name()
		default:
			p.error("anonymous field expected")
		}
	}
	field = types.NewField(token.NoPos, pkg, name, typ, anon)
	if p.tok == scanner.String {
		tag = p.parseString()
	}
	return
}

// Param = Name ["..."] Type .
func (p *parser) parseParam(pkg *types.Package) (param *types.Var, isVariadic bool) {
	name := p.parseName()
	if p.tok == '.' {
		p.next()
		p.expect('.')
		p.expect('.')
		isVariadic = true
	}
	typ := p.parseType(pkg)
	if isVariadic {
		typ = types.NewSlice(typ)
	}
	param = types.NewParam(token.NoPos, pkg, name, typ)
	return
}

// Var = Name Type .
func (p *parser) parseVar(pkg *types.Package) *types.Var {
	name := p.parseName()
	return types.NewVar(token.NoPos, pkg, name, p.parseType(pkg))
}

// ConstValue     = string | "false" | "true" | ["-"] (int ["'"] | FloatOrComplex) .
// FloatOrComplex = float ["i" | ("+"|"-") float "i"] .
func (p *parser) parseConstValue() (val exact.Value, typ types.Type) {
	switch p.tok {
	case scanner.String:
		str := p.parseString()
		val = exact.MakeString(str)
		typ = types.Typ[types.UntypedString]
		return

	case scanner.Ident:
		b := false
		switch p.lit {
		case "false":
		case "true":
			b = true

		default:
			p.errorf("expected const value, got %s (%q)", scanner.TokenString(p.tok), p.lit)
		}

		p.next()
		val = exact.MakeBool(b)
		typ = types.Typ[types.UntypedBool]
		return
	}

	sign := ""
	if p.tok == '-' {
		p.next()
		sign = "-"
	}

	switch p.tok {
	case scanner.Int:
		val = exact.MakeFromLiteral(sign+p.lit, token.INT)
		if val == nil {
			p.error("could not parse integer literal")
		}

		p.next()
		if p.tok == '\'' {
			p.next()
			typ = types.Typ[types.UntypedRune]
		} else {
			typ = types.Typ[types.UntypedInt]
		}

	case scanner.Float:
		re := sign + p.lit
		p.next()

		var im string
		switch p.tok {
		case '+':
			p.next()
			im = p.expect(scanner.Float)

		case '-':
			p.next()
			im = "-" + p.expect(scanner.Float)

		case scanner.Ident:
			// re is in fact the imaginary component. Expect "i" below.
			im = re
			re = "0"

		default:
			val = exact.MakeFromLiteral(re, token.FLOAT)
			if val == nil {
				p.error("could not parse float literal")
			}
			typ = types.Typ[types.UntypedFloat]
			return
		}

		p.expectKeyword("i")
		reval := exact.MakeFromLiteral(re, token.FLOAT)
		if reval == nil {
			p.error("could not parse real component of complex literal")
		}
		imval := exact.MakeFromLiteral(im+"i", token.IMAG)
		if imval == nil {
			p.error("could not parse imag component of complex literal")
		}
		val = exact.BinaryOp(reval, token.ADD, imval)
		typ = types.Typ[types.UntypedComplex]

	default:
		p.errorf("expected const value, got %s (%q)", scanner.TokenString(p.tok), p.lit)
	}

	return
}

// Const = Name [Type] "=" ConstValue .
func (p *parser) parseConst(pkg *types.Package) *types.Const {
	name := p.parseName()
	var typ types.Type
	if p.tok == '<' {
		typ = p.parseType(pkg)
	}
	p.expect('=')
	val, vtyp := p.parseConstValue()
	if typ == nil {
		typ = vtyp
	}
	return types.NewConst(token.NoPos, pkg, name, typ, val)
}

// TypeName = ExportedName .
func (p *parser) parseTypeName() *types.TypeName {
	pkg, name := p.parseExportedName()
	scope := pkg.Scope()
	if obj := scope.Lookup(name); obj != nil {
		return obj.(*types.TypeName)
	}
	obj := types.NewTypeName(token.NoPos, pkg, name, nil)
	// a named type may be referred to before the underlying type
	// is known - set it up
	types.NewNamed(obj, nil, nil)
	scope.Insert(obj)
	return obj
}

// NamedType = TypeName Type { Method } .
// Method    = "func" "(" Param ")" Name ParamList ResultList ";" .
func (p *parser) parseNamedType(n int) types.Type {
	obj := p.parseTypeName()

	pkg := obj.Pkg()
	typ := obj.Type()
	p.typeMap[n] = typ

	nt, ok := typ.(*types.Named)
	if !ok {
		// This can happen for unsafe.Pointer, which is a TypeName holding a Basic type.
		pt := p.parseType(pkg)
		if pt != typ {
			p.error("unexpected underlying type for non-named TypeName")
		}
		return typ
	}

	underlying := p.parseType(pkg)
	if nt.Underlying() == nil {
		nt.SetUnderlying(underlying.Underlying())
	}

	for p.tok == scanner.Ident {
		// collect associated methods
		p.expectKeyword("func")
		p.expect('(')
		receiver, _ := p.parseParam(pkg)
		p.expect(')')
		name := p.parseName()
		params, isVariadic := p.parseParamList(pkg)
		results := p.parseResultList(pkg)
		p.expect(';')

		sig := types.NewSignature(pkg.Scope(), receiver, params, results, isVariadic)
		nt.AddMethod(types.NewFunc(token.NoPos, pkg, name, sig))
	}

	return nt
}

// ArrayOrSliceType = "[" [ int ] "]" Type .
func (p *parser) parseArrayOrSliceType(pkg *types.Package) types.Type {
	p.expect('[')
	if p.tok == ']' {
		p.next()
		return types.NewSlice(p.parseType(pkg))
	}

	lit := p.expect(scanner.Int)
	n, err := strconv.ParseInt(lit, 10, 0)
	if err != nil {
		p.error(err)
	}
	p.expect(']')
	return types.NewArray(p.parseType(pkg), n)
}

// MapType = "map" "[" Type "]" Type .
func (p *parser) parseMapType(pkg *types.Package) types.Type {
	p.expectKeyword("map")
	p.expect('[')
	key := p.parseType(pkg)
	p.expect(']')
	elem := p.parseType(pkg)
	return types.NewMap(key, elem)
}

// ChanType = "chan" ["<-" | "-<"] Type .
func (p *parser) parseChanType(pkg *types.Package) types.Type {
	p.expectKeyword("chan")
	dir := types.SendRecv
	switch p.tok {
	case '-':
		p.next()
		p.expect('<')
		dir = types.SendOnly

	case '<':
		// don't consume '<' if it belongs to Type
		if p.scanner.Peek() == '-' {
			p.next()
			p.expect('-')
			dir = types.RecvOnly
		}
	}

	return types.NewChan(dir, p.parseType(pkg))
}

// StructType = "struct" "{" { Field } "}" .
func (p *parser) parseStructType(pkg *types.Package) types.Type {
	p.expectKeyword("struct")

	var fields []*types.Var
	var tags []string

	p.expect('{')
	for p.tok != '}' && p.tok != scanner.EOF {
		field, tag := p.parseField(pkg)
		p.expect(';')
		fields = append(fields, field)
		tags = append(tags, tag)
	}
	p.expect('}')

	return types.NewStruct(fields, tags)
}

// ParamList = "(" [ { Parameter "," } Parameter ] ")" .
func (p *parser) parseParamList(pkg *types.Package) (*types.Tuple, bool) {
	var list []*types.Var
	isVariadic := false

	p.expect('(')
	for p.tok != ')' && p.tok != scanner.EOF {
		if len(list) > 0 {
			p.expect(',')
		}
		par, variadic := p.parseParam(pkg)
		list = append(list, par)
		if variadic {
			if isVariadic {
				p.error("... not on final argument")
			}
			isVariadic = true
		}
	}
	p.expect(')')

	return types.NewTuple(list...), isVariadic
}

// ResultList = Type | ParamList .
func (p *parser) parseResultList(pkg *types.Package) *types.Tuple {
	switch p.tok {
	case '<':
		return types.NewTuple(types.NewParam(token.NoPos, pkg, "", p.parseType(pkg)))

	case '(':
		params, _ := p.parseParamList(pkg)
		return params

	default:
		return nil
	}
}

// FunctionType = ParamList ResultList .
func (p *parser) parseFunctionType(pkg *types.Package) *types.Signature {
	params, isVariadic := p.parseParamList(pkg)
	results := p.parseResultList(pkg)
	return types.NewSignature(pkg.Scope(), nil, params, results, isVariadic)
}

// Func = Name FunctionType .
func (p *parser) parseFunc(pkg *types.Package) *types.Func {
	name := p.parseName()
	if strings.ContainsRune(name, '$') {
		// This is a Type$equal or Type$hash function, which we don't want to parse,
		// except for the types.
		p.discardDirectiveWhileParsingTypes(pkg)
		return nil
	}
	return types.NewFunc(token.NoPos, pkg, name, p.parseFunctionType(pkg))
}

// InterfaceType = "interface" "{" { ("?" Type | Func) ";" } "}" .
func (p *parser) parseInterfaceType(pkg *types.Package) types.Type {
	p.expectKeyword("interface")

	var methods []*types.Func
	var typs []*types.Named

	p.expect('{')
	for p.tok != '}' && p.tok != scanner.EOF {
		if p.tok == '?' {
			p.next()
			typs = append(typs, p.parseType(pkg).(*types.Named))
		} else {
			method := p.parseFunc(pkg)
			methods = append(methods, method)
		}
		p.expect(';')
	}
	p.expect('}')

	return types.NewInterface(methods, typs)
}

// PointerType = "*" ("any" | Type) .
func (p *parser) parsePointerType(pkg *types.Package) types.Type {
	p.expect('*')
	if p.tok == scanner.Ident {
		p.expectKeyword("any")
		return types.Typ[types.UnsafePointer]
	}
	return types.NewPointer(p.parseType(pkg))
}

// TypeDefinition = NamedType | MapType | ChanType | StructType | InterfaceType | PointerType | ArrayOrSliceType | FunctionType .
func (p *parser) parseTypeDefinition(pkg *types.Package, n int) types.Type {
	var t types.Type
	switch p.tok {
	case scanner.String:
		t = p.parseNamedType(n)

	case scanner.Ident:
		switch p.lit {
		case "map":
			t = p.parseMapType(pkg)

		case "chan":
			t = p.parseChanType(pkg)

		case "struct":
			t = p.parseStructType(pkg)

		case "interface":
			t = p.parseInterfaceType(pkg)
		}

	case '*':
		t = p.parsePointerType(pkg)

	case '[':
		t = p.parseArrayOrSliceType(pkg)

	case '(':
		t = p.parseFunctionType(pkg)
	}

	p.typeMap[n] = t
	return t
}

const (
	// From gofrontend/antha/export.h
	// Note that these values are negative in the gofrontend and have been made positive
	// in the gccgoimporter.
	gccgoBuiltinINT8       = 1
	gccgoBuiltinINT16      = 2
	gccgoBuiltinINT32      = 3
	gccgoBuiltinINT64      = 4
	gccgoBuiltinUINT8      = 5
	gccgoBuiltinUINT16     = 6
	gccgoBuiltinUINT32     = 7
	gccgoBuiltinUINT64     = 8
	gccgoBuiltinFLOAT32    = 9
	gccgoBuiltinFLOAT64    = 10
	gccgoBuiltinINT        = 11
	gccgoBuiltinUINT       = 12
	gccgoBuiltinUINTPTR    = 13
	gccgoBuiltinBOOL       = 15
	gccgoBuiltinSTRING     = 16
	gccgoBuiltinCOMPLEX64  = 17
	gccgoBuiltinCOMPLEX128 = 18
	gccgoBuiltinERROR      = 19
	gccgoBuiltinBYTE       = 20
	gccgoBuiltinRUNE       = 21
)

func lookupBuiltinType(typ int) types.Type {
	return [...]types.Type{
		gccgoBuiltinINT8:       types.Typ[types.Int8],
		gccgoBuiltinINT16:      types.Typ[types.Int16],
		gccgoBuiltinINT32:      types.Typ[types.Int32],
		gccgoBuiltinINT64:      types.Typ[types.Int64],
		gccgoBuiltinUINT8:      types.Typ[types.Uint8],
		gccgoBuiltinUINT16:     types.Typ[types.Uint16],
		gccgoBuiltinUINT32:     types.Typ[types.Uint32],
		gccgoBuiltinUINT64:     types.Typ[types.Uint64],
		gccgoBuiltinFLOAT32:    types.Typ[types.Float32],
		gccgoBuiltinFLOAT64:    types.Typ[types.Float64],
		gccgoBuiltinINT:        types.Typ[types.Int],
		gccgoBuiltinUINT:       types.Typ[types.Uint],
		gccgoBuiltinUINTPTR:    types.Typ[types.Uintptr],
		gccgoBuiltinBOOL:       types.Typ[types.Bool],
		gccgoBuiltinSTRING:     types.Typ[types.String],
		gccgoBuiltinCOMPLEX64:  types.Typ[types.Complex64],
		gccgoBuiltinCOMPLEX128: types.Typ[types.Complex128],
		gccgoBuiltinERROR:      types.Universe.Lookup("error").Type(),
		gccgoBuiltinBYTE:       types.Typ[types.Byte],
		gccgoBuiltinRUNE:       types.Typ[types.Rune],
	}[typ]
}

// Type = "<" "type" ( "-" int | int [ TypeDefinition ] ) ">" .
func (p *parser) parseType(pkg *types.Package) (t types.Type) {
	p.expect('<')
	p.expectKeyword("type")

	switch p.tok {
	case scanner.Int:
		n, err := strconv.ParseInt(p.lit, 10, 0)
		if err != nil {
			p.error(err)
		}
		p.next()

		if p.tok == '>' {
			t = p.typeMap[int(n)]
		} else {
			t = p.parseTypeDefinition(pkg, int(n))
		}

	case '-':
		p.next()
		lit := p.expect(scanner.Int)
		n, err := strconv.ParseInt(lit, 10, 0)
		if err != nil {
			p.error(err)
		}
		t = lookupBuiltinType(int(n))

	default:
		p.errorf("expected type number, got %s (%q)", scanner.TokenString(p.tok), p.lit)
		return nil
	}

	p.expect('>')
	return
}

// Throw away tokens until we see a ';'. If we see a '<', attempt to parse as a type.
func (p *parser) discardDirectiveWhileParsingTypes(pkg *types.Package) {
	for {
		switch p.tok {
		case ';':
			return
		case '<':
			p.parseType(p.pkg)
		case scanner.EOF:
			p.error("unexpected EOF")
		default:
			p.next()
		}
	}
}

// Create the package if we have parsed both the package path and package name.
func (p *parser) maybeCreatePackage() {
	if p.pkgname != "" && p.pkgpath != "" {
		p.pkg = p.getPkg(p.pkgpath, p.pkgname)
	}
}

// Directive = ("v1" | "priority" | "init" | "checksum") { <any token> } ";" |
//             "package" unquotedString ";" |
//             "pkgpath" unquotedString ";" |
//             "import" unquotedString unquotedString string ";" |
//             "func" Func ";" |
//             "type" Type ";" |
//             "var" Var ";" |
//             "const" Const ";" .
func (p *parser) parseDirective() {
	if p.tok != scanner.Ident {
		p.expect(scanner.Ident)
	}

	switch p.lit {
	case "v1", "priority", "init":
		// We can't parse these yet.
		p.discardDirectiveWhileParsingTypes(p.pkg)
		p.next()

	case "package":
		p.next()
		p.pkgname = p.parseUnquotedString()
		p.maybeCreatePackage()
		p.expect(';')

	case "pkgpath":
		p.next()
		p.pkgpath = p.parseUnquotedString()
		p.maybeCreatePackage()
		p.expect(';')

	case "import":
		p.next()
		pkgname := p.parseUnquotedString()
		pkgpath := p.parseUnquotedString()
		p.getPkg(pkgpath, pkgname)
		p.parseString()
		p.expect(';')

	case "func":
		p.next()
		fun := p.parseFunc(p.pkg)
		if fun != nil {
			p.pkg.Scope().Insert(fun)
		}
		p.expect(';')

	case "type":
		p.next()
		p.parseType(p.pkg)
		p.expect(';')

	case "var":
		p.next()
		v := p.parseVar(p.pkg)
		p.pkg.Scope().Insert(v)
		p.expect(';')

	case "const":
		p.next()
		c := p.parseConst(p.pkg)
		p.pkg.Scope().Insert(c)
		p.expect(';')

	case "checksum":
		// Don't let the scanner try to parse the checksum as a number.
		p.scanner.Mode &^= scanner.ScanInts | scanner.ScanFloats
		p.next()
		p.parseUnquotedString()
		p.expect(';')
		p.scanner.Mode |= scanner.ScanInts | scanner.ScanFloats

	default:
		p.errorf("unexpected identifier: %q", p.lit)
	}
}

// Package = { Directive } .
func (p *parser) parsePackage() *types.Package {
	for p.tok != scanner.EOF {
		p.parseDirective()
	}
	p.pkg.MarkComplete()
	return p.pkg
}