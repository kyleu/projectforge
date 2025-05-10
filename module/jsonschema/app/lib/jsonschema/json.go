package jsonschema

import "{{{ .Package }}}/app/util"

func FromJSON(b []byte) (*Schema, error) {
	ret := &Schema{}
	if err := util.FromJSONStrict(b, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
