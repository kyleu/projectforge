package numeric

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/kyleu/idlingengine/app/util"
)

var prefixes = []string{
	"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion", "sextillion", "septillion", "octillion", "nonillion",
	"decillion", "undecillion", "duodecillion", "tredecillion", "quattuordecillion", "quindecillion", "sedecillion|sexdecillion", "septendecillion",
	"octodecillion", "novemdecillion ", "vigintillion", "unvigintillion", "duovigintillion", "trevigintillion|tresvigintillion", "quattuorvigintillion",
	"quinvigintillion", "sexvigintillion", "septenvigintillion", "octovigintillion", "novemvigintillion", "trigintillion", "untrigintillion",
	"duotrigintillion", "tretrigintillion", "quattuortrigintillion", "quintrigintillion", "sextrigintillion", "septentrigintillion",
	"octotrigintillion", "novemtrigintillion", "quadragintillion", "unquadragintillion", "duoquadragintillion", "trequadragintillion",
	"quattuorquadragintillion", "quinquadragintillion", "sexquadragintillion", "septenquadragintillion|septquadragintillion", "octoquadragintillion",
	"novemquadragintillion", "quinquagintillion", "unquinquagintillion", "duoquinquagintillion", "trequinquagintillion", "quattuorquinquagintillion",
	"quinquinquagintillion", "sexquinquagintillion", "septenquinquagintillion|septquinquagintillion", "octoquinquagintillion", "novemquinquagintillion",
	"sexagintillion", "unsexagintillion", "duosexagintillion", "tresexagintillion", "quattuorsexagintillion", "quinsexagintillion", "sexsexagintillion",
	"septsexagintillion", "octosexagintillion", "octosexagintillion", "septuagintillion", "unseptuagintillion", "duoseptuagintillion", "treseptuagintillion",
	"quinseptuagintillion", "sexseptuagintillion", "septseptuagintillion", "octoseptuagintillion", "novemseptuagintillion", "octogintillion", "unoctogintillion",
	"duooctogintillion", "treoctogintillion", "quattuoroctogintillion", "quinoctogintillion", "sexoctogintillion", "septoctogintillion", "octooctogintillion",
	"novemoctogintillion", "nonagintillion", "unnonagintillion", "duononagintillion", "trenonagintillion", "quattuornonagintillion|quattornonagintillion",
	"quinnonagintillion", "sexnonagintillion", "septnonagintillion", "octononagintillion", "novemnonagintillion", "centillion",
}

var prefixMap = func() map[string]int {
	ret := make(map[string]int, len(prefixes))
	for idx, x := range prefixes {
		for _, s := range util.StringSplitAndTrim(x, "|") {
			ret[s] = idx * 3
		}
	}
	return ret
}()

func Pow10FromEnglish(word string) (int, error) {
	word = strings.ToLower(word)
	if power, exists := prefixMap[word]; exists {
		return power, nil
	}
	return 0, errors.Errorf("unknown number word: %s", word)
}
