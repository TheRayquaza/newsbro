import streamlit as st
import numpy as np
import pandas as pd
from preprocessing import preprocessing
import re
import os
import sys
import base64
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Import the recommendation models
sys.path.append("../")
sys.path.append("../models")
from models.SBertModel import SBertModel
from models.TF_IDFModel import TfidfModel

# Custom CSS for better styling
MAIN_PAGE_CSS = """
<style>
    /* Main container styling */
    .main {
        padding-top: 1rem;
    }
    
    /* Custom card styling */
    .news-card {
        border: 1px solid #e8eaed;
        border-radius: 16px;
        padding: 20px;
        margin: 12px 0;
        background: #ffffff;
        box-shadow: 0 2px 8px rgba(60, 64, 67, 0.1);
        transition: all 0.2s ease;
        height: 100%;
        display: flex;
        flex-direction: column;
    }
    
    .news-card:hover {
        box-shadow: 0 4px 16px rgba(60, 64, 67, 0.15);
        border-color: #1976d2;
    }
    
    /* Header styling */
    .main-header {
        background: linear-gradient(135deg, #2196f3 0%, #21cbf3 100%);
        padding: 40px 20px;
        border-radius: 20px;
        color: white;
        text-align: center;
        margin-bottom: 40px;
    }
    
    /* Search section styling */
    .search-section {
        background: #f8f9fa;
        padding: 24px;
        border-radius: 16px;
        margin-bottom: 30px;
        border: 1px solid #e8eaed;
    }
    
    /* Button styling */
    .stButton > button {
        border-radius: 12px;
        border: none;
        background: #1976d2;
        color: white;
        font-weight: 500;
        padding: 6px 12px;
        transition: all 0.2s ease;
        font-size: 0.85em;
    }
    
    .stButton > button:hover {
        background: #1565c0;
        transform: translateY(-1px);
    }
    
    /* Category badges */
    .category-badge {
        background: #e3f2fd;
        color: #1976d2;
        padding: 4px 12px;
        border-radius: 16px;
        font-size: 0.85em;
        font-weight: 500;
        display: inline-block;
        margin: 2px 4px 2px 0;
    }
    
    /* Link styling */
    .news-title {
        color: #1976d2;
        text-decoration: none;
        font-weight: 600;
        font-size: 1.1em;
        margin-bottom: 8px;
        display: block;
        line-height: 1.3;
    }
    
    .news-title:hover {
        color: #1565c0;
        text-decoration: underline;
    }
    
    /* Stats cards */
    .stat-card {
        background: white;
        padding: 16px;
        border-radius: 12px;
        border: 1px solid #e8eaed;
        text-align: center;
        margin-bottom: 12px;
    }
    
    .stat-number {
        font-size: 1.5em;
        font-weight: 600;
        color: #1976d2;
    }
    
    .stat-label {
        font-size: 0.85em;
        color: #666;
        margin-top: 4px;
    }
    
    .logout-btn {
        position: fixed;
        top: 60px;
        right: 20px;
        z-index: 1000;
    }
</style>
"""

# Initialize recommendation models
@st.cache_resource
def initialize_recommendation_models():
    """Initialize and return the recommendation models"""
    models = {
        "S-BERT Model": SBertModel("S-BERT", "Semantic similarity using sentence transformers"),
        "TF-IDF Model": TfidfModel("TF-IDF", "Traditional TF-IDF based content similarity")
    }
    return models

# Load data and initialize models
@st.cache_data
def load_data():
    """Load the news data"""
    if not os.path.exists('../data/news_cleaned.csv'):
        with st.spinner("Processing news data... This may take a while."):
            preprocessing()
    
    news = pd.read_csv('../data/news_cleaned.csv', index_col=0)
    return news

# Setup models
@st.cache_resource
def setup_models(_models, _news):
    """Setup and train the recommendation models"""
    trained_models = {}
    
    for name, model in _models.items():
        try:
            with st.spinner(f"Training {name}..."):
                model.fit(_news)
                trained_models[name] = model
                logger.info(f"Successfully trained {name}")
        except Exception as e:
            logger.error(f"Error training {name}: {e}")
            trained_models[name] = model
    
    return trained_models

def get_base64_image(img_path):
    try:
        with open(img_path, "rb") as f:
            data = f.read()
        return base64.b64encode(data).decode()
    except FileNotFoundError:
        return None

def get_search_suggestions(query, news_df):
    """Get search suggestions based on query"""
    if not query or len(query) < 2:
        return []
    
    query = query.lower()
    suggestions = []
    
    title_matches = news_df[news_df['title'].str.lower().str.contains(query, na=False, regex=False)]['title'].tolist()
    suggestions.extend([f"📰 {title[:60]}..." if len(title) > 60 else f"📰 {title}" for title in title_matches[:3]])
    
    abstract_matches = news_df[news_df['abstract'].str.lower().str.contains(query, na=False, regex=False)]['title'].tolist()
    suggestions.extend([f"📄 {title[:60]}..." if len(title) > 60 else f"📄 {title}" for title in abstract_matches[:2]])
    
    return suggestions[:8]

def search_news(query, news_df):
    """Enhanced search news by title or abstract"""
    if not query:
        return news_df.sample(n=min(12, len(news_df)))
    
    clean_query = re.sub(r'^[📰📄]\s*', '', query).lower()
    query_terms = clean_query.split()
    masks = []
    
    for term in query_terms:
        term_mask = (
            news_df['title'].str.lower().str.contains(term, na=False, regex=False) |
            news_df['abstract'].str.lower().str.contains(term, na=False, regex=False)
        )
        masks.append(term_mask)
    
    if masks:
        combined_mask = masks[0]
        for mask in masks[1:]:
            combined_mask &= mask
    else:
        combined_mask = pd.Series([True] * len(news_df), index=news_df.index)
    
    return news_df[combined_mask]

def recommend_articles(news_title, top_n, model_name, trained_models, news_df):
    """Get recommended articles using the selected model"""
    try:
        if model_name in trained_models and trained_models[model_name].is_trained:
            return trained_models[model_name].recommend(news_title, top_n)
        else:
            return news_df.sample(n=min(top_n, len(news_df)))
    except Exception as e:
        logger.error(f"Error getting recommendations: {e}")
        return news_df.sample(n=min(top_n, len(news_df)))

def display_news_card(article, show_recommendations, news_index):
    """Display a news article in a card-like format"""
    if 'link' in article and pd.notna(article['link']):
        title_html = f'<a href="{article["link"]}" target="_blank" class="news-title">{article["title"]}</a>'
    else:
        title_html = f'<div class="news-title">{article["title"]}</div>'
    
    similarity_html = ""
    if isinstance(article, pd.Series) and 'similarity_score' in article and pd.notna(article['similarity_score']):
        score = float(article['similarity_score'])
        similarity_html = f'<div style="font-size: 0.8em; color: #666; margin-top: 8px;">📊 Similarity: {score:.3f}</div>'
    
    card_html = f"""
    <div class="news-card">
        {title_html}
        {similarity_html}
    </div>
    """
    
    st.markdown(card_html, unsafe_allow_html=True)
    
    col1, col2 = st.columns(2)
    
    with col1:
        if 'abstract' in article and pd.notna(article['abstract']):
            if st.button("📖 Abstract", key=f"abs_{news_index}", help="View abstract"):
                st.session_state[f"show_abstract_{news_index}"] = not st.session_state.get(f"show_abstract_{news_index}", False)
                st.rerun()
    
    with col2:
        if show_recommendations and news_index is not None:
            if st.button("🔍 Similar", key=f"rec_{news_index}", help="Find similar news"):
                st.session_state.show_recommendations = True
                st.session_state.selected_news_title = article['title']
                st.rerun()
    
    if st.session_state.get(f"show_abstract_{news_index}", False) and 'abstract' in article and pd.notna(article['abstract']):
        with st.container():
            abstract_text = str(article['abstract'])
            if len(abstract_text) > 400:
                st.markdown(f"**Abstract:** {abstract_text[:400]}...")
            else:
                st.markdown(f"**Abstract:** {abstract_text}")

def display_news_grid(news_df_display, show_recommendations=True, items_per_row=3):
    """Display news in a grid layout"""
    news_list = news_df_display.to_dict('records')
    
    for i in range(0, len(news_list), items_per_row):
        row_news = news_list[i:i + items_per_row]
        cols = st.columns(items_per_row)
        
        for j, article in enumerate(row_news):
            if j < len(row_news):
                with cols[j]:
                    display_news_card(article, show_recommendations, news_index=news_df_display.index[i + j])

def render_main_page(set_page):
    """Main rendering function for the news discovery page"""
    
    # Apply CSS
    st.markdown(MAIN_PAGE_CSS, unsafe_allow_html=True)
    
    # Load data and models
    news_df = load_data()
    recommendation_models = initialize_recommendation_models()
    trained_models = setup_models(recommendation_models, news_df)
    
    # Initialize session state
    if 'show_recommendations' not in st.session_state:
        st.session_state.show_recommendations = False
    if 'selected_news_title' not in st.session_state:
        st.session_state.selected_news_title = ''
    if 'search_query' not in st.session_state:
        st.session_state.search_query = ''
    if 'selected_model' not in st.session_state:
        st.session_state.selected_model = "S-BERT Model"
    
    # Logout button (top right)
    col1, col2 = st.columns([6, 1])
    with col2:
        if st.button("🚪 Logout", key="logout_btn"):
            st.session_state.authenticated = False
            st.session_state.auth_token = None
            st.session_state.login_user = None
            st.query_params.clear()
            set_page("login")
    
    # Header with image
    img_base64 = get_base64_image("./img/Albert.png")
    
    if img_base64:
        st.markdown(f"""
        <div class="main-header">
            <div style="display: flex; align-items: center; justify-content: center; gap: 30px;">
                <img src="data:image/png;base64,{img_base64}" 
                     style="height:200px; border-radius:50%;">
                <div style="text-align: left;">
                    <h1>📰 News Discovery Platform</h1>
                    <p style="font-size: 16px; opacity: 0.9;">
                        Discover and explore news articles with AI-powered recommendations
                    </p>
                    <p style="font-size: 25px; opacity: 1; margin: 5px 0;">
                        with Albert the NewsBro
                    </p>
                </div>
            </div>
        </div>
        """, unsafe_allow_html=True)
    else:
        st.markdown("""
        <div class="main-header">
            <h1>📰 News Discovery Platform</h1>
            <p style="font-size: 18px; opacity: 0.9;">
                Discover and explore news articles with AI-powered recommendations
            </p>
            <p style="font-size: 28px; opacity: 1; margin: 10px 0;">
                with Albert the NewsBro
            </p>
        </div>
        """, unsafe_allow_html=True)
    
    # Model selection
    st.markdown("### 🤖 Recommendation Model")
    available_models = list(trained_models.keys())
    col1, col2 = st.columns([2, 1])
    with col1:
        selected_model = st.selectbox(
            "Choose Recommendation Model",
            available_models,
            index=available_models.index(st.session_state.selected_model) if st.session_state.selected_model in available_models else 0,
            key="model_selector"
        )
        
        if selected_model != st.session_state.selected_model:
            st.session_state.selected_model = selected_model
            st.rerun()
    
    # Search section
    st.markdown('<div class="search-section">', unsafe_allow_html=True)
    search_col1, search_col2 = st.columns([4, 1])
    with search_col1:
        search_query = st.text_input(
            "Search",
            value=st.session_state.get('search_query', ''),
            placeholder="Search news articles, topics...",
            label_visibility="collapsed",
            key="main_search"
        )
    with search_col2:
        if st.button("🔍 Search", type="primary", use_container_width=True):
            st.session_state.search_query = search_query
            st.rerun()
    
    # Search suggestions
    if search_query and search_query != st.session_state.get('search_query', ''):
        suggestions = get_search_suggestions(search_query, news_df)
        if suggestions:
            st.markdown("**Suggestions:**")
            suggestion_cols = st.columns(min(len(suggestions), 4))
            for i, suggestion in enumerate(suggestions[:4]):
                with suggestion_cols[i % 4]:
                    if st.button(suggestion, key=f"sug_{i}", help="Click to search"):
                        st.session_state.search_query = suggestion
                        st.rerun()
    
    st.markdown('</div>', unsafe_allow_html=True)
    
    # Update session state
    if search_query != st.session_state.get('search_query', ''):
        st.session_state.search_query = search_query
    
    # Search results or featured news
    if st.session_state.search_query:
        with st.spinner("🔍 Searching..."):
            results = search_news(st.session_state.search_query, news_df)
        
        if len(results) == 0:
            st.info("🔍 No news articles found. Try different keywords.")
            st.markdown("### 🎲 You might be interested in:")
            sample_news = news_df.sample(n=6)
            display_news_grid(sample_news)
        else:
            st.markdown(f"### 📊 Found {len(results)} articles")
            display_results = results.head(12)
            display_news_grid(display_results, items_per_row=3)
    else:
        st.markdown("### ✨ Featured Articles")
        
        if 'random_seed' not in st.session_state:
            st.session_state.random_seed = np.random.randint(0, 1000)
        
        sample_news = news_df.sample(n=9, random_state=st.session_state.random_seed)
        display_news_grid(sample_news, items_per_row=3)
        
        col1, col2, col3 = st.columns([1, 1, 1])
        with col2:
            if st.button("🔄 Refresh", use_container_width=True):
                st.session_state.random_seed = np.random.randint(0, 1000)
                st.rerun()
    
    # Recommendations section
    if st.session_state.show_recommendations and st.session_state.selected_news_title:
        st.markdown("---")
        st.markdown("### 🎯 Similar Articles")
        st.info(f"📌 Based on: {st.session_state.selected_news_title[:80]}{'...' if len(st.session_state.selected_news_title) > 80 else ''}  \n🤖 Using: {st.session_state.selected_model}")
        
        with st.spinner("🤖 Finding similar articles..."):
            try:
                recommendations = recommend_articles(
                    st.session_state.selected_news_title, 
                    top_n=9, 
                    model_name=st.session_state.selected_model,
                    trained_models=trained_models,
                    news_df=news_df
                )
            except Exception as e:
                st.error(f"Error generating recommendations: {e}")
                recommendations = news_df.sample(n=9)
        
        display_news_grid(recommendations, show_recommendations=False, items_per_row=3)
        
        col1, col2, col3 = st.columns([1, 1, 1])
        with col2:
            if st.button("✖️ Close", use_container_width=True):
                st.session_state.show_recommendations = False
                st.rerun()
    
    # Sidebar
    with st.sidebar:
        st.markdown("### 📊 Statistics")
        st.markdown(f"""
        <div class="stat-card">
            <div class="stat-number">{len(news_df):,}</div>
            <div class="stat-label">Total Articles</div>
        </div>
        """, unsafe_allow_html=True)
        
        st.markdown("### 🤖 Model Details")
        if st.session_state.selected_model in trained_models:
            current_model = trained_models[st.session_state.selected_model]
            model_info = current_model.get_info()
            st.markdown(f"""
            **Active Model:** {model_info['name']}  
            **Status:** {'✅ Ready' if model_info['is_trained'] else '❌ Not Ready'}  
            **Description:** {model_info['description']}
            """)
        
        st.markdown("### 🚀 Quick Actions")
        if st.button("🎲 Random Articles", use_container_width=True):
            st.session_state.random_seed = np.random.randint(0, 1000)
            st.session_state.search_query = ''
            st.session_state.show_recommendations = False
            st.rerun()
        
        if st.button("🔄 Reset", use_container_width=True):
            st.session_state.search_query = ''
            st.session_state.show_recommendations = False
            st.rerun()
    
    # Footer
    st.markdown("---")
    st.markdown(f"""
    <div style="text-align: center; color: #666; padding: 20px;">
        <p>📰 <strong>News Discovery Platform</strong> • Powered by AI</p>
        <p style="font-size: 0.8em;">Current Model: {st.session_state.selected_model} | Articles: {len(news_df):,}</p>
    </div>
    """, unsafe_allow_html=True)


# If running standalone (for testing)
if __name__ == "__main__":
    st.set_page_config(page_title="News Discovery Platform", page_icon="📰", layout="wide")
    
    def dummy_set_page(page):
        st.session_state.page = page
    
    render_main_page(dummy_set_page)
