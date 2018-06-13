package responses

// RedirectResponse represents the url for redirection
type RedirectResponse struct {
	// required: true
	URL string `json:"url"`
}

// see RedirectResponse
// swagger:response Redirect
type swaggerRedirectResponse struct {
	//in: body
	Body RedirectResponse
}
