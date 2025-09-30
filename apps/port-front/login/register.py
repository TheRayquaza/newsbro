import streamlit as st

def register_page(set_page):
    st.markdown('<div class="login-header"><h1>🎓 Academic Paper Platform</h1><p>Create a new account</p></div>', unsafe_allow_html=True)
    st.markdown('<div class="login-container">', unsafe_allow_html=True)

    if "reg_user_input" not in st.session_state:
        st.session_state.reg_user_input = ""
    if "reg_pass_input" not in st.session_state:
        st.session_state.reg_pass_input = ""

    username = st.text_input("New Username", key="reg_user_input")
    password = st.text_input("New Password", type="password", key="reg_pass_input")

    col1, col2 = st.columns(2)
    with col1:
        if st.button("Register"):
            if username in st.session_state.users:
                st.error("Username already exists")
            elif username == "" or password == "":
                st.error("Please fill all fields")
            else:
                st.session_state.users[username] = password
                st.success("Registration successful! Please log in.")
                set_page("login")
                st.rerun() 
    with col2:
        if st.button("Back to Login"):
            set_page("login")
            st.rerun()

    st.markdown('</div>', unsafe_allow_html=True)
