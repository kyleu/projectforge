// $PF_GENERATE_ONCE$
package gql

import "context"

func (s *Schema) Hello(ctx context.Context) (string, error) {
	return "Howdy!", nil
}

func (s *Schema) Poke(ctx context.Context) (string, error) {
	return "OK!", nil
}
