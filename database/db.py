import databases
import os

from datetime import datetime
from sqlalchemy import MetaData, Boolean, Column, Integer, String, DateTime, Text, ForeignKey, Table
from sqlalchemy.ext.declarative import declarative_base, DeferredReflection
from sqlalchemy import create_engine


DATABASE_URL = os.getenv('db_url')

database = databases.Database(DATABASE_URL)

metadata = MetaData()

engine = create_engine(DATABASE_URL)

Base = declarative_base(metadata=metadata)
