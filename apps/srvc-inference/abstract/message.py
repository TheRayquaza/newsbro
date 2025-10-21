from pydantic import BaseModel
from datetime import datetime

class ArticleAggregate(BaseModel):
    id: int
    category: str
    subcategory: str
    title: str
    abstract: str
    link: str
    rss_link: str
    published_at: datetime
    is_active: bool

class FeedbackAggregate(BaseModel):
    user_id: int
    news_id: int
    value: int
    is_active: bool
    date: datetime

class RSSAggregate(BaseModel):
    id: int
    name: str
    description: str
    link: str
    active: bool

class InferenceCommand(BaseModel):
    user_id: int
    model_name: str
    article: ArticleAggregate
    date: datetime
