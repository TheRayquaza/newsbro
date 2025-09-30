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

    # --- Input fields ---
    email = st.text_input("Email", key="login_email_input")
    password = st.text_input("Password", type="password", key="login_pass_input")

    col1, col2 = st.columns(2)

    # --- Login button ---
    with col1:
        if st.button("Login"):
            if email and password:
                try:
                    resp = requests.post(
                        API_LOGIN,
                        json={"email": email, "password": password},
                        timeout=5
                    )
                    if resp.status_code == 200:
                        data = resp.json()
                        # Save authenticated state and token
                        st.session_state.authenticated = True
                        st.session_state.login_user = email
                        st.session_state.auth_token = data.get("token")  # assuming API returns {"token": "..."}
                        st.success(f"Welcome, {email}!")
                        st.rerun()
                    else:
                        # Show detailed error if available
                        try:
                            err = resp.json().get("message", "")
                            st.error(f"Login failed ({resp.status_code}): {err}")
                        except:
                            st.error(f"Login failed ({resp.status_code})")
                except Exception as e:
                    st.error(f"Error: {e}")
            else:
                st.error("Please enter email and password")

    # --- Register button ---
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

# Example usage:
# if 'authenticated' not in st.session_state:
#     st.session_state.authenticated = False
# 
# if not st.session_state.authenticated:
#     login_page(lambda page: st.session_state.update(page=page))
# else:
#     st.write("You are logged in!")
