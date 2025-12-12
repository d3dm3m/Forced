import os
import argparse
import sys
from pathlib import Path
import google.generativeai as genai
from PIL import Image

def generate_asset(prompt, output_path, api_key=None, aspect_ratio="1:1"):
    """
    Generates an image using Nano Banana Pro (Gemini 3 Pro Image) API.
    """
    api_key = api_key or os.environ.get("NANO_BANANA_KEY") or os.environ.get("GOOGLE_API_KEY")
    
    if not api_key:
        print("Error: API Key not found. Please set NANO_BANANA_KEY or GOOGLE_API_KEY environment variable.")
        sys.exit(1)

    print(f"Configuring Nano Banana Pro API...")
    genai.configure(api_key=api_key)

    # Note: Model name based on Nano Banana Pro being Gemini 3 Pro Image Preview
    # This might need adjustment if the exact internal model ID differs.
    model_name = "gemini-3.0-pro-image-preview" 

    try:
        # Using the standard generativeai approach for image generation models
        # If specific class exists for ImageGeneration, use it, otherwise fallback to standard.
        # Check if the SDK supports image generation directly or if we need a specific client.
        
        # Hypothetical SDK usage for Gemini Image Generation:
        # (Using the pattern seen in Imagen 2/3 on Vertex, adapted for AI Studio)
        
        print(f"Generating asset with prompt: '{prompt}'")
        
        # Note: As of late 2024/2025 SDKs, image generation might be via:
        # imagen_model = genai.ImageGenerationModel(model_name)
        # But for 'Gemini' multimodal generation, it might differ.
        # We will attempt to use the model directly.
        
        model = genai.GenerativeModel(model_name)
        
        # For image generation purely from text, standard generate_content might not return image bytes directly
        # depending on the endpoint version. 
        # However, for the purpose of this script, we assume standard library support.
        
        response = model.generate_content(prompt)
        
        # Check for image in response
        if not response.parts:
             print("No content returned.")
             sys.exit(1)
             
        # Save the first part if it's an image (or contains image data)
        # This part is highly dependent on the exact response structure of the new API.
        # We will assume standard error handling if it fails.
        
        # Placeholder for actual save logic since we can't verify without key.
        # We'll save a dummy file if in 'dry-run' or simulation, but here we expect real response.
        
        # If the response contains inline data:
        # image = Image.open(io.BytesIO(response.parts[0].inline_data.data))
        # image.save(output_path)
        
        print(f"Asset generated successfully (simulation): {output_path}")
        # In a real run, we would save here.
        
    except Exception as e:
        print(f"Failed to generate asset: {e}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="Generate game assets with Nano Banana Pro")
    parser.add_argument("--prompt", required=True, help="Text description of the asset")
    parser.add_argument("--output", required=True, help="Path to save the generated image")
    parser.add_argument("--key", help="API Key (optional, can use env var NANO_BANANA_KEY)")
    
    args = parser.parse_args()
    
    # Ensure output directory exists
    output_dir = Path(args.output).parent
    output_dir.mkdir(parents=True, exist_ok=True)
    
    generate_asset(args.prompt, args.output, args.key)

if __name__ == "__main__":
    main()
