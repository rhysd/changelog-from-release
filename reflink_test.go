package main

import "testing"

func TestLinkRefs(t *testing.T) {
	tests := []struct {
		what    string
		input   string
		want    string
		repoURL string
	}{
		{
			what:  "issue",
			input: "#123",
			want:  "[#123](https://gh/u/r/issues/123)",
		},
		{
			what:  "user",
			input: "@foo",
			want:  "[@foo](https://gh/foo)",
		},
		{
			what:  "user includes hypen",
			input: "@a-B-2",
			want:  "[@a-B-2](https://gh/a-B-2)",
		},
		{
			what:  "issue number in codeblock",
			input: "`#123`",
			want:  "`#123`",
		},
		{
			what:  "user name in codeblock",
			input: "`@foo`",
			want:  "`@foo`",
		},
		{
			what:  "multiple issues",
			input: "#1 #2 #3",
			want:  "[#1](https://gh/u/r/issues/1) [#2](https://gh/u/r/issues/2) [#3](https://gh/u/r/issues/3)",
		},
		{
			what:  "multiple users",
			input: "@a @b @c",
			want:  "[@a](https://gh/a) [@b](https://gh/b) [@c](https://gh/c)",
		},
		{
			what:  "issues list",
			input: "- #1\n- #2\n- #3",
			want:  "- [#1](https://gh/u/r/issues/1)\n- [#2](https://gh/u/r/issues/2)\n- [#3](https://gh/u/r/issues/3)",
		},
		{
			what:  "issues list",
			input: "- #1\n- #2\n- #3",
			want:  "- [#1](https://gh/u/r/issues/1)\n- [#2](https://gh/u/r/issues/2)\n- [#3](https://gh/u/r/issues/3)",
		},
		{
			what:  "users list",
			input: "- @a\n- @b\n- @c",
			want:  "- [@a](https://gh/a)\n- [@b](https://gh/b)\n- [@c](https://gh/c)",
		},
		{
			what:  "issue follows alphabet",
			input: "a#123",
			want:  "a#123",
		},
		{
			what:  "issue followed by alphabet",
			input: "#123a",
			want:  "#123a",
		},
		{
			// Note: This behavior is different from GitHub (#123 should not be linked). But aligning
			// to GitHub's behavior is hard and it is very edge case for us. So we determined not to
			// align the behavior
			what:  "issue follows _",
			input: "_#123",
			want:  "_[#123](https://gh/u/r/issues/123)",
		},
		{
			// Note: This behavior is different from GitHub (#123 should not be linked). But aligning
			// to GitHub's behavior is hard and it is very edge case for us. So we determined not to
			// align the behavior
			what:  "issue followed by _",
			input: "#123_",
			want:  "[#123](https://gh/u/r/issues/123)_",
		},
		{
			what:  "issue as italic text",
			input: "_#123_",
			want:  "_[#123](https://gh/u/r/issues/123)_", // Linked because it's an italic text
		},
		{
			what:  "issue as part of italic text",
			input: "_foo #123_",
			want:  "_foo [#123](https://gh/u/r/issues/123)_", // Linked because it's an italic text
		},
		{
			what:  "issue next to italic",
			input: "_foo_#123",
			want:  "_foo_[#123](https://gh/u/r/issues/123)",
		},
		{
			what:  "italic next to issue",
			input: "#123_foo_",
			want:  "[#123](https://gh/u/r/issues/123)_foo_",
		},
		{
			what:  "issue follows number",
			input: "1#123",
			want:  "1#123",
		},
		{
			what:  "issue followed by alphabet",
			input: "#123a",
			want:  "#123a",
		},
		{
			what:  "issue followed by sharp",
			input: "#123#456",
			want:  "[#123](https://gh/u/r/issues/123)#456",
		},
		{
			what:  "issue surrounded by punctuations",
			input: "!#123?",
			want:  "![#123](https://gh/u/r/issues/123)?",
		},
		{
			what:  "issue among multibyte characters",
			input: "い#1🐶#2ぬ",
			want:  "い[#1](https://gh/u/r/issues/1)🐶[#2](https://gh/u/r/issues/2)ぬ",
		},
		{
			what:  "user follows alphabet",
			input: "a@foo",
			want:  "a@foo",
		},
		{
			// Note: This behavior is different from GitHub (@foo should not be linked). But aligning
			// to GitHub's behavior is hard and it is very edge case for us. So we determined not to
			// align the behavior
			what:  "user follows _",
			input: "_@foo",
			want:  "_[@foo](https://gh/foo)",
		},
		{
			// Note: This behavior is different from GitHub (@foo should not be linked). But aligning
			// to GitHub's behavior is hard and it is very edge case for us. So we determined not to
			// align the behavior
			what:  "user followed by _",
			input: "@foo_",
			want:  "[@foo](https://gh/foo)_",
		},
		{
			what:  "user as italic text",
			input: "_@foo_",
			want:  "_[@foo](https://gh/foo)_", // Linked because of italic text
		},
		{
			what:  "user as part of italic text",
			input: "_foo @foo_",
			want:  "_foo [@foo](https://gh/foo)_", // Linked because of italic text
		},
		{
			what:  "user follows hyphen",
			input: "-@foo",
			want:  "-[@foo](https://gh/foo)",
		},
		{
			what:  "user ends with hyphen",
			input: "@foo-",
			want:  "@foo-",
		},
		{
			what:  "user starts with hyphen",
			input: "@-foo",
			want:  "@-foo",
		},
		{
			what:  "user follows number",
			input: "1@foo",
			want:  "1@foo",
		},
		{
			what:  "user followed by other user",
			input: "@a@b",
			want:  "[@a](https://gh/a)@b",
		},
		{
			what:  "user surrounded by punctuations",
			input: "!@a?",
			want:  "![@a](https://gh/a)?",
		},
		{
			what:  "user among multibyte characters",
			input: "い@X🐶@Yぬ",
			want:  "い[@X](https://gh/X)🐶[@Y](https://gh/Y)ぬ",
		},
		{
			what:  "users and issues are mixed",
			input: "#1 @a #2 @b",
			want:  "[#1](https://gh/u/r/issues/1) [@a](https://gh/a) [#2](https://gh/u/r/issues/2) [@b](https://gh/b)",
		},
		{
			what:  "single sharp",
			input: "#",
			want:  "#",
		},
		{
			what:  "single at",
			input: "@",
			want:  "@",
		},
		{
			what:  "text ends with sharp",
			input: "#123 foo #",
			want:  "[#123](https://gh/u/r/issues/123) foo #",
		},
		{
			what:  "text ends with at",
			input: "@foo bar @",
			want:  "[@foo](https://gh/foo) bar @",
		},
		{
			what:  "text starts with sharp",
			input: "# 123",
			want:  "# 123",
		},
		{
			what:  "text starts with at",
			input: "@ foo",
			want:  "@ foo",
		},
		{
			what:  "empty",
			input: "",
			want:  "",
		},
		{
			what:  "quote",
			input: "> @foo\n> #1",
			want:  "> [@foo](https://gh/foo)\n> [#1](https://gh/u/r/issues/1)",
		},
		{
			what:  "issue in link",
			input: "[oops #1](https://example.com/foo/bar?a=b#frag)",
			want:  "[oops #1](https://example.com/foo/bar?a=b#frag)",
		},
		{
			what:  "user name in link",
			input: "[@foo woo](https://example.com/foo/bar?a=b#frag)",
			want:  "[@foo woo](https://example.com/foo/bar?a=b#frag)",
		},
		{
			what:  "italic",
			input: "*@foo* *#1*",
			want:  "*[@foo](https://gh/foo)* *[#1](https://gh/u/r/issues/1)*",
		},
		{
			what:  "italic with _",
			input: "_@foo_ _#1_",
			want:  "_[@foo](https://gh/foo)_ _[#1](https://gh/u/r/issues/1)_",
		},
		{
			what:  "bold",
			input: "**@foo** **#1**",
			want:  "**[@foo](https://gh/foo)** **[#1](https://gh/u/r/issues/1)**",
		},
		{
			what:  "bold with _",
			input: "__@foo__ __#1__",
			want:  "__[@foo](https://gh/foo)__ __[#1](https://gh/u/r/issues/1)__",
		},
		{
			what:  "code fence",
			input: "```\n#123\n@foo\n```",
			want:  "```\n#123\n@foo\n```",
		},
		{
			what:  "<pre> html element",
			input: "<pre>hi #123 @foo</pre>",
			want:  "<pre>hi #123 @foo</pre>",
		},
		{
			what:    "issue with GHE URL",
			input:   "#123",
			want:    "[#123](https://github.some-company.com/user/repo/issues/123)",
			repoURL: "https://github.some-company.com/user/repo",
		},
		{
			what:    "user name with GHE URL",
			input:   "@foo",
			want:    "[@foo](https://github.some-company.com/foo)",
			repoURL: "https://github.some-company.com/user/repo",
		},
		{
			what:  "slash after user name",
			input: "@foo/",
			want:  "@foo/",
		},
		{
			what:  "slash and something after user name",
			input: "@foo/bar",
			want:  "@foo/bar",
		},
		{
			what:  "slash before user name",
			input: "/@foo",
			want:  "/[@foo](https://gh/foo)",
		},
		{
			what:  "commit sha",
			input: "41608e5f4109208a6ab995c58266554e6071c5b2",
			want:  "[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)",
		},
		{
			what:  "mutiple commit sha",
			input: "41608e5f4109208a6ab995c58266554e6071c5b2 41608e5f4109208a6ab995c58266554e6071c5b2 f7b60f34e0a60a0e67f2864f6cebdacc7e247e29",
			want:  "[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2) [`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2) [`f7b60f34e0`](https://gh/u/r/commit/f7b60f34e0a60a0e67f2864f6cebdacc7e247e29)",
		},
		{
			what:  "commit sha shorter than 40 characters",
			input: "41608e5f4109208a6ab995c58266554e6071c5b",
			want:  "41608e5f4109208a6ab995c58266554e6071c5b",
		},
		{
			what:  "commit sha longer than 40 characters",
			input: "41608e5f4109208a6ab995c58266554e6071c5b2f",
			want:  "41608e5f4109208a6ab995c58266554e6071c5b2f",
		},
		{
			what:  "italic commit sha",
			input: "_41608e5f4109208a6ab995c58266554e6071c5b2_ *41608e5f4109208a6ab995c58266554e6071c5b2*",
			want:  "_[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)_ *[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)*",
		},
		{
			what:  "bold commit sha",
			input: "__41608e5f4109208a6ab995c58266554e6071c5b2__ **41608e5f4109208a6ab995c58266554e6071c5b2**",
			want:  "__[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)__ **[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)**",
		},
		{
			what:  "commit sha follows alphabets",
			input: "z41608e5f4109208a6ab995c58266554e6071c5b2",
			want:  "z41608e5f4109208a6ab995c58266554e6071c5b2",
		},
		{
			what:  "commit sha followed by alphabets",
			input: "41608e5f4109208a6ab995c58266554e6071c5b2z",
			want:  "41608e5f4109208a6ab995c58266554e6071c5b2z",
		},
		{
			what:  "commit sha in link",
			input: "[41608e5f4109208a6ab995c58266554e6071c5b2 is awesome commit](https://example.com)",
			want:  "[41608e5f4109208a6ab995c58266554e6071c5b2 is awesome commit](https://example.com)",
		},
		{
			what:  "commit sha in code span",
			input: "`41608e5f4109208a6ab995c58266554e6071c5b2`",
			want:  "`41608e5f4109208a6ab995c58266554e6071c5b2`",
		},
		{
			what:  "commit sha in code fence",
			input: "```\n41608e5f4109208a6ab995c58266554e6071c5b2\n```",
			want:  "```\n41608e5f4109208a6ab995c58266554e6071c5b2\n```",
		},
		{
			what:  "commit sha in <pre> html element",
			input: "<pre>41608e5f4109208a6ab995c58266554e6071c5b2</pre>",
			want:  "<pre>41608e5f4109208a6ab995c58266554e6071c5b2</pre>",
		},
		{
			what:  "commit sha among multiple characters",
			input: "いぬ41608e5f4109208a6ab995c58266554e6071c5b2🐶",
			want:  "いぬ[`41608e5f41`](https://gh/u/r/commit/41608e5f4109208a6ab995c58266554e6071c5b2)🐶",
		},
	}

	for _, tc := range tests {
		t.Run(tc.what, func(t *testing.T) {
			u := "https://gh/u/r"
			if tc.repoURL != "" {
				u = tc.repoURL
			}
			have := LinkRefs(tc.input, u)
			if have != tc.want {
				t.Fatalf("wanted %q but got %q", tc.want, have)
			}
		})
	}
}
