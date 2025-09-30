import streamlit as st
from login.login import login_page
from login.register import register_page
from main_pages.admin_page import admin_page

st.set_page_config(page_title="Academic Paper Platform", page_icon="🎓", layout="centered")

if "page" not in st.session_state:
    st.session_state.page = "login"
if "authenticated" not in st.session_state:
    st.session_state.authenticated = False
if "users" not in st.session_state:
    st.session_state.users = {"admin": "password"}

query_params = st.query_params

if "page" in query_params:
    st.session_state.page = query_params["page"]

def set_page(page: str):
    st.session_state.page = page
    st.query_params["page"] = page
    st.rerun()

# ------------------------------
# CSS styling
# ------------------------------
st.markdown("""
<style>
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
    if st.session_state.get("login_user") == "admin@admin.com":
        admin_page(set_page)
    else:
        st.markdown('<div class="login-container">', unsafe_allow_html=True)
        st.markdown(f"### 🎉 Welcome, {st.session_state.get('login_user', 'User')}!")
        if st.button("Logout"):
            st.session_state.authenticated = False
            set_page("login")
        st.markdown('</div>', unsafe_allow_html=True)

else:
    if st.session_state.page == "login":
        login_page(set_page)
    elif st.session_state.page == "register":
        register_page(set_page)
