/*
Copyright (C) 2022-2023 Traefik Labs

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package v1alpha1_test

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestAPI_Validation(t *testing.T) {
	t.Parallel()

	tests := []validationTestCase{
		{
			desc: "valid: minimal",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  currentVersion: my-api-v2`),
		},
		{
			desc: "valid: full",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  cors:
    allowCredentials: true
    allowHeaders: ["X-API-Name"]
    allowMethods: ["GET"]
    allowOriginList: ["*"]
    allowOriginListRegex: [".*"]
    exposeHeaders: ["Content-Encoding"]
    maxAge: 10
  headers:
    request:
      set:
        X-API-Name: my-api
      delete: ["X-Old-API-Name"]
    response:
      set:
        X-API-Name: my-api
      delete: ["X-Old-API-Name"]
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: /openapi.json`),
		},
		{
			desc: "missing resource namespace",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: "my-api"
spec:
  pathPrefix: /api
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeRequired, Field: "metadata.namespace", BadValue: ""}},
		},
		{
			desc: "invalid resource name",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: .non-dns-compliant-api
  namespace: my-ns
spec:
  pathPrefix: /api
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "metadata.name", BadValue: ".non-dns-compliant-api", Detail: "a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?')"}},
		},
		{
			desc: "missing resource name",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: ""
  namespace: my-ns
spec:
  pathPrefix: /api
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeRequired, Field: "metadata.name", BadValue: "", Detail: "name or generateName is required"}},
		},
		{
			desc: "resource name is too long",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: api-with-a-way-toooooooooooooooooooooooooooooooooooooo-long-name
  namespace: my-ns
spec:
  pathPrefix: /api
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "metadata.name", BadValue: "api-with-a-way-toooooooooooooooooooooooooooooooooooooo-long-name", Detail: "must be no more than 63 characters"}},
		},
		{
			desc: "path prefix is required",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeRequired, Field: "spec.pathPrefix", BadValue: ""}},
		},
		{
			desc: "path prefix must start with a /",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: something
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.pathPrefix", BadValue: "string", Detail: "must start with a '/'"}},
		},
		{
			desc: "path prefix cannot contains ../",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /foo/../bar
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.pathPrefix", BadValue: "string", Detail: "cannot contains '../'"}},
		},
		{
			desc: "path prefix cannot ends with /..",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /foo/..
  currentVersion: my-api-v2`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.pathPrefix", BadValue: "string", Detail: "cannot contains '../'"}},
		},
		{
			desc: "valid: pathPrefix with segment starting with ..",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /foo/..bar
  currentVersion: my-api-v2`),
		},
		{
			desc: "valid: pathPrefix with segment starting with .well-known",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /foo/.well-known
  currentVersion: my-api-v2`),
		},
		{
			desc: "versioned API cannot define cors, headers and service",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  currentVersion: my-api-v2
  cors:
    allowMethods: ["GET"]
  headers:
    request:
      set:
        X-API-Name: my-api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: /openapi.json`),
			wantErrs: field.ErrorList{
				{Type: field.ErrorTypeInvalid, Field: "spec", BadValue: "object", Detail: "currentVersion and service are mutually exclusive"},
				{Type: field.ErrorTypeInvalid, Field: "spec", BadValue: "object", Detail: "currentVersion and cors are mutually exclusive"},
				{Type: field.ErrorTypeInvalid, Field: "spec", BadValue: "object", Detail: "currentVersion and headers are mutually exclusive"},
			},
		},
		{
			desc: "currentVersion or service must be defined",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec", BadValue: "object", Detail: "currentVersion or service must be defined"}},
		},
		{
			desc: "service name, port and openApiSpec are required",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service: {}`),
			wantErrs: field.ErrorList{
				{Type: field.ErrorTypeRequired, Field: "spec.service.name", BadValue: ""},
				{Type: field.ErrorTypeRequired, Field: "spec.service.openApiSpec", BadValue: ""},
				{Type: field.ErrorTypeRequired, Field: "spec.service.port", BadValue: ""},
			},
		},
		{
			desc: "openApiSpec must have a path or an url",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec: {}`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.service.openApiSpec", BadValue: "object", Detail: "path or url must be defined"}},
		},
		{
			desc: "openApiSpec url must be a valid URL",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      url: ../invalid-spec-url.json`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.service.openApiSpec.url", BadValue: "string", Detail: "must be a valid URL"}},
		},
		{
			desc: "openApiSpec path must start with a /",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: something`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.service.openApiSpec.path", BadValue: "string", Detail: "must start with a '/'"}},
		},
		{
			desc: "path prefix cannot contains ../",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: /foo/../bar`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.service.openApiSpec.path", BadValue: "string", Detail: "cannot contains '../'"}},
		},
		{
			desc: "path prefix cannot ends with /..",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: /foo/..`),
			wantErrs: field.ErrorList{{Type: field.ErrorTypeInvalid, Field: "spec.service.openApiSpec.path", BadValue: "string", Detail: "cannot contains '../'"}},
		},
		{
			desc: "valid: pathPrefix with segment starting with ..",
			manifest: []byte(`
apiVersion: hub.traefik.io/v1alpha1
kind: API
metadata:
  name: my-api
  namespace: my-ns
spec:
  pathPrefix: /api
  service:
    name: my-svc
    port:
      number: 8080
    openApiSpec:
      path: /foo/..bar`),
		},
	}

	checkValidationTestCases(t, tests)
}