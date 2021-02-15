from datetime import datetime
from typing import List
from sqlalchemy import Boolean, Column, Integer, String, DateTime, Text, ForeignKey, Table
from .db import metadata
from sqlalchemy.orm import relationship
from sqlalchemy.ext.declarative import declared_attr

from .db import Base, engine


class IdMixin(object):
    """Миксин ID"""

    id = Column(Integer, primary_key=True, index=True)


class ModifiedAtMixin(object):
    """Миксин времени создания и времени изменения"""

    created_at = Column(DateTime, default=datetime.utcnow())
    updated_at = Column(DateTime, default=datetime.utcnow())


class User(Base, IdMixin, ModifiedAtMixin):
    """Модель пользователя"""

    __tablename__ = 'users'

    username = Column(String, unique=True)
    email = Column(String, unique=True)
    password = Column(String)
    is_active = Column(Boolean, default=True)

    profile = relationship('Profile', back_populates='user', uselist=False)


class Profile(Base, IdMixin):
    """Модель профиля пользователя"""

    __tablename__ = 'profiles'

    avatar = Column(String(length=255))
    activity = Column(String(length=255))
    bio = Column(String(length=500))
    social_networks = Column(String(length=255))
    user_id = Column(Integer, ForeignKey('users.id'))

    user = relationship('User', back_populates='profile')


metadata.create_all(engine)

'''
class Composition(Base, IdMixin, ModifiedAtMixin, PathToFileMixin):
    __abstract__ = True

    name = Column(String(length=255))
    genre = Column(String(length=255))
    description = Column(Text(length=500))
    author_id = Column(Integer, ForeignKey('users.id'))

    author = relationship('Profile', back_populates='composition')


class Art(Base, Composition):
    __tablename__ = 'arts'


class Audio(Base, Composition):
    __tablename__ = 'audios'

    cover = Column(String(length=255))
    producer = Column(String(length=255))
    featured = Column(String(length=255))
'''