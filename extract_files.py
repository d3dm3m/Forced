import os
import re

def extract_files(md_file):
    with open(md_file, 'r', encoding='utf-8') as f:
        content = f.read()

    # Regex to find file blocks
    # Format: ### ./path/to/file
    # ```lang
    # content
    # ```

    # We will split by "### ./" and then process each chunk
    chunks = content.split("### ./")

    # Skip the first chunk (header stuff)
    for chunk in chunks[1:]:
        # First line is the filepath
        lines = chunk.splitlines()
        filepath = lines[0].strip()

        # Determine content
        # Find the first ``` and the last ```
        # Note: Code blocks might be indented or just start.
        # But usually in the provided format they are strictly formatted.

        # Check if binary/asset file placeholder
        if "*Binary/Asset File*" in chunk:
            print(f"Skipping binary file: {filepath}")
            continue

        # Extract code block
        try:
            start_idx = chunk.find("```")
            if start_idx == -1:
                print(f"No code block found for {filepath}")
                continue

            # Find end of start line (e.g. ```go)
            code_start_newline = chunk.find("\n", start_idx)

            # Find closing ```
            end_idx = chunk.rfind("```")

            if end_idx <= start_idx:
                print(f"Malformed block for {filepath}")
                continue

            file_content = chunk[code_start_newline+1:end_idx]

            # Remove trailing newline if it looks like artifact?
            # Usually keep as is.

            # Write file
            # Ensure directory exists
            dirpath = os.path.dirname(filepath)
            if dirpath and not os.path.exists(dirpath):
                os.makedirs(dirpath)

            # If file already exists, I should check if I should overwrite.
            # The plan says "Restore Missing Files".
            # If I already moved it from Force, maybe I shouldn't overwrite?
            # Or maybe PROJECT_CONTEXT.md is the source of truth?
            # The user said "restore them correctly", "Force" directory has files.
            # I should probably prioritize "Force" files if they exist, and only write if not exists.

            if os.path.exists(filepath):
                print(f"File exists, skipping: {filepath}")
            else:
                with open(filepath, 'w', encoding='utf-8') as out:
                    out.write(file_content)
                print(f"Restored: {filepath}")

        except Exception as e:
            print(f"Error processing {filepath}: {e}")

if __name__ == "__main__":
    extract_files("PROJECT_CONTEXT.md")
