// Package enum - Content managed by Project Forge, see [projectforge.md] for details.
package enum

var (
	ExampleValues = Values{{Key: "x", Name: "X", Description: "It's an x", Icon: "x-icon", Default: false, Simple: false}}

	Examples = map[string]any{
		"values": ExampleValues,
	}
)
