package workspace_service

// UpdateWorkspaceParams contains parameters for updating a workspace
type UpdateWorkspaceParams struct {
	Name         string
	Logo         string
	Website      string
	About        string
	BillingEmail string
}
