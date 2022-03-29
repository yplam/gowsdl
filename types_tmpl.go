// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var typesTmpl = `
{{define "SimpleType"}}
	// SimpleTypeT
	{{$typeName := replaceReservedWords .Name | wrapNS | makePublic}}
	{{if .Doc}} {{.Doc | comment}} {{end}}
	{{if ne .List.ItemType ""}}
		type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
	{{else if ne .Union.MemberTypes ""}}
		type {{$typeName}} string
	{{else if .Union.SimpleType}}
		type {{$typeName}} string
	{{else if .Restriction.Base}}
		type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
    {{else}}
		type {{$typeName}} interface{}
	{{end}}

	{{if .Restriction.Enumeration}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
		{{end}}
	)
	{{end}}
{{end}}

{{define "ComplexContent"}}
	// ComplexContent
	{{$baseType := toGoType .Extension.Base false}}
	{{ if $baseType }}
		{{$baseType}}
	{{end}}

	{{template "Elements" .Extension.Sequence}}
	{{template "Elements" .Extension.Choice}}
	{{template "Elements" .Extension.SequenceChoice}}
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "Attributes"}}
	// Attributes
    {{ $targetNamespace := getNS }}
	{{range .}}
		{{if .Doc}} {{.Doc | comment}} {{end}}
		{{ if ne .Type "" }}
			{{ normalize .Name | makeFieldPublic}} {{toGoType .Type false}} ` + "`" + `xml:"{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ else }}
			{{ normalize .Name | makeFieldPublic}} string ` + "`" + `xml:"{{.Name}},attr,omitempty" json:"{{.Name}},omitempty"` + "`" + `
		{{ end }}
	{{end}}
{{end}}

{{define "SimpleContent"}}
	// SimpleContent
	Value {{toGoType .Extension.Base false}} ` + "`xml:\",chardata\" json:\"-,\"`" + `
	{{template "Attributes" .Extension.Attributes}}
{{end}}

{{define "ComplexTypeInline"}}
	// ComplexTypeInline
	{{replaceReservedWords .Name | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}struct {
	{{with .ComplexType}}
		{{if ne .ComplexContent.Extension.Base ""}}
			{{template "ComplexContent" .ComplexContent}}
		{{else if ne .SimpleContent.Extension.Base ""}}
			{{template "SimpleContent" .SimpleContent}}
		{{else}}
			{{template "Elements" .Sequence}}
			{{template "Elements" .Choice}}
			{{template "Elements" .SequenceChoice}}
			{{template "Elements" .All}}
			{{template "Attributes" .Attributes}}
		{{end}}
	{{end}}
	} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
{{end}}

{{define "Elements"}}
	// ElementsT
	{{range .}}
		{{if ne .Ref ""}}
			{{removeNS .Ref | replaceReservedWords  | makePublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Ref .Nillable }} ` + "`" + `xml:"{{.Ref | removeNS}},omitempty" json:"{{.Ref | removeNS}},omitempty"` + "`" + `
		{{else}}
		{{if not .Type}}
			{{if .SimpleType}}
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .SimpleType.List.ItemType ""}}
					{{ normalize .Name | makeFieldPublic}} []{{toGoType .SimpleType.List.ItemType false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{else}}
					{{ normalize .Name | makeFieldPublic}} {{toGoType .SimpleType.Restriction.Base false}} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + `
				{{end}}
			{{else}}
				{{template "ComplexTypeInline" .}}
			{{end}}
		{{else}}
			{{if .Doc}}{{.Doc | comment}} {{end}}
			{{replaceAttrReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{toGoType .Type .Nillable }} ` + "`" + `xml:"{{.Name}},omitempty" json:"{{.Name}},omitempty"` + "`" + ` {{end}}
		{{end}}
	{{end}}
{{end}}

{{define "Any"}}
	// Any
	{{range .}}
		Items     []string ` + "`" + `xml:",any" json:"items,omitempty"` + "`" + `
	{{end}}
{{end}}

{{range .Schemas}}
	{{ $targetNamespace := setNS .TargetNamespace }}
	// Schema {{$targetNamespace}}
	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}

	{{range .Elements}}
		{{$name := .Name}}
		// Elements {{$targetNamespace}}
		{{$typeName := replaceReservedWords $name | wrapNS | makePublic}}
		{{if not .Type}}
			{{/* ComplexTypeLocal */}}
			{{with .ComplexType}}
				// ComplexTypeLocal $targetNamespace
				type {{$typeName}} struct {
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$name}}\"`" + `
					{{if ne .ComplexContent.Extension.Base ""}}
						{{template "ComplexContent" .ComplexContent}}
					{{else if ne .SimpleContent.Extension.Base ""}}
						{{template "SimpleContent" .SimpleContent}}
					{{else}}
						{{template "Elements" .Sequence}}
						{{template "Any" .Any}}
						{{template "Elements" .Choice}}
						{{template "Elements" .SequenceChoice}}
						{{template "Elements" .All}}
						{{template "Attributes" .Attributes}}
					{{end}}
				}
			{{end}}
			{{/* SimpleTypeLocal */}}
			{{with .SimpleType}}
				// SimpleType
				{{if .Doc}} {{.Doc | comment}} {{end}}
				{{if ne .List.ItemType ""}}
					type {{$typeName}} []{{toGoType .List.ItemType false | removePointerFromType}}
				{{else if ne .Union.MemberTypes ""}}
					type {{$typeName}} string
				{{else if .Union.SimpleType}}
					type {{$typeName}} string
				{{else if .Restriction.Base}}
					type {{$typeName}} {{toGoType .Restriction.Base false | removePointerFromType}}
				{{else}}
					type {{$typeName}} interface{}
				{{end}}

				{{if .Restriction.Enumeration}}
				const (
					{{with .Restriction}}
						{{range .Enumeration}}
							{{if .Doc}} {{.Doc | comment}} {{end}}
							{{$typeName}}{{$value := replaceReservedWords .Value}}{{$value | makePublic}} {{$typeName}} = "{{goString .Value}}" {{end}}
					{{end}}
				)
				{{end}}
			{{end}}
		{{else}}
			{{$type := toGoType .Type .Nillable | removePointerFromType}}
			{{if ne ($typeName) ($type)}}
				type {{$typeName}} {{$type}}
				{{if eq ($type) ("soap.XSDDateTime")}}
					func (xdt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDateTime(xdt).MarshalXML(e, start)
					}

					func (xdt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDateTime)(xdt).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDDate")}}
					func (xd {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDDate(xd).MarshalXML(e, start)
					}

					func (xd *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDDate)(xd).UnmarshalXML(d, start)
					}
				{{else if eq ($type) ("soap.XSDTime")}}
					func (xt {{$typeName}}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
						return soap.XSDTime(xt).MarshalXML(e, start)
					}

					func (xt *{{$typeName}}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
						return (*soap.XSDTime)(xt).UnmarshalXML(d, start)
					}
				{{end}}
			{{end}}
		{{end}}
	{{end}}

	{{range .ComplexTypes}}
		{{/* ComplexTypeGlobal */}}
		// ComplexTypeGlobal {{ $targetNamespace }}
		{{$typeName := replaceReservedWords .Name | wrapNS | makePublic}}
		{{if and (eq (len .SimpleContent.Extension.Attributes) 0) (eq (toGoType .SimpleContent.Extension.Base false) "string") }}
			type {{$typeName}} string
		{{else}}
			type {{$typeName}} struct {
				{{$type := findNameByType .Name}}
				{{if ne .Name $type}}
					XMLName xml.Name ` + "`xml:\"{{$targetNamespace}} {{$type}}\"`" + `
				{{end}}

				{{if ne .ComplexContent.Extension.Base ""}}
					{{template "ComplexContent" .ComplexContent}}
				{{else if ne .SimpleContent.Extension.Base ""}}
					{{template "SimpleContent" .SimpleContent}}
				{{else}}
					{{template "Elements" .Sequence}}
					{{template "Any" .Any}}
					{{template "Elements" .Choice}}
					{{template "Elements" .SequenceChoice}}
					{{template "Elements" .All}}
					{{template "Attributes" .Attributes}}
				{{end}}
			}
		{{end}}
	{{end}}
{{end}}
`
