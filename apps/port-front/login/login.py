import streamlit as st
import requests

API_LOGIN = "https://account.newsbro.cc/api/v1/auth/login"
FORGE_OAUTH = "https://account.newsbro.cc/api/v1/auth/oauth/login"

def login_page(set_page):
    st.markdown(
        '<div class="login-header"><h1>🎓 Academic Paper Platform</h1>'
        '<p>Please log in to access the platform</p></div>',
        unsafe_allow_html=True
    )
    st.markdown('<div class="login-container">', unsafe_allow_html=True)

    username = st.text_input("Username", key="login_user_input")
    password = st.text_input("Password", type="password", key="login_pass_input")

    col1, col2 = st.columns(2)
    with col1:
        if st.button("Login"):
            if username and password:
                try:
                    resp = requests.post(API_LOGIN, json={"username": username, "password": password}, timeout=5)
                    if resp.status_code == 200:
                        st.session_state.authenticated = True
                        st.session_state.login_user = username
                        st.success(f"Welcome, {username}!")
                        st.rerun()
                    else:
                        st.error(f"Login failed ({resp.status_code})")
                except Exception as e:
                    st.error(f"Error: {e}")
            else:
                st.error("Please enter username and password")
    with col2:
        if st.button("Register"):
            set_page("register")
            st.rerun()

    # --- Forge login ---
    st.markdown("---")
    if st.button("🔑 Connect via Forge"):
        st.markdown(
            f'<meta http-equiv="refresh" content="0; url={FORGE_OAUTH}" />',
            unsafe_allow_html=True
        )

    st.markdown('</div>', unsafe_allow_html=True)
