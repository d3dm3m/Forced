import os
import argparse
import sys
import time
import requests
from pathlib import Path

MESHY_API_URL = "https://api.meshy.ai/openapi/v2/text-to-3d"

def generate_asset_meshy(prompt, output_dir, api_key=None):
    """
    Generates a 3D asset using Meshy API (Preview -> Refine workflow).
    """
    api_key = api_key or os.environ.get("MESHY_API_KEY")
    
    if not api_key:
        print("Error: API Key not found. Please set MESHY_API_KEY environment variable.")
        sys.exit(1)

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json"
    }

    # 1. Create Preview Task
    print(f"Starting Meshy Preview for: '{prompt}'")
    payload = {
        "mode": "preview",
        "prompt": prompt,
        "art_style": "realistic", # Default to realistic, can be parameterized
        "negative_prompt": "low quality, low resolution, low poly, ugly, blurring, pixelated" 
    }
    
    try:
        response = requests.post(MESHY_API_URL, headers=headers, json=payload)
        response.raise_for_status()
        task_data = response.json()
        preview_task_id = task_data.get("result")
        print(f"Preview Task ID: {preview_task_id}")
    except Exception as e:
        print(f"Failed to create preview task: {e}")
        if hasattr(e, 'response') and e.response:
             print(e.response.text)
        sys.exit(1)

    # 2. Poll for Preview Completion
    print("Waiting for Preview...")
    while True:
        try:
            status_response = requests.get(f"https://api.meshy.ai/openapi/v2/text-to-3d/{preview_task_id}", headers=headers)
            status_response.raise_for_status()
            status_data = status_response.json()
            status = status_data.get("status")
            progress = status_data.get("progress", 0)
            
            if status == "SUCCEEDED":
                print("Preview Complete.")
                break
            elif status in ["FAILED", "EXPIRED"]:
                print(f"Preview Failed with status: {status}")
                sys.exit(1)
            
            print(f"Preview Progress: {progress}%")
            time.sleep(2)
        except Exception as e:
            print(f"Error polling preview: {e}")
            time.sleep(5)

    # 3. Create Refine Task
    print("Starting Refine Task...")
    refine_payload = {
        "mode": "refine",
        "preview_task_id": preview_task_id
    }
    
    try:
        response = requests.post(MESHY_API_URL, headers=headers, json=refine_payload)
        response.raise_for_status()
        refine_task_data = response.json()
        refine_task_id = refine_task_data.get("result")
        print(f"Refine Task ID: {refine_task_id}")
    except Exception as e:
         print(f"Failed to create refine task: {e}")
         if hasattr(e, 'response') and e.response:
             print(e.response.text)
         sys.exit(1)

    # 4. Poll for Refine Completion
    print("Waiting for Texturing (Refine)...")
    final_model_url = None
    
    while True:
        try:
            status_response = requests.get(f"https://api.meshy.ai/openapi/v2/text-to-3d/{refine_task_id}", headers=headers)
            status_response.raise_for_status()
            status_data = status_response.json()
            status = status_data.get("status")
            progress = status_data.get("progress", 0)
            
            if status == "SUCCEEDED":
                print("Refine Complete.")
                model_urls = status_data.get("model_urls", {})
                final_model_url = model_urls.get("glb") # Prefer GLB
                break
            elif status in ["FAILED", "EXPIRED"]:
                print(f"Refine Failed with status: {status}")
                sys.exit(1)
            
            print(f"Refine Progress: {progress}%")
            time.sleep(2)
        except Exception as e:
            print(f"Error polling refine: {e}")
            time.sleep(5)

    # 5. Download Model
    if final_model_url:
        print(f"Downloading model from: {final_model_url}")
        try:
            model_response = requests.get(final_model_url)
            model_response.raise_for_status()
            
            filename = f"{prompt.replace(' ', '_')[:20]}_{int(time.time())}.glb"
            output_file = Path(output_dir) / filename
            output_file.parent.mkdir(parents=True, exist_ok=True)
            
            with open(output_file, 'wb') as f:
                f.write(model_response.content)
            
            print(f"Asset saved to: {output_file}")
            
        except Exception as e:
            print(f"Failed to download model: {e}")
            sys.exit(1)
    else:
        print("No model URL found in response.")

def main():
    parser = argparse.ArgumentParser(description="Generate 3D assets with Meshy AI")
    parser.add_argument("--prompt", required=True, help="Text description of the 3D asset")
    parser.add_argument("--output_dir", required=True, help="Directory to save the generated asset")
    parser.add_argument("--key", help="API Key (optional, can use env var MESHY_API_KEY)")
    
    args = parser.parse_args()
    
    generate_asset_meshy(args.prompt, args.output_dir, args.key)

if __name__ == "__main__":
    main()
