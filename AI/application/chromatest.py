import chromadb

# Connect to ChromaDB (client)
client = chromadb.PersistentClient(path="./resumes_collection")  # or chromadb.Client() for in-memory

# Select your collection
collection = client.get_collection(name="resumes")

# The job_id you want to fetch
target_job_id = "42d4dbb7-9f17-4362-b7b0-dcd292c5269f"

# Query ChromaDB filtering by metadata
results = collection.get(
    where={"job_id": target_job_id}  # Filter by metadata key
)

# Display results
for i in range(len(results["ids"])):
    print(f"ID: {results['ids'][i]}")
    print(f"Document: {results['documents'][i]}")
    print(f"Metadata: {results['metadatas'][i]}")
    print("-" * 40)