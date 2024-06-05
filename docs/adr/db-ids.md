# Title: Decision Between Binary and Hexadecimal IDs for Database Entries

# Date: 2024-06-05

# Status: Accepted

# Context:

The project needs to generate unique identifiers (IDs) for database entries.
The options considered include using binary or hexadecimal formats for these IDs.
The IDs need to be unique, easy to generate, and efficient to store and retrieve from the database.

# Assumptions:

- The project prioritizes storage efficiency and performance over human readability for IDs;

# Decision:

The project will use binary format for the IDs.

# Consequences:

|Positive|Negatives|
|--|--|
|Storage Efficiency: Binary IDs are more space-efficient. Each byte in binary is directly stored, reducing the overall size of the data. |Readability: Binary IDs are not human-readable, making them difficult to debug and log without conversion tools.
|Performance: Using binary can reduce the overhead associated with encoding and decoding, leading to potentially faster operations in databases that support binary formats natively.|Tooling and Integration: Some libraries and tools may have limited support for binary data, requiring additional handling or custom solutions.|
|Consistency: Binary representation ensures consistency in the size and structure of the IDs, which can be beneficial for indexing and retrieval in the database.|URL Safety: Binary data needs to be encoded (e.g., Base64) to be safely used in URLs, adding an extra step for web-based applications.|

# Alternatives Considered:

## Hexadecimal:
|Pros|Cons|
|--|--|
|More human-readable and easier to debug. Hexadecimal is also URL-safe and widely supported by libraries and tools.|Less storage-efficient, as each byte of binary data is represented by two hex characters, increasing the storage requirements.|

