"""Regression coverage for static web asset references."""

import unittest

from check_web_links import WEB_ROOT, normalize, target_for


class WebLinkTests(unittest.TestCase):
    def test_decodes_encoded_local_path_after_stripping_query(self):
        self.assertEqual(normalize("images/test%20signal.svg?v=2#plot"), "images/test signal.svg")

    def test_preserves_encoded_hash_in_filename(self):
        self.assertEqual(normalize("images/gain%23plot.svg#caption"), "images/gain#plot.svg")

    def test_skips_external_and_non_file_links_case_insensitively(self):
        for value in ("HTTPS://example.org", "MAILTO:person@example.org",
                      "tel:+12345678", "SMS:+12345678", "//cdn.example.org/lib.js",
                      "data:image/png;base64,aabb", "#section"):
            with self.subTest(value=value):
                self.assertIsNone(normalize(value))

    def test_resolves_local_assets_under_web_root(self):
        source = WEB_ROOT / "reports" / "index.php"
        self.assertEqual(
            target_for(source, "../img/logo.svg"),
            (WEB_ROOT / "img" / "logo.svg").resolve(),
        )
        self.assertEqual(
            target_for(source, "/img/logo.svg"),
            (WEB_ROOT / "img" / "logo.svg").resolve(),
        )

    def test_rejects_paths_escaping_web_root(self):
        source = WEB_ROOT / "index.php"
        escaped = normalize("%2e%2e/private.txt")
        self.assertIsNotNone(escaped)
        self.assertIsNone(target_for(source, escaped))


if __name__ == "__main__":
    unittest.main()
