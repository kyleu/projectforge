package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

var (
	FmtCode       = Format{Key: "code", Title: "Code", Description: ""}
	FmtCodeHidden = Format{Key: "codehidden", Title: "Code (hidden)", Description: ""}
	FmtColor      = Format{Key: "color", Title: "Color", Description: ""}
	FmtCountry    = Format{Key: "country", Title: "Country", Description: ""}
	FmtHTML       = Format{Key: "html", Title: "HTML", Description: ""}
	FmtIcon       = Format{Key: "icon", Title: "Icon", Description: ""}
	FmtImage      = Format{Key: "image", Title: "Image", Description: ""}
	FmtJSON       = Format{Key: "json", Title: "JSON", Description: ""}
	FmtLinebreaks = Format{Key: "linebreaks", Title: "Linebreaks", Description: ""}
	FmtSelect     = Format{Key: "select", Title: "Select Box", Description: ""}
	FmtSeconds    = Format{Key: "seconds", Title: "Seconds", Description: ""}
	FmtSI         = Format{Key: "si", Title: "SI Units", Description: ""}
	FmtSQL        = Format{Key: "sql", Title: "SQL", Description: ""}
	FmtTags       = Format{Key: "tags", Title: "Tags", Description: ""}
	FmtURL        = Format{Key: "url", Title: "URL", Description: ""}

	AllFormats = Formats{FmtCode, FmtCodeHidden, FmtColor, FmtCountry, FmtHTML, FmtIcon, FmtImage, FmtJSON, FmtSelect, FmtSeconds, FmtSI, FmtSQL, FmtURL}
)

type Format struct {
	Key         string
	Title       string
	Description string
}

func (x *Format) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(x.Key, false), nil
}

func (x *Format) UnmarshalJSON(data []byte) error {
	str, err := util.FromJSONString(data)
	if err != nil {
		return err
	}
	curr := AllFormats.Get(str, nil)
	x.Key = curr.Key
	x.Title = curr.Title
	x.Description = curr.Description
	return nil
}

type Formats []Format

func (l Formats) Get(s string, logger util.Logger) Format {
	for _, x := range l {
		if x.Key == s {
			return x
		}
	}
	msg := fmt.Sprintf("unable to find [Format] enum with key [%s]", s)
	if logger != nil {
		logger.Warn(msg)
	}
	return Format{Key: "_error", Title: "error: " + msg}
}
