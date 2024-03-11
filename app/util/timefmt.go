// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"time"
)

const (
	dateFmtFull    = "2006-01-02 15:04:05"
	dateFmtFullMS  = "2006-01-02 15:04:05.000000"
	dateFmtHours   = "15:04:05"
	dateFmtHTML    = "2006-01-02T15:04:05"
	dateFmtJS      = "2006-01-02T15:04:05Z"
	dateFmtRFC3339 = "2006-01-02T15:04:05.000Z07:00"
	dateFmtVerbose = "Mon Jan 2 15:04:05 2006 -0700"
	dateFmtYMD     = "2006-01-02"
)

func TimeToFull(d *time.Time) string {
	return TimeToString(d, dateFmtFull)
}

func TimeToFullMS(d *time.Time) string {
	return TimeToString(d, dateFmtFullMS)
}

func TimeToHours(d *time.Time) string {
	return TimeToString(d, dateFmtHours)
}

func TimeToHTML(d *time.Time) string {
	return TimeToString(d, dateFmtHTML)
}

func TimeToJS(d *time.Time) string {
	return TimeToString(d, dateFmtJS)
}

func TimeToRFC3339(d *time.Time) string {
	return TimeToString(d, dateFmtRFC3339)
}

func TimeToVerbose(d *time.Time) string {
	return TimeToString(d, dateFmtVerbose)
}

func TimeToYMD(d *time.Time) string {
	return TimeToString(d, dateFmtYMD)
}

func TimeFromFull(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtFull)
}

func TimeFromFullMS(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtFullMS)
}

func TimeFromHTML(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtHTML)
}

func TimeFromJS(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtJS)
}

func TimeFromRFC3339(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtRFC3339)
}

func TimeFromVerbose(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtVerbose)
}

func TimeFromYMD(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtYMD)
}
