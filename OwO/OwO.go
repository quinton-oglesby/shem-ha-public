// Package to OwOify a given string.
package OwO

import (
	"math/rand"
	"strings"
	"time"
)

var prefixes = []string{
	"<3 ",
	"0w0 ",
	"H-hewwo?? ",
	"HIIII! ",
	"Haiiii! ",
	"Huohhhh. ",
	"OWO ",
	"OwO ",
	"UwU ",
}

var suffixes = []string{
	" ( ͡° ᴥ ͡°)",
	" (´・ω・｀)",
	" (இωஇ )",
	" (๑•́ ₃ •̀๑)",
	" (• o •)",
	" (⁎˃ᆺ˂)",
	" (╯﹏╰）",
	" (●´ω｀●)",
	" (◠‿◠✿)",
	" (✿ ♡‿♡)",
	" (❁´◡`❁)",
	" (人◕ω◕)",
	" (；ω；)",
	" (｀へ´)",
	" ._.",
	" :3",
	" :D",
	" :P",
	" ;-;",
	" ;3",
	" ;_;",
	" <{^v^}>",
	" >_<",
	" >_>",
	" UwU",
	" XDDD",
	" ^-^",
	" ^_^",
	" x3",
	" x3",
	" xD",
	" ÙωÙ",
	" ʕʘ‿ʘʔ",
	" ㅇㅅㅇ",
	", fwendo",
	"（＾ｖ＾）",
}

var substitutions = map[string]string{
	"r":    "w",
	"l":    "w",
	"R":    "W",
	"L":    "W",
	"no":   "nu",
	"has":  "haz",
	"have": "haz",
	"you":  "uu",
	"the ": "da ",
	"The ": "Da ",
	"THE ": "DA ",
}

/*UwU Convewts da specified stwing into OwO speak ʕʘ‿ʘʔ
//:param text: Huohhhh. Da text uu want to convewt..
:return: OWO Da convewted stwing (人◕ω◕)*/
func whats_this(text string) string {
	for key, value := range substitutions {
		text = strings.Replace(text, key, value, -1)
	}

	rand.Seed(time.Now().Unix())
	text = prefixes[rand.Intn(len(prefixes))] + text + suffixes[rand.Intn(len(suffixes))]

	return text
}
