import streamlit as st
import requests

API_URL = "https://account.newsbro.cc/api/v1/auth/register"

def register_page(set_page):
    st.markdown(
        '<div class="login-header"><h1>🎓 Academic Paper Platform</h1>'
        '<p>Create a new account</p></div>',
        unsafe_allow_html=True
    )
    st.markdown('<div class="login-container">', unsafe_allow_html=True)

    # --- Initialize session state for inputs ---
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

    # --- Input fields ---
    email = st.text_input("Email", key="reg_email_input")
    first_name = st.text_input("First Name", key="reg_first_name_input")
    last_name = st.text_input("Last Name", key="reg_last_name_input")
    username = st.text_input("Username", key="reg_user_input")
    password = st.text_input("Password", type="password", key="reg_pass_input")

    col1, col2 = st.columns(2)

    # --- Register button ---
    with col1:
        if st.button("Register"):
            if not all([email, first_name, last_name, username, password]):
                st.error("Please fill all fields")
            else:
                try:
                    payload = {
                        "email": email,
                        "first_name": first_name,
                        "last_name": last_name,
                        "username": username,
                        "password": password
                    }
                    response = requests.post(API_URL, json=payload, timeout=10)
                    if response.status_code == 201:  # success
                        st.success("Registration successful! Please log in.")
                        set_page("login")
                        st.rerun()
                    else:
                        # Show API error message if available
                        try:
                            error_msg = response.json()
                        except Exception:
                            error_msg = response.text
                        st.error(f"Registration failed: {error_msg}")
                except requests.exceptions.RequestException as e:
                    st.error(f"Network error: {e}")

    # --- Back to login button ---
    with col2:
        if st.button("Back to Login"):
            set_page("login")
            st.rerun()

    st.markdown('</div>', unsafe_allow_html=True)
