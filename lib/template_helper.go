package lib
import "errors"
// renderIndex generates the HTML response for this route.
func renderIndex(displayUser interface{}) ([]byte,error) {
	if username == nil {
	markup = 	executeTemplate("login" vars map[string]interface{})
	}
	str, ok := displayUser.(string)
	if !ok {
		return nil, errors.New("Session data corrupted.")
	}
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
	markup =

	return nil,nil
}
