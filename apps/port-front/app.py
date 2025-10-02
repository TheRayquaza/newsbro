import streamlit as st
from login.login import login_page
from register.register import register_page
from main_pages.admin_page import admin_page
import base64
import json

st.set_page_config(page_title="Newsbro", page_icon="📰", layout="wide")

# ------------------------------
# Helper functions for session persistence
# ------------------------------

def encode_session(email, token):
    """Encode session data to base64 for URL parameter"""
    data = {"email": email, "token": token}
    json_str = json.dumps(data)
    encoded = base64.b64encode(json_str.encode()).decode()
    return encoded

def decode_session(encoded_str):
    """Decode session data from base64 URL parameter"""
    try:
        decoded = base64.b64decode(encoded_str.encode()).decode()
        data = json.loads(decoded)
        return data.get("email"), data.get("token")
    except:
        return None, None

# ------------------------------
# Initialize session state from URL params first
# ------------------------------

if "session" in st.query_params:
    session_param = st.query_params["session"]
    if isinstance(session_param, list):
        session_param = session_param[0]
    
    email, token = decode_session(session_param)
    if email and token:
        st.session_state.authenticated = True
        st.session_state.login_user = email
        st.session_state.auth_token = token

# Initialize other session state variables
if "page" not in st.session_state:
    st.session_state.page = "login"
if "authenticated" not in st.session_state:
    st.session_state.authenticated = False
if "login_user" not in st.session_state:
    st.session_state.login_user = None
if "auth_token" not in st.session_state:
    st.session_state.auth_token = None
if "users" not in st.session_state:
    st.session_state.users = {"admin@admin.com": "password"}

# ------------------------------
# Helper function to change page
# ------------------------------

def set_page(page: str):
    st.session_state.page = page
    params = {"page": page}
    
    # Always include session data in URL if authenticated
    if st.session_state.authenticated and st.session_state.auth_token:
        session_data = encode_session(st.session_state.login_user, st.session_state.auth_token)
        params["session"] = session_data
    
    st.query_params.update(params)
    st.rerun()

# ------------------------------
# CSS styling
# ------------------------------

@st.cache_data
def load_css():
    s = ""
    with open("./styles.css") as f:
        s = f.read()
    return s

st.markdown(f"""
<style>
{load_css()}
</style>
""", unsafe_allow_html=True)

# ------------------------------
# Sync URL params with session state
# ------------------------------

if "page" in st.query_params:
    query_page = st.query_params["page"]
    if isinstance(query_page, list):
        query_page = query_page[0]
    st.session_state.page = query_page

# If authenticated, ensure session is in URL
if st.session_state.authenticated and st.session_state.auth_token:
    if "session" not in st.query_params:
        session_data = encode_session(st.session_state.login_user, st.session_state.auth_token)
        st.query_params["session"] = session_data

# ------------------------------
# Admin access control
# ------------------------------

if st.session_state.page == "admin":
    if not st.session_state.authenticated or not st.session_state.auth_token or st.session_state.login_user != "admin@admin.com":
        st.warning("⚠️ You must be an admin to access this page.")
        st.session_state.page = "login"
        st.query_params.clear()
        st.query_params["page"] = "login"
        st.rerun()

# ------------------------------
# Page rendering
# ------------------------------

if st.session_state.authenticated and st.session_state.auth_token:
    # Admin dashboard
    if st.session_state.login_user == "admin@admin.com":
        if st.session_state.page != "admin":
            st.session_state.page = "admin"
            set_page("admin")
        else:
            admin_page(set_page)
    else:
        # Normal user - show main news discovery page
        # Import main page functionality here to avoid circular imports
        import sys
        import os
        
        # Add necessary paths
        current_dir = os.path.dirname(os.path.abspath(__file__))
        if current_dir not in sys.path:
            sys.path.insert(0, current_dir)
        
        # Import and run the main page
        try:
            # Import main page components
            from main import render_main_page
            render_main_page(set_page)
        except ImportError:
            # Fallback if main page is structured differently
            st.markdown(f"### 🎉 Welcome, {st.session_state.login_user}!")
            st.info("📰 Loading News Discovery Platform...")
            
            # Show logout button
            if st.button("Logout", key="logout_main"):
                st.session_state.authenticated = False
                st.session_state.auth_token = None
                st.session_state.login_user = None
                st.query_params.clear()
                set_page("login")
else:
    # Login or register pages
    if st.session_state.page == "login":
        login_page(set_page)
    elif st.session_state.page == "register":
        register_page(set_page)
