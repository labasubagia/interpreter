import streamlit as st
import subprocess
import os

dir_path = os.path.dirname(os.path.realpath(__file__))

if 'result' not in st.session_state:
    st.session_state['result'] = ''


def handle_run():
    bin_path = os.path.join(dir_path, "./interpreter")
    args = (bin_path, "string", code)
    popen = subprocess.Popen(args, stdout=subprocess.PIPE)
    popen.wait()
    if popen.stdout:
        val = popen.stdout.read().decode('utf-8')
        st.session_state['result'] = val if val else '~ No output'

st.header("INTERPRETER")

example_dir = os.path.join(dir_path,'../example')
example_file = st.selectbox('Select Example', ['none'] + os.listdir(example_dir))

example_code = "let x = 5;\nx = 11;\nputs(x * x);"
if example_file and example_file != 'none':
    with open(os.path.join(example_dir, example_file), 'r') as file:
        example_code = file.read()

code = st.text_area("Code", example_code, height=300)
st.button('Run', on_click=handle_run)
st.text(st.session_state['result'])