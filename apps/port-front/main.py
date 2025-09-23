import streamlit as st
import numpy as np
import pandas as pd
from preprocessing import preprocessing
import re
import os

# Import the new recommendation system
import sys
sys.path.append("../")

from apps.models.BaseRecommendationModel import create_default_engine, RecommendationEngine
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Configure page
st.set_page_config(
    page_title="Academic Paper Search",
    page_icon="🎓",
    layout="wide"
)

# Custom CSS for better styling
st.markdown("""
<style>
    /* Main container styling */
    .main {
        padding-top: 1rem;
    }
    
    /* Custom card styling */
    .paper-card {
        border: 1px solid #e8eaed;
        border-radius: 16px;
        padding: 20px;  /* Slightly reduced padding */
        margin: 12px 0;
        background: #ffffff;
        box-shadow: 0 2px 8px rgba(60, 64, 67, 0.1);
        transition: all 0.2s ease;
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    
    .paper-card:hover {
        box-shadow: 0 4px 16px rgba(60, 64, 67, 0.15);
        border-color: #1976d2;
    }
    
    /* Header styling */
    .main-header {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
        padding: 6px 12px;  /* Slightly smaller padding */
        transition: all 0.2s ease;
        font-size: 0.85em;  /* Slightly smaller font */
    }
    
    .stButton > button:hover {
        background: #1565c0;
        transform: translateY(-1px);
    }
    
    /* Secondary button styling */
    .secondary-btn {
        background: #f5f5f5 !important;
        color: #666 !important;
        border: 1px solid #ddd !important;
    }
    
    .secondary-btn:hover {
        background: #eeeeee !important;
        color: #333 !important;
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
    .paper-title {
        color: #1976d2;
        text-decoration: none;
        font-weight: 600;
        font-size: 1.1em;
        margin-bottom: 8px;
        display: block;
        line-height: 1.3;
    }
    
    .paper-title:hover {
        color: #1565c0;
        text-decoration: underline;
    }
    
    /* Author styling */
    .paper-authors {
        color: #666;
        font-size: 0.9em;
        margin-bottom: 12px;
    }
    
    /* Search suggestions */
    .search-suggestions {
        background: white;
        border: 1px solid #e8eaed;
        border-radius: 8px;
        margin-top: 4px;
        max-height: 200px;
        overflow-y: auto;
        box-shadow: 0 2px 8px rgba(60, 64, 67, 0.1);
    }
    
    .suggestion-item {
        padding: 8px 16px;
        border-bottom: 1px solid #f0f0f0;
        cursor: pointer;
        font-size: 0.9em;
    }
    
    .suggestion-item:hover {
        background: #f8f9fa;
    }
    
    .suggestion-item:last-child {
        border-bottom: none;
    }
    
    /* Remove extra spacing */
    .block-container {
        padding-top: 2rem;
    }
    
    /* Sidebar styling */
    .sidebar .sidebar-content {
        background: #f8f9fa;
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

    /* Model selector styling */
    .model-info {
        background: #f0f8ff;
        padding: 12px;
        border-radius: 8px;
        border: 1px solid #b3d9ff;
        margin: 8px 0;
        font-size: 0.85em;
    }
</style>
""", unsafe_allow_html=True)

# Initialize recommendation engine
@st.cache_resource
def initialize_recommendation_engine():
    """Initialize and return the recommendation engine"""
    return create_default_engine()

# Load data and initialize models
@st.cache_data
def load_data():
    recompute = False

    try:
        similarity_matrix = np.load('../data/similarity_matrix.npy', allow_pickle=True)
        if not isinstance(similarity_matrix, np.ndarray):
            raise ValueError("Loaded similarity matrix is not a NumPy array.")
    except (FileNotFoundError, ValueError):
        recompute = True
        similarity_matrix = None

    if not os.path.exists('../data/papers.csv'):
        recompute = True

    # Run preprocessing if needed
    if recompute:
        with st.spinner("Computing similarity matrix and papers data... This may take a while."):
            preprocessing()
        similarity_matrix = np.load('../data/similarity_matrix.npy', allow_pickle=True)

    # Load papers.csv after ensuring it exists
    papers = pd.read_csv('../data/papers.csv')

    return similarity_matrix, papers

@st.cache_resource
def setup_models(_engine, _papers, _similarity_matrix):
    """Setup and train the recommendation models"""
    try:
        # Train the similarity matrix model
        _engine.fit(_papers, model_name="Similarity Matrix", similarity_matrix=_similarity_matrix)
        
        # Try to train content-based model if available
        try:
            if "Content-Based TF-IDF" in _engine.get_available_models():
                with st.spinner("Training content-based model..."):
                    _engine.fit(_papers, model_name="Content-Based TF-IDF")
        except Exception as e:
            logger.warning(f"Could not train content-based model: {e}")
        
        return _engine
    except Exception as e:
        logger.error(f"Error setting up models: {e}")
        return _engine

# Load data and setup models
similarity_matrix, papers = load_data()
recommendation_engine = initialize_recommendation_engine()
recommendation_engine = setup_models(recommendation_engine, papers, similarity_matrix)

def get_search_suggestions(query):
    """Get search suggestions based on query"""
    if not query or len(query) < 2:
        return []
    
    query = query.lower()
    suggestions = []
    
    # Search in titles
    title_matches = papers[papers['Title'].str.lower().str.contains(query, na=False, regex=False)]['Title'].tolist()
    suggestions.extend([f"📄 {title[:60]}..." if len(title) > 60 else f"📄 {title}" for title in title_matches[:3]])
    
    # Search in authors
    author_matches = papers[papers['Author'].str.lower().str.contains(query, na=False, regex=False)]['Author'].tolist()
    unique_authors = []
    for authors in author_matches:
        if pd.notna(authors):
            author_list = re.split(r'[,;&]', str(authors))
            for author in author_list:
                author = author.strip()
                if query in author.lower() and author not in unique_authors and len(unique_authors) < 3:
                    unique_authors.append(author)
    suggestions.extend([f"👤 {author}" for author in unique_authors])
    
    # Search in categories
    if 'Category' in papers.columns:
        cat_matches = papers[papers['Category'].str.lower().str.contains(query, na=False, regex=False)]['Category'].unique().tolist()
        suggestions.extend([f"🏷️ {cat}" for cat in cat_matches[:2]])
    
    return suggestions[:8]  # Limit to 8 suggestions

def search_papers(query, papers_df):
    """Enhanced search papers by title, author, or category"""
    if not query:
        return papers_df.sample(n=min(12, len(papers_df)))  # Show random sample
    
    # Clean the query (remove emoji prefixes from suggestions)
    clean_query = re.sub(r'^[📄👤🏷️]\s*', '', query).lower()
    
    # Split query into individual terms for better matching
    query_terms = clean_query.split()
    
    # Create masks for different fields
    masks = []
    
    for term in query_terms:
        term_mask = (
            papers_df['Title'].str.lower().str.contains(term, na=False, regex=False) |
            papers_df['Author'].str.lower().str.contains(term, na=False, regex=False) |
            papers_df['Category'].str.lower().str.contains(term, na=False, regex=False) |
            papers_df['Category'].str.lower().str.contains(term, na=False, regex=False)
        )
        
        # Also search in abstract if available
        if 'Abstract' in papers_df.columns:
            term_mask |= papers_df['Abstract'].str.lower().str.contains(term, na=False, regex=False)
        
        masks.append(term_mask)
    
    # Combine masks (AND operation for multiple terms)
    if masks:
        combined_mask = masks[0]
        for mask in masks[1:]:
            combined_mask &= mask
    else:
        combined_mask = pd.Series([True] * len(papers_df), index=papers_df.index)
    
    return papers_df[combined_mask]

def recommend_articles(paper_index, top_n=6, model_name=None):
    """Get recommended articles using the recommendation engine"""
    try:
        if model_name:
            return recommendation_engine.recommend(paper_index, top_n, model_name=model_name)
        else:
            return recommendation_engine.recommend(paper_index, top_n)
    except Exception as e:
        logger.error(f"Error getting recommendations: {e}")
        # Fallback to showing random papers
        return papers.sample(n=min(top_n, len(papers)))

def format_authors(authors):
    """Format author names nicely"""
    if pd.isna(authors):
        return "Unknown"
    # Split by common separators and take first few authors
    author_list = re.split(r'[,;&]', str(authors))
    if len(author_list) > 2:
        return f"{', '.join(author_list[:2])} et al."
    return ', '.join(author_list)

def display_paper_grid(papers_df, show_recommendations=True, items_per_row=3):
    """Display papers in a grid layout"""
    papers_list = papers_df.to_dict('records')
    
    # Group papers by rows
    for i in range(0, len(papers_list), items_per_row):
        row_papers = papers_list[i:i + items_per_row]
        cols = st.columns(items_per_row)
        
        for j, paper in enumerate(row_papers):
            if j < len(row_papers):
                with cols[j]:
                    display_paper_card(paper, show_recommendations, paper_index=papers_df.index[i + j])

def display_paper_card(paper, show_recommendations=True, paper_index=None):
    """Display a paper in a card-like format"""
    
    # Paper title with link if available
    if 'Link' in paper and pd.notna(paper['Link']):
        title_html = f'<a href="{paper["Link"]}" target="_blank" class="paper-title">{paper["Title"]}</a>'
    else:
        title_html = f'<div class="paper-title">{paper["Title"]}</div>'
    
    # Authors
    authors_html = f'<div class="paper-authors">👤 {format_authors(paper["Author"])}</div>'
    
    # Categories
    categories_html = ""
    if 'Category' in paper and pd.notna(paper['Category']):
        categories_html = f'<span class="category-badge">{paper["Category"]}</span>'
    
    # Show similarity score if available
    similarity_html = ""
    if isinstance(paper, pd.Series) and 'similarity_score' in paper and pd.notna(paper['similarity_score']):
        score = float(paper['similarity_score'])
        similarity_html = f'<div style="font-size: 0.8em; color: #666; margin-top: 8px;">📊 Similarity: {score:.3f}</div>'
    
    card_html = f"""
    <div class="paper-card">
        {title_html}
        {authors_html}
        {categories_html}
        {similarity_html}
    </div>
    """
    
    st.markdown(card_html, unsafe_allow_html=True)
    
    # Action buttons
    col1, col2 = st.columns(2)
    
    with col1:
        if 'Abstract' in paper and pd.notna(paper['Abstract']):
            if st.button("📖 Abstract", key=f"abs_{paper_index}", help="View abstract"):
                st.session_state[f"show_abstract_{paper_index}"] = not st.session_state.get(f"show_abstract_{paper_index}", False)
                st.rerun()
    
    with col2:
        if show_recommendations and paper_index is not None:
            if st.button("🔍 Similar", key=f"rec_{paper_index}", help="Find similar papers"):
                st.session_state.show_recommendations = True
                st.session_state.selected_paper_idx = paper_index
                st.rerun()
    
    # Show abstract if toggled
    if st.session_state.get(f"show_abstract_{paper_index}", False) and 'Abstract' in paper and pd.notna(paper['Abstract']):
        with st.container():
            st.markdown(f"**Abstract:** {paper['Abstract'][:400]}..." if len(str(paper['Abstract'])) > 400 else f"**Abstract:** {paper['Abstract']}")

# Initialize session state
if 'show_recommendations' not in st.session_state:
    st.session_state.show_recommendations = False
if 'selected_paper_idx' not in st.session_state:
    st.session_state.selected_paper_idx = 0
if 'search_query' not in st.session_state:
    st.session_state.search_query = ''
if 'selected_model' not in st.session_state:
    st.session_state.selected_model = recommendation_engine.active_model.name if recommendation_engine.active_model else "Similarity Matrix"

# Header
import base64

def get_base64_image(img_path):
    try:
        with open(img_path, "rb") as f:
            data = f.read()
        return base64.b64encode(data).decode()
    except FileNotFoundError:
        return None

img_base64 = get_base64_image("./img/Albert.png")

# Header with conditional image
if img_base64:
    st.markdown(f"""
    <div class="main-header">
        <div style="display: flex; align-items: center; justify-content: center; gap: 30px;">
            <img src="data:image/png;base64,{img_base64}" 
                 style="height:200px; border-radius:50%;">
            <div style="text-align: left;">
                <h1>🎓 Academic Paper Discovery Platform</h1>
                <p style="font-size: 16px; opacity: 0.9;">
                    Discover and explore research papers with AI-powered recommendations
                </p>
                <p style="font-size: 25px; opacity: 1; margin: 5px 0;">
                    with Albert the ScientificPapersBro
                </p>
            </div>
        </div>
    </div>
    """, unsafe_allow_html=True)
else:
    st.markdown("""
    <div class="main-header">
        <h1>🎓 Academic Paper Discovery Platform</h1>
        <p style="font-size: 18px; opacity: 0.9;">
            Discover and explore research papers with AI-powered recommendations
        </p>
        <p style="font-size: 28px; opacity: 1; margin: 10px 0;">
            with Albert the ScientificPapersBro
        </p>
    </div>
    """, unsafe_allow_html=True)

# Model selection section
st.markdown("### 🤖 Recommendation Model")
available_models = recommendation_engine.get_available_models()

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
        recommendation_engine.set_active_model(selected_model)
        st.rerun()

# Search section
st.markdown('<div class="search-section">', unsafe_allow_html=True)

# Search input with suggestions
search_col1, search_col2 = st.columns([4, 1])

with search_col1:
    search_query = st.text_input(
        "Search",
        value=st.session_state.get('search_query', ''),
        placeholder="Search papers, authors, or topics...",
        label_visibility="collapsed",
        key="main_search"
    )

with search_col2:
    if st.button("🔍 Search", type="primary", use_container_width=True):
        st.session_state.search_query = search_query
        st.rerun()

# Show search suggestions
if search_query and search_query != st.session_state.get('search_query', ''):
    suggestions = get_search_suggestions(search_query)
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

# Search results or featured papers
if st.session_state.search_query:
    with st.spinner("🔍 Searching..."):
        results = search_papers(st.session_state.search_query, papers)
    
    if len(results) == 0:
        st.info("🔍 No papers found. Try different keywords.")
        # Show random sample as fallback
        st.markdown("### 🎲 You might be interested in:")
        sample_papers = papers.sample(n=6)
        display_paper_grid(sample_papers)
    else:
        st.markdown(f"### 📊 Found {len(results)} papers")
        
        # Show first 12 results
        display_results = results.head(12)
        display_paper_grid(display_results, items_per_row=3)
        
        # Show more button if there are more results
        if len(results) > 12:
            col1, col2, col3 = st.columns([1, 1, 1])
            with col2:
                if st.button(f"📄 Show {min(12, len(results) - 12)} more papers", use_container_width=True):
                    # This would typically implement pagination
                    pass

else:
    # Show featured papers when no search
    st.markdown("### ✨ Featured Papers")
    
    # Show different samples each time
    if 'random_seed' not in st.session_state:
        st.session_state.random_seed = np.random.randint(0, 1000)
    
    sample_papers = papers.sample(n=9, random_state=st.session_state.random_seed)
    display_paper_grid(sample_papers, items_per_row=3)
    
    col1, col2, col3 = st.columns([1, 1, 1])
    with col2:
        if st.button("🔄 Refresh", use_container_width=True):
            st.session_state.random_seed = np.random.randint(0, 1000)
            st.rerun()

# Recommendations section
if st.session_state.show_recommendations:
    st.markdown("---")
    selected_paper = papers.iloc[st.session_state.selected_paper_idx]
    
    st.markdown("### 🎯 Similar Papers")
    
    # Show which model is being used
    active_model_name = recommendation_engine.active_model.name
    st.info(f"📌 Based on: **{selected_paper['Title'][:80]}{'...' if len(selected_paper['Title']) > 80 else ''}**  \n🤖 Using: **{active_model_name}**")
    
    with st.spinner("🤖 Finding similar papers..."):
        try:
            recommendations = recommend_articles(st.session_state.selected_paper_idx, top_n=9)
        except Exception as e:
            st.error(f"Error generating recommendations: {e}")
            recommendations = papers.sample(n=9)  # Fallback
    
    display_paper_grid(recommendations, show_recommendations=False, items_per_row=3)    
    
    col1, col2, col3 = st.columns([1, 1, 1])
    with col2:
        if st.button("✖️ Close", use_container_width=True):
            st.session_state.show_recommendations = False
            st.rerun()

# Enhanced sidebar
with st.sidebar:
    st.markdown("### 📊 Statistics")
    
    st.markdown(f"""
    <div class="stat-card">
        <div class="stat-number">{len(papers):,}</div>
        <div class="stat-label">Total Papers</div>
    </div>
    """, unsafe_allow_html=True)
    
    if 'Category' in papers.columns:
        unique_categories = papers['Category'].nunique()
        st.markdown(f"""
        <div class="stat-card">
            <div class="stat-number">{unique_categories}</div>
            <div class="stat-label">Research Areas</div>
        </div>
        """, unsafe_allow_html=True)
    
    st.markdown("### 🎛️ Filters")
    
    # Simple category filter
    if 'Category' in papers.columns:
        categories = ['All Categories'] + sorted(papers['Category'].dropna().unique().tolist()[:10])  # Limit to top 10
        selected_category = st.selectbox("Research Area", categories)
    
    st.markdown("### 🤖 Model Details")
    
    # Show current model info
    current_model_info = recommendation_engine.get_model_info(st.session_state.selected_model)
    st.markdown(f"""
    **Active Model:** {current_model_info['name']}  
    **Status:** {'✅ Ready' if current_model_info['is_trained'] else '❌ Not Ready'}  
    **Description:** {current_model_info['description']}
    """)
    
    st.markdown("### 🚀 Quick Actions")
    
    if st.button("🎲 Random Papers", use_container_width=True):
        st.session_state.random_seed = np.random.randint(0, 1000)
        st.session_state.search_query = ''
        st.session_state.show_recommendations = False
        st.rerun()
    
    if st.button("🔄 Reset", use_container_width=True):
        st.session_state.search_query = ''
        st.session_state.show_recommendations = False
        st.rerun()

# Clean footer
st.markdown("---")
st.markdown(f"""
<div style="text-align: center; color: #666; padding: 20px;">
    <p>🎓 <strong>Academic Paper Discovery Platform</strong> • Powered by AI</p>
    <p style="font-size: 0.8em;">Current Model: {st.session_state.selected_model} | Papers: {len(papers):,}</p>
</div>
""", unsafe_allow_html=True)
