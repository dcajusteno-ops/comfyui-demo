import json
import random

input_file = r'F:\demo\comfyui-demo-main\data\prompt-library\imported_prompts.json'
with open(input_file, 'r', encoding='utf-8') as f:
    prompts = json.load(f)

styles = [p for p in prompts if p.get('category') == '风格']
print(f'Total in 风格: {len(styles)}')
for p in random.sample(styles, min(50, len(styles))):
    print(f\"{p.get('text_en')} - {p.get('text_zh')}\")
