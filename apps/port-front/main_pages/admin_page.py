import streamlit as st
import requests
import validators

API_USERS = "https://account.newsbro.cc/api/v1/users"
API_ARTICLES = "https://article.newsbro.cc/api/v1/articles"
API_CATEGORIES = "https://article.newsbro.cc/api/v1/articles/categories"
API_SUBCATEGORIES = "https://article.newsbro.cc/api/v1/articles/subcategories"


def admin_page(set_page):
    st.markdown(
        '<div class="login-header"><h1>🛠 Admin Dashboard</h1>'
        '<p>Manage users and articles</p></div>',
        unsafe_allow_html=True
    )
    st.markdown('<div class="login-container">', unsafe_allow_html=True)

    # ----------------------------
    # Authentication check
    # ----------------------------
    token = st.session_state.get("auth_token")
    if not token:
        st.error("⚠️ No access token found. Please log in again.")
        if st.button("Go to Login"):
            st.session_state.authenticated = False
            set_page("login")
            st.experimental_set_query_params(page="login")
            st.rerun()
        st.markdown('</div>', unsafe_allow_html=True)
        return

    headers = {"Authorization": f"Bearer {token}"}

    # ----------------------------
    # Tabs
    # ----------------------------
    tabs = st.tabs(["Users", "Articles", "Feedback"])

    # ----------------------------
    # Users tab
    # ----------------------------
    with tabs[0]:
        st.subheader("All Users")
        try:
            resp = requests.get(API_USERS, headers=headers, timeout=10)
            if resp.status_code == 200:
                users = resp.json()
                user_table = [{"ID": u['id'], "Username": u['username'], "Email": u['email'], "Active": u['is_active']} for u in users]
                if user_table:
                    st.table(user_table)
                else:
                    st.info("No users found.")
            else:
                st.error(f"Failed to fetch users ({resp.status_code}): {resp.text}")
        except Exception as e:
            st.error(f"Error fetching users: {e}")

    # ----------------------------
    # Articles tab
    # ----------------------------
    with tabs[1]:
        article_tabs = st.tabs(["Browse Articles", "Create Article", "Update Article", "Delete Article"])

        # Browse Articles
        with article_tabs[0]:
            st.subheader("Browse Articles")
            categories = get_categories(headers)
            subcategories = get_subcategories(headers)

            category = st.selectbox("Select category", ["All"] + categories)
            subcategory = st.selectbox("Select subcategory", ["All"] + subcategories)

            params = {}
            if category != "All":
                params["category"] = category
            if subcategory != "All":
                params["subcategory"] = subcategory

            try:
                resp = requests.get(API_ARTICLES, headers=headers, params=params, timeout=10)
                if resp.status_code == 200:
                    articles_resp = resp.json()
                    articles_list = articles_resp.get("data", articles_resp if isinstance(articles_resp, list) else [])
                    for a in articles_list:
                        st.write(f"**{a.get('title', 'No Title')}** (ID: {a.get('id', 'N/A')}) - {a.get('summary', '')}")
                else:
                    st.error(f"Failed to fetch articles ({resp.status_code})")
            except Exception as e:
                st.error(f"Error fetching articles: {e}")

        # Create Article
        with article_tabs[1]:
            st.subheader("Create New Article")
            with st.form("create_article_form"):
                title = st.text_input("Title")
                summary = st.text_area("Summary")
                content = st.text_area("Content")

                categories = get_categories(headers)
                category_create = st.selectbox("Category", ["Select a category"] + categories, key="create_cat")
                if category_create == "Select a category":
                    category_create = None

                subcategories = get_subcategories(headers)
                subcategory_create = st.selectbox("Subcategory", ["Select a subcategory"] + subcategories, key="create_subcat")
                if subcategory_create == "Select a subcategory":
                    subcategory_create = None

                link = st.text_input("Link")

                submit = st.form_submit_button("Create Article")
                if submit:
                    if not title or not category_create or not subcategory_create or not validators.url(link):
                        st.error("All fields are required and Link must be a valid URL.")
                    else:
                        payload = {
                            "title": title,
                            "summary": summary,
                            "content": content,
                            "category": category_create,
                            "subcategory": subcategory_create,
                            "link": link
                        }
                        try:
                            resp = requests.post(API_ARTICLES, headers=headers, json=payload, timeout=10)
                            if resp.status_code in [200, 201]:
                                st.success("✅ Article created successfully!")
                            else:
                                st.error(f"Failed to create article ({resp.status_code}): {resp.text}")
                        except Exception as e:
                            st.error(f"Error creating article: {e}")

        # Update Article
        with article_tabs[2]:
            st.subheader("Update Article")
            article_id_update = st.text_input("Article ID to update")

            if article_id_update:
                with st.form("update_article_form"):
                    title_update = st.text_input("Title")
                    summary_update = st.text_area("Summary")
                    content_update = st.text_area("Content")

                    categories = get_categories(headers)
                    category_update = st.selectbox("Category", ["Select a category"] + categories, key="update_cat")
                    if category_update == "Select a category":
                        category_update = None

                    subcategories = get_subcategories(headers)
                    subcategory_update = st.selectbox("Subcategory", ["Select a subcategory"] + subcategories, key="update_subcat")
                    if subcategory_update == "Select a subcategory":
                        subcategory_update = None

                    link_update = st.text_input("Link")

                    submit_update = st.form_submit_button("Update Article")
                    if submit_update:
                        if not title_update or not category_update or not subcategory_update or not validators.url(link_update):
                            st.error("All fields are required and Link must be a valid URL.")
                        else:
                            payload = {
                                "title": title_update,
                                "summary": summary_update,
                                "content": content_update,
                                "category": category_update,
                                "subcategory": subcategory_update,
                                "link": link_update
                            }
                            try:
                                resp = requests.put(f"{API_ARTICLES}/{article_id_update}", headers=headers, json=payload, timeout=10)
                                if resp.status_code == 200:
                                    st.success("✅ Article updated successfully!")
                                else:
                                    st.error(f"Failed to update article ({resp.status_code}): {resp.text}")
                            except Exception as e:
                                st.error(f"Error updating article: {e}")

        # Delete Article
        with article_tabs[3]:
            st.subheader("Delete Article")
            article_id_delete = st.text_input("Article ID to delete", key="delete_id")
            if st.button("Delete Article", type="primary"):
                if article_id_delete:
                    try:
                        resp = requests.delete(f"{API_ARTICLES}/{article_id_delete}", headers=headers, timeout=10)
                        if resp.status_code in [200, 204]:
                            st.success("✅ Article deleted successfully!")
                        else:
                            st.error(f"Failed to delete article ({resp.status_code}): {resp.text}")
                    except Exception as e:
                        st.error(f"Error deleting article: {e}")
                else:
                    st.warning("Please enter an article ID")

    # ----------------------------
    # Feedback tab
    # ----------------------------
    with tabs[2]:
        st.subheader("Feedback Stats (Admin)")
        try:
            resp = requests.get("https://article.newsbro.cc/api/v1/feedback/stats", headers=headers)
            if resp.status_code == 200:
                stats = resp.json()
                st.json(stats)
            else:
                st.error(f"Failed to fetch feedback stats ({resp.status_code})")
        except Exception as e:
            st.error(f"Error fetching feedback stats: {e}")

    st.markdown('</div>', unsafe_allow_html=True)

    # ----------------------------
    # Logout button
    # ----------------------------
    if st.button("Logout"):
        # Clear all authentication data
        st.session_state.authenticated = False
        st.session_state.login_user = None
        st.session_state.auth_token = None
        
        # Force page to login
        st.session_state.page = "login"
        st.query_params = {"page": "login"}
        
        # Stop current run to trigger immediate rerun
        st.stop()

# ----------------------------
# Helper functions
# ----------------------------
def get_categories(headers):
    try:
        resp = requests.get(API_CATEGORIES, headers=headers, timeout=10)
        if resp.status_code == 200:
            data = resp.json()
            if isinstance(data, dict) and "data" in data:
                return [c['name'] for c in data['data']]
            elif isinstance(data, list):
                return [c['name'] for c in data]
    except Exception as e:
        st.error(f"Error fetching categories: {e}")
    return []


def get_subcategories(headers):
    try:
        resp = requests.get(API_SUBCATEGORIES, headers=headers, timeout=10)
        if resp.status_code == 200:
            data = resp.json()
            if isinstance(data, dict) and "data" in data:
                return [c['name'] for c in data['data']]
            elif isinstance(data, list):
                return [c['name'] for c in data]
    except Exception as e:
        st.error(f"Error fetching subcategories: {e}")
    return []
