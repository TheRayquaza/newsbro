import streamlit as st
import requests
import os

REGISTER_API_URL = f"{os.getenv('ACCOUNT_BASE_URL', 'https://account.newsbro.cc')}/api/v1/auth/register"

@st.cache_data
def register_css():
    s = ""
    cw = os.path.dirname(os.path.realpath(__file__))
    print("Current working directory:", cw)
    with open("./register/register.css") as f:
        s = f.read()
    return s

def register_page(set_page):
    st.markdown(
        f"""
        <style>
            {register_css()}
        </style>
        """,
        unsafe_allow_html=True
    )
    
    # Initialize session state for inputs
    if "reg_email_input" not in st.session_state:
        st.session_state.reg_email_input = ""
    if "reg_first_name_input" not in st.session_state:
        st.session_state.reg_first_name_input = ""
    if "reg_last_name_input" not in st.session_state:
        st.session_state.reg_last_name_input = ""
    if "reg_user_input" not in st.session_state:
        st.session_state.reg_user_input = ""
    if "reg_pass_input" not in st.session_state:
        st.session_state.reg_pass_input = ""
    
    # Create centered column layout
    col1, col2, col3 = st.columns([1, 1.5, 1])
    
    with col2:
        with st.container():
            # Header
            st.markdown("# 📰 NewsBro")
            st.markdown('<p class="subtitle">Create your account to get started</p>', unsafe_allow_html=True)
            
            # Input fields
            email = st.text_input("Email", placeholder="Enter your email", key="reg_email_input")
            
            # First and Last name in two columns
            name_col1, name_col2 = st.columns(2)
            with name_col1:
                first_name = st.text_input("First Name", placeholder="First name", key="reg_first_name_input")
            with name_col2:
                last_name = st.text_input("Last Name", placeholder="Last name", key="reg_last_name_input")
            
            username = st.text_input("Username", placeholder="Choose a username", key="reg_user_input")
            password = st.text_input("Password", type="password", placeholder="Create a password", key="reg_pass_input")
            
            col1, col2, col3 = st.columns([2, 2, 1])

            with col2:
                # Centered Register button
                if st.button("Register", key="register_btn"):
                    if not all([email, first_name, last_name, username, password]):
                        st.error("⚠️ Please fill in all fields")
                    else:
                        try:
                            payload = {
                                "email": email,
                                "first_name": first_name,
                                "last_name": last_name,
                                "username": username,
                                "password": password
                            }
                            response = requests.post(REGISTER_API_URL, json=payload, timeout=10)
                            if response.status_code == 201:
                                st.success("✅ Registration successful! Please log in.")
                                set_page("login")
                                st.rerun()
                            else:
                                try:
                                    st.error(f"❌ Registration failed: {response.json()['error']}")
                                except:
                                    st.error(f"❌ Registration failed: {response.text}")
                        except requests.exceptions.RequestException as e:
                            st.error(f"⚠️ Network error: {e}")

            # Divider
            st.divider()

            # Fallback button (hidden) for actual functionality
            if st.button("↩ Back to Login", key="back_btn", help="Return to login page"):
                set_page("login")
                st.rerun()
