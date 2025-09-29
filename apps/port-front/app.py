import streamlit as st
from login.login import login_page
from login.register import register_page

# ------------------------------
# Page configuration (must be first)
# ------------------------------
st.set_page_config(page_title="Academic Paper Platform", page_icon="🎓", layout="centered")

# ------------------------------
# Session state initialization
# ------------------------------
if "page" not in st.session_state:
    st.session_state.page = "login"
if "authenticated" not in st.session_state:
    st.session_state.authenticated = False
if "users" not in st.session_state:
    st.session_state.users = {"admin": "password"}

# ------------------------------
# CSS styling (remove extra white space)
# ------------------------------
st.markdown("""
<style>
    /* Remove top padding */
    .css-18e3th9 {padding-top: 0rem !important;}
    
    .login-header {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        padding: 30px 20px;
        border-radius: 20px;
        color: white;
        text-align: center;
        margin-bottom: 20px;
    }

    .login-container {
        background: #f8f9fa;
        padding: 20px 20px;
        border-radius: 16px;
        max-width: 400px;
        margin: auto;
        box-shadow: 0 2px 12px rgba(60, 64, 67, 0.1);
    }

    .stButton > button {
        border-radius: 12px;
        border: none;
        background: #1976d2;
        color: white;
        font-weight: 500;
        padding: 10px 20px;
        font-size: 1em;
        transition: all 0.2s ease;
        width: 100%;
    }

    .stButton > button:hover {
        background: #1565c0;
        transform: translateY(-1px);
    }

    .secondary-btn > button {
        background: #f5f5f5 !important;
        color: #333 !important;
        border: 1px solid #ddd !important;
    }

    .secondary-btn > button:hover {
        background: #eeeeee !important;
        color: #000 !important;
    }
</style>
""", unsafe_allow_html=True)

# ------------------------------
# Page navigation
# ------------------------------
if st.session_state.authenticated:
    st.markdown('<div class="login-container">', unsafe_allow_html=True)
    st.markdown(f"### 🎉 Welcome, {st.session_state.get('login_user', 'User')}!")
    if st.button("Logout"):
        st.session_state.authenticated = False
        st.session_state.page = "login"
    st.markdown('</div>', unsafe_allow_html=True)
else:
    if st.session_state.page == "login":
        login_page()
    elif st.session_state.page == "register":
        register_page()
