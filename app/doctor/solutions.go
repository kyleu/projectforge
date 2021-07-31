package doctor

type Solution struct {
	Code  string   `json:"code"`
	Args  []string `json:"args,omitempty"`
	MacOS string   `json:"macOS"`
}

func NewSolution(code string, args []string) *Solution {
	return &Solution{Code: code, Args: args}
}

func (s *Solution) WithMacOS(macOS string) *Solution {
	s.MacOS = macOS
	return s
}

func FindSolution(code string, args ...string) *Solution {
	switch code {
	case "missing":
		switch args[0] {
		case "rsvg":
			return NewSolution(code, args).WithMacOS("brew install librsvg")
		default:
			return nil
		}
	default:
		return nil
	}
}
