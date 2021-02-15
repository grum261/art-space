import databases

from sqlalchemy import MetaData
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine


DATABASE_URL = 'postgresql://grum231:jeJiqNxu6nqC@localhost:5432/art_space'

database = databases.Database(DATABASE_URL)

metadata = MetaData()

engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})
metadata.bind(engine)

Base = declarative_base()
