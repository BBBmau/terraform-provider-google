package tpgresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
)

// HandleListError is a helper function that handles the common error reporting pattern
// in list resource implementations. It adds an error to diagnostics, creates a result,
// pushes it to the stream, and sets the stream results to diagnostics.
//
// This function should be called when an error occurs within a list resource's
// stream.Results function. After calling this function, the caller should return
// from the function to stop processing.
func HandleListError(
	ctx context.Context,
	req list.ListRequest,
	diags *diag.Diagnostics,
	push func(list.ListResult) bool,
	stream *list.ListResultsStream,
	summary, message string,
) {
	diags.AddError(summary, message)
	result := req.NewListResult(ctx)
	result.Diagnostics = *diags
	push(result)
	stream.Results = list.ListResultsStreamDiagnostics(*diags)
}
