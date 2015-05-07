package lib

// renderIndex generates the HTML response for this route.
func renderIndex(fv map[string]interface{}, user tinyUser) []byte {
	// Generate the markup for the results template.
	// if user != nil {
	// 	vars := map[string]interface{}{"User": user}
	// 	markup := executeTemplate("results", vars)
	// 	fv["Results"] = template.HTML(string(markup))
	// }
	//
	// // Generate the markup for the index template.
	// markup := executeTemplate("index", fv)
	//
	// // Generate the final markup with the layout template.
	// vars := map[string]interface{}{"LayoutContent": template.HTML(string(markup))}
	// return executeTemplate("layout", vars)
	return nil
}
