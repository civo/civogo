package main

import (
	"go/format"
	"strings"
	"testing"
)

func TestSplitEnumDoc(t *testing.T) {
	overview, values := splitEnumDoc("Reasons overview.\n\nMore text.\n\n - A_NOT_FOUND: No A exists.\n - B_TAKEN: B is taken\nby someone else.\n ")
	if overview != "Reasons overview.\n\nMore text." {
		t.Errorf("overview = %q", overview)
	}
	if values["A_NOT_FOUND"] != "No A exists." || values["B_TAKEN"] != "B is taken\nby someone else." {
		t.Errorf("values = %q", values)
	}
}

func TestRenderReasons(t *testing.T) {
	g := &generator{protoCommit: "abc", specName: "compute.openapi3.json", spec: &spec{}}
	if g.renderReasons() != nil {
		t.Fatal("spec without ErrorReason must render nothing")
	}

	g.spec.Components.Schemas = map[string]*schema{reasonsSchema: {
		Type:        "string",
		Description: "Overview.\n\n - INSTANCE_NOT_FOUND: No instance.\n - SSH_KEY_ALREADY_EXISTS: Key exists.",
		Enum:        []string{"ERROR_REASON_UNSPECIFIED", "INSTANCE_NOT_FOUND", "SSH_KEY_ALREADY_EXISTS"},
	}}
	src, err := format.Source(g.renderReasons())
	if err != nil {
		t.Fatal(err)
	}
	out := string(src)
	for _, want := range []string{
		"package compute",
		`ReasonInstanceNotFound = "INSTANCE_NOT_FOUND"`,
		`ReasonSSHKeyAlreadyExists = "SSH_KEY_ALREADY_EXISTS"`,
		"// ReasonSSHKeyAlreadyExists — Key exists.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "UNSPECIFIED") {
		t.Errorf("UNSPECIFIED value emitted:\n%s", out)
	}
}

func TestFieldTypeNullable(t *testing.T) {
	g := &generator{spec: &spec{}}
	g.spec.Components.Schemas = map[string]*schema{
		"v3alpha1PlacementRule":  {Type: "object", Properties: map[string]*schema{"x": {Type: "string"}}},
		"v3alpha1InstanceStatus": {Type: "string", Enum: []string{"A"}},
	}
	for _, tc := range []struct {
		s    *schema
		want string
	}{
		{&schema{Type: "boolean", Nullable: true}, "*bool"},
		{&schema{Type: "string", Nullable: true}, "*string"},
		{&schema{Type: "boolean"}, "bool"},
		{&schema{Type: "array", Items: &schema{Type: "string"}, Nullable: true}, "[]string"},
		{&schema{Ref: "#/components/schemas/v3alpha1PlacementRule"}, "*PlacementRule"},
		{&schema{Ref: "#/components/schemas/v3alpha1InstanceStatus"}, "InstanceStatus"},
		{&schema{Ref: "#/components/schemas/v1PageInfo"}, "api.PageInfo"},
	} {
		if got := g.fieldType(tc.s); got != tc.want {
			t.Errorf("fieldType(%+v) = %s, want %s", tc.s, got, tc.want)
		}
	}
}
