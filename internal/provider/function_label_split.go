package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var labelSplitReturnAttrTypes = map[string]attr.Type{"label": types.StringType}

var _ function.Function = &LabelSplitFunction{}

type LabelSplitFunction struct{}

func NewLabelSplitFunction() function.Function {
	return &LabelSplitFunction{}
}

func (f *LabelSplitFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "label_split"
}

func (f *LabelSplitFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Breakout nested tags",
		Description: "Given a Terraform tag with /s, will parse and return a comma separated string of tags at each level of /.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "label",
				Description: "Terraform string to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: labelSplitReturnAttrTypes,
		},
	}
}

func (f *LabelSplitFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var label string
	var AssembledChunks string

	resp.Error = req.Arguments.Get(ctx, &label)
	if resp.Error != nil {
		return
	}

	prevChunks := ""
	chunks := strings.Split(label, "/")

	for _, chunk := range chunks {
		if chunk == "" {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error parsing label: %q cannot have blank label section", label))
			return
		}
		AssembledChunks += prevChunks + chunk + ","
		prevChunks += chunk + "/"
	}

	output, diags := types.ObjectValue(
		labelSplitReturnAttrTypes,
		map[string]attr.Value{"label": types.StringValue(AssembledChunks)},
	)

	resp.Error = function.FuncErrorFromDiags(ctx, diags)
	if resp.Error != nil {
		return
	}

	resp.Error = resp.Result.Set(ctx, &output)
}
