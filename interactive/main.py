import streamlit as st
import subprocess
import urllib.request
import os

st.header("INTERPRETER")

dir_path = os.path.dirname(os.path.realpath(__file__))
bin_path = os.path.join(dir_path, "./interpreter")

if 'result' not in st.session_state:
    st.session_state['result'] = ''

def download_bin():
    if os.path.exists(bin_path):
        return
    bin_url = "https://github.com/labasubagia/interpreter/releases/download/latest/interpreter"
    with st.spinner(text="Download binary file"):
        urllib.request.urlretrieve(bin_url, bin_path)
        os.chmod(bin_path, 755)

def handle_run():
    if not os.path.exists(bin_path):
        return
    args = (bin_path, "string", code)
    popen = subprocess.Popen(args, stdout=subprocess.PIPE)
    popen.wait()
    if not popen.stdout:
        return
    val = popen.stdout.read().decode('utf-8')
    st.session_state['result'] = val if val else '~ No output'

download_bin()

example_dir = os.path.join(dir_path,'../example')
example_file = st.selectbox('Select Example', ['none'] + os.listdir(example_dir))

example_code = "let x = 5;\nx = 11;\nputs(x * x);"
if example_file and example_file != 'none':
    with open(os.path.join(example_dir, example_file), 'r') as file:
        example_code = file.read()

code = st.text_area("Code", example_code, height=300)
st.button('Run', on_click=handle_run)
st.text(st.session_state['result'])
