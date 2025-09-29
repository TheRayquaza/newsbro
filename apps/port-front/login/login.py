import streamlit as st

def login_page():
    st.markdown('<div class="login-header"><h1>🎓 Academic Paper Platform</h1><p style="font-size:16px;opacity:0.9;">Please log in to access the platform</p></div>', unsafe_allow_html=True)
    st.markdown('<div class="login-container">', unsafe_allow_html=True)

    # Use session_state to fix double-click bug
    if "login_user_input" not in st.session_state:
        st.session_state.login_user_input = ""
    if "login_pass_input" not in st.session_state:
        st.session_state.login_pass_input = ""

    username = st.text_input("Username", key="login_user_input")
    password = st.text_input("Password", type="password", key="login_pass_input")

    col1, col2 = st.columns(2)
    with col1:
        if st.button("Login"):
            if username in st.session_state.users and st.session_state.users[username] == password:
                st.session_state.authenticated = True
                st.session_state.login_user = username
                st.success(f"Welcome, {username}!")
            else:
                st.error("Invalid username or password")
    with col2:
        if st.button("Register"):
            st.session_state.page = "register"

    st.markdown('</div>', unsafe_allow_html=True)
