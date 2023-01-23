import hashlib


def rolling_hash_diff(original_bytes: bytes, updated_bytes: bytes, hash_func: callable, chunk_size: int = 1024) -> dict:
    """
    Compares the original and updated version of a file using a rolling hash and returns a description of the chunks
    that can be reused and the chunks that have been added, modified, or removed.

    Parameters:
      original_bytes (bytes): The original data
      updated_bytes (bytes): The updated data
      hash_func (callable): A callable for creating a new rolling hash
      chunk_size (int): The size of the chunks, defaults to 1024 bytes

    Returns:
      dict: A dictionary containing the reusables chunks as a list of tuple (index, data) and the modified chunks
      as a list of tuple (index, data)
    """
    # initialize variables
    original_hash = rolling_hash(hash_func)
    updated_hash = rolling_hash(hash_func)
    delta = []
    reusables = []
    removed = []

    # split the data into chunks
    original_chunks = chunk_data(original_bytes, chunk_size)
    updated_chunks = chunk_data(updated_bytes, chunk_size)

    # Edge case: if original_chunks length is lower than updated, we need to fill blank chunks to effectively compare
    while len(original_chunks) < len(updated_chunks):
        original_chunks.append(b'')

    # Iterate over the original and updated hash objects
    idx = 0
    for original_chunk, updated_chunk in zip(original_chunks, updated_chunks):
        if original_chunk == updated_chunk:
            reusables.append((idx, original_chunk))
        else:
            delta.append((idx, updated_chunk))
            removed.append((idx, original_chunk))
        idx += 1
    return {"reusables": reusables, "modified": delta, "removed": removed}


# Helper function to create a new rolling hash based on a given hashing algorithm
def rolling_hash(hash_func):
    return hash_func()


# Helper function to split the data into chunks
def chunk_data(data, chunk_size):
    return [data[i:i+chunk_size] for i in range(0, len(data), chunk_size)]


# For debugging
if __name__ == '__main__':
    result1 = rolling_hash_diff(b"Testx", b"Testx", hashlib.sha256, chunk_size=1)
    print(result1)
    result2 = rolling_hash_diff(b"Testx", b"Testz", hashlib.sha256, chunk_size=1)
    print(result2)
    result3 = rolling_hash_diff(b"Testx", b"Tessz", hashlib.sha256, chunk_size=1)
    print(result3)

    with open('original.txt', 'rb') as f1, open('updated.txt', 'rb') as f2:
        original_data = f1.read()
        updated_data = f2.read()

        result = rolling_hash_diff(original_data, updated_data, hashlib.sha256, chunk_size=1)
        print("Comparison finished")
        print(f"Result: {result}")
