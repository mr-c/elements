package lib

import "github.com/antha-lang/antha/component"

var (
	components []component.Component
)

// Helper function to add appropriate component information to the component
// library
func addComponent(desc component.Component) error {
	if err := component.UpdateParamTypes(&desc); err != nil {
		return err
	}
	components = append(components, desc)
	return nil
}

func GetComponents() []component.Component {
	return components
}
