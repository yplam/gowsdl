// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package gowsdl

var opsTmpl = `
{{range .}}
	{{$privateType := .Name | makePrivate}}
	{{$exportType := .Name | makePublic}}

	type {{$exportType}} interface {
		{{range .Operations}}
			{{$faults := len .Faults}}
			{{$soapAction := findSOAPAction .Name $privateType}}
			{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
			{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}

			{{/*if ne $soapAction ""*/}}
			{{if gt $faults 0}}
			// Error can be either of the following types:
			// {{range .Faults}}
			//   - {{.Name}} {{.Doc}}{{end}}{{end}}
			{{if ne .Doc ""}}/* {{.Doc}} */{{end}}
			{{makePublic .Name | replaceReservedWords}} (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error)
			{{/*end*/}}
		{{end}}
	}

	type {{$privateType}} struct {
		client *soap.Client
	}

	func New{{$exportType}}(client *soap.Client) {{$exportType}} {
		return &{{$privateType}}{
			client: client,
		}
	}

	{{range .Operations}}
		{{$requestType := findType .Input.Message | replaceReservedWords | makePublic}}
		{{$soapAction := findSOAPAction .Name $privateType}}
		{{$responseType := findType .Output.Message | replaceReservedWords | makePublic}}

		type {{$requestType}}Body struct {
			{{$requestType}} {{$requestType}}
		}

		type {{$responseType}}Body struct {
			soap.EnvelopeResponseBody
			{{$responseType}} {{$responseType}}
		}

		func (service *{{$privateType}}) {{makePublic .Name | replaceReservedWords}} (ctx context.Context, {{if ne $requestType ""}}request *{{$requestType}}{{end}}) ({{if ne $responseType ""}}*{{$responseType}}, {{end}}error) {

			envelope := soap.NewEnvelope()
			envelope.Body = &request

			response := new({{$responseType}}Body)
			envelopeResp := &soap.EnvelopeResponse{
				Body: response,
			}
			err := service.client.Call(ctx,
				"{{if ne $soapAction ""}}{{$soapAction}}{{else}}''{{end}}",
				envelope, envelopeResp)
			if err != nil {
				return {{if ne $responseType ""}}nil, {{end}}err
			}
			return {{if ne $responseType ""}}&response.{{$responseType}}, {{end}}nil

		}
	{{end}}
{{end}}
`
