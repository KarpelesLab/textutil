package textutil

import (
	"strings"
	"testing"
)

func TestWrapString(t *testing.T) {
	cases := []struct {
		Input, Output string
		Opts          *WrapOptions
	}{
		// A simple word passes through.
		{
			Input:  "foo",
			Output: "foo",
			Opts:   &WrapOptions{Limit: 4},
		},
		// A single word that is too long passes through.
		// We do not break words.
		{
			Input:  "foobarbaz",
			Output: "foobarbaz",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Lines are broken at whitespace.
		{
			Input:  "foo bar baz",
			Output: "foo\nbar\nbaz",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Lines are broken at whitespace, even if words
		// are too long. We do not break words.
		{
			Input:  "foo bars bazzes",
			Output: "foo\nbars\nbazzes",
			Opts:   &WrapOptions{Limit: 4},
		},
		// A word that would run beyond the width is wrapped.
		{
			Input:  "fo sop",
			Output: "fo\nsop",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Do not break on non-breaking space.
		{
			Input:  "foo bar\u00A0baz",
			Output: "foo\nbar\u00A0baz",
			Opts:   &WrapOptions{Limit: 10},
		},
		// Whitespace that trails a line and fits the width
		// are trimmped, as does whitespace prefixing an
		// explicit line break. A tab counts as one character.
		{
			Input:  "foo\nb\t r\n baz",
			Output: "foo\nb\t r\nbaz",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Trailing whitespace is removed if it doesn't fit the width.
		// Runs of whitespace on which a line is broken are removed.
		{
			Input:  "foo    \nb   ar   ",
			Output: "foo\nb\nar",
			Opts:   &WrapOptions{Limit: 4},
		},
		// An explicit line break at the end of the input is preserved.
		{
			Input:  "foo bar baz\n",
			Output: "foo\nbar\nbaz\n",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Explicit break are always preserved.
		{
			Input:  "\nfoo bar\n\n\nbaz\n",
			Output: "\nfoo\nbar\n\n\nbaz\n",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Complete example:
		{
			Input:  " This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* baz  \nBAM    ",
			Output: "This\nis a\nlist:\n\n* foo\n* bar\n\n\n* baz\nBAM",
			Opts:   &WrapOptions{Limit: 6},
		},
		// Multi-byte characters
		{
			Input:  strings.Repeat("\u2584 ", 4),
			Output: "\u2584 \u2584\n\u2584 \u2584",
			Opts:   &WrapOptions{Limit: 4},
		},
		// Email header test
		{
			Input:  "Received: from a25-34.smtp-out.us-west-2.amazonses.com (a25-34.smtp-out.us-west-2.amazonses.com. [54.240.25.34]) by mx.google.com with ESMTPS id 189d4az89d4az98d.dazfze8fz.Fezfzefez.2021.10.25 for <test@example.com> (version=TLS1_2 cipher=ECDHE-ECDSA-AES128-SHA bits=128/128); Mon, 25 Oct 2021 08:53:19 -0700 (PDT)",
			Output: "Received: from a25-34.smtp-out.us-west-2.amazonses.com\r\n\t(a25-34.smtp-out.us-west-2.amazonses.com. [54.240.25.34]) by mx.google.com\r\n\twith ESMTPS id 189d4az89d4az98d.dazfze8fz.Fezfzefez.2021.10.25 for\r\n\t<test@example.com> (version=TLS1_2 cipher=ECDHE-ECDSA-AES128-SHA\r\n\tbits=128/128); Mon, 25 Oct 2021 08:53:19 -0700 (PDT)",
			Opts: &WrapOptions{
				Limit:     76,
				Prefix:    "\t",
				Linebreak: CRLF,
			},
		},
		// ARC test
		{
			Input:  "ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=google.com; s=arc-20160816; h=feedback-id:date:bounces-to:mime-version:subject:message-id:to:reply-to:from:dkim-signature:dkim-signature; bh=jQeY2dlYpkluPbrBBFicWp/Jx7XMgQUiI6R8I7mXbd4=; b=takw5mTuZV9nYb/GiPlNsA3QrrYeJC3E+wchH/KHCeXBoiy/j/fxlHdTN4GNFflVJJo3tVpKeWyk8nqGp3OIYRGGNEtZ2xWj8/I+9QxzE4J657uAdMM11Wg7J7CyZFXKGFAKvpVYDlBBUsbbXnUGSEmjLX2vgVvidMppLTpqO7Gtzjej09NBr7T1dPTk/B/FBiONTb6Mgxby2/JOqKlPe8ZPSPZvJTNaD9wdI6YXHTUwGcuaWm5U4ZThmn3G9lhmwXY/eYP5mppTtJR7Dpf2JlLMBA+G0+VjEC7/qR6974PeJEI32QGS7RkLijFWGq6u23ALTrym5nzztH8WOzHscw==",
			Output: "ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=google.com;\r\n        s=arc-20160816;\r\n        h=feedback-id:date:bounces-to:mime-version:subject:message-id:to:reply-to:from:dkim-signature:dkim-signature;\r\n         bh=jQeY2dlYpkluPbrBBFicWp/Jx7XMgQUiI6R8I7mXbd4=;\r\n        b=takw5mTuZV9nYb/GiPlNsA3QrrYeJC3E+wchH/KHCeXBoiy/j/fxlHdTN4GNFflVJJo3tVpKeWyk8nqGp3OIYRGGNEtZ2xWj8/I+9QxzE4J657uAdMM11Wg7J7CyZFXKGFAKvpVYDlBBUsbbXnUGSEmjLX2vgVvidMppLTpqO7Gtzjej09NBr7T1dPTk/B/FBiONTb6Mgxby2/JOqKlPe8ZPSPZvJTNaD9wdI6YXHTUwGcuaWm5U4ZThmn3G9lhmwXY/eYP5mppTtJR7Dpf2JlLMBA+G0+VjEC7/qR6974PeJEI32QGS7RkLijFWGq6u23ALTrym5nzztH8WOzHscw==",
			Opts: &WrapOptions{
				Limit:     76,
				Prefix:    "        ", // 8 spaces
				Linebreak: CRLF,
			},
		},
		// ARC test with break words
		{
			Input:  "ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=google.com; s=arc-20160816; h=feedback-id:date:bounces-to:mime-version:subject:message-id:to:reply-to:from:dkim-signature:dkim-signature; bh=jQeY2dlYpkluPbrBBFicWp/Jx7XMgQUiI6R8I7mXbd4=; b=takw5mTuZV9nYb/GiPlNsA3QrrYeJC3E+wchH/KHCeXBoiy/j/fxlHdTN4GNFflVJJo3tVpKeWyk8nqGp3OIYRGGNEtZ2xWj8/I+9QxzE4J657uAdMM11Wg7J7CyZFXKGFAKvpVYDlBBUsbbXnUGSEmjLX2vgVvidMppLTpqO7Gtzjej09NBr7T1dPTk/B/FBiONTb6Mgxby2/JOqKlPe8ZPSPZvJTNaD9wdI6YXHTUwGcuaWm5U4ZThmn3G9lhmwXY/eYP5mppTtJR7Dpf2JlLMBA+G0+VjEC7/qR6974PeJEI32QGS7RkLijFWGq6u23ALTrym5nzztH8WOzHscw==",
			Output: "ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=google.com;\r\n        s=arc-20160816;\r\n        h=feedback-id:date:bounces-to:mime-version:subject:message-id:to:repl\r\n        y-to:from:dkim-signature:dkim-signature;\r\n        bh=jQeY2dlYpkluPbrBBFicWp/Jx7XMgQUiI6R8I7mXbd4=;\r\n        b=takw5mTuZV9nYb/GiPlNsA3QrrYeJC3E+wchH/KHCeXBoiy/j/fxlHdTN4GNFflVJJo\r\n        3tVpKeWyk8nqGp3OIYRGGNEtZ2xWj8/I+9QxzE4J657uAdMM11Wg7J7CyZFXKGFAKvpVY\r\n        DlBBUsbbXnUGSEmjLX2vgVvidMppLTpqO7Gtzjej09NBr7T1dPTk/B/FBiONTb6Mgxby2\r\n        /JOqKlPe8ZPSPZvJTNaD9wdI6YXHTUwGcuaWm5U4ZThmn3G9lhmwXY/eYP5mppTtJR7Dp\r\n        f2JlLMBA+G0+VjEC7/qR6974PeJEI32QGS7RkLijFWGq6u23ALTrym5nzztH8WOzHscw=\r\n        =",
			Opts: &WrapOptions{
				Limit:      76,
				Prefix:     "        ", // 8 spaces
				Linebreak:  CRLF,
				BreakWords: true,
			},
		},
	}

	for i, tc := range cases {
		actual := WrapString(tc.Input, tc.Opts)
		if actual != tc.Output {
			t.Fatalf("Case %d Input:\n\n%q\n\nExpected Output:\n\n%q\n\nActual Output:\n\n%q", i, tc.Input, tc.Output, actual)
		}
	}
}
