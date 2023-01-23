import unittest
import hashlib
from rolling_hash import rolling_hash_diff


class TestFileDiff(unittest.TestCase):
    def test_same_file(self):
        original_data = b"Same files"
        updated_data = b"Same files"
        delta = rolling_hash_diff(original_data, updated_data, hashlib.sha256)
        self.assertEqual(delta, {'modified': [], 'removed': [], 'reusables': [(0, b'Same files')]})

    def test_different_files(self):
        original_data = b"Same files"
        updated_data = b"Not Same files"
        delta = rolling_hash_diff(original_data, updated_data, hashlib.sha256)
        self.assertEqual(delta, {'modified': [(0, b'Not Same files')], 'removed': [(0, b'Same files')],  'reusables': []})

    def test_different_chunk_size(self):
        original_data = b"Same files"
        updated_data = b"Emas files"
        delta = rolling_hash_diff(original_data, updated_data, hashlib.sha256, chunk_size=5)
        self.assertEqual(delta, {'modified': [(0, b'Emas ')], 'removed': [(0, b'Same ')],  'reusables': [(1, b'files')]})

    def test_removed_chunk(self):
        original_data = b"Same sure files"
        updated_data = b"Same files"
        delta = rolling_hash_diff(original_data, updated_data, hashlib.sha256, chunk_size=5)
        self.assertEqual(delta['reusables'], [(0, b"Same ")])
        self.assertEqual(delta['modified'], [(1, b"files")])
        self.assertEqual(delta['removed'], [(1, b"sure ")])

    def test_added_chunk(self):
        original_data = b"Same files"
        updated_data = b"Same files, larger!!"
        delta = rolling_hash_diff(original_data, updated_data, hashlib.sha256, chunk_size=10)
        self.assertEqual(delta['reusables'], [(0, b"Same files")])
        self.assertEqual(delta['modified'], [(1, b", larger!!")])
        self.assertEqual(delta['removed'], [(1, b"")])


if __name__ == '__main__':
    unittest.main()
