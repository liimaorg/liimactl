// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package japanese_test

import (
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/internal"
	"golang.org/x/text/encoding/internal/enctest"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func dec(e encoding.Encoding) (dir string, t transform.Transformer, err error) {
	return "Decode", e.NewDecoder(), nil
}
func enc(e encoding.Encoding) (dir string, t transform.Transformer, err error) {
	return "Encode", e.NewEncoder(), internal.ErrASCIIReplacement
}

func TestNonRepertoire(t *testing.T) {
	testCases := []struct {
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	}{
		{dec, japanese.EUCJP, "\xfe\xfc", "\ufffd"},
		{dec, japanese.ISO2022JP, "\x1b$B\x7e\x7e", "\ufffd"},
		{dec, japanese.ShiftJIS, "\xef\xfc", "\ufffd"},

		{enc, japanese.EUCJP, "갂", ""},
		{enc, japanese.EUCJP, "a갂", "a"},
		{enc, japanese.EUCJP, "丌갂", "\x8f\xb0\xa4"},

		{enc, japanese.ISO2022JP, "갂", ""},
		{enc, japanese.ISO2022JP, "a갂", "a"},
		{enc, japanese.ISO2022JP, "朗갂", "\x1b$BzF\x1b(B"}, // switch back to ASCII mode at end

		{enc, japanese.ShiftJIS, "갂", ""},
		{enc, japanese.ShiftJIS, "a갂", "a"},
		{enc, japanese.ShiftJIS, "\u2190갂", "\x81\xa9"},
	}
	for _, tc := range testCases {
		dir, tr, wantErr := tc.init(tc.e)

		dst, _, err := transform.String(tr, tc.src)
		if err != wantErr {
			t.Errorf("%s %v(%q): got %v; want %v", dir, tc.e, tc.src, err, wantErr)
		}
		if got := string(dst); got != tc.want {
			t.Errorf("%s %v(%q):\ngot  %q\nwant %q", dir, tc.e, tc.src, got, tc.want)
		}
	}
}

func TestCorrect(t *testing.T) {
	testCases := []struct {
		init      func(e encoding.Encoding) (string, transform.Transformer, error)
		e         encoding.Encoding
		src, want string
	}{
		{dec, japanese.ShiftJIS, "\x9f\xfc", "滌"},
		{dec, japanese.ShiftJIS, "\xfb\xfc", "髙"},
		{dec, japanese.ShiftJIS, "\xfa\xb1", "﨑"},
		{enc, japanese.ShiftJIS, "滌", "\x9f\xfc"},
		{enc, japanese.ShiftJIS, "﨑", "\xed\x95"},
	}
	for _, tc := range testCases {
		dir, tr, _ := tc.init(tc.e)

		dst, _, err := transform.String(tr, tc.src)
		if err != nil {
			t.Errorf("%s %v(%q): got %v; want %v", dir, tc.e, tc.src, err, nil)
		}
		if got := string(dst); got != tc.want {
			t.Errorf("%s %v(%q):\ngot  %q\nwant %q", dir, tc.e, tc.src, got, tc.want)
		}
	}
}

func TestBasics(t *testing.T) {
	// The encoded forms can be verified by the iconv program:
	// $ echo 月日は百代 | iconv -f UTF-8 -t SHIFT-JIS | xxd
	testCases := []struct {
		e         encoding.Encoding
		encPrefix string
		encSuffix string
		encoded   string
		utf8      string
	}{{
		// "A｡ｶﾟ 0208: etc 0212: etc" is a nonsense string that contains ASCII, half-width
		// kana, JIS X 0208 (including two near the kink in the Shift JIS second byte
		// encoding) and JIS X 0212 encodable codepoints.
		//
		// "月日は百代の過客にして、行かふ年も又旅人也。" is from the 17th century poem
		// "Oku no Hosomichi" and contains both hiragana and kanji.
		e: japanese.EUCJP,
		encoded: "A\x8e\xa1\x8e\xb6\x8e\xdf " +
			"0208: \xa1\xa1\xa1\xa2\xa1\xdf\xa1\xe0\xa1\xfd\xa1\xfe\xa2\xa1\xa2\xa2\xf4\xa6 " +
			"0212: \x8f\xa2\xaf\x8f\xed\xe3",
		utf8: "A｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199 " +
			"0212: \u02d8\u9fa5",
	}, {
		e: japanese.EUCJP,
		encoded: "\xb7\xee\xc6\xfc\xa4\xcf\xc9\xb4\xc2\xe5\xa4\xce\xb2\xe1\xb5\xd2" +
			"\xa4\xcb\xa4\xb7\xa4\xc6\xa1\xa2\xb9\xd4\xa4\xab\xa4\xd5\xc7\xaf" +
			"\xa4\xe2\xcb\xf4\xce\xb9\xbf\xcd\xcc\xe9\xa1\xa3",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	}, {
		e:         japanese.ISO2022JP,
		encSuffix: "\x1b\x28\x42",
		encoded: "\x1b\x28\x49\x21\x36\x5f\x1b\x28\x42 " +
			"0208: \x1b\x24\x42\x21\x21\x21\x22\x21\x5f\x21\x60\x21\x7d\x21\x7e\x22\x21\x22\x22\x74\x26",
		utf8: "｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199",
	}, {
		e:         japanese.ISO2022JP,
		encPrefix: "\x1b\x24\x42",
		encSuffix: "\x1b\x28\x42",
		encoded: "\x37\x6e\x46\x7c\x24\x4f\x49\x34\x42\x65\x24\x4e\x32\x61\x35\x52" +
			"\x24\x4b\x24\x37\x24\x46\x21\x22\x39\x54\x24\x2b\x24\x55\x47\x2f" +
			"\x24\x62\x4b\x74\x4e\x39\x3f\x4d\x4c\x69\x21\x23",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	}, {
		e: japanese.ShiftJIS,
		encoded: "A\xa1\xb6\xdf " +
			"0208: \x81\x40\x81\x41\x81\x7e\x81\x80\x81\x9d\x81\x9e\x81\x9f\x81\xa0\xea\xa4",
		utf8: "A｡ｶﾟ " +
			"0208: \u3000\u3001\u00d7\u00f7\u25ce\u25c7\u25c6\u25a1\u7199",
	}, {
		e: japanese.ShiftJIS,
		encoded: "\x8c\x8e\x93\xfa\x82\xcd\x95\x53\x91\xe3\x82\xcc\x89\xdf\x8b\x71" +
			"\x82\xc9\x82\xb5\x82\xc4\x81\x41\x8d\x73\x82\xa9\x82\xd3\x94\x4e" +
			"\x82\xe0\x96\x94\x97\xb7\x90\x6c\x96\xe7\x81\x42",
		utf8: "月日は百代の過客にして、行かふ年も又旅人也。",
	}}

	for _, tc := range testCases {
		enctest.TestEncoding(t, tc.e, tc.encoded, tc.utf8, tc.encPrefix, tc.encSuffix)
	}
}

func TestFiles(t *testing.T) {
	enctest.TestFile(t, japanese.EUCJP)
	enctest.TestFile(t, japanese.ISO2022JP)
	enctest.TestFile(t, japanese.ShiftJIS)
}

func BenchmarkEncoding(b *testing.B) {
	enctest.Benchmark(b, japanese.EUCJP)
	enctest.Benchmark(b, japanese.ISO2022JP)
	enctest.Benchmark(b, japanese.ShiftJIS)
}
