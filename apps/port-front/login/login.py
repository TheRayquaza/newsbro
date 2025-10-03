import streamlit as st
import requests
import os

API_LOGIN = f"{os.getenv('API_BASE_URL', 'https://account.newsbro.cc')}/api/v1/auth/login"
FORGE_OAUTH = f"{os.getenv('API_BASE_URL', 'https://account.newsbro.cc')}/api/v1/auth/oauth/login"

@st.cache_data
def login_css():
    s = ""
    with open("./login/login.css") as f:
        s = f.read()
    return s

def login_page(set_page):
    st.markdown(
        f"""
        <style>
            {login_css()}
        </style>
        """,
        unsafe_allow_html=True
    )
    
    # Create centered column layout
    col1, col2, col3 = st.columns([1, 1.5, 1])
    
    with col2:
        with st.container():
            # Header
            st.markdown("# 📰 NewsBro")
            st.markdown('<p class="subtitle">Welcome back! Please login to continue</p>', unsafe_allow_html=True)
            
            # Input fields
            email = st.text_input("Email", placeholder="Enter your email", key="login_email_input")
            password = st.text_input("Password", type="password", placeholder="Enter your password", key="login_pass_input")
            
            # Centered Login button
            if st.button("Login", key="login_btn"):
                if email and password:
                    try:
                        resp = requests.post(
                            API_LOGIN,
                            json={"email": email, "password": password},
                            timeout=5
                        )
                        if resp.status_code == 200:
                            data = resp.json()
                            token = data.get("access_token") or data.get("token")
                            if token:
                                st.session_state.authenticated = True
                                st.session_state.login_user = email
                                st.session_state.auth_token = token
                                if email == "admin@admin.com":
                                    set_page("admin")
                                else:
                                    set_page("home")
                                st.rerun()
                            else:
                                st.error("🔒 Login succeeded but no token received")
                        else:
                            try:
                                err = resp.json().get("message", "Invalid credentials")
                                st.error(f"❌ {err}")
                            except Exception:
                                st.error(f"❌ Login failed ({resp.status_code})")
                    except Exception as e:
                        st.error(f"⚠️ Connection error: {str(e)}")
                else:
                    st.error("⚠️ Please enter both email and password")
            
            # Divider
            st.markdown('<div class="divider">OR</div>', unsafe_allow_html=True)
            
            # Forge OAuth button
            st.markdown(
                f'''
                <a href="{FORGE_OAUTH}" class="forge-button" target="_blank">
                    <img src="https://docs.forge.epita.fr/assets/images/forge_logo_white-38d02f6ef9b56a34f60b46030ee8b0de.svg" alt="Forge"/>
                    <span>Continue with ForgeID</span>
                </a>
                ''',
                unsafe_allow_html=True
            )

            # Divider
            st.divider()
            
            if st.button("Create an account", key="register_btn", help="Sign up for a new account"):
                set_page("register")
                st.rerun()