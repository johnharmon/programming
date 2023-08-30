import sys
sys.path.append('/home/jharmon/.local/lib/python3.9/site-packages')
import openai

openai.api_key_path = '../../../openai/api_key'

response = openai.Completion.create(
model="text-davinci-003",
prompt="write a jinja template that creates a mvn pom"
)

'''
import sys
sys.path.append('/home/jharmon/.local/lib/python3.9/site-packages')
import openai

openai.api_key_path = '../../../openai/api_key'



def main():






'''