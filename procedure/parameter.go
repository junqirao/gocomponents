package procedure

type (
	// Parameter defines parameter
	Parameter struct {
		// Name of parameter
		Name string `json:"name"`
		// Validate validator define of parameter
		Validate *ParameterValidate `json:"validate,omitempty"`
		// Meta of parameter
		Meta Meta `json:"meta"`
	}
	// ParameterValidate defines parameter
	// validate rules and error messages.
	// see https://goframe.org/docs/core/gvalid
	ParameterValidate struct {
		// validator rule
		// e.g. "required|length:6,16|same:password2"
		Rule string `json:"rule"`
		// error message accept map or string
		// e.g. "message": "xxx"
		// e.g. "message": {"required": "xxx", "same": "xxx"}
		Message any `json:"message"`
	}
)
